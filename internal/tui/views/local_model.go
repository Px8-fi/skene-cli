package views

import (
	"fmt"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// LocalModelStatus represents detection state
type LocalModelStatus int

const (
	LocalModelDetecting LocalModelStatus = iota
	LocalModelFound
	LocalModelNotFound
)

// LocalModelView handles local model runtime detection
type LocalModelView struct {
	width        int
	height       int
	status       LocalModelStatus
	providerName string // "ollama" or "lmstudio"
	models       []string
	selectedIdx  int
	baseURL      string
	spinner      *components.Spinner
	header       *components.WizardHeader
	errorMsg     string
}

// NewLocalModelView creates a new local model view
func NewLocalModelView(providerName string) *LocalModelView {
	var baseURL string
	switch providerName {
	case "ollama":
		baseURL = "http://localhost:11434/v1"
	case "lmstudio":
		baseURL = "http://localhost:1234/v1"
	}

	return &LocalModelView{
		status:       LocalModelDetecting,
		providerName: providerName,
		baseURL:      baseURL,
		spinner:      components.NewSpinner(),
		header:       components.NewWizardHeader(4, "Local Model Setup"),
	}
}

// SetSize updates dimensions
func (v *LocalModelView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)
}

// SetStatus updates the detection status
func (v *LocalModelView) SetStatus(status LocalModelStatus) {
	v.status = status
}

// SetModels sets available models
func (v *LocalModelView) SetModels(models []string) {
	v.models = models
	if len(models) > 0 {
		v.status = LocalModelFound
	} else {
		v.status = LocalModelNotFound
	}
}

// SetError sets error message
func (v *LocalModelView) SetError(msg string) {
	v.errorMsg = msg
	v.status = LocalModelNotFound
}

// TickSpinner advances spinner
func (v *LocalModelView) TickSpinner() {
	v.spinner.Tick()
}

// HandleUp moves model selection up
func (v *LocalModelView) HandleUp() {
	if v.selectedIdx > 0 {
		v.selectedIdx--
	}
}

// HandleDown moves model selection down
func (v *LocalModelView) HandleDown() {
	if v.selectedIdx < len(v.models)-1 {
		v.selectedIdx++
	}
}

// GetSelectedModel returns the selected model name
func (v *LocalModelView) GetSelectedModel() string {
	if v.selectedIdx >= 0 && v.selectedIdx < len(v.models) {
		return v.models[v.selectedIdx]
	}
	return ""
}

// GetBaseURL returns the base URL for the local provider
func (v *LocalModelView) GetBaseURL() string {
	return v.baseURL
}

// IsDetecting returns true if still detecting
func (v *LocalModelView) IsDetecting() bool {
	return v.status == LocalModelDetecting
}

// IsFound returns true if local model runtime was found
func (v *LocalModelView) IsFound() bool {
	return v.status == LocalModelFound
}

// Render the local model view
func (v *LocalModelView) Render() string {
	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	// Wizard header
	wizHeader := v.header.Render()

	// Content based on status
	var mainContent string
	switch v.status {
	case LocalModelDetecting:
		mainContent = v.renderDetecting(sectionWidth)
	case LocalModelFound:
		mainContent = v.renderModelList(sectionWidth)
	case LocalModelNotFound:
		mainContent = v.renderNotFound(sectionWidth)
	}

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.WizardSelectHelp())

	// Combine
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		wizHeader,
		"",
		mainContent,
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

func (v *LocalModelView) renderDetecting(width int) string {
	displayName := v.providerName
	if displayName == "lmstudio" {
		displayName = "LM Studio"
	} else if displayName == "ollama" {
		displayName = "Ollama"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		v.spinner.SpinnerWithText(fmt.Sprintf("Detecting %s runtime...", displayName)),
		"",
		styles.Muted.Render(fmt.Sprintf("Checking %s", v.baseURL)),
	)

	return styles.Box.Width(width).Render(content)
}

func (v *LocalModelView) renderModelList(width int) string {
	header := styles.SectionHeader.Render("Select a local model")

	var items []string
	for i, model := range v.models {
		if i == v.selectedIdx {
			items = append(items, styles.ListItemSelected.Render(model))
		} else {
			items = append(items, styles.ListItem.Render(model))
		}
	}

	list := lipgloss.JoinVertical(lipgloss.Left, items...)

	statusLine := styles.SuccessText.Render(fmt.Sprintf("✓ %d model(s) available at %s", len(v.models), v.baseURL))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		statusLine,
		"",
		list,
	)

	return styles.Box.Width(width).Render(content)
}

func (v *LocalModelView) renderNotFound(width int) string {
	displayName := v.providerName
	var installGuide string

	switch v.providerName {
	case "ollama":
		displayName = "Ollama"
		installGuide = "1. Install: curl -fsSL https://ollama.com/install.sh | sh\n" +
			"2. Start:   ollama serve\n" +
			"3. Pull:    ollama pull llama3.3"
	case "lmstudio":
		displayName = "LM Studio"
		installGuide = "1. Download from: https://lmstudio.ai\n" +
			"2. Load a model in the Developer tab\n" +
			"3. Start the local server"
	}

	header := styles.Error.Render(fmt.Sprintf("✗ %s not detected", displayName))

	errDetail := ""
	if v.errorMsg != "" {
		errDetail = styles.Muted.Render(v.errorMsg) + "\n"
	}

	guideHeader := styles.SectionHeader.Render("Setup Guide")
	guide := styles.Body.Render(installGuide)

	retryHint := styles.Accent.Render("Press 'r' to retry detection or 'esc' to go back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		errDetail,
		guideHeader,
		"",
		guide,
		"",
		retryHint,
	)

	return styles.Box.Width(width).Render(content)
}

// GetHelpItems returns context-specific help
func (v *LocalModelView) GetHelpItems() []components.HelpItem {
	switch v.status {
	case LocalModelDetecting:
		return []components.HelpItem{
			{Key: "ctrl+c", Desc: "quit"},
		}
	case LocalModelNotFound:
		return []components.HelpItem{
			{Key: "r", Desc: "retry detection"},
			{Key: "esc", Desc: "go back"},
			{Key: "ctrl+c", Desc: "quit"},
		}
	default:
		return []components.HelpItem{
			{Key: "↑/↓", Desc: "select model"},
			{Key: "enter", Desc: "confirm"},
			{Key: "esc", Desc: "go back"},
			{Key: "ctrl+c", Desc: "quit"},
		}
	}
}
