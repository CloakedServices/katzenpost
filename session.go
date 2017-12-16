// session.go - mixnet session client
// Copyright (C) 2017  Yawning Angel, Ruben Pollan, David Stainton
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

// Package client provides the Katzenpost midclient
package client

import (
	"errors"
	"fmt"

	"github.com/katzenpost/core/crypto/ecdh"
	"github.com/katzenpost/core/crypto/rand"
	"github.com/katzenpost/core/log"
	"github.com/katzenpost/core/sphinx/constants"
	"github.com/katzenpost/minclient"
	"github.com/katzenpost/minclient/block"
	"github.com/op/go-logging"
)

// UserKeyDiscovery interface for user key discovery
type UserKeyDiscovery interface {
	Get(identity string) (*ecdh.PublicKey, error)
}

// IngressBlock is used for storing decrypted
// blocked received from remote clients.
type IngressBlock struct {
	SenderPubKey *ecdh.PublicKey
	Block        *block.Block
}

// ToBytes serializes an IngressBlock into bytes
func (i *IngressBlock) ToBytes() ([]byte, error) {
	raw := []byte{}
	rawSenderPubKey, err := i.SenderPubKey.MarshalBinary()
	if err != nil {
		return nil, err
	}
	raw = append(raw, rawSenderPubKey...)
	rawBlock, err := i.Block.ToBytes()
	if err != nil {
		return nil, err
	}
	raw = append(raw, rawBlock...)
	return raw, nil
}

func (i *IngressBlock) FromBytes(raw []byte) error {
	pubKey := new(ecdh.PublicKey)
	err := pubKey.FromBytes(raw[:ecdh.PublicKeySize])
	if err != nil {
		return err
	}
	i.SenderPubKey = pubKey
	i.Block = new(block.Block)
	return i.Block.FromBytes(raw[ecdh.PublicKeySize:])
}

// Storage is an interface user for persisting
// ARQ and fragmentation/reassembly state
type Storage interface {
	GetBlocks(*[block.MessageIDLength]byte) ([][]byte, error)
	PutBlock(*[block.MessageIDLength]byte, []byte) error
}

// MessageConsumer is an interface used for
// processing received messages
type MessageConsumer interface {
	ReceivedMessage(senderPubKey *ecdh.PublicKey, message []byte)
	ReceivedACK(messageID *[block.MessageIDLength]byte, message []byte)
}

// SessionConfig is specifies the configuration for a new session
type SessionConfig struct {
	User             string
	Provider         string
	IdentityPrivKey  *ecdh.PrivateKey
	LinkPrivKey      *ecdh.PrivateKey
	MessageConsumer  MessageConsumer
	Storage          Storage
	UserKeyDiscovery UserKeyDiscovery
}

// Session holds the client session
type Session struct {
	cfg             *SessionConfig
	client          *minclient.Client
	queue           chan string
	log             *logging.Logger
	logBackend      *log.Backend
	messageConsumer MessageConsumer
	connected       chan bool
	identityPrivKey *ecdh.PrivateKey
}

// NewSession stablishes a session with provider using key.
// This method will block until session is connected to the Provider.
// This method takes the following arguments:
// user: the username of the account
// provider: the Provider name indicates which Provider the user account is on
// identityKeyPriv: the private messaging key for end to end message exchanges with other users
// linkKeyPriv: the private link layer key for our noise wire protocol
// consumer: the message consumer consumes received messages
func (c *Client) NewSession(cfg *SessionConfig) (*Session, error) {
	var err error
	session := new(Session)
	clientCfg := &minclient.ClientConfig{
		User:        cfg.User,
		Provider:    cfg.Provider,
		LinkKey:     cfg.LinkPrivKey,
		LogBackend:  c.logBackend,
		PKIClient:   c.cfg.PKIClient,
		OnConnFn:    session.onConnection,
		OnMessageFn: session.onMessage,
		OnACKFn:     session.onACK,
	}
	session.identityPrivKey = cfg.IdentityPrivKey
	session.connected = make(chan bool, 0)
	session.messageConsumer = cfg.MessageConsumer
	session.log = c.logBackend.GetLogger(fmt.Sprintf("%s@%s_session", cfg.User, cfg.Provider))
	session.client, err = minclient.New(clientCfg)
	if err != nil {
		return nil, err
	}
	err = session.waitForConnection()
	if err != nil {
		return nil, err
	}
	return session, nil
}

