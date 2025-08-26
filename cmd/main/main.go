package main

import (
	"fmt"
	"os"
	"term-p2p/internals/components/app"
	"term-p2p/internals/p2p"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	host, peerChan, evntChan := p2p.StartConnection()
	app := app.InitApp(&host, peerChan, evntChan)

	_, err := tea.NewProgram(&app, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
