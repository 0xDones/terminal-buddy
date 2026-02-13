package main

import (
	"fmt"
	"os"

	"clo/internal/config"
	"clo/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	m := ui.New(cfg.Commands, cfg.Categories())
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if result, ok := finalModel.(ui.Model); ok {
		if sel := result.Selected(); sel != nil {
			fmt.Print(sel.Command)
		}
	}
}
