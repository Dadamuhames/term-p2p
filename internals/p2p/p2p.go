package p2p

import (
	"encoding/binary"
	"fmt"
	"term-p2p/internals/config"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

func (n *discoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	n.PeerChan <- peerInfo
}

func InitDiscoveryServer(node host.Host) *(chan peerstore.AddrInfo) {
	notifee := discoveryNotifee{}
	notifee.PeerChan = make(chan peer.AddrInfo)

	// look for peers
	discoveryService := mdns.NewMdnsService(
		node,
		config.DiscoveryNamespace,
		&notifee,
	)

	if err := discoveryService.Start(); err != nil {
		panic(err)
	}

	return &notifee.PeerChan
}

// Custom messages
type NewStreamMsg struct {
	id     peer.ID
	stream network.Stream
}

func (n NewStreamMsg) Id() peer.ID {
	return n.id
}

func (n NewStreamMsg) Stream() network.Stream {
	return n.stream
}

func StartConnection() (host.Host, *(chan peerstore.AddrInfo), *(chan NewStreamMsg)) {
	var eventCh = make(chan NewStreamMsg)

	node, err := libp2p.New(
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}

	node.SetStreamHandler(protocol.ID(config.ProtocolID), func(s network.Stream) {
		eventCh <- NewStreamMsg{id: s.Conn().RemotePeer(), stream: s}
	})

	// decover peers
	return node, InitDiscoveryServer(node), &eventCh
}

func writeData(s network.Stream) {
	var counter uint64 = 29

	for {
		<-time.After(time.Second)
		counter++

		err := binary.Write(s, binary.BigEndian, counter)

		fmt.Printf("Sent %d to %s\n", counter, s.ID())
		if err != nil {
			fmt.Println("Error writing: ", err)
		}
	}
}

func readData(s network.Stream) {
	for {
		var counter uint64

		err := binary.Read(s, binary.BigEndian, &counter)
		if err != nil {
			fmt.Println("Error reading: ", err)
		}

		fmt.Printf("Received %d from %s\n", counter, s.ID())
	}
}
