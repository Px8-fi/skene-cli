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
		anim: components.NewASCIIMotionWithDefaults(),
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
	// Title
	title := styles.Title.Copy().
		MarginBottom(1).
		Render("Welcome to Skene Growth")

	// Animated logo
	logo := v.anim.View()

	// Subtitle
	subtitle := styles.Subtitle.Render("Product-Led Growth analysis for your codebase")

	// Features
	features := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.Body.Render("  Detect tech stacks & growth opportunities"),
		styles.Body.Render("  Generate growth manifests & documentation"),
		styles.Body.Render("  AI-powered growth loop analysis"),
	)

	// Call to action
	enterKey := styles.Accent.Bold(true).Render(">ENTER<")
	cta := styles.Body.Render("Press ") + enterKey + styles.Body.Render(" to begin setup")

	// Version info
	version := styles.Muted.Render("v0.1.8 • github.com/SkeneTechnologies/skene-growth")

	// Footer help
	footer := components.FooterHelp([]components.HelpItem{
		{Key: "enter", Desc: "start"},
		{Key: "?", Desc: "help"},
		{Key: "ctrl+c", Desc: "quit"},
	})

	// Combine elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		logo,
		"",
		subtitle,
		"",
		features,
		"",
		"",
		cta,
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

	// Add footer at bottom
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
