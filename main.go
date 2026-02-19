package main

import (
	"fmt"
	"os"
	"runtime"

	"tb/internal/config"
	"tb/internal/shell"
	"tb/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
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
			if isatty.IsTerminal(os.Stdout.Fd()) {
				// Stdout is a terminal — no shell integration is capturing output.
				// Show a friendly summary on stderr instead of bare stdout.
				r := lipgloss.DefaultRenderer()
				label := r.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
				cmd := r.NewStyle().Foreground(lipgloss.Color("114"))
				dim := r.NewStyle().Foreground(lipgloss.Color("242"))
				hint := r.NewStyle().Foreground(lipgloss.Color("39"))
				fmt.Fprintf(os.Stderr, "\n%s\n\n%s\n\n%s %s\n",
					label.Render(sel.Name),
					cmd.Render("$ "+sel.Command),
					dim.Render("Tip: add"),
					hint.Render("eval \"$(tb init <shell>)\"")+dim.Render(" to your shell config so commands prefill your prompt."),
				)
			} else {
				// Stdout is a pipe — shell integration is capturing output.
				fmt.Print(sel.Command)
			}
		}
	}
}
