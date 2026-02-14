package main

import (
	"fmt"
	"os"
	"runtime"

	"tb/internal/clipboard"
	"tb/internal/config"
	"tb/internal/shell"
	"tb/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var version = "dev"

func main() {
	// CLI routing
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "version":
			fmt.Printf("tb %s (%s)\n", version, runtime.Version())
			return
		case "init":
			if len(os.Args) != 3 {
				fmt.Fprintln(os.Stderr, "Usage: tb init <bash|zsh|fish>")
				os.Exit(1)
			}
			script, err := shell.InitScript(os.Args[2])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(script)
			return
		}
	}

	// Default: launch TUI
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Use stderr for all TUI rendering so it displays correctly even when
	// stdout is captured by a shell widget (e.g., selected="$(tb)").
	// Also tell lipgloss to detect color support from stderr, not stdout.
	lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stderr))

	m := ui.New(cfg.Commands, cfg.Categories(), cfg.Keybindings)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithOutput(os.Stderr))

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if result, ok := finalModel.(ui.Model); ok {
		if sel := result.Selected(); sel != nil {
			// Best-effort clipboard copy (ignore errors for headless/no-xclip)
			_ = clipboard.Write(sel.Command)
			// Print to stdout for shell integration / piping
			fmt.Print(sel.Command)
		}
	}
}
