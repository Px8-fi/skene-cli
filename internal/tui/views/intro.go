package views

import (
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// IntroView renders the welcome screen
type IntroView struct {
	width  int
	height int
	time   float64
}

// NewIntroView creates a new intro view
func NewIntroView() *IntroView {
	return &IntroView{}
}

// SetSize updates dimensions
func (v *IntroView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// SetTime updates animation time
func (v *IntroView) SetTime(t float64) {
	v.time = t
}

// Render the intro view
func (v *IntroView) Render() string {
	// Title
	title := styles.Title.Copy().
		MarginBottom(2).
		Render("Welcome to Skene")

	// Animated logo
	logo := components.RenderAnimatedLogo(v.time)

	// Call to action
	enterKey := styles.Accent.Bold(true).Render(">ENTER<")
	cta := styles.Body.Render("Press ") + enterKey + styles.Body.Render(" to start")

	// Footer help
	footer := components.IntroHelp()

	// Combine elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		logo,
		"",
		"",
		cta,
	)

	// Center in viewport
	centered := lipgloss.Place(
		v.width,
		v.height-3, // Leave room for footer
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	// Add footer at bottom
	footerStyled := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		MarginTop(1).
		Render(footer)

	return centered + "\n" + footerStyled
}

// GetHelpItems returns context-specific help
func (v *IntroView) GetHelpItems() []components.HelpItem {
	return []components.HelpItem{
		{Key: "enter", Desc: "start setup"},
		{Key: "?", Desc: "toggle help"},
		{Key: "q", Desc: "quit"},
	}
}
