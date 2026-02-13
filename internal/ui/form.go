package ui

import (
	"fmt"
	"strings"

	"tb/internal/config"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	fieldName = iota
	fieldDesc
	fieldCmd
	fieldCat
	numFields
)

var fieldLabels = [numFields]string{"Name (*)", "Description", "Command (*)", "Category"}

func (m Model) initCreateForm() Model {
	m.mode = modeForm
	m.formEditing = false
	m.formErr = ""
	for i := 0; i < numFields; i++ {
		ti := textinput.New()
		ti.Prompt = "  "
		ti.CharLimit = 256
		ti.PromptStyle = lipgloss.NewStyle()
		ti.TextStyle = lipgloss.NewStyle().Foreground(clrTextPri)
		ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(clrTextSec)
		ti.Cursor.Style = lipgloss.NewStyle().Foreground(clrAccent)
		m.formFields[i] = ti
	}
	m.formFocused = 0
	m.formFields[0].Focus()
	return m
}

func (m Model) initEditForm() Model {
	m = m.initCreateForm()
	m.formEditing = true

	// Find the index of the selected command in m.commands by name match
	target := m.filtered[m.cursor]
	for i, c := range m.commands {
		if c.Name == target.Name {
			m.formEditIdx = i
			break
		}
	}

	cmd := m.commands[m.formEditIdx]
	m.formFields[fieldName].SetValue(cmd.Name)
	m.formFields[fieldDesc].SetValue(cmd.Description)
	m.formFields[fieldCmd].SetValue(cmd.Command)
	m.formFields[fieldCat].SetValue(cmd.Category)
	return m
}

func (m Model) handleFormKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.ClearEsc):
		m.mode = modeBrowse
		return m, nil
	case key.Matches(msg, keys.Select):
		return m.saveForm()
	case key.Matches(msg, keys.FormTab):
		m.formFields[m.formFocused].Blur()
		m.formFocused = (m.formFocused + 1) % numFields
		m.formFields[m.formFocused].Focus()
		return m, textinput.Blink
	case key.Matches(msg, keys.FormBackTab):
		m.formFields[m.formFocused].Blur()
		m.formFocused = (m.formFocused - 1 + numFields) % numFields
		m.formFields[m.formFocused].Focus()
		return m, textinput.Blink
	}

	var cmd tea.Cmd
	m.formFields[m.formFocused], cmd = m.formFields[m.formFocused].Update(msg)
	return m, cmd
}

func (m Model) saveForm() (Model, tea.Cmd) {
	name := strings.TrimSpace(m.formFields[fieldName].Value())
	cmdText := strings.TrimSpace(m.formFields[fieldCmd].Value())

	if name == "" {
		m.formErr = "Name is required"
		return m, nil
	}
	if cmdText == "" {
		m.formErr = "Command is required"
		return m, nil
	}
	for i, c := range m.commands {
		if c.Name == name && !(m.formEditing && i == m.formEditIdx) {
			m.formErr = fmt.Sprintf("Name %q already exists", name)
			return m, nil
		}
	}

	newCmd := config.Command{
		Name:        name,
		Description: strings.TrimSpace(m.formFields[fieldDesc].Value()),
		Command:     cmdText,
		Category:    strings.TrimSpace(m.formFields[fieldCat].Value()),
	}

	if m.formEditing {
		m.commands[m.formEditIdx] = newCmd
		m.statusMsg = "Command updated"
	} else {
		m.commands = append(m.commands, newCmd)
		m.statusMsg = "Command created"
	}

	if err := config.Save(m.commands); err != nil {
		m.formErr = fmt.Sprintf("Save failed: %v", err)
		return m, nil
	}

	m.mode = modeBrowse
	m = m.refreshAfterMutation()
	return m, nil
}

func (m Model) renderForm() string {
	title := " New Command "
	if m.formEditing {
		title = " Edit Command "
	}
	header := formHeaderStyle.Render(title)

	var rows []string
	for i := 0; i < numFields; i++ {
		labelText := "  " + fieldLabels[i]

		var label string
		if i == m.formFocused {
			label = formFocusedLabelStyle.Render(labelText)
		} else {
			label = formLabelStyle.Render(labelText)
		}

		input := m.formFields[i].View()

		if i == m.formFocused {
			underline := formUnderlineStyle.Render("  " + strings.Repeat("─", 30))
			rows = append(rows, label+"\n"+input+"\n"+underline)
		} else {
			rows = append(rows, label+"\n"+input)
		}
	}
	body := strings.Join(rows, "\n\n")

	var errLine string
	if m.formErr != "" {
		errLine = "\n" + formErrStyle.Render("  ✗ "+m.formErr)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		header, "", body, errLine)
	card := formContainerStyle.Render(content)

	helpLine := helpKeyStyle.Render("tab") + helpDescStyle.Render(" next") +
		helpSepStyle.Render(" · ") +
		helpKeyStyle.Render("enter") + helpDescStyle.Render(" save") +
		helpSepStyle.Render(" · ") +
		helpKeyStyle.Render("esc") + helpDescStyle.Render(" cancel")

	cardWithHelp := lipgloss.JoinVertical(lipgloss.Center,
		card, "", helpLine)

	return lipgloss.Place(m.innerWidth(), m.innerHeight(),
		lipgloss.Center, lipgloss.Center, cardWithHelp)
}

func (m Model) handleDeleteConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, keys.ClearEsc) {
		m.mode = modeBrowse
		return m, nil
	}

	switch msg.String() {
	case "y":
		// Find and remove the command by name match
		target := m.filtered[m.cursor]
		for i, c := range m.commands {
			if c.Name == target.Name {
				m.commands = append(m.commands[:i], m.commands[i+1:]...)
				break
			}
		}

		if err := config.Save(m.commands); err != nil {
			m.statusMsg = fmt.Sprintf("Delete failed: %v", err)
			m.mode = modeBrowse
			return m, nil
		}

		m.statusMsg = "Command deleted"
		m.mode = modeBrowse
		m = m.refreshAfterMutation()
	case "n":
		m.mode = modeBrowse
	}
	return m, nil
}

func (m Model) renderDeleteConfirm(areaWidth, areaHeight int) string {
	if m.cursor >= len(m.filtered) {
		return ""
	}
	name := m.filtered[m.cursor].Name

	title := deleteConfirmStyle.Render(fmt.Sprintf("Delete %q?", name))
	options := helpKeyStyle.Render("y") + helpDescStyle.Render(" yes") +
		helpSepStyle.Render("  ") +
		helpKeyStyle.Render("n") + helpDescStyle.Render(" no") +
		helpSepStyle.Render("  ") +
		helpKeyStyle.Render("esc") + helpDescStyle.Render(" cancel")

	content := lipgloss.JoinVertical(lipgloss.Center, title, "", options)
	box := deleteBoxStyle.Render(content)

	return lipgloss.Place(areaWidth, areaHeight,
		lipgloss.Center, lipgloss.Center, box)
}
