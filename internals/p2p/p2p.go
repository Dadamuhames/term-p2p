package p2p

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

const (
	protocolID         = "/term-p2p/1.0.0"
	discoveryNamespace = "term-p2p"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

func (n *discoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	n.PeerChan <- peerInfo
}

func InitDiscoveryServer(node host.Host) chan peerstore.AddrInfo {
	notifee := discoveryNotifee{}
	notifee.PeerChan = make(chan peer.AddrInfo)

	// look for peers
	discoveryService := mdns.NewMdnsService(
		node,
		discoveryNamespace,
		&notifee,
	)

	if err := discoveryService.Start(); err != nil {
		panic(err)
	}

	return notifee.PeerChan
}

func handleStream(stream network.Stream) {
	go readData(stream)
	go writeData(stream)
}

func StartConnection() chan peerstore.AddrInfo {
	node, err := libp2p.New(
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}

	node.SetStreamHandler(protocol.ID(protocolID), handleStream)

	// print the node's listening addresses
	// peerInfo := peerstore.AddrInfo{
	// 	ID:    node.ID(),
	// 	Addrs: node.Addrs(),
	// }
	// addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	// fmt.Println("libp2p node address:", addrs)
	//
	// decover peers
	return InitDiscoveryServer(node)
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
