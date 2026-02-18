package views

import (
	"skene/internal/tui/components"
	"skene/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// NextStepAction represents an available next step
type NextStepAction struct {
	ID          string
	Name        string
	Description string
	Command     string
}

// NextStepsView presents available next steps after analysis
type NextStepsView struct {
	width       int
	height      int
	actions     []NextStepAction
	selectedIdx int
	header      *components.WizardHeader
}

// NewNextStepsView creates a new next steps view
func NewNextStepsView() *NextStepsView {
	return &NextStepsView{
		selectedIdx: 0,
		header:      components.NewWizardHeader(7, "Next Steps"),
		actions: []NextStepAction{
			{
				ID:          "plan",
				Name:        "Generate Growth Plan",
				Description: "Create a prioritized growth plan with implementation roadmap",
				Command:     "uvx skene-growth plan",
			},
			{
				ID:          "build",
				Name:        "Build Implementation Prompt",
				Description: "Generate a ready-to-use prompt for Cursor, Claude, or other AI tools",
				Command:     "uvx skene-growth build",
			},
			{
				ID:          "validate",
				Name:        "Validate Manifest",
				Description: "Validate the growth manifest against the schema",
				Command:     "uvx skene-growth validate",
			},
			{
				ID:          "rerun",
				Name:        "Re-run Analysis",
				Description: "Analyze the codebase again with the current configuration",
				Command:     "uvx skene-growth analyze .",
			},
			{
				ID:          "open",
				Name:        "Open Generated Files",
				Description: "View the analysis output in ./skene-context/",
				Command:     "",
			},
			{
				ID:          "config",
				Name:        "Change Configuration",
				Description: "Modify provider, model, or project settings",
				Command:     "",
			},
			{
				ID:          "exit",
				Name:        "Exit",
				Description: "Close Skene CLI",
				Command:     "",
			},
		},
	}
}

// SetSize updates dimensions
func (v *NextStepsView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)
}

// HandleUp moves selection up
func (v *NextStepsView) HandleUp() {
	if v.selectedIdx > 0 {
		v.selectedIdx--
	}
}

// HandleDown moves selection down
func (v *NextStepsView) HandleDown() {
	if v.selectedIdx < len(v.actions)-1 {
		v.selectedIdx++
	}
}

// GetSelectedAction returns the selected action
func (v *NextStepsView) GetSelectedAction() *NextStepAction {
	if v.selectedIdx >= 0 && v.selectedIdx < len(v.actions) {
		return &v.actions[v.selectedIdx]
	}
	return nil
}

// Render the next steps view
func (v *NextStepsView) Render() string {
	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	// Wizard header
	wizHeader := lipgloss.NewStyle().Width(sectionWidth).Render(v.header.Render())

	// Success message
	successMsg := styles.SuccessText.Render("Analysis complete! What would you like to do next?")

	// Actions list
	actionsSection := v.renderActions(sectionWidth)

	// Command preview
	commandPreview := v.renderCommandPreview(sectionWidth)

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.WizardSelectHelp())

	// Combine
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		wizHeader,
		"",
		successMsg,
		"",
		actionsSection,
		"",
		commandPreview,
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

func (v *NextStepsView) renderActions(width int) string {
	var items []string

	for i, action := range v.actions {
		isSelected := i == v.selectedIdx

		var name, desc string
		if isSelected {
			name = styles.ListItemSelected.Render(action.Name)
			desc = styles.ListDescriptionSelected.Render(action.Description)
		} else {
			name = styles.ListItem.Render(action.Name)
			desc = styles.ListDescription.Render(action.Description)
		}

		item := name + "\n" + desc
		items = append(items, item)

		// Add spacing between items (but not after last)
		if i < len(v.actions)-1 {
			items = append(items, "")
		}
	}

	list := lipgloss.JoinVertical(lipgloss.Left, items...)
	return styles.Box.Width(width).Render(list)
}

func (v *NextStepsView) renderCommandPreview(width int) string {
	action := v.GetSelectedAction()
	if action == nil || action.Command == "" {
		return ""
	}

	cmdLabel := styles.Muted.Render("Command: ")
	cmdValue := styles.Accent.Render(action.Command)
	preview := cmdLabel + cmdValue
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(preview)
}

// GetHelpItems returns context-specific help
func (v *NextStepsView) GetHelpItems() []components.HelpItem {
	return []components.HelpItem{
		{Key: "↑/↓", Desc: "navigate"},
		{Key: "enter", Desc: "select"},
		{Key: "esc", Desc: "back to results"},
		{Key: "ctrl+c", Desc: "quit"},
	}
}
