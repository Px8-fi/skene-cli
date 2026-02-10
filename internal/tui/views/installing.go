package views

import (
	"fmt"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// InstallTask represents a single installation task
type InstallTask struct {
	Name     string
	Progress float64
	Done     bool
	Error    string
	Active   bool
}

// InstallingView shows installation progress
type InstallingView struct {
	width       int
	height      int
	tasks       []InstallTask
	elapsedTime float64
	header      *components.WizardHeader
	spinner     *components.Spinner
	failed      bool
	failMessage string
}

// NewInstallingView creates a new installing progress view
func NewInstallingView(method string) *InstallingView {
	var tasks []InstallTask
	if method == "uvx" {
		tasks = []InstallTask{
			{Name: "Resolving ephemeral environment", Active: true},
			{Name: "Downloading skene-growth", Active: false},
			{Name: "Verifying installation", Active: false},
		}
	} else {
		tasks = []InstallTask{
			{Name: "Setting up environment", Active: true},
			{Name: "Installing skene-growth", Active: false},
			{Name: "Installing dependencies", Active: false},
			{Name: "Verifying installation", Active: false},
		}
	}

	return &InstallingView{
		tasks:   tasks,
		header:  components.NewWizardHeader(2, "Installing"),
		spinner: components.NewSpinner(),
	}
}

// SetSize updates dimensions
func (v *InstallingView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)
}

// SetElapsedTime sets the elapsed time
func (v *InstallingView) SetElapsedTime(t float64) {
	v.elapsedTime = t
}

// TickSpinner advances spinner animation
func (v *InstallingView) TickSpinner() {
	v.spinner.Tick()
}

// UpdateTask updates task progress by index
func (v *InstallingView) UpdateTask(index int, progress float64) {
	if index >= 0 && index < len(v.tasks) {
		v.tasks[index].Progress = progress
		v.tasks[index].Active = progress < 1.0
		if progress >= 1.0 {
			v.tasks[index].Done = true
			v.tasks[index].Active = false
			// Activate next task
			if index+1 < len(v.tasks) {
				v.tasks[index+1].Active = true
			}
		}
	}
}

// SetTaskError marks a task as failed
func (v *InstallingView) SetTaskError(index int, errMsg string) {
	if index >= 0 && index < len(v.tasks) {
		v.tasks[index].Error = errMsg
		v.tasks[index].Active = false
		v.failed = true
		v.failMessage = errMsg
	}
}

// AllTasksDone returns true if all tasks are complete
func (v *InstallingView) AllTasksDone() bool {
	for _, t := range v.tasks {
		if !t.Done {
			return false
		}
	}
	return true
}

// HasFailed returns true if installation failed
func (v *InstallingView) HasFailed() bool {
	return v.failed
}

// GetFailMessage returns the failure message
func (v *InstallingView) GetFailMessage() string {
	return v.failMessage
}

// Render the installing view
func (v *InstallingView) Render() string {
	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	// Wizard header
	wizHeader := v.header.Render()

	// Installation title
	installTitle := styles.Accent.Render(components.StaticLogo)

	// Progress section
	progressSection := v.renderProgress(sectionWidth)

	// Elapsed time
	elapsed := styles.Muted.Render(fmt.Sprintf("Elapsed: %.1fs", v.elapsedTime))

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.WizardProgressHelp())

	// Combine
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		wizHeader,
		"",
		installTitle,
		"",
		progressSection,
		"",
		elapsed,
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

func (v *InstallingView) renderProgress(width int) string {
	// Bar fills the container: subtract border (2), padding (4), indent (2), " 100%" (5)
	barWidth := width - 13
	if barWidth < 20 {
		barWidth = 20
	}
	var lines []string

	for _, task := range v.tasks {
		// Status icon
		var icon string
		var labelStyle lipgloss.Style

		if task.Error != "" {
			icon = styles.Error.Render("✗")
			labelStyle = styles.Error
		} else if task.Done {
			icon = styles.SuccessText.Render("✓")
			labelStyle = styles.Body
		} else if task.Active {
			icon = v.spinner.Render()
			labelStyle = styles.Body
		} else {
			icon = styles.Muted.Render("○")
			labelStyle = styles.Muted
		}

		// Task name
		taskLine := icon + " " + labelStyle.Render(task.Name)
		lines = append(lines, taskLine)

		// Progress bar for active/done tasks
		if task.Active || task.Done || task.Error != "" {
			filledWidth := int(float64(barWidth) * task.Progress)
			emptyWidth := barWidth - filledWidth

			filled := strings.Repeat("█", filledWidth)
			empty := strings.Repeat("░", emptyWidth)

			var barColor lipgloss.Color
			if task.Error != "" {
				barColor = styles.Coral
			} else {
				barColor = styles.Amber
			}

			bar := lipgloss.NewStyle().Foreground(barColor).Render(filled) +
				lipgloss.NewStyle().Foreground(styles.MidGray).Render(empty)

			percent := fmt.Sprintf("%3.0f%%", task.Progress*100)
			barLine := "  " + bar + " " + styles.Muted.Render(percent)
			lines = append(lines, barLine)
		}

		// Error message
		if task.Error != "" {
			lines = append(lines, "  "+styles.Error.Render(task.Error))
		}

		lines = append(lines, "")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return styles.Box.Width(width).Render(content)
}

// GetHelpItems returns context-specific help
func (v *InstallingView) GetHelpItems() []components.HelpItem {
	if v.failed {
		return []components.HelpItem{
			{Key: "r", Desc: "retry"},
			{Key: "ctrl+c", Desc: "quit"},
		}
	}
	return []components.HelpItem{
		{Key: "g", Desc: "play mini game"},
		{Key: "ctrl+c", Desc: "quit"},
	}
}
