package components

import (
	"fmt"
	"skene-terminal-v2/internal/tui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ProgressBar renders a custom progress bar matching the design
type ProgressBar struct {
	Label    string
	Progress float64 // 0.0 to 1.0
	Width    int
}

// NewProgressBar creates a new progress bar
func NewProgressBar(label string, width int) *ProgressBar {
	return &ProgressBar{
		Label:    label,
		Progress: 0.0,
		Width:    width,
	}
}

// SetProgress updates progress (0.0 to 1.0)
func (p *ProgressBar) SetProgress(progress float64) {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	p.Progress = progress
}

// Render the progress bar
func (p *ProgressBar) Render() string {
	// Calculate filled width
	barWidth := p.Width - 10 // Leave space for percentage
	filledWidth := int(float64(barWidth) * p.Progress)
	emptyWidth := barWidth - filledWidth

	// Build the bar
	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", emptyWidth)

	bar := lipgloss.NewStyle().Foreground(styles.Amber).Render(filled) +
		lipgloss.NewStyle().Foreground(styles.MidGray).Render(empty)

	// Percentage
	percent := fmt.Sprintf("%3.0f%%", p.Progress*100)

	// Combine label and bar
	labelStyle := styles.Body
	percentStyle := styles.Muted

	return fmt.Sprintf("%s\n%s %s",
		labelStyle.Render(p.Label),
		bar,
		percentStyle.Render(percent),
	)
}

// RenderCompact renders just the bar without label
func (p *ProgressBar) RenderCompact() string {
	barWidth := p.Width
	filledWidth := int(float64(barWidth) * p.Progress)
	emptyWidth := barWidth - filledWidth

	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", emptyWidth)

	return lipgloss.NewStyle().Foreground(styles.Amber).Render(filled) +
		lipgloss.NewStyle().Foreground(styles.MidGray).Render(empty)
}

// TaskProgress represents a task with progress
type TaskProgress struct {
	Name     string
	Progress float64
	Done     bool
	Error    error
}

// ProgressGroup renders multiple progress bars
type ProgressGroup struct {
	Tasks       []TaskProgress
	BarWidth    int
	ElapsedTime float64
}

// NewProgressGroup creates a new progress group
func NewProgressGroup(barWidth int) *ProgressGroup {
	return &ProgressGroup{
		Tasks:    make([]TaskProgress, 0),
		BarWidth: barWidth,
	}
}

// AddTask adds a task to track
func (pg *ProgressGroup) AddTask(name string) {
	pg.Tasks = append(pg.Tasks, TaskProgress{
		Name:     name,
		Progress: 0.0,
		Done:     false,
	})
}

// UpdateTask updates a task's progress
func (pg *ProgressGroup) UpdateTask(index int, progress float64) {
	if index >= 0 && index < len(pg.Tasks) {
		pg.Tasks[index].Progress = progress
		if progress >= 1.0 {
			pg.Tasks[index].Done = true
		}
	}
}

// SetError marks a task as failed
func (pg *ProgressGroup) SetError(index int, err error) {
	if index >= 0 && index < len(pg.Tasks) {
		pg.Tasks[index].Error = err
	}
}

// AllDone checks if all tasks are complete
func (pg *ProgressGroup) AllDone() bool {
	for _, task := range pg.Tasks {
		if !task.Done {
			return false
		}
	}
	return true
}

// Render the progress group
func (pg *ProgressGroup) Render() string {
	var lines []string

	for _, task := range pg.Tasks {
		// Status indicator
		var statusIcon string
		var labelStyle lipgloss.Style

		if task.Error != nil {
			statusIcon = styles.Error.Render("✗")
			labelStyle = styles.Error
		} else if task.Done {
			statusIcon = styles.SuccessText.Render("✓")
			labelStyle = styles.Body
		} else {
			statusIcon = styles.Muted.Render("○")
			labelStyle = styles.Body
		}

		// Task label
		label := fmt.Sprintf("%s %s", statusIcon, labelStyle.Render(task.Name))
		lines = append(lines, label)

		// Progress bar
		bar := NewProgressBar("", pg.BarWidth)
		bar.SetProgress(task.Progress)
		lines = append(lines, "  "+bar.RenderCompact()+" "+styles.Muted.Render(fmt.Sprintf("%3.0f%%", task.Progress*100)))
		lines = append(lines, "")
	}

	// Elapsed time
	elapsed := fmt.Sprintf("Elapsed time: %.3fs", pg.ElapsedTime)
	lines = append(lines, styles.Muted.Render(elapsed))

	return strings.Join(lines, "\n")
}

// Spinner component
type Spinner struct {
	frames []string
	index  int
}

// NewSpinner creates a new spinner
func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		index:  0,
	}
}

// Tick advances the spinner
func (s *Spinner) Tick() {
	s.index = (s.index + 1) % len(s.frames)
}

// Render the spinner
func (s *Spinner) Render() string {
	return styles.Accent.Render(s.frames[s.index])
}

// SpinnerWithText renders spinner with text
func (s *Spinner) SpinnerWithText(text string) string {
	return s.Render() + " " + styles.Body.Render(text)
}
