package views

import (
	"fmt"
	"skene-terminal-v2/internal/services/config"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// ProviderView handles provider selection
type ProviderView struct {
	width         int
	height        int
	providers     []config.Provider
	selectedIndex int
	scrollOffset  int
	maxVisible    int
	buttonGroup   *components.ButtonGroup
	buttonFocused bool
	header        *components.WizardHeader
}

// NewProviderView creates a new provider view
func NewProviderView() *ProviderView {
	return &ProviderView{
		providers:     config.GetProviders(),
		selectedIndex: 0,
		maxVisible:    7,
		buttonGroup:   components.NavigationButtons(false),
		buttonFocused: false,
		header:        components.NewWizardHeader(3, "AI Provider"),
	}
}

// SetSize updates dimensions
func (v *ProviderView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)
	// Adjust max visible based on height
	v.maxVisible = (height - 18) / 3
	if v.maxVisible < 3 {
		v.maxVisible = 3
	}
	if v.maxVisible > 10 {
		v.maxVisible = 10
	}
}

// HandleUp moves selection up
func (v *ProviderView) HandleUp() {
	if v.buttonFocused {
		v.buttonFocused = false
		return
	}
	if v.selectedIndex > 0 {
		v.selectedIndex--
		if v.selectedIndex < v.scrollOffset {
			v.scrollOffset = v.selectedIndex
		}
	}
}

// HandleDown moves selection down
func (v *ProviderView) HandleDown() {
	if !v.buttonFocused && v.selectedIndex < len(v.providers)-1 {
		v.selectedIndex++
		if v.selectedIndex >= v.scrollOffset+v.maxVisible {
			v.scrollOffset = v.selectedIndex - v.maxVisible + 1
		}
	} else if !v.buttonFocused {
		v.buttonFocused = true
	}
}

// HandleLeft moves button focus
func (v *ProviderView) HandleLeft() {
	if v.buttonFocused {
		v.buttonGroup.Previous()
	}
}

// HandleRight moves button focus
func (v *ProviderView) HandleRight() {
	if v.buttonFocused {
		v.buttonGroup.Next()
	}
}

// IsButtonFocused returns if buttons are focused
func (v *ProviderView) IsButtonFocused() bool {
	return v.buttonFocused
}

// GetSelectedProvider returns the selected provider
func (v *ProviderView) GetSelectedProvider() *config.Provider {
	if v.selectedIndex >= 0 && v.selectedIndex < len(v.providers) {
		return &v.providers[v.selectedIndex]
	}
	return nil
}

// GetButtonLabel returns selected button label
func (v *ProviderView) GetButtonLabel() string {
	return v.buttonGroup.GetActiveLabel()
}

// Render the provider view
func (v *ProviderView) Render() string {
	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	// Wizard header
	wizHeader := v.header.Render()

	// Provider list section
	listSection := v.renderProviderList(sectionWidth)

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

func (v *ProviderView) renderProviderList(width int) string {
	header := styles.SectionHeader.Render("Select AI Provider")

	// Provider count
	count := styles.Muted.Render(fmt.Sprintf("%d / %d providers", v.selectedIndex+1, len(v.providers)))

	// Provider list
	var items []string

	endIdx := v.scrollOffset + v.maxVisible
	if endIdx > len(v.providers) {
		endIdx = len(v.providers)
	}

	for i := v.scrollOffset; i < endIdx; i++ {
		p := v.providers[i]
		isSelected := i == v.selectedIndex && !v.buttonFocused

		var item string
		if isSelected {
			name := styles.ListItemSelected.Render(p.Name)
			desc := styles.ListDescriptionSelected.Render(p.Description)
			item = name + "\n" + desc
		} else {
			name := styles.ListItem.Render(p.Name)
			desc := styles.ListDescription.Render(p.Description)
			item = name + "\n" + desc
		}

		// Add badges
		if p.IsLocal {
			item += "  " + styles.Muted.Render("[local]")
		}
		if p.IsGeneric {
			item += "  " + styles.Muted.Render("[custom endpoint]")
		}

		items = append(items, item)

		// Add spacing between items (but not after last)
		if i < endIdx-1 {
			items = append(items, "")
		}
	}

	// Scroll indicators
	if v.scrollOffset > 0 {
		items = append([]string{styles.Muted.Render("  ↑ more above")}, items...)
	}
	if endIdx < len(v.providers) {
		items = append(items, styles.Muted.Render("  ↓ more below"))
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
func (v *ProviderView) GetHelpItems() []components.HelpItem {
	return []components.HelpItem{
		{Key: "↑/↓", Desc: "select provider"},
		{Key: "enter", Desc: "confirm selection"},
		{Key: "esc", Desc: "go back"},
		{Key: "ctrl+c", Desc: "quit"},
	}
}
