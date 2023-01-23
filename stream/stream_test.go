// stream_test.go - map service stream tests
// Copyright (C) 2022  Masala
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
//go:build docker_test
// +build docker_test

package client

import (
	"encoding/base64"
	"github.com/katzenpost/katzenpost/client"
	"github.com/katzenpost/katzenpost/client/config"
	"github.com/katzenpost/katzenpost/core/crypto/rand"
	mClient "github.com/katzenpost/katzenpost/map/client"
	"github.com/stretchr/testify/require"
	"io"
	"sync"
	"testing"
	"time"
)

// newStreams returns an initialized pair of Streams
func newStreams() (*Stream, *Stream) {
	a := new(Stream)
	b := new(Stream)

	// allocate memory for keys
	a.writekey = &[keySize]byte{}
	a.readkey = &[keySize]byte{}
	b.writekey = &[keySize]byte{}
	b.readkey = &[keySize]byte{}
	asecret := &[keySize]byte{}
	bsecret := &[keySize]byte{}

	// initialize handshake secrets from random
	io.ReadFull(rand.Reader, asecret[:])
	io.ReadFull(rand.Reader, bsecret[:])

	// Stream keys should now be initialized
	a.exchange(asecret[:], bsecret[:])
	b.exchange(bsecret[:], asecret[:])
	return a, b
}

func TestFrameKey(t *testing.T) {
	require := require.New(t)

	// the same key should be returned for every idx
	a, b := newStreams()
	for i := 0; i < 4096; i++ {
		i := uint64(i)
		// require sender/receiver frame ID match
		require.Equal(a.rxFrameID(i), b.txFrameID(i))
		require.Equal(a.txFrameID(i), b.rxFrameID(i))

		// require sender/receiver frame keys match
		require.Equal(a.rxFrameKey(i), b.txFrameKey(i))
		require.Equal(a.txFrameKey(i), b.rxFrameKey(i))
	}
	a.Halt()
	b.Halt()
}

func TestCreateStream(t *testing.T) {
	require := require.New(t)

	cfg, err := config.LoadFile("testdata/client.toml")
	require.NoError(err)

	cc, err := client.New(cfg)
	require.NoError(err)
	require.NotNil(cc)

	session, err := cc.NewTOFUSession()
	require.NoError(err)
	require.NotNil(session)
	session.WaitForDocument()

	c, err := mClient.NewClient(session)
	require.NoError(err)
	require.NotNil(c)

	asecret := &[keySize]byte{}
	bsecret := &[keySize]byte{}

	// initialize handshake secrets from random
	io.ReadFull(rand.Reader, asecret[:])
	io.ReadFull(rand.Reader, bsecret[:])

	// our view of stream
	s := NewStream(c, asecret[:], bsecret[:])
	// "other end" of stream
	r := NewStream(c, bsecret[:], asecret[:])

	msg := []byte("Hello World")
	n, err := s.Write(msg)
	require.NoError(err)
	require.Equal(n, len(msg))

	yolo := make([]byte, len(msg))
	for {
		// XXX: the tricky part is that we don't have a convenience method that will handle spinning on Read() for us and
		// ReadAtLeast payload
		// I thought io.ReadAtLeast would do this, but we get EOF too soon
		// because we are just proxying the calls through bytes.Buffer and whatever it does
		n, err = r.Read(yolo)
		if n == len(msg) {
			break
		}
	}
	require.NoError(err)
	require.Equal(n, len(msg))
	require.Equal(yolo, msg)

	msg = []byte("Goodbye World")
	n, err = s.Write(msg)
	require.NoError(err)
	require.Equal(n, len(msg))

	yolo = make([]byte, len(msg))
	n, err = io.ReadAtLeast(r, yolo, len(msg))
	require.NoError(err)
	require.Equal(n, len(msg))
	require.Equal(yolo, msg)
	s.Halt()
	r.Halt()

	// initialize handshake secrets from random
	io.ReadFull(rand.Reader, asecret[:])
	io.ReadFull(rand.Reader, bsecret[:])

	wg := new(sync.WaitGroup)
	wg.Add(2)
	sidechannel := make(chan string, 0)
	// worker A
	go func() {
		s := NewStream(c, asecret[:], bsecret[:])
		for i := 0; i < 4; i++ {
			entropic := make([]byte, 4242) // ensures fragmentation
			io.ReadFull(rand.Reader, entropic)
			message := base64.StdEncoding.EncodeToString(entropic)
			// tell the other worker what message we're going to try and send
			t.Logf("Sending %d bytes", len(message))
			sidechannel <- message
			s.Write([]byte(message))
		}
		close(sidechannel)
		t.Logf("SendWorker Done()")
		wg.Done()
		// wait until reader has finished reading
		// before halting the writer()
		wg.Wait()
		s.Halt()
	}()

	// worker B
	go func() {
		s := NewStream(c, bsecret[:], asecret[:])
		for {
			msg, ok := <-sidechannel
			// channel was closed by writer, we're done
			if !ok {
				t.Logf("ReadWorker Done()")
				wg.Done()
				s.Halt()
				return
			}
			b := make([]byte, len(msg))
			// Read() data until we have received the message
			for readOff := 0; readOff < len(msg); {
				n, err := s.Read(b[readOff:])
				if err != nil {
					panic(err)
				}
				t.Logf("Read %d bytes", n)

				readOff += n
				if n == 0 {
					// XXX retry a sensible time later, like the average round trip time
					<-time.After(time.Second * 2)
				}
			}
			t.Logf("Read total %d", len(b))
			require.Equal([]byte(msg), b)
		}
	}()
	wg.Wait()
}