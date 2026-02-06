package views

import (
	"fmt"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// GeneratingTask represents a generation task
type GeneratingTask struct {
	Name     string
	Progress float64
	Done     bool
}

// GeneratingView shows the generating/installing progress
type GeneratingView struct {
	width       int
	height      int
	tasks       []GeneratingTask
	elapsedTime float64
	showGame    bool
}

// NewGeneratingView creates a new generating view
func NewGeneratingView() *GeneratingView {
	return &GeneratingView{
		tasks: []GeneratingTask{
			{Name: "Generating config", Progress: 0},
			{Name: "Building prompt context", Progress: 0},
		},
		elapsedTime: 0,
	}
}

// SetSize updates dimensions
func (v *GeneratingView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// SetTasks sets the tasks to display
func (v *GeneratingView) SetTasks(tasks []GeneratingTask) {
	v.tasks = tasks
}

// UpdateTask updates a task's progress
func (v *GeneratingView) UpdateTask(index int, progress float64) {
	if index >= 0 && index < len(v.tasks) {
		v.tasks[index].Progress = progress
		if progress >= 1.0 {
			v.tasks[index].Done = true
		}
	}
}

// SetElapsedTime updates elapsed time
func (v *GeneratingView) SetElapsedTime(t float64) {
	v.elapsedTime = t
}

// AllTasksDone returns true if all tasks complete
func (v *GeneratingView) AllTasksDone() bool {
	for _, t := range v.tasks {
		if !t.Done {
			return false
		}
	}
	return true
}

// ToggleGame toggles game mode
func (v *GeneratingView) ToggleGame() {
	v.showGame = !v.showGame
}

// IsGameMode returns if game is active
func (v *GeneratingView) IsGameMode() bool {
	return v.showGame
}

// Render the generating view
func (v *GeneratingView) Render() string {
	// Page title
	title := styles.PageTitle("Configuration", v.width)

	// Big GENERATING text
	generatingText := v.renderGeneratingText()

	// Progress bars
	progressSection := v.renderProgress()

	// Elapsed time
	elapsed := styles.Muted.Render(fmt.Sprintf("Elapsed time: %.3fs", v.elapsedTime))

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.LoadingHelp())

	// Combine
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		generatingText,
		"",
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

func (v *GeneratingView) renderGeneratingText() string {
	// Simplified stylized text that fits better
	text := `
 ██████╗ ███████╗███╗   ██╗███████╗██████╗  █████╗ ████████╗██╗███╗   ██╗ ██████╗ 
██╔════╝ ██╔════╝████╗  ██║██╔════╝██╔══██╗██╔══██╗╚══██╔══╝██║████╗  ██║██╔════╝ 
██║  ███╗█████╗  ██╔██╗ ██║█████╗  ██████╔╝███████║   ██║   ██║██╔██╗ ██║██║  ███╗
██║   ██║██╔══╝  ██║╚██╗██║██╔══╝  ██╔══██╗██╔══██║   ██║   ██║██║╚██╗██║██║   ██║
╚██████╔╝███████╗██║ ╚████║███████╗██║  ██║██║  ██║   ██║   ██║██║ ╚████║╚██████╔╝
 ╚═════╝ ╚══════╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝╚═╝  ╚═══╝ ╚═════╝ 
`
	return styles.Accent.Render(text)
}

func (v *GeneratingView) renderProgress() string {
	var lines []string
	barWidth := 40

	for _, task := range v.tasks {
		// Task name centered above bar
		taskLabel := lipgloss.NewStyle().
			Width(barWidth + 10).
			Align(lipgloss.Center).
			Render(styles.Body.Render(task.Name))

		// Progress bar
		filledWidth := int(float64(barWidth) * task.Progress)
		emptyWidth := barWidth - filledWidth

		filled := strings.Repeat("█", filledWidth)
		empty := strings.Repeat("░", emptyWidth)

		bar := lipgloss.NewStyle().Foreground(styles.Amber).Render(filled) +
			lipgloss.NewStyle().Foreground(styles.MidGray).Render(empty)

		// Percentage
		percent := fmt.Sprintf("%3.0f%%", task.Progress*100)

		barLine := lipgloss.JoinHorizontal(
			lipgloss.Center,
			bar,
			" ",
			styles.Muted.Render(percent),
		)

		barLineCentered := lipgloss.NewStyle().
			Width(barWidth + 10).
			Align(lipgloss.Center).
			Render(barLine)

		lines = append(lines, taskLabel)
		lines = append(lines, barLineCentered)
		lines = append(lines, "")
	}

	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}

// GetHelpItems returns context-specific help
func (v *GeneratingView) GetHelpItems() []components.HelpItem {
	return []components.HelpItem{
		{Key: "g", Desc: "play mini game"},
		{Key: "q", Desc: "quit"},
	}
}
