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

var (
	// ── Application frame ───────────────────────────────────────────
	frameStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(clrBorder).
			Padding(0, 1)

	// ── Thin rule (horizontal separator) ────────────────────────────
	thinRuleStyle = lipgloss.NewStyle().
			Foreground(clrTextMuted)

	// ── Tab styles ──────────────────────────────────────────────────
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(clrNearBlack).
			Background(clrAccent).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(clrTextSec).
				Padding(0, 1)

	// ── List item styles ────────────────────────────────────────────
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

	// ── Detail pane styles ──────────────────────────────────────────
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

	// ── Search bar styles ───────────────────────────────────────────
	searchPromptStyle = lipgloss.NewStyle().
				Foreground(clrAccent)

	// ── Help bar styles ─────────────────────────────────────────────
	helpKeyStyle = lipgloss.NewStyle().
			Foreground(clrAccentDim)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(clrTextSec)

	helpSepStyle = lipgloss.NewStyle().
			Foreground(clrTextMuted)

	// ── Form styles ─────────────────────────────────────────────────
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
			Bold(true).
			Foreground(clrAccentDim).
			Width(14).
			Align(lipgloss.Right)

	formFocusedLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(clrAccent).
				Width(14).
				Align(lipgloss.Right)

	formErrStyle = lipgloss.NewStyle().
			Foreground(clrError)

	// ── Status / confirmation styles ────────────────────────────────
	statusMsgStyle = lipgloss.NewStyle().
			Foreground(clrSuccess)

	deleteConfirmStyle = lipgloss.NewStyle().
				Foreground(clrWarning).
				Bold(true)

	deleteBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(clrWarning).
			Padding(1, 3)
)
