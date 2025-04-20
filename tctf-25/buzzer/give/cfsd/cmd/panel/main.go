package main

import (
	"fmt"
	"log"

	"cfsd/tui/panel"

	tea "github.com/charmbracelet/bubbletea"
)

func runTUI() error {
	m, err := panel.NewModel()
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	m.ThisProgram = p

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to init tea program: %w", err)
	}

	return nil
}

func main() {
	if err := runTUI(); err != nil {
		log.Fatalf("bruh: %v", err)
	}
}