// Shutdown the session
func (s *Session) Shutdown() {
	s.client.Shutdown()
}

// waitForConnection blocks until the client is
// connected to the Provider
func (s *Session) waitForConnection() error {
	isConnected := <-s.connected
	if !isConnected {
		return errors.New("status is not connected even with status change")
	}
	return nil
}

// Send reliably delivers the message to the recipient's queue
// on the destination provider or returns an error
func (s *Session) Send(recipient, provider string, message []byte) (*[block.MessageIDLength]byte, error) {
	s.log.Debugf("Send")
	return nil, errors.New("failure: Send is not yet implemented")
}

// SendUnreliable unreliably sends a message to the recipient's queue
// on the destination provider or returns an error
func (s *Session) SendUnreliable(recipient, provider string, message []byte) error {
	s.log.Debugf("SendUnreliable")
	messageID := [block.MessageIDLength]byte{}
	_, err := rand.Reader.Read(messageID[:])
	if err != nil {
		return err
	}
	recipientPubKey, err := s.cfg.UserKeyDiscovery.Get(recipient)
	if err != nil {
		return err
	}
	blocks, err := block.EncryptMessage(&messageID, message, s.identityPrivKey, recipientPubKey)
	if err != nil {
		return err
	}
	for _, block := range blocks {
		err = s.client.SendUnreliableCiphertext(recipient, provider, block)
		if err != nil {
			break
		}
	}
	return err
}

// OnConnection will be called by the minclient api
// upon connecting to the Provider
func (s *Session) onConnection(isConnected bool) {
	s.log.Debugf("OnConnection")
	s.connected <- isConnected
}

// OnMessage will be called by the minclient api
// upon receiving a message
func (s *Session) onMessage(ciphertextBlock []byte) error {
	s.log.Debugf("OnMessage")
	rBlock, senderPubKey, err := block.DecryptBlock(ciphertextBlock, s.identityPrivKey)
	if err != nil {
		return err
	}
	if rBlock.TotalBlocks == 1 {
		s.messageConsumer.ReceivedMessage(senderPubKey, rBlock.Payload)
		return nil
	}
	ingressBlock := IngressBlock{
		SenderPubKey: senderPubKey,
		Block:        rBlock,
	}
	rawStoredBlocks, err := s.cfg.Storage.GetBlocks(&rBlock.MessageID)
	if err != nil {
		return err
	}
	rawBlock, err := ingressBlock.ToBytes()
	if err != nil {
		return err
	}
	rawBlocks := append(rawStoredBlocks, rawBlock)
	ingressBlocks := make([]*IngressBlock, len(rawBlocks))
	for i, b := range rawBlocks {
		ingressBlock := &IngressBlock{}
		err := ingressBlock.FromBytes(b)
		if err != nil {
			return err
		}
		ingressBlocks[i] = ingressBlock
	}
	message, err := reassemble(ingressBlocks)
	if err != nil {
		err = s.cfg.Storage.PutBlock(&ingressBlock.Block.MessageID, rawBlock)
		if err != nil {
			return err
		}
	}
	s.messageConsumer.ReceivedMessage(senderPubKey, message)
	return nil
}

// OnACK is called by the minclient api whe
// we receive an ACK message
func (s *Session) onACK(surbid *[constants.SURBIDLength]byte, message []byte) error {
	s.log.Debugf("OnACK")
	return nil
}
