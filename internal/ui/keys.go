package ui

import (
	"strings"

	"tb/internal/config"

	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	NextTab  key.Binding
	PrevTab  key.Binding
	Search   key.Binding
	ClearEsc key.Binding
	Select   key.Binding
	Copy     key.Binding
	Quit     key.Binding
	// Command management
	Create      key.Binding
	Edit        key.Binding
	Delete      key.Binding
	FormTab     key.Binding
	FormBackTab key.Binding
}

var keys keyMap

// keyDisplayNames maps BubbleTea key identifiers to display characters.
var keyDisplayNames = map[string]string{
	"up": "↑", "down": "↓", "left": "←", "right": "→",
}

func helpKeyLabel(keyList []string) string {
	parts := make([]string, len(keyList))
	for i, k := range keyList {
		if d, ok := keyDisplayNames[k]; ok {
			parts[i] = d
		} else {
			parts[i] = k
		}
	}
	return strings.Join(parts, "/")
}

func buildBinding(configured []string, defaults []string, helpDesc string) key.Binding {
	k := defaults
	if len(configured) > 0 {
		k = configured
	}
	return key.NewBinding(key.WithKeys(k...), key.WithHelp(helpKeyLabel(k), helpDesc))
}

func initKeys(kb config.Keybindings) {
	keys = keyMap{
		Up:          buildBinding(kb.Up, []string{"up", "k"}, "up"),
		Down:        buildBinding(kb.Down, []string{"down", "j"}, "down"),
		NextTab:     buildBinding(kb.NextTab, []string{"tab"}, "next tab"),
		PrevTab:     buildBinding(kb.PrevTab, []string{"shift+tab"}, "prev tab"),
		Search:      buildBinding(kb.Search, []string{"/"}, "search"),
		ClearEsc:    buildBinding(kb.ClearEsc, []string{"esc"}, "clear/exit search"),
		Select:      buildBinding(kb.Select, []string{"enter"}, "select"),
		Copy:        buildBinding(kb.Copy, []string{"c"}, "copy"),
		Quit:        buildBinding(kb.Quit, []string{"q", "ctrl+c"}, "quit"),
		Create:      buildBinding(kb.Create, []string{"n"}, "new"),
		Edit:        buildBinding(kb.Edit, []string{"e"}, "edit"),
		Delete:      buildBinding(kb.Delete, []string{"d"}, "delete"),
		FormTab:     key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
		FormBackTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev field")),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.NextTab, k.Search, k.Select, k.Copy, k.Create, k.Edit, k.Delete, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.NextTab, k.PrevTab},
		{k.Search, k.ClearEsc},
		{k.Select, k.Copy, k.Quit},
		{k.Create, k.Edit, k.Delete},
	}
}
