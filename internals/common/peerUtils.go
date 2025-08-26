package common

import (
	"term-p2p/internals/p2p"

	tea "github.com/charmbracelet/bubbletea"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
)

func GetNextNewPeer(evtChan *(chan p2p.NewStreamMsg)) tea.Cmd {
	return func() tea.Msg {
		return <-*evtChan
	}
}

func GetNextPeer(c <-chan peerstore.AddrInfo) tea.Cmd {
	return func() tea.Msg {
		return <-c
	}
}

func GetNextMessage(m <-chan Message) tea.Cmd {
	return func() tea.Msg {
		return <-m
	}
}
