package common

import (
	tea "github.com/charmbracelet/bubbletea"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
)

func GetNextPeer(c <-chan peerstore.AddrInfo) tea.Cmd {
	return func() tea.Msg {
		return <-c
	}
}
