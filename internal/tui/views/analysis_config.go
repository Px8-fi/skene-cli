package views

import (
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// AnalysisConfigView allows configuration before analysis
type AnalysisConfigView struct {
	width        int
	height       int
	useDefaults  bool
	selectedIdx  int
	options      []AnalysisOption
	header       *components.WizardHeader
	buttonGroup  *components.ButtonGroup
	providerName string
	modelName    string
	projectDir   string
}

// AnalysisOption represents a configurable option
type AnalysisOption struct {
	Name        string
	Description string
	Enabled     bool
	Value       string
}

// NewAnalysisConfigView creates a new analysis configuration view
func NewAnalysisConfigView(provider, model, projectDir string) *AnalysisConfigView {
	return &AnalysisConfigView{
		useDefaults:  true,
		selectedIdx:  0,
		providerName: provider,
		modelName:    model,
		projectDir:   projectDir,
		header:       components.NewWizardHeader(6, "Analysis Configuration"),
		buttonGroup:  components.YesNoButtons(true),
		options: []AnalysisOption{
			{Name: "Product docs", Description: "Generate product documentation", Enabled: false},
			{Name: "Business type", Description: "Specify business type for template", Enabled: false, Value: "auto-detect"},
			{Name: "Verbose output", Description: "Show detailed analysis progress", Enabled: true},
		},
	}
}

// SetSize updates dimensions
func (v *AnalysisConfigView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)
}

// HandleUp moves selection up
func (v *AnalysisConfigView) HandleUp() {
	if !v.useDefaults && v.selectedIdx > 0 {
		v.selectedIdx--
	}
}

// HandleDown moves selection down
func (v *AnalysisConfigView) HandleDown() {
	if !v.useDefaults && v.selectedIdx < len(v.options)-1 {
		v.selectedIdx++
	}
}

// HandleLeft handles left for button group
func (v *AnalysisConfigView) HandleLeft() {
	if v.useDefaults {
		v.buttonGroup.Previous()
	}
}

// HandleRight handles right for button group
func (v *AnalysisConfigView) HandleRight() {
	if v.useDefaults {
		v.buttonGroup.Next()
	}
}

// HandleSpace toggles option
func (v *AnalysisConfigView) HandleSpace() {
	if !v.useDefaults && v.selectedIdx >= 0 && v.selectedIdx < len(v.options) {
		v.options[v.selectedIdx].Enabled = !v.options[v.selectedIdx].Enabled
	}
}

// GetButtonLabel returns selected button label
func (v *AnalysisConfigView) GetButtonLabel() string {
	return v.buttonGroup.GetActiveLabel()
}

// SetCustomMode enables custom configuration mode
func (v *AnalysisConfigView) SetCustomMode() {
	v.useDefaults = false
}

// IsDefaultMode returns if using default settings
func (v *AnalysisConfigView) IsDefaultMode() bool {
	return v.useDefaults
}

// GetProductDocs returns if product docs should be generated
func (v *AnalysisConfigView) GetProductDocs() bool {
	return v.options[0].Enabled
}

// GetVerbose returns if verbose output is enabled
func (v *AnalysisConfigView) GetVerbose() bool {
	return v.options[2].Enabled
}

// Render the analysis config view
func (v *AnalysisConfigView) Render() string {
	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	// Wizard header
	wizHeader := lipgloss.NewStyle().Width(sectionWidth).Render(v.header.Render())

	// Summary section
	summarySection := v.renderSummary(sectionWidth)

	var actionSection string
	if v.useDefaults {
		actionSection = v.renderDefaultQuestion(sectionWidth)
	} else {
		actionSection = v.renderCustomOptions(sectionWidth)
	}

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.WizardHelp())

	// Combine
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		wizHeader,
		"",
		summarySection,
		"",
		actionSection,
	)

	padded := lipgloss.NewStyle().PaddingTop(2).Render(content)

	centered := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Top,
		padded,
	)

	return centered + "\n" + footer
}

func (v *AnalysisConfigView) renderSummary(width int) string {
	header := styles.SectionHeader.Render("Analysis Summary")

	rows := []string{
		styles.Label.Render("Provider:   ") + styles.Body.Render(v.providerName),
		styles.Label.Render("Model:      ") + styles.Body.Render(v.modelName),
		styles.Label.Render("Directory:  ") + styles.Body.Render(v.projectDir),
		styles.Label.Render("Output:     ") + styles.Body.Render("./skene-context/"),
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		lipgloss.JoinVertical(lipgloss.Left, rows...),
	)

	return styles.Box.Width(width).Render(content)
}

func (v *AnalysisConfigView) renderDefaultQuestion(width int) string {
	question := styles.Body.Render("Use recommended settings?")
	buttons := v.buttonGroup.Render()

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		question,
		"",
		buttons,
	)

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(content)
}

func (v *AnalysisConfigView) renderCustomOptions(width int) string {
	header := styles.SectionHeader.Render("Advanced Configuration")

	var items []string
	for i, opt := range v.options {
		isSelected := i == v.selectedIdx

		// Checkbox
		var checkbox string
		if opt.Enabled {
			checkbox = styles.SuccessText.Render("[✓]")
		} else {
			checkbox = styles.Muted.Render("[ ]")
		}

		var nameStyle, descStyle lipgloss.Style
		if isSelected {
			nameStyle = lipgloss.NewStyle().Foreground(styles.Amber)
			descStyle = lipgloss.NewStyle().Foreground(styles.Sand)
		} else {
			nameStyle = styles.Body
			descStyle = styles.Muted
		}

		line := checkbox + " " + nameStyle.Render(opt.Name) + "  " + descStyle.Render(opt.Description)
		items = append(items, line)
	}

	list := lipgloss.JoinVertical(lipgloss.Left, items...)
	hint := styles.Muted.Render("Space to toggle • Enter to start analysis")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		list,
		"",
		hint,
	)

	return styles.Box.Width(width).Render(content)
}

// GetHelpItems returns context-specific help
func (v *AnalysisConfigView) GetHelpItems() []components.HelpItem {
	if v.useDefaults {
		return []components.HelpItem{
			{Key: "←/→", Desc: "select option"},
			{Key: "enter", Desc: "confirm"},
			{Key: "esc", Desc: "go back"},
		}
	}
	return []components.HelpItem{
		{Key: "↑/↓", Desc: "navigate"},
		{Key: "space", Desc: "toggle option"},
		{Key: "enter", Desc: "start analysis"},
		{Key: "esc", Desc: "go back"},
	}
}
