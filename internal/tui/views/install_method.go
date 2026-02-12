package views

import (
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// InstallMethodOption represents an installation method choice
type InstallMethodOption struct {
	ID          string
	Name        string
	Description string
	Detail      string
	Command     string
	Recommended bool
}

// InstallMethodView allows selection of installation method
type InstallMethodView struct {
	width         int
	height        int
	options       []InstallMethodOption
	selectedIndex int
	header        *components.WizardHeader
	uvAvailable   bool
}

// NewInstallMethodView creates a new install method selection view
func NewInstallMethodView(uvAvailable bool) *InstallMethodView {
	options := []InstallMethodOption{}

	return &InstallMethodView{
		options:       options,
		selectedIndex: 0,
		header:        components.NewWizardHeader(2, "Install Method"),
		uvAvailable:   uvAvailable,
	}
}

// SetSize updates dimensions
func (v *InstallMethodView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)
}

// HandleUp moves selection up
func (v *InstallMethodView) HandleUp() {
	if v.selectedIndex > 0 {
		v.selectedIndex--
	}
}

// HandleDown moves selection down
func (v *InstallMethodView) HandleDown() {
	if v.selectedIndex < len(v.options)-1 {
		v.selectedIndex++
	}
}

// GetSelectedMethod returns the selected install method
func (v *InstallMethodView) GetSelectedMethod() string {
	if v.selectedIndex >= 0 && v.selectedIndex < len(v.options) {
		return v.options[v.selectedIndex].ID
	}
	return "uvx"
}

// Render the install method selection view
func (v *InstallMethodView) Render() string {
	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	// Wizard header
	wizHeader := lipgloss.NewStyle().Width(sectionWidth).Render(v.header.Render())

	// Options list
	optionsSection := v.renderOptions(sectionWidth)

	// Info box
	infoSection := v.renderInfo(sectionWidth)

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
		optionsSection,
		"",
		infoSection,
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

func (v *InstallMethodView) renderOptions(width int) string {
	header := styles.SectionHeader.Render("How would you like to run skene-growth?")

	var items []string
	for i, opt := range v.options {
		isSelected := i == v.selectedIndex

		var name, desc, detail string
		if isSelected {
			name = styles.ListItemSelected.Render(opt.Name)
			desc = styles.ListDescriptionSelected.Render(opt.Description)
			detail = styles.ListDescriptionSelected.Render(opt.Detail)
		} else {
			name = styles.ListItem.Render(opt.Name)
			desc = styles.ListDescription.Render(opt.Description)
			detail = styles.ListDescription.Render(opt.Detail)
		}

		// Add recommended badge
		if opt.Recommended {
			name += "  " + styles.SuccessText.Render("[recommended]")
		}

		// UV not available warning
		if opt.ID == "uvx" && !v.uvAvailable {
			name += "  " + styles.Error.Render("[uv not installed]")
		}

		item := name + "\n" + desc + "\n" + detail
		items = append(items, item)

		// Add spacing between items (but not after last)
		if i < len(v.options)-1 {
			items = append(items, "")
		}
	}

	list := lipgloss.JoinVertical(lipgloss.Left, items...)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		list,
	)

	return styles.Box.Width(width).Render(content)
}

func (v *InstallMethodView) renderInfo(width int) string {
	selected := v.options[v.selectedIndex]

	info := styles.Muted.Render("Command: ") + styles.Accent.Render(selected.Command)

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(info)
}

// GetHelpItems returns context-specific help
func (v *InstallMethodView) GetHelpItems() []components.HelpItem {
	return []components.HelpItem{
		{Key: "↑/↓", Desc: "select method"},
		{Key: "enter", Desc: "confirm"},
		{Key: "esc", Desc: "go back"},
		{Key: "ctrl+c", Desc: "quit"},
	}
}
