package views

import (
	"fmt"
	"skene-terminal-v2/internal/services/config"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// AuthView handles Skene auth with magic link and fallback
type AuthView struct {
	width        int
	height       int
	provider     *config.Provider
	countdown    int // seconds remaining
	authURL      string
	showFallback bool
	header       *components.WizardHeader
	spinner      *components.Spinner
	authState    AuthState
}

// AuthState represents the authentication state
type AuthState int

const (
	AuthStateCountdown AuthState = iota
	AuthStateBrowserOpen
	AuthStateWaiting
	AuthStateSuccess
	AuthStateFallback
)

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
		header:       components.NewWizardHeader(4, "Authentication"),
		spinner:      components.NewSpinner(),
		authState:    AuthStateCountdown,
	}
}

// SetSize updates dimensions
func (v *AuthView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)
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

// SetAuthState updates the auth state
func (v *AuthView) SetAuthState(state AuthState) {
	v.authState = state
}

// ShowFallback enables fallback mode
func (v *AuthView) ShowFallback() {
	v.showFallback = true
	v.authState = AuthStateFallback
}

// IsFallbackShown returns if fallback is shown
func (v *AuthView) IsFallbackShown() bool {
	return v.showFallback
}

// TickSpinner advances the spinner
func (v *AuthView) TickSpinner() {
	v.spinner.Tick()
}

// Render the auth view
func (v *AuthView) Render() string {
	if v.showFallback {
		return v.renderFallback()
	}

	sectionWidth := 60

	// Wizard header
	wizHeader := v.header.Render()

	// Auth content based on state
	var authContent string
	switch v.authState {
	case AuthStateCountdown:
		authContent = v.renderCountdown(sectionWidth)
	case AuthStateBrowserOpen, AuthStateWaiting:
		authContent = v.renderWaiting(sectionWidth)
	case AuthStateSuccess:
		authContent = v.renderSuccess(sectionWidth)
	default:
		authContent = v.renderCountdown(sectionWidth)
	}

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
		wizHeader,
		"",
		"",
		authContent,
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

func (v *AuthView) renderCountdown(width int) string {
	message := styles.Body.Render("Opening browser for Skene authentication")
	url := styles.Accent.Render(v.authURL)

	countdownText := fmt.Sprintf("Redirecting in %ds...", v.countdown)
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

	return styles.Box.
		Width(width).
		Align(lipgloss.Center).
		Render(content)
}

func (v *AuthView) renderWaiting(width int) string {
	message := v.spinner.SpinnerWithText("Waiting for authentication...")
	subMessage := styles.Muted.Render("Complete the login in your browser")
	url := styles.Accent.Render(v.authURL)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		message,
		"",
		subMessage,
		"",
		url,
	)

	return styles.Box.
		Width(width).
		Align(lipgloss.Center).
		Render(content)
}

func (v *AuthView) renderSuccess(width int) string {
	message := styles.SuccessText.Render("✓ Authentication successful!")
	subMessage := styles.Muted.Render("API key received and saved")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		message,
		"",
		subMessage,
	)

	return styles.Box.
		Width(width).
		Align(lipgloss.Center).
		Render(content)
}

func (v *AuthView) renderFallback() string {
	sectionWidth := 60

	wizHeader := v.header.Render()

	message := styles.Body.Render("Browser auth cancelled.")
	subMessage := styles.Muted.Render("You can enter your Skene API key manually.")
	hint := styles.Accent.Render("Press Enter to continue to manual entry")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		message,
		"",
		subMessage,
		"",
		hint,
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
		wizHeader,
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
