package peerview

import (
	"fmt"
	"strings"
	"term-p2p/internals/components/mainview/list"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/libp2p/go-libp2p/core/network"
)

const gap = "\n\n"

type (
	errMsg error
)

type PeerMessage struct {
	message string
	sender  string
}

type PeerView struct {
	peerInfo     list.CustomListItem
	peerStream   *network.Stream
	viewport     viewport.Model
	messages     []PeerMessage
	textarea     textarea.Model
	senderStyle  lipgloss.Style
	windowWidth  int
	windowHeight int
	err          error
}

func InitialModel(item list.CustomListItem, peerStream *network.Stream) PeerView {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 28)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return PeerView{
		textarea:    ta,
		messages:    []PeerMessage{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		peerInfo:    item,
		peerStream:  peerStream,
	}
}

func (m PeerView) messagesToString() string {
	builder := strings.Builder{}

	for i := 0; i < len(m.messages); i++ {
		msg := m.messages[i]

		builder.WriteString(fmt.Sprintf("%s: %s\n\n", m.senderStyle.Render(msg.sender), msg.message))
	}

	return builder.String()
}

func (m PeerView) Init() tea.Cmd {
	return textarea.Blink
}

func (m PeerView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.viewport.Width = msg.Width / 2
		m.textarea.SetWidth(msg.Width / 2)
		m.viewport.Height = msg.Height
		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(
				lipgloss.NewStyle().Width(m.viewport.Width).Render())
		}
		m.viewport.GotoBottom()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textarea.Value() == "" {
				break
			}

			m.messages = append(
				m.messages,
				PeerMessage{message: m.textarea.Value(), sender: "You"},
			)
			m.viewport.SetContent(
				lipgloss.NewStyle().Width(m.viewport.Width).Render(m.messagesToString()))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m PeerView) View() string {
	marginLeft := 40

	contentStyle := lipgloss.NewStyle().
		MarginLeft(marginLeft).
		Align(lipgloss.Center, lipgloss.Center)

	composedView := lipgloss.JoinVertical(
		lipgloss.Bottom,
		contentStyle.Render(m.viewport.View()),
		contentStyle.Render(m.textarea.View()),
	)

	containerStyle := lipgloss.NewStyle().
		Width(m.windowWidth).
		Height(m.windowHeight).
		Align(lipgloss.Center, lipgloss.Bottom)

	return containerStyle.Render(composedView)
}
