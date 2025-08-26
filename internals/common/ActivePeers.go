package common

import (
	"encoding/binary"

	"github.com/libp2p/go-libp2p/core/network"
)

type ActivePeers struct {
	peers map[string]network.Stream
}

func NewActivePeers() ActivePeers {
	return ActivePeers{
		peers: make(map[string]network.Stream),
	}
}

type Message struct {
	FromPeerId string
	Text       string
}

func startReading(s network.Stream, messageChan *(chan Message)) {
	go func() {
		for {
			var message string

			err := binary.Read(s, binary.BigEndian, &message)
			if err != nil {
				panic(err)
			}

			*messageChan <- Message{FromPeerId: s.ID(), Text: message}
		}
	}()
}

func (a *ActivePeers) AddPeer(peerId string, s network.Stream, messageChan *(chan Message)) {
	a.peers[peerId] = s

	startReading(s, messageChan)
}

func (a ActivePeers) GetPeers() map[string]network.Stream {
	return a.peers
}

func (a ActivePeers) GetPeerById(peerId string) (network.Stream, bool) {
	stream, exists := a.peers[peerId]

	return stream, exists
}
