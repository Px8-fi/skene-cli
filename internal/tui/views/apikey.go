package views

import (
	"fmt"
	"skene-terminal-v2/internal/services/config"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// APIKeyView handles API key entry
type APIKeyView struct {
	width       int
	height      int
	provider    *config.Provider
	model       *config.Model
	textInput   textinput.Model
	buttonGroup *components.ButtonGroup
	inputFocus  bool
	error       string
}

// NewAPIKeyView creates a new API key view
func NewAPIKeyView(provider *config.Provider, model *config.Model) *APIKeyView {
	ti := textinput.New()
	ti.Placeholder = "Enter API Key"
	ti.CharLimit = 128
	ti.Width = 40
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	ti.Focus()

	return &APIKeyView{
		provider:    provider,
		model:       model,
		textInput:   ti,
		buttonGroup: components.NavigationButtons(false),
		inputFocus:  true,
	}
}

// SetProvider updates provider and model
func (v *APIKeyView) SetProvider(provider *config.Provider, model *config.Model) {
	v.provider = provider
	v.model = model
}

// SetSize updates dimensions
func (v *APIKeyView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// Update handles text input updates
func (v *APIKeyView) Update(msg interface{}) {
	if v.inputFocus {
		var cmd interface{}
		v.textInput, cmd = v.textInput.Update(msg)
		_ = cmd
	}
}

// HandleTab toggles between input and buttons
func (v *APIKeyView) HandleTab() {
	v.inputFocus = !v.inputFocus
	if v.inputFocus {
		v.textInput.Focus()
	} else {
		v.textInput.Blur()
	}
}

// HandleLeft moves button focus
func (v *APIKeyView) HandleLeft() {
	if !v.inputFocus {
		v.buttonGroup.Previous()
	}
}

// HandleRight moves button focus
func (v *APIKeyView) HandleRight() {
	if !v.inputFocus {
		v.buttonGroup.Next()
	}
}

// IsInputFocused returns if input is focused
func (v *APIKeyView) IsInputFocused() bool {
	return v.inputFocus
}

// GetAPIKey returns the entered API key
func (v *APIKeyView) GetAPIKey() string {
	return v.textInput.Value()
}

// GetButtonLabel returns selected button label
func (v *APIKeyView) GetButtonLabel() string {
	return v.buttonGroup.GetActiveLabel()
}

// Validate checks if the API key is valid
func (v *APIKeyView) Validate() bool {
	key := v.textInput.Value()
	if len(key) < 10 {
		v.error = "API key is too short"
		return false
	}
	v.error = ""
	return true
}

// GetTextInput returns the text input model
func (v *APIKeyView) GetTextInput() *textinput.Model {
	return &v.textInput
}

// Render the API key view
func (v *APIKeyView) Render() string {
	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	// Page title
	title := styles.PageTitle("Configuration", v.width)

	// Main content section
	contentSection := v.renderContent(sectionWidth)

	// Step indicator and buttons
	stepIndicator := components.StepIndicator(3, 4)
	buttons := v.buttonGroup.Render()

	bottomBar := lipgloss.JoinHorizontal(
		lipgloss.Center,
		stepIndicator,
		"          ",
		buttons,
	)

	bottomBarCentered := lipgloss.NewStyle().
		Width(sectionWidth).
		Align(lipgloss.Right).
		Render(bottomBar)

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.InputHelp())

	// Combine
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		contentSection,
		"",
		bottomBarCentered,
	)

	centered := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	return centered + "\n" + footer
}

func (v *APIKeyView) renderContent(width int) string {
	header := styles.SectionHeader.Render("Selected Values")

	// Show selected provider and model
	providerName := ""
	modelName := ""
	if v.provider != nil {
		providerName = v.provider.Name
	}
	if v.model != nil {
		modelName = v.model.Name
	}

	// Table header
	tableHeader := styles.TableHeader.Render(
		fmt.Sprintf("%-12s %s", "Type", "Value"),
	)
	sep := styles.TableSeparator.Render(strings.Repeat("─", width-8))

	rows := []string{
		fmt.Sprintf("%-12s %s", styles.Body.Render("provider"), styles.Body.Render(providerName)),
		fmt.Sprintf("%-12s %s", styles.Body.Render("model"), styles.Body.Render(modelName)),
	}

	// API Key input
	apiKeyLabel := styles.Label.Render("[API Key]:")
	inputPrompt := styles.Muted.Render("> ")
	inputField := v.textInput.View()

	// Error message if any
	errorMsg := ""
	if v.error != "" {
		errorMsg = "\n" + styles.Error.Render(v.error)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		tableHeader,
		sep,
		strings.Join(rows, "\n"),
		"",
		"",
		apiKeyLabel+"     "+inputPrompt+inputField+errorMsg,
	)

	return styles.Box.Width(width).Render(content)
}

// GetHelpItems returns context-specific help
func (v *APIKeyView) GetHelpItems() []components.HelpItem {
	return []components.HelpItem{
		{Key: "enter", Desc: "submit key"},
		{Key: "tab", Desc: "switch focus"},
		{Key: "esc", Desc: "go back"},
		{Key: "q", Desc: "quit"},
	}
}
