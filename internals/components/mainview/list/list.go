package list

import (
	"term-p2p/internals/common"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
)

var (
	highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	appStyle = lipgloss.NewStyle().Padding(0, 0).Align(
		lipgloss.Left,
	)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
	peerChan     chan peerstore.AddrInfo
}

type customListItem struct {
	id       string
	title    string
	subtitle string
}

func (c customListItem) peerSelectCmd() tea.Msg {
	return PeerSelectMsg(c)
}

func (c customListItem) Id() string          { return c.id }
func (c customListItem) Title() string       { return c.title }
func (c customListItem) Description() string { return c.subtitle }
func (c customListItem) FilterValue() string { return c.title }

func NewListModel(peerChan chan peerstore.AddrInfo) model {
	var (
		listKeys     = newListKeyMap()
		delegateKeys = newDelegateKeyMap()
	)

	items := make([]list.Item, 0)

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	groceryList := list.New(
		items,
		delegate,
		0, 0)

	groceryList.SetShowTitle(false)
	groceryList.SetShowStatusBar(false)
	groceryList.SetShowHelp(false)
	groceryList.SetShowPagination(false)
	groceryList.Styles.Title = titleStyle
	groceryList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return model{
		list:         groceryList,
		keys:         listKeys,
		delegateKeys: delegateKeys,
		peerChan:     peerChan,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height/2)

	case peerstore.AddrInfo:
		itemsList := m.list.Items()

		newItem := customListItem{id: msg.ID.String(), title: msg.ID.String(), subtitle: "new peer"}

		itemsList = append(itemsList, newItem)

		cmd := m.list.SetItems(itemsList)

		cmds = append(cmds, cmd)
		cmds = append(cmds, common.GetNextPeer(m.peerChan))

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}
