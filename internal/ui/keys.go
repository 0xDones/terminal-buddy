package ui

import "github.com/charmbracelet/bubbles/key"

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

var keys = keyMap{
	Up:          key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:        key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	NextTab:     key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
	PrevTab:     key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev tab")),
	Search:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	ClearEsc:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "clear/exit search")),
	Select:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Copy:        key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy")),
	Quit:        key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Create:      key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
	Edit:        key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Delete:      key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	FormTab:     key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
	FormBackTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev field")),
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
