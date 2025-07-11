package mainview

import (
	"strings"
	"term-p2p/internals/common"
	"term-p2p/internals/components/mainview/list"
	"term-p2p/internals/components/peerview"
	"term-p2p/internals/components/tab"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
)

type MainView struct {
	tab           tea.Model
	content       []tea.Model
	currentWindow int
	windowWidth   int
	windowHeight  int
	peerChan      chan peerstore.AddrInfo
}

func NewMainView(peerChan chan peerstore.AddrInfo) MainView {
	tabs := []string{"Peers", "Chats"}

	tab := tab.NewTabModel(tabs)

	peers := make(chan peerstore.AddrInfo)

	go func() {
		for {
			peer := <-peerChan
			peers <- peer
		}
	}()

	content := []tea.Model{
		list.NewListModel(peers),
		list.NewListModel(peers),
	}

	return MainView{tab: tab, content: content, currentWindow: 0, peerChan: peers}
}

func (c MainView) Init() tea.Cmd {
	return common.GetNextPeer(c.peerChan)
}

func (c MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tab.ActiveTabMsg:
		c.currentWindow = int(msg)

	case list.PeerSelectMsg:
		peerChat := peerview.InitialModel()

		return peerChat.Update(msg)

	case tea.WindowSizeMsg:
		c.windowWidth = msg.Width
		c.windowHeight = msg.Height
		for i := 0; i < len(c.content); i++ {
			updated, _ := c.content[i].Update(msg)
			c.content[i] = updated
		}
	}

	// update tab
	tabUpdated, cmd := c.tab.Update(msg)

	cmds = append(cmds, cmd)

	c.tab = tabUpdated

	// update content
	contentUpdated, contentCmd := c.content[c.currentWindow].Update(msg)

	c.content[c.currentWindow] = contentUpdated

	cmds = append(cmds, contentCmd)

	return c, tea.Batch(cmds...)
}

func (c MainView) View() string {
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
