package views

import (
	"fmt"
	"skene-terminal-v2/internal/services/config"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// AuthView handles Skene auth simulation with countdown
type AuthView struct {
	width     int
	height    int
	provider  *config.Provider
	countdown int // seconds remaining
	authURL   string
	showFallback bool
}

// NewAuthView creates a new auth view
func NewAuthView(provider *config.Provider) *AuthView {
	authURL := "https://www.skene.ai/login?=retrieve-api-key"
	if provider != nil && provider.AuthURL != "" {
		authURL = provider.AuthURL
	}

	return &AuthView{
		provider:     provider,
		countdown:    3,
		authURL:      authURL,
		showFallback: false,
	}
}

// SetSize updates dimensions
func (v *AuthView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// SetCountdown updates countdown value
func (v *AuthView) SetCountdown(seconds int) {
	v.countdown = seconds
}

// GetCountdown returns current countdown
func (v *AuthView) GetCountdown() int {
	return v.countdown
}

// GetAuthURL returns the auth URL
func (v *AuthView) GetAuthURL() string {
	return v.authURL
}

// ShowFallback enables fallback mode
func (v *AuthView) ShowFallback() {
	v.showFallback = true
}

// IsFallbackShown returns if fallback is shown
func (v *AuthView) IsFallbackShown() bool {
	return v.showFallback
}

// Render the auth view
func (v *AuthView) Render() string {
	if v.showFallback {
		return v.renderFallback()
	}

	sectionWidth := 60

	// Page title
	title := styles.PageTitle("Configuration", v.width)

	// Auth message
	message := styles.Body.Render("Redirecting you to")
	url := styles.Accent.Render(v.authURL)
	
	countdownText := fmt.Sprintf("in %ds", v.countdown)
	countdownStyled := styles.Muted.Render(countdownText)

	// Countdown visual
	var countdownVisual string
	switch v.countdown {
	case 3:
		countdownVisual = styles.Accent.Render("● ● ●")
	case 2:
		countdownVisual = styles.Accent.Render("● ●") + styles.Muted.Render(" ○")
	case 1:
		countdownVisual = styles.Accent.Render("●") + styles.Muted.Render(" ○ ○")
	default:
		countdownVisual = styles.Muted.Render("○ ○ ○")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		message,
		"",
		url,
		"",
		countdownStyled,
		"",
		countdownVisual,
	)

	box := styles.Box.
		Width(sectionWidth).
		Align(lipgloss.Center).
		Render(content)

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.FooterHelp([]components.HelpItem{
			{Key: "m", Desc: "manual entry"},
			{Key: "esc", Desc: "cancel"},
		}))

	// Combine
	fullContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		box,
	)

	centered := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Center,
		fullContent,
	)

	return centered + "\n" + footer
}

func (v *AuthView) renderFallback() string {
	sectionWidth := 60

	title := styles.PageTitle("Configuration", v.width)

	message := styles.Body.Render("Browser auth cancelled.")
	subMessage := styles.Muted.Render("You can enter your API key manually.")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		message,
		"",
		subMessage,
		"",
		styles.Muted.Render("Press Enter to continue to manual entry"),
	)

	box := styles.Box.
		Width(sectionWidth).
		Align(lipgloss.Center).
		Render(content)

	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.FooterHelp([]components.HelpItem{
			{Key: "enter", Desc: "continue"},
			{Key: "esc", Desc: "go back"},
		}))

	fullContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		box,
	)

	centered := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Center,
		fullContent,
	)

	return centered + "\n" + footer
}

// GetHelpItems returns context-specific help
func (v *AuthView) GetHelpItems() []components.HelpItem {
	if v.showFallback {
		return []components.HelpItem{
			{Key: "enter", Desc: "continue to manual entry"},
			{Key: "esc", Desc: "go back to provider selection"},
		}
	}
	return []components.HelpItem{
		{Key: "m", Desc: "skip to manual entry"},
		{Key: "esc", Desc: "cancel and go back"},
	}
}
