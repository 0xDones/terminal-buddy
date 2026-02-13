package ui

import (
	"fmt"
	"sort"
	"strings"

	"clo/internal/config"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

const (
	modeBrowse = iota
	modeForm
	modeDeleteConfirm
)

// Model is the main BubbleTea model for the command browser TUI.
type Model struct {
	commands     []config.Command // all commands from config
	filtered     []config.Command // after category + search filter
	tabs         []string         // "All" + category names
	activeTab    int
	cursor       int
	scrollOffset int
	search       textinput.Model
	searchActive bool // true when search input is focused
	width        int
	height       int
	selected     *config.Command // set on Enter, nil otherwise

	// Command management state
	mode        int
	formFields  [numFields]textinput.Model
	formFocused int
	formEditing bool
	formEditIdx int
	formErr     string
	statusMsg   string
}

// New creates the TUI model from loaded commands.
func New(commands []config.Command, categories []string) Model {
	ti := textinput.New()
	ti.Placeholder = "type to search..."
	ti.Prompt = "/ "
	ti.PromptStyle = searchPromptStyle
	ti.TextStyle = lipgloss.NewStyle().Foreground(clrTextPri)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(clrTextSec)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(clrAccent)
	ti.CharLimit = 100

	tabs := append([]string{"All"}, categories...)

	m := Model{
		commands: commands,
		tabs:     tabs,
		search:   ti,
	}
	m = m.filterCommands()
	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Selected returns the command chosen by the user, or nil if they quit.
func (m Model) Selected() *config.Command {
	return m.selected
}

// ── Layout helpers ──────────────────────────────────────────────────

func (m Model) innerWidth() int {
	// frame border(2) + frame horizontal padding(2)
	return m.width - 4
}

func (m Model) innerHeight() int {
	// frame border(2)
	return m.height - 2
}

func (m Model) bodyHeight() int {
	// inner height minus tabs(1) + rule(1) + search/status(1) + help(1)
	return m.innerHeight() - 4
}

func (m Model) listHeight() int {
	// body height minus scroll indicator line(1)
	return m.bodyHeight() - 1
}

// ── Update ──────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, keys.Quit) && msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		switch m.mode {
		case modeForm:
			return m.handleFormKeys(msg)
		case modeDeleteConfirm:
			return m.handleDeleteConfirmKeys(msg)
		default:
			// Clear status message on any keypress in browse mode
			m.statusMsg = ""

			if m.searchActive {
				return m.handleSearchKeys(msg)
			}
			return m.handleNavigationKeys(msg)
		}
	}
	return m, nil
}

func (m Model) handleSearchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, keys.ClearEsc) {
		m.searchActive = false
		m.search.Blur()
		return m, nil
	}
	if key.Matches(msg, keys.Select) {
		if len(m.filtered) > 0 {
			cmd := m.filtered[m.cursor]
			m.selected = &cmd
			return m, tea.Quit
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.search, cmd = m.search.Update(msg)
	m = m.filterCommands()
	m.cursor = 0
	m.scrollOffset = 0
	return m, cmd
}

func (m Model) handleNavigationKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, keys.Search):
		m.searchActive = true
		m.search.Focus()
		return m, textinput.Blink
	case key.Matches(msg, keys.Select):
		if len(m.filtered) > 0 {
			cmd := m.filtered[m.cursor]
			m.selected = &cmd
			return m, tea.Quit
		}
	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, keys.Down):
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
	case key.Matches(msg, keys.NextTab):
		m.activeTab = (m.activeTab + 1) % len(m.tabs)
		m = m.filterCommands()
		m.cursor = 0
		m.scrollOffset = 0
	case key.Matches(msg, keys.PrevTab):
		m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
		m = m.filterCommands()
		m.cursor = 0
		m.scrollOffset = 0
	case key.Matches(msg, keys.Create):
		m = m.initCreateForm()
		return m, textinput.Blink
	case key.Matches(msg, keys.Edit):
		if len(m.filtered) > 0 {
			m = m.initEditForm()
			return m, textinput.Blink
		}
	case key.Matches(msg, keys.Delete):
		if len(m.filtered) > 0 {
			m.mode = modeDeleteConfirm
		}
	}

	m = m.adjustScroll()
	return m, nil
}

