package main

import (
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"os"
	"os/signal"
	"syscall"
)

const (
	protocolID         = "/term-p2p/1.0.0"
	discoveryNamespace = "term-p2p"
)

func main() {
	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/192.168.100.13/tcp/2000"),
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
	fmt.Println("found peer", peerInfo.String())
}
