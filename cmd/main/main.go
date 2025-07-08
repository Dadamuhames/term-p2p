package main

import (
	"fmt"
	"os"
	"term-p2p/internals/components/container"
	"term-p2p/internals/p2p"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	peerChan := p2p.StartConnection()
	model := container.NewContainer(peerChan)

	_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
