package views

import (
	"fmt"
	"skene-terminal-v2/internal/services/config"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ConfigView displays existing configuration and asks to edit
type ConfigView struct {
	width        int
	height       int
	configMgr    *config.Manager
	buttonGroup  *components.ButtonGroup
	configStatus []config.ConfigStatus
}

// NewConfigView creates a new config view
func NewConfigView(mgr *config.Manager) *ConfigView {
	return &ConfigView{
		configMgr:    mgr,
		buttonGroup:  components.YesNoButtons(true),
		configStatus: mgr.CheckConfigs(),
	}
}

// SetSize updates dimensions
func (v *ConfigView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// HandleLeft moves button focus left
func (v *ConfigView) HandleLeft() {
	v.buttonGroup.Previous()
}

// HandleRight moves button focus right
func (v *ConfigView) HandleRight() {
	v.buttonGroup.Next()
}

// GetSelectedButton returns the selected button label
func (v *ConfigView) GetSelectedButton() string {
	return v.buttonGroup.GetActiveLabel()
}

// Render the config view
func (v *ConfigView) Render() string {
	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	// Page title
	title := styles.PageTitle("Configuration", v.width)

	// Config files section
	configFilesSection := v.renderConfigFilesSection(sectionWidth)

	// Current values section
	currentValuesSection := v.renderCurrentValuesSection(sectionWidth)

	// Modal with question and buttons
	modalSection := v.renderModal(60)

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.NavHelp())

	// Combine all sections
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		configFilesSection,
		"",
		currentValuesSection,
		"",
		"",
		modalSection,
	)

	// Center content and add footer
	centered := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	return centered + "\n" + footer
}

func (v *ConfigView) renderConfigFilesSection(width int) string {
	// Header
	header := styles.SectionHeader.Render("Config files")

	// Table header
	tableHeader := styles.TableHeader.Render(
		fmt.Sprintf("%-10s %-40s %s", "Type", "Path", "Status"),
	)

	// Separator
	sep := styles.TableSeparator.Render(strings.Repeat("─", width-8))

	// Rows
	var rows []string
	for _, status := range v.configStatus {
		path := config.GetShortenedPath(status.Path, 38)
		statusText := "not found"
		statusStyle := styles.Muted
		if status.Found {
			statusText = "found"
			statusStyle = styles.SuccessText
		}

		row := fmt.Sprintf("%-10s %-40s %s",
			styles.Body.Render(status.Type),
			styles.Body.Render(path),
			statusStyle.Render(statusText),
		)
		rows = append(rows, row)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		tableHeader,
		sep,
		strings.Join(rows, "\n"),
	)

	return styles.Box.Width(width).Render(content)
}

func (v *ConfigView) renderCurrentValuesSection(width int) string {
	header := styles.SectionHeader.Render("Current Values")

	// Table header
	tableHeader := styles.TableHeader.Render(
		fmt.Sprintf("%-12s %s", "Type", "Value"),
	)

	sep := styles.TableSeparator.Render(strings.Repeat("─", width-8))

	// Get current config values
	cfg := v.configMgr.Config

	rows := []string{
		fmt.Sprintf("%-12s %s", styles.Body.Render("api_key"), styles.Body.Render(v.configMgr.GetMaskedAPIKey())),
		fmt.Sprintf("%-12s %s", styles.Body.Render("provider"), styles.Body.Render(cfg.Provider)),
		fmt.Sprintf("%-12s %s", styles.Body.Render("model"), styles.Body.Render(cfg.Model)),
		fmt.Sprintf("%-12s %s", styles.Body.Render("output_dir"), styles.Body.Render(cfg.OutputDir)),
		fmt.Sprintf("%-12s %s", styles.Body.Render("verbose"), styles.Body.Render(fmt.Sprintf("%v", cfg.Verbose))),
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		tableHeader,
		sep,
		strings.Join(rows, "\n"),
	)

	return styles.Box.Width(width).Render(content)
}

func (v *ConfigView) renderModal(width int) string {
	question := styles.Body.Copy().
		Width(width - 8).
		Align(lipgloss.Center).
		Render("Do you want to edit this configuration?")

	buttons := lipgloss.NewStyle().
		Width(width - 8).
		Align(lipgloss.Center).
		Render(v.buttonGroup.Render())

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		question,
		"",
		buttons,
	)

	return styles.Box.Width(width).Render(content)
}

// GetHelpItems returns context-specific help
func (v *ConfigView) GetHelpItems() []components.HelpItem {
	return []components.HelpItem{
		{Key: "←/→", Desc: "select option"},
		{Key: "enter", Desc: "confirm"},
		{Key: "?", Desc: "toggle help"},
		{Key: "q", Desc: "quit"},
	}
}
