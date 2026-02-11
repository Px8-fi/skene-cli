package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// WelcomeView renders the wizard welcome screen
type WelcomeView struct {
	width  int
	height int
	time   float64
	anim   components.ASCIIMotionModel
}

// NewWelcomeView creates a new welcome view
func NewWelcomeView() *WelcomeView {
	return &WelcomeView{
		anim: components.NewASCIIMotion(styles.IsDarkBackground),
	}
}

// SetSize updates dimensions
func (v *WelcomeView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.anim.SetSize(width, height)
}

// SetTime updates animation time
func (v *WelcomeView) SetTime(t float64) {
	v.time = t
}

// UpdateAnimation updates the animation model with a message
func (v *WelcomeView) UpdateAnimation(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	updatedModel, cmd := v.anim.Update(msg)
	v.anim = updatedModel.(components.ASCIIMotionModel)
	return cmd
}

// InitAnimation returns the initialization command for the animation
func (v *WelcomeView) InitAnimation() tea.Cmd {
	return v.anim.Init()
}

// Render the welcome view
func (v *WelcomeView) Render() string {
	// Content width for consistent centering
	contentWidth := 60
	if v.width > 0 && v.width < contentWidth {
		contentWidth = v.width - 4
	}

	center := lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center)

	// Animated logo
	logo := v.anim.View()

	// Subtitle
	subtitle := center.Render(styles.Subtitle.Render("Product-Led Growth analysis for your codebase"))

	// Call to action
	enterKey := styles.Accent.Bold(true).Render(">ENTER<")
	cta := center.Render(enterKey)

	// Version info
	version := center.Render(styles.Muted.Render("v0.1.8 • github.com/SkeneTechnologies/skene-growth"))

	// Footer help
	footer := components.FooterHelp([]components.HelpItem{
		{Key: "enter", Desc: "start"},
		{Key: "?", Desc: "help"},
		{Key: "ctrl+c", Desc: "quit"},
	})

	// Combine elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		"",
		"",
		cta,
		"",
		subtitle,
		"",
		version,
	)

	// Center in viewport
	centered := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	// Footer pinned at bottom
	footerStyled := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		MarginTop(1).
		Render(footer)

	return centered + "\n" + footerStyled
}

// GetHelpItems returns context-specific help
func (v *WelcomeView) GetHelpItems() []components.HelpItem {
	return []components.HelpItem{
		{Key: "enter", Desc: "start wizard"},
		{Key: "?", Desc: "toggle help"},
		{Key: "ctrl+c", Desc: "quit"},
	}
}
