// main.go - Reunion server using the Katzenpost mix server cbor plugin system.
// Copyright (C) 2020  David Stainton.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/katzenpost/reunion/commands"
	"github.com/katzenpost/reunion/epochtime/katzenpost"
	"github.com/katzenpost/reunion/server"
	"github.com/katzenpost/server/cborplugin"
	"github.com/ugorji/go/codec"
	"gopkg.in/op/go-logging.v1"
)

func parametersHandler(response http.ResponseWriter, req *http.Request) {
	params := new(cborplugin.Parameters)
	var serialized []byte
	enc := codec.NewEncoderBytes(&serialized, new(codec.CborHandle))
	if err := enc.Encode(params); err != nil {
		panic(err)
	}
	_, err := response.Write(serialized)
	if err != nil {
		panic(err)
	}
}

func requestHandler(log *logging.Logger, server *server.Server, response http.ResponseWriter, req *http.Request) {
	cborHandle := new(codec.CborHandle)
	request := cborplugin.Request{
		Payload: make([]byte, 0),
	}
	err := codec.NewDecoder(req.Body, new(codec.CborHandle)).Decode(&request)
	if err != nil {
		log.Errorf("query command must be of type cborplugin.Request: %s", err.Error())
		return
	}
	cmd, err := commands.FromBytes(request.Payload)
	if err != nil {
		log.Errorf("invalid Reunion query command found in request Payload len %d: %s", len(request.Payload), err.Error())
		return
	}
	replyCmd, err := server.ProcessQuery(cmd)
	if err != nil {
		log.Errorf("reunion HTTP server invalid reply command: %s", err.Error())
		return
	}
	rawReply := replyCmd.ToBytes()
	reply := cborplugin.Response{
		Payload: rawReply,
	}
	var serialized []byte
	enc := codec.NewEncoderBytes(&serialized, cborHandle)
	if err := enc.Encode(reply); err != nil {
		log.Error(err.Error())
		return
	}
	log.Debugf("serialized response is len %d", len(serialized))
	_, err = response.Write(serialized)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Debug("sent response")
}

func main() {
	var logLevel string
	logPath := flag.String("log", "", "Log file path. Default STDOUT.")
	flag.StringVar(&logLevel, "log_level", "DEBUG", "logging level could be set to: DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL")
	stateFilePath := flag.String("s", "statefile", "State file path.")
	epochClockName := flag.String("epochClock", "katzenpost", "The epoch-clock to use.")
	flag.Parse()

	if *epochClockName != "katzenpost" {
		panic("Thus far only the Katzenpost epoch clock is supported in this server implementation.")
	}

	// start service
	httpServer := http.Server{}

	tmpDir, err := ioutil.TempDir("", "reunion_server")
	if err != nil {
		panic(err)
	}
	socketFile := filepath.Join(tmpDir, fmt.Sprintf("%d.reunion.socket", os.Getpid()))

	unixListener, err := net.Listen("unix", socketFile)
	if err != nil {
		panic(err)
	}

	// XXX make a reunionServer.
	clock := new(katzenpost.Clock)
	reunionServer, err := server.NewServer(clock, *stateFilePath, *logPath, logLevel)
	if err != nil {
		panic(err)
	}
	httpLog := reunionServer.GetNewLogger("reunion_http_server")

	httpLog.Debug("Starting up...")
	_requestHandler := func(response http.ResponseWriter, request *http.Request) {
		requestHandler(httpLog, reunionServer, response, request)
	}
	http.HandleFunc("/request", _requestHandler)
	http.HandleFunc("/parameters", parametersHandler)
	fmt.Printf("%s\n", socketFile)
	httpServer.Serve(unixListener)
	os.Remove(socketFile)
}
