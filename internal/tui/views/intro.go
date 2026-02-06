package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// IntroView renders the welcome screen
type IntroView struct {
	width  int
	height int
	time   float64
	anim   components.ASCIIMotionModel
}

// NewIntroView creates a new intro view
func NewIntroView() *IntroView {
	return &IntroView{
		anim: components.NewASCIIMotionWithDefaults(),
	}
}

// SetSize updates dimensions
func (v *IntroView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.anim.SetSize(width, height)
}

// SetTime updates animation time (kept for compatibility, but animation now self-manages)
func (v *IntroView) SetTime(t float64) {
	v.time = t
}

// UpdateAnimation updates the animation model with a message
func (v *IntroView) UpdateAnimation(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	updatedModel, cmd := v.anim.Update(msg)
	v.anim = updatedModel.(components.ASCIIMotionModel)
	return cmd
}

// InitAnimation returns the initialization command for the animation
func (v *IntroView) InitAnimation() tea.Cmd {
	return v.anim.Init()
}

// Render the intro view
func (v *IntroView) Render() string {
	// Title
	title := styles.Title.Copy().
		MarginBottom(2).
		Render("Welcome to Skene")

	// Animated logo from ASCII motion component
	logo := v.anim.View()

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
