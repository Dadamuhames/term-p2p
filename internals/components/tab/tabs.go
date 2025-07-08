package tab

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActiveTabMsg int

type TabModel struct {
	Tabs         []string
	activeTab    int
	TabSeparator lipgloss.Style
	TabInactive  lipgloss.Style
	TabActive    lipgloss.Style
}

func NewTabModel(tabs []string) TabModel {
	return TabModel{
		Tabs:        tabs,
		activeTab:   0,
		TabInactive: lipgloss.NewStyle().BorderForeground(highlightColor).Padding(0, 0),
		TabActive: lipgloss.NewStyle().BorderForeground(
			highlightColor).Padding(0, 0).Foreground(lipgloss.Color("210")),
		TabSeparator: lipgloss.NewStyle().SetString("│").Padding(
			0, 1).Foreground(lipgloss.Color("238")),
	}
}

func (m TabModel) Init() tea.Cmd {
	return nil
}

func (m TabModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, m.activeTabCmd
		case "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, m.activeTabCmd
		}
	}

	return m, nil
}

var (
	docStyle       = lipgloss.NewStyle().Padding(1, 0)
	highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
)

func (m TabModel) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		tbString := t

		var style lipgloss.Style
		if i == m.activeTab {
			style = m.TabActive
		} else {
			style = m.TabInactive
		}

		if i != len(m.Tabs)-1 {
			tbString = fmt.Sprintf("%s%s", tbString, m.TabSeparator.String())
		}

		renderedTabs = append(renderedTabs, style.Render(tbString))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)

	return docStyle.Render(doc.String())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (t *TabModel) activeTabCmd() tea.Msg {
	return ActiveTabMsg(t.activeTab)
}
