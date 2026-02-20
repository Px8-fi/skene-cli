package views

import (
	"skene/internal/constants"
	"skene/internal/tui/components"
	"skene/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// SkenePackage represents a Skene ecosystem package
type SkenePackage struct {
	ID          string
	Name        string
	Description string
	URL         string
	Enabled     bool
}

// AnalysisConfigView allows configuration before analysis
type AnalysisConfigView struct {
	width        int
	height       int
	useDefaults  bool
	selectedIdx  int
	packages     []SkenePackage
	header       *components.WizardHeader
	buttonGroup  *components.ButtonGroup
	providerName string
	modelName    string
	projectDir   string
}

// NewAnalysisConfigView creates a new analysis configuration view
func NewAnalysisConfigView(provider, model, projectDir string) *AnalysisConfigView {
	return &AnalysisConfigView{
		useDefaults:  true,
		selectedIdx:  0,
		providerName: provider,
		modelName:    model,
		projectDir:   projectDir,
		header:       components.NewWizardHeader(3, constants.StepNameAnalysisConfig),
		buttonGroup:  components.YesNoButtons(true),
		packages: func() []SkenePackage {
			var pkgs []SkenePackage
			for _, p := range constants.SkenePackages {
				pkgs = append(pkgs, SkenePackage{
					ID:          p.ID,
					Name:        p.Name,
					Description: p.Description,
					URL:         p.URL,
					Enabled:     true,
				})
			}
			return pkgs
		}(),
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
	if !v.useDefaults {
		if v.selectedIdx > 0 {
			v.selectedIdx--
		}
	}
}

// HandleDown moves selection down
func (v *AnalysisConfigView) HandleDown() {
	if !v.useDefaults {
		if v.selectedIdx < len(v.packages)-1 {
			v.selectedIdx++
		}
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
	if !v.useDefaults {
		if v.selectedIdx >= 0 && v.selectedIdx < len(v.packages) {
			v.packages[v.selectedIdx].Enabled = !v.packages[v.selectedIdx].Enabled
		}
	}
}

// GetButtonLabel returns selected button label
func (v *AnalysisConfigView) GetButtonLabel() string {
	return v.buttonGroup.GetActiveLabel()
}

// SetCustomMode enables custom configuration mode
func (v *AnalysisConfigView) SetCustomMode() {
	v.useDefaults = false
	v.selectedIdx = 0
}

// IsDefaultMode returns if using default settings
func (v *AnalysisConfigView) IsDefaultMode() bool {
	return v.useDefaults
}

// GetUseGrowth returns if Skene Growth is enabled
func (v *AnalysisConfigView) GetUseGrowth() bool {
	if len(v.packages) > 0 {
		return v.packages[0].Enabled
	}
	return true
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
	header := styles.SectionHeader.Render(constants.AnalysisConfigSummary)

	rows := []string{
		styles.Label.Render("Provider:   ") + styles.Body.Render(v.providerName),
		styles.Label.Render("Model:      ") + styles.Body.Render(v.modelName),
		styles.Label.Render("Directory:  ") + styles.Body.Render(v.projectDir),
		styles.Label.Render("Output:     ") + styles.Body.Render(constants.DefaultOutputDir+"/"),
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
	question := styles.Body.Render(constants.AnalysisConfigQuestion)
	desc := styles.Muted.Render(constants.AnalysisConfigDefault)
	buttons := v.buttonGroup.Render()

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		question,
		desc,
		"",
		buttons,
	)

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(content)
}

func (v *AnalysisConfigView) renderCustomOptions(width int) string {
	header := styles.Body.Render(constants.AnalysisConfigSelectPkgs)

	var pkgItems []string
	for i, pkg := range v.packages {
		isSelected := i == v.selectedIdx

		var checkbox string
		if pkg.Enabled {
			checkbox = styles.SuccessText.Render("✓")
		} else {
			checkbox = styles.Muted.Render("✗")
		}

		var nameStyle, descStyle lipgloss.Style
		if isSelected {
			nameStyle = styles.ListItemSelected
			descStyle = styles.ListDescriptionSelected
		} else {
			nameStyle = styles.ListItem
			descStyle = styles.ListDescription
		}

		// Use the styles.ListItemSelected which has a left border
		// We need to be careful not to double up on padding if the style already has it
		// styles.ListItemSelected has PaddingLeft(1)
		// styles.ListItem has PaddingLeft(2)

		line := checkbox + "  " + pkg.Name
		nameLine := nameStyle.Render(line)
		descLine := descStyle.Render(pkg.Description)

		item := lipgloss.JoinVertical(lipgloss.Left, nameLine, descLine)
		if i < len(v.packages)-1 {
			item += "\n"
		}
		pkgItems = append(pkgItems, item)
	}

	pkgList := lipgloss.JoinVertical(lipgloss.Left, pkgItems...)

	// Start hint
	hint := styles.Muted.Render(constants.AnalysisConfigToggleHint)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		header,
		"",
		pkgList,
		"",
		hint,
	)

	return styles.Box.Width(width).Render(content)
}

// GetHelpItems returns context-specific help
func (v *AnalysisConfigView) GetHelpItems() []components.HelpItem {
	if v.useDefaults {
		return []components.HelpItem{
			{Key: constants.HelpKeyLeftRight, Desc: constants.HelpDescSelectOption},
			{Key: constants.HelpKeyEnter, Desc: constants.HelpDescConfirm},
			{Key: constants.HelpKeyEsc, Desc: constants.HelpDescGoBack},
		}
	}
	return []components.HelpItem{
		{Key: constants.HelpKeyUpDown, Desc: constants.HelpDescNavigate},
		{Key: constants.HelpKeySpace, Desc: constants.HelpDescToggleOption},
		{Key: constants.HelpKeyEnter, Desc: constants.HelpDescStartAnalysis},
		{Key: constants.HelpKeyEsc, Desc: constants.HelpDescGoBack},
	}
}
