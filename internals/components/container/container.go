package container

import (
	"strings"
	"term-p2p/internals/common"
	"term-p2p/internals/components/list"
	"term-p2p/internals/components/tab"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
)

type Container struct {
	tab           tea.Model
	content       []tea.Model
	currentWindow int
	windowWidth   int
	windowHeight  int
	peerChan      chan peerstore.AddrInfo
}

func NewContainer(peerChan chan peerstore.AddrInfo) Container {
	tabs := []string{"Peers", "Chats"}

	tab := tab.NewTabModel(tabs)

	content := []tea.Model{
		list.NewListModel(),
		list.NewListModel(),
	}

	peers := make(chan peerstore.AddrInfo)

	go func() {
		for {
			peer := <-peerChan
			peers <- peer
		}
	}()

	return Container{tab: tab, content: content, currentWindow: 0, peerChan: peers}
}

func (c Container) Init() tea.Cmd {
	return common.GetNextPeer(c.peerChan)
}

func (c Container) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tab.ActiveTabMsg:
		c.currentWindow = int(msg)

	case tea.WindowSizeMsg:
		c.windowWidth = msg.Width
		c.windowHeight = msg.Height
		for i := 0; i < len(c.content); i++ {
			updated, _ := c.content[i].Update(msg)
			c.content[i] = updated
		}

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return c, tea.Quit
		}
	}

	tabUpdated, cmd := c.tab.Update(msg)

	c.tab = tabUpdated

	c.content[c.currentWindow], _ = c.content[c.currentWindow].Update(msg)

	return c, cmd
}

func (c Container) View() string {
	var stringBuilder strings.Builder
	contentModel := c.content[c.currentWindow]

	tabView := c.tab.View()
	contentView := contentModel.View()

	contentStyle := lipgloss.NewStyle().
		Width(c.windowWidth/2).
		Align(lipgloss.Left, lipgloss.Top)

	composedView := lipgloss.JoinVertical(
		lipgloss.Top,
		contentStyle.Render(tabView),
		contentStyle.Render(contentView),
	)

	containerStyle := lipgloss.NewStyle().
		Width(c.windowWidth).
		Height(c.windowHeight).
		Align(lipgloss.Center, lipgloss.Center)

	stringBuilder.WriteString(containerStyle.Render(composedView))
	return stringBuilder.String()
}