func (m Model) adjustScroll() Model {
	listHeight := m.listHeight()
	if listHeight <= 0 {
		return m
	}
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
	if m.cursor >= m.scrollOffset+listHeight {
		m.scrollOffset = m.cursor - listHeight + 1
	}
	return m
}

func (m Model) filterCommands() Model {
	// Step 1: filter by category
	var pool []config.Command
	if m.activeTab == 0 { // "All"
		pool = m.commands
	} else {
		cat := m.tabs[m.activeTab]
		for _, cmd := range m.commands {
			if cmd.Category == cat {
				pool = append(pool, cmd)
			}
		}
	}

	// Step 2: fuzzy search
	query := m.search.Value()
	if query == "" {
		m.filtered = pool
		return m
	}

	source := commandSource(pool)
	matches := fuzzy.FindFrom(query, source)
	result := make([]config.Command, len(matches))
	for i, match := range matches {
		result[i] = pool[match.Index]
	}
	m.filtered = result
	return m
}

// commandSource adapts []config.Command for sahilm/fuzzy.
type commandSource []config.Command

func (s commandSource) String(i int) string {
	return s[i].Name + " " + s[i].Description
}

func (s commandSource) Len() int {
	return len(s)
}

// ── View ────────────────────────────────────────────────────────────

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}
	if m.width < 40 || m.height < 10 {
		return "Terminal too small. Please resize to at least 40x10."
	}

	iw := m.innerWidth()
	ih := m.innerHeight()
	var inner string

	switch m.mode {
	case modeForm:
		inner = m.renderForm()
	case modeDeleteConfirm:
		tabBar := m.renderTabs()
		rule := m.thinRule()
		confirmAreaH := m.bodyHeight() + 1 // body + search line
		confirm := m.renderDeleteConfirm(iw, confirmAreaH)
		helpView := m.renderHelp()
		inner = lipgloss.JoinVertical(lipgloss.Left,
			tabBar, rule, confirm, helpView)
	default:
		tabBar := m.renderTabs()
		rule := m.thinRule()
		body := m.renderBody()
		var searchOrStatus string
		if m.statusMsg != "" {
			searchOrStatus = statusMsgStyle.Render(" " + m.statusMsg)
		} else {
			searchOrStatus = m.renderSearchBar()
		}
		helpView := m.renderHelp()
		inner = lipgloss.JoinVertical(lipgloss.Left,
			tabBar, rule, body, searchOrStatus, helpView)
	}

	return frameStyle.
		Width(iw).
		Height(ih).
		Render(inner)
}

func (m Model) thinRule() string {
	return thinRuleStyle.Render(strings.Repeat("─", m.innerWidth()))
}

func (m Model) renderTabs() string {
	var tabs []string
	for i, name := range m.tabs {
		if i == m.activeTab {
			tabs = append(tabs, activeTabStyle.Render(name))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(name))
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...)
	// Clamp to inner width so overflowing tabs don't break layout geometry.
	if lipgloss.Width(row) > m.innerWidth() {
		row = lipgloss.NewStyle().MaxWidth(m.innerWidth()).Render(row)
	}
	return row
}

func (m Model) renderBody() string {
	iw := m.innerWidth()
	bh := m.bodyHeight()

	listWidth := iw * 2 / 5
	detailOuter := iw - listWidth - 1 // 1 for gap column

	list := m.renderList(listWidth)
	listBox := lipgloss.NewStyle().
		Width(listWidth).
		Height(bh).
		Render(list)

	detailInner := detailOuter - 4 // border(2) + padding(2)
	detail := m.renderDetail(detailInner)
	detailBox := detailBoxStyle.
		Width(detailInner).
		Height(bh - 2). // border(2)
		Render(detail)

	return lipgloss.JoinHorizontal(lipgloss.Top, listBox, " ", detailBox)
}

