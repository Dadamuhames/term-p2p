package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
)

const (
	protocolID         = "/term-p2p/1.0.0"
	discoveryNamespace = "term-p2p"
)

func main() {
	node, err := libp2p.New(
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}

	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	// print the node's listening addresses
	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	fmt.Println("libp2p node address:", addrs[0])

	notifee := discoveryNotifee{h: node}

	// look for peers
	discoveryService := mdns.NewMdnsService(
		node,
		discoveryNamespace,
		&notifee,
	)

	if err := discoveryService.Start(); err != nil {
		fmt.Printf("Failed to start mDNS service: %v\n", err)
		return
	}

	// wait for CTRL+C
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

	// shut the node down
	if err := node.Close(); err != nil {
		panic(err)
	}
}

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	fmt.Println("found peer", peerInfo.Addrs)

	peerAddrInfos, err := peer.AddrInfosFromP2pAddrs(peerInfo.Addrs...)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(peerAddrInfos); i++ {
		peerAddrInfo := peerAddrInfos[i]

		if !strings.HasPrefix(peerAddrInfo.String(), "/ip4/192.168.100.160/tcp/") {
			continue
		}

		if err := n.h.Connect(context.Background(), *&peerAddrInfo); err != nil {
			panic(err)
		}
		fmt.Println("Connected to", peerAddrInfo.String())

		s, err := n.h.NewStream(context.Background(), peerAddrInfo.ID, protocolID)
		if err != nil {
			panic(err)
		}

		go writeCounter(s)
		go readCounter(s)
	}
}

func writeCounter(s network.Stream) {
	var counter uint64

	for {
		<-time.After(time.Second)
		counter++

		err := binary.Write(s, binary.BigEndian, counter)
		if err != nil {
			panic(err)
		}
	}
}

func readCounter(s network.Stream) {
	for {
		var counter uint64

		err := binary.Read(s, binary.BigEndian, &counter)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Received %d from %s\n", counter, s.ID())
	}
}
