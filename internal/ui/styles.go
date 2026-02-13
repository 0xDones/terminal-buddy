package ui

import "github.com/charmbracelet/lipgloss"

// ── Palette ─────────────────────────────────────────────────────────
// All colors use 256-color codes for broad terminal compatibility.
const (
	clrAccent    = lipgloss.Color("39")  // deep sky blue
	clrAccentDim = lipgloss.Color("75")  // steel blue
	clrTextPri   = lipgloss.Color("253") // bright gray
	clrTextSec   = lipgloss.Color("242") // dim gray
	clrTextMuted = lipgloss.Color("238") // dark gray
	clrBorder    = lipgloss.Color("238") // same as muted
	clrHighBg    = lipgloss.Color("236") // subtle dark bg
	clrNearBlack = lipgloss.Color("232") // near black
	clrCmdText   = lipgloss.Color("114") // light green
	clrSuccess   = lipgloss.Color("78")  // green
	clrError     = lipgloss.Color("196") // red
	clrWarning   = lipgloss.Color("214") // orange
)

// Styles are declared here and initialized lazily via initStyles() so that
// they pick up the correct lipgloss renderer (which main configures to use
// stderr before calling ui.New).
var (
	// ── Application frame ───────────────────────────────────────────
	frameStyle lipgloss.Style

	// ── Thin rule (horizontal separator) ────────────────────────────
	thinRuleStyle lipgloss.Style

	// ── Tab styles ──────────────────────────────────────────────────
	activeTabStyle   lipgloss.Style
	inactiveTabStyle lipgloss.Style

	// ── List item styles ────────────────────────────────────────────
	cursorStyle          lipgloss.Style
	selectedItemStyle    lipgloss.Style
	normalItemStyle      lipgloss.Style
	scrollIndicatorStyle lipgloss.Style

	// ── Detail pane styles ──────────────────────────────────────────
	detailBoxStyle    lipgloss.Style
	detailTitleStyle  lipgloss.Style
	detailLabelStyle  lipgloss.Style
	detailValueStyle  lipgloss.Style
	commandValueStyle lipgloss.Style
	categoryTagStyle  lipgloss.Style

	// ── Search bar styles ───────────────────────────────────────────
	searchPromptStyle lipgloss.Style

	// ── Help bar styles ─────────────────────────────────────────────
	helpKeyStyle  lipgloss.Style
	helpDescStyle lipgloss.Style
	helpSepStyle  lipgloss.Style

	// ── Form styles ─────────────────────────────────────────────────
	formContainerStyle    lipgloss.Style
	formHeaderStyle       lipgloss.Style
	formLabelStyle        lipgloss.Style
	formFocusedLabelStyle lipgloss.Style
	formUnderlineStyle    lipgloss.Style
	formErrStyle          lipgloss.Style

	// ── Status / confirmation styles ────────────────────────────────
	statusMsgStyle     lipgloss.Style
	deleteConfirmStyle lipgloss.Style
	deleteBoxStyle     lipgloss.Style
)

// initStyles creates all lipgloss styles. Must be called after the default
// lipgloss renderer has been configured (e.g., to use stderr).
func initStyles() {
	frameStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(clrBorder).
		Padding(0, 1)

	thinRuleStyle = lipgloss.NewStyle().
		Foreground(clrTextMuted)

	activeTabStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(clrNearBlack).
		Background(clrAccent).
		Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
		Foreground(clrTextSec).
		Padding(0, 1)

	cursorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(clrAccent)

	selectedItemStyle = lipgloss.NewStyle().
		Foreground(clrAccent).
		Background(clrHighBg).
		Bold(true)

	normalItemStyle = lipgloss.NewStyle().
		Foreground(clrTextPri)

	scrollIndicatorStyle = lipgloss.NewStyle().
		Foreground(clrTextSec)

	detailBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(clrBorder).
		Padding(0, 1)

	detailTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(clrAccent)

	detailLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(clrAccentDim)

	detailValueStyle = lipgloss.NewStyle().
		Foreground(clrTextPri)

	commandValueStyle = lipgloss.NewStyle().
		Foreground(clrCmdText)

	categoryTagStyle = lipgloss.NewStyle().
		Foreground(clrTextSec).
		Italic(true)

	searchPromptStyle = lipgloss.NewStyle().
		Foreground(clrAccent)

	helpKeyStyle = lipgloss.NewStyle().
		Foreground(clrAccentDim)

	helpDescStyle = lipgloss.NewStyle().
		Foreground(clrTextSec)

	helpSepStyle = lipgloss.NewStyle().
		Foreground(clrTextMuted)

	formContainerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(clrAccent).
		Padding(1, 2)

	formHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(clrNearBlack).
		Background(clrAccent).
		Padding(0, 2)

	formLabelStyle = lipgloss.NewStyle().
		Foreground(clrTextSec)

	formFocusedLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(clrAccent)

	formUnderlineStyle = lipgloss.NewStyle().
		Foreground(clrAccent)

	formErrStyle = lipgloss.NewStyle().
		Foreground(clrError)

	statusMsgStyle = lipgloss.NewStyle().
		Foreground(clrSuccess)

	deleteConfirmStyle = lipgloss.NewStyle().
		Foreground(clrWarning).
		Bold(true)

	deleteBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(clrWarning).
		Padding(1, 3)
}
