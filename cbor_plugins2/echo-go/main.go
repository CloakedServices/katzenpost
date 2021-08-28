// main.go - echo service using cbor plugin system
// Copyright (C) 2018  David Stainton.
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
	"os"
	"path/filepath"

	"github.com/katzenpost/core/log"

	"github.com/fxamacker/cbor/v2"
	corecborplugin "github.com/katzenpost/core/cborplugin"
	"golang.org/x/sync/errgroup"
)

type Payload struct {
	Payload []byte
}

func (p *Payload) Marshal() ([]byte, error) {
	return cbor.Marshal(p)
}

func (p *Payload) Unmarshal(b []byte) error {
	return cbor.Unmarshal(b, p)
}

type payloadFactory struct{}

func (p *payloadFactory) Build() corecborplugin.Command {
	return new(Payload)
}

type Echo struct{}

func (e *Echo) OnCommand(cmd corecborplugin.Command) (corecborplugin.Command, error) {
	return cmd, nil
}

func (e *Echo) RegisterConsumer(s *corecborplugin.Server) {
	// noop
}

func main() {
	var logLevel string
	var logDir string
	flag.StringVar(&logDir, "log_dir", "", "logging directory")
	flag.StringVar(&logLevel, "log_level", "DEBUG", "logging level could be set to: DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL")
	flag.Parse()

	// Ensure that the log directory exists.
	s, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		fmt.Printf("Log directory '%s' doesn't exist.", logDir)
		os.Exit(1)
	}
	if !s.IsDir() {
		fmt.Println("Log directory must actually be a directory.")
		os.Exit(1)
	}

	logBackend, err := log.New("", "DEBUG", false)
	serverLog := logBackend.GetLogger("server")

	// start service
	tmpDir, err := ioutil.TempDir("", "echo_server")
	if err != nil {
		panic(err)
	}
	socketFile := filepath.Join(tmpDir, fmt.Sprintf("%d.echo.socket", os.Getpid()))
	commandFactory := new(payloadFactory)
	echo := new(Echo)

	var server *corecborplugin.Server
	g := new(errgroup.Group)
	g.Go(func() error {
		server = corecborplugin.NewServer(serverLog, socketFile, commandFactory, echo)
		return nil

	})
	err = g.Wait()
	fmt.Printf("%s\n", socketFile)
	server.Accept()
	server.Wait()
	os.Remove(socketFile)
}
