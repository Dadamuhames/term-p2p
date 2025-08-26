package app

import (
	"term-p2p/internals/common"
	"term-p2p/internals/components/mainview"
	"term-p2p/internals/p2p"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/libp2p/go-libp2p/core/host"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
)

type App struct {
	host        *host.Host
	activePeers common.ActivePeers
	currentView tea.Model
	evtChan     *(chan p2p.NewStreamMsg)
	messageChan (chan common.Message)
}

func InitApp(host *host.Host, peerChan *(chan peerstore.AddrInfo), evtChan *(chan p2p.NewStreamMsg)) App {
	peers := common.NewActivePeers()

	messageChan := make(chan common.Message)

	mainView := mainview.NewMainView(host, peerChan, &peers, &messageChan)

	return App{
		host:        host,
		currentView: mainView,
		activePeers: peers,
		evtChan:     evtChan,
		messageChan: messageChan,
	}
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.currentView.Init(),
		common.GetNextNewPeer(a.evtChan),
		common.GetNextMessage(a.messageChan),
	)
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return a, tea.Quit
		}

	case p2p.NewStreamMsg:
		a.activePeers.AddPeer(msg.Id().String(), msg.Stream(), &a.messageChan)
	}

	nextView, cmd := a.currentView.Update(msg)

	a.currentView = nextView

	return a, tea.Batch(
		cmd,
		common.GetNextNewPeer(a.evtChan),
		common.GetNextMessage(a.messageChan),
	)
}

func (a *App) View() string {
	return a.currentView.View()
}