func (m Model) renderList(width int) string {
	if len(m.filtered) == 0 {
		empty := lipgloss.NewStyle().Foreground(clrTextSec)
		return empty.Render("  No commands found")
	}

	lh := m.listHeight()
	var lines []string
	end := min(m.scrollOffset+lh, len(m.filtered))
	for i := m.scrollOffset; i < end; i++ {
		name := m.filtered[i].Name
		if i == m.cursor {
			cursor := cursorStyle.Render(" > ")
			name = selectedItemStyle.
				Width(width - lipgloss.Width(cursor)).
				Render(name)
			lines = append(lines, cursor+name)
		} else {
			lines = append(lines, "   "+normalItemStyle.Render(name))
		}
	}

	// Pad to fill visible item area
	for len(lines) < lh {
		lines = append(lines, "")
	}

	// Scroll indicator on last line, right-aligned
	indicator := scrollIndicatorStyle.Render(
		fmt.Sprintf("%d/%d", m.cursor+1, len(m.filtered)))
	padLen := width - lipgloss.Width(indicator)
	lines = append(lines, strings.Repeat(" ", max(0, padLen))+indicator)

	return strings.Join(lines, "\n")
}

func (m Model) renderDetail(width int) string {
	if len(m.filtered) == 0 || m.cursor >= len(m.filtered) {
		return ""
	}

	cmd := m.filtered[m.cursor]

	title := detailTitleStyle.Render("Details")

	desc := detailLabelStyle.Render("Description") + "\n" +
		detailValueStyle.Width(width).Render(cmd.Description)

	cmdText := detailLabelStyle.Render("Command") + "\n" +
		commandValueStyle.Width(width).Render(cmd.Command)

	parts := []string{title, "", desc, "", cmdText}

	if cmd.Category != "" {
		parts = append(parts, "", categoryTagStyle.Render(cmd.Category))
	}

	return strings.Join(parts, "\n")
}

func (m Model) renderSearchBar() string {
	return " " + m.search.View()
}

func (m Model) renderHelp() string {
	bindings := keys.ShortHelp()
	var parts []string
	for _, b := range bindings {
		k := helpKeyStyle.Render(b.Help().Key)
		d := helpDescStyle.Render(b.Help().Desc)
		parts = append(parts, k+" "+d)
	}
	sep := helpSepStyle.Render(" · ")
	line := " " + strings.Join(parts, sep)
	// Clamp to inner width so wrapping doesn't consume extra rows.
	return lipgloss.NewStyle().MaxWidth(m.innerWidth()).Render(line)
}

// refreshAfterMutation rebuilds tabs, filters, and clamps cursor after a command list change.
func (m Model) refreshAfterMutation() Model {
	// Preserve current tab name so we can re-find it after rebuild
	currentTab := ""
	if m.activeTab < len(m.tabs) {
		currentTab = m.tabs[m.activeTab]
	}

	seen := make(map[string]struct{})
	for _, cmd := range m.commands {
		if cmd.Category != "" {
			seen[cmd.Category] = struct{}{}
		}
	}
	cats := make([]string, 0, len(seen))
	for cat := range seen {
		cats = append(cats, cat)
	}
	sort.Strings(cats)
	m.tabs = append([]string{"All"}, cats...)

	// Re-find the tab by name; fall back to "All" if the category was removed
	m.activeTab = 0
	for i, t := range m.tabs {
		if t == currentTab {
			m.activeTab = i
			break
		}
	}

	m = m.filterCommands()

	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
	m.scrollOffset = 0
	m = m.adjustScroll()
	return m
}
