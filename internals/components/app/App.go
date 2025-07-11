package app

import (
	"term-p2p/internals/components/mainview"

	tea "github.com/charmbracelet/bubbletea"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
)

type App struct {
	currentView tea.Model
}

func InitApp(peerChan chan peerstore.AddrInfo) App {
	mainView := mainview.NewMainView(peerChan)

	return App{currentView: mainView}
}

func (a *App) Init() tea.Cmd {
	return a.currentView.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return a, tea.Quit
		}
	}

	nextView, cmd := a.currentView.Update(msg)

	a.currentView = nextView

	return a, cmd
}

func (a *App) View() string {
	return a.currentView.View()
}
