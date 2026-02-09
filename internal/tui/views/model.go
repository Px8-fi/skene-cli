package views

import (
	"fmt"
	"skene-terminal-v2/internal/services/config"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// ModelView handles model selection for a provider
type ModelView struct {
	width         int
	height        int
	provider      *config.Provider
	selectedIndex int
	buttonGroup   *components.ButtonGroup
	buttonFocused bool
	header        *components.WizardHeader
}

// NewModelView creates a new model view
func NewModelView(provider *config.Provider) *ModelView {
	return &ModelView{
		provider:      provider,
		selectedIndex: 0,
		buttonGroup:   components.NavigationButtons(false),
		buttonFocused: false,
		header:        components.NewWizardHeader(3, "Select Model"),
	}
}

// SetProvider updates the provider
func (v *ModelView) SetProvider(provider *config.Provider) {
	v.provider = provider
	v.selectedIndex = 0
	v.buttonFocused = false
}

// SetSize updates dimensions
func (v *ModelView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)
}

// HandleUp moves selection up
func (v *ModelView) HandleUp() {
	if v.buttonFocused {
		v.buttonFocused = false
		return
	}
	if v.selectedIndex > 0 {
		v.selectedIndex--
	}
}

// HandleDown moves selection down
func (v *ModelView) HandleDown() {
	if v.provider == nil {
		return
	}
	if !v.buttonFocused && v.selectedIndex < len(v.provider.Models)-1 {
		v.selectedIndex++
	} else if !v.buttonFocused {
		v.buttonFocused = true
	}
}

// HandleLeft moves button focus
func (v *ModelView) HandleLeft() {
	if v.buttonFocused {
		v.buttonGroup.Previous()
	}
}

// HandleRight moves button focus
func (v *ModelView) HandleRight() {
	if v.buttonFocused {
		v.buttonGroup.Next()
	}
}

// IsButtonFocused returns if buttons are focused
func (v *ModelView) IsButtonFocused() bool {
	return v.buttonFocused
}

// GetSelectedModel returns the selected model
func (v *ModelView) GetSelectedModel() *config.Model {
	if v.provider == nil || v.selectedIndex < 0 || v.selectedIndex >= len(v.provider.Models) {
		return nil
	}
	return &v.provider.Models[v.selectedIndex]
}

// GetButtonLabel returns selected button label
func (v *ModelView) GetButtonLabel() string {
	return v.buttonGroup.GetActiveLabel()
}

// Render the model view
func (v *ModelView) Render() string {
	if v.provider == nil {
		return "No provider selected"
	}

	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	// Wizard header
	wizHeader := v.header.Render()

	// Model list section
	listSection := v.renderModelList(sectionWidth)

	// Buttons
	buttons := v.buttonGroup.Render()
	buttonsCentered := lipgloss.NewStyle().
		Width(sectionWidth).
		Align(lipgloss.Right).
		Render(buttons)

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
		listSection,
		"",
		buttonsCentered,
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

func (v *ModelView) renderModelList(width int) string {
	header := styles.SectionHeader.Render(fmt.Sprintf("Select Model for %s", v.provider.Name))

	// Model count
	count := styles.Muted.Render(fmt.Sprintf("%d / %d models", v.selectedIndex+1, len(v.provider.Models)))

	// Model list
	var items []string
	for i, m := range v.provider.Models {
		isSelected := i == v.selectedIndex && !v.buttonFocused

		var item string
		if isSelected {
			name := styles.ListItemSelected.Render(m.Name)
			desc := styles.ListDescriptionSelected.Render(m.Description)
			item = name + "\n" + desc
		} else {
			name := styles.ListItem.Render(m.Name)
			desc := styles.ListDescription.Render(m.Description)
			item = name + "\n" + desc
		}
		items = append(items, item)
	}

	list := lipgloss.JoinVertical(lipgloss.Left, items...)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		count,
		"",
		list,
	)

	return styles.Box.Width(width).Render(content)
}

// GetHelpItems returns context-specific help
func (v *ModelView) GetHelpItems() []components.HelpItem {
	return []components.HelpItem{
		{Key: "↑/↓", Desc: "select model"},
		{Key: "enter", Desc: "confirm selection"},
		{Key: "esc", Desc: "go back"},
		{Key: "ctrl+c", Desc: "quit"},
	}
}
