package mainview

import (
	"context"
	"term-p2p/internals/common"
	"term-p2p/internals/components/mainview/list"
	"term-p2p/internals/components/peerview"
	"term-p2p/internals/components/tab"
	"term-p2p/internals/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/libp2p/go-libp2p/core/host"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type MainView struct {
	host          *host.Host
	tab           tea.Model
	content       []tea.Model
	currentWindow int
	windowWidth   int
	windowHeight  int
	peerChan      chan peerstore.AddrInfo
	activePeers   *common.ActivePeers
	messageChan   *(chan common.Message)
}

func NewMainView(
	host *host.Host,
	peerChan *(chan peerstore.AddrInfo),
	activePeers *common.ActivePeers,
	messageChan *(chan common.Message)) MainView {

	tabs := []string{"Discover", "Requests"}

	tab := tab.NewTabModel(tabs)

	content := []tea.Model{
		list.NewListModel(*peerChan),
		list.NewListModel(*peerChan),
	}

	return MainView{
		host:          host,
		tab:           tab,
		content:       content,
		currentWindow: 0,
		peerChan:      *peerChan,
		activePeers:   activePeers,
		messageChan:   messageChan,
	}
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
		peerId := msg.Item.Id()

		peerStream, exists := (*c.activePeers).GetPeerById(peerId)

		if !exists {
			ctx := context.Background()

			peerAddrs := msg.Item.Addrs()

			for i := 0; i < len(peerAddrs); i++ {
				addr := peerAddrs[i]

				peer, err := peerstore.AddrInfoFromP2pAddr(addr)

				if err != nil {
					continue
				}

				if err := (*c.host).Connect(ctx, *peer); err != nil {
					panic(err)
				}

				peerStream, err = (*c.host).NewStream(ctx, peer.ID, protocol.ID(config.ProtocolID))

				(*c.activePeers).AddPeer(peerId, peerStream, c.messageChan)

				break
			}
		}

		peerChat := peerview.InitialModel(msg.Item, &peerStream)

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

	return containerStyle.Render(composedView)
}
