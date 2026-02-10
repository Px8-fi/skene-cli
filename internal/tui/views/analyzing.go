package views

import (
	"fmt"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// AnalysisPhase represents a phase of the analysis
type AnalysisPhase struct {
	Name     string
	Progress float64
	Done     bool
	Active   bool
	Error    string
}

// AnalyzingView shows analysis progress with detailed phases
type AnalyzingView struct {
	width       int
	height      int
	phases      []AnalysisPhase
	elapsedTime float64
	header      *components.WizardHeader
	spinner     *components.Spinner
	failed      bool
	currentIdx  int
}

// NewAnalyzingView creates a new analysis progress view
func NewAnalyzingView() *AnalyzingView {
	phases := []AnalysisPhase{
		{Name: "Scanning codebase", Active: true},
		{Name: "Detecting product features", Active: false},
		{Name: "Growth loop analysis", Active: false},
		{Name: "Monetisation analysis", Active: false},
		{Name: "Opportunity modelling", Active: false},
		{Name: "Generating manifests & docs", Active: false},
	}

	return &AnalyzingView{
		phases:  phases,
		header:  components.NewWizardHeader(6, "Running Analysis"),
		spinner: components.NewSpinner(),
	}
}

// SetSize updates dimensions
func (v *AnalyzingView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)
}

// SetElapsedTime updates elapsed time
func (v *AnalyzingView) SetElapsedTime(t float64) {
	v.elapsedTime = t
}

// TickSpinner advances spinner animation
func (v *AnalyzingView) TickSpinner() {
	v.spinner.Tick()
}

// UpdatePhase updates a phase's progress
func (v *AnalyzingView) UpdatePhase(index int, progress float64) {
	if index >= 0 && index < len(v.phases) {
		v.phases[index].Progress = progress
		v.phases[index].Active = progress < 1.0
		if progress >= 1.0 {
			v.phases[index].Done = true
			v.phases[index].Active = false
			// Activate next phase
			if index+1 < len(v.phases) {
				v.phases[index+1].Active = true
				v.currentIdx = index + 1
			}
		}
	}
}

// SetPhaseError marks a phase as failed
func (v *AnalyzingView) SetPhaseError(index int, errMsg string) {
	if index >= 0 && index < len(v.phases) {
		v.phases[index].Error = errMsg
		v.phases[index].Active = false
		v.failed = true
	}
}

// AllPhasesDone returns true if all phases are complete
func (v *AnalyzingView) AllPhasesDone() bool {
	for _, p := range v.phases {
		if !p.Done {
			return false
		}
	}
	return true
}

// HasFailed returns true if analysis failed
func (v *AnalyzingView) HasFailed() bool {
	return v.failed
}

// GetOverallProgress returns overall progress 0.0-1.0
func (v *AnalyzingView) GetOverallProgress() float64 {
	done := 0
	for _, p := range v.phases {
		if p.Done {
			done++
		}
	}
	return float64(done) / float64(len(v.phases))
}

// Render the analyzing view
func (v *AnalyzingView) Render() string {
	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	// Wizard header
	wizHeader := v.header.Render()

	// Big title
	title := styles.Accent.Render(`
 █████╗ ███╗   ██╗ █████╗ ██╗  ██╗   ██╗███████╗██╗███╗   ██╗ ██████╗
██╔══██╗████╗  ██║██╔══██╗██║  ╚██╗ ██╔╝╚══███╔╝██║████╗  ██║██╔════╝
███████║██╔██╗ ██║███████║██║   ╚████╔╝   ███╔╝ ██║██╔██╗ ██║██║  ███╗
██╔══██║██║╚██╗██║██╔══██║██║    ╚██╔╝   ███╔╝  ██║██║╚██╗██║██║   ██║
██║  ██║██║ ╚████║██║  ██║███████╗██║   ███████╗██║██║ ╚████║╚██████╔╝
╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝╚═╝   ╚══════╝╚═╝╚═╝  ╚═══╝ ╚═════╝`)

	// Overall progress bar
	overallBar := v.renderOverallProgress(sectionWidth)

	// Phase details
	phaseSection := v.renderPhases(sectionWidth)

	// Elapsed
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
		title,
		"",
		overallBar,
		"",
		phaseSection,
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

func (v *AnalyzingView) renderOverallProgress(width int) string {
	// Bar fills the section: subtract " 100%" (5) and some margin (2)
	barWidth := width - 7
	if barWidth < 20 {
		barWidth = 20
	}
	progress := v.GetOverallProgress()
	filledWidth := int(float64(barWidth) * progress)
	emptyWidth := barWidth - filledWidth

	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", emptyWidth)

	bar := lipgloss.NewStyle().Foreground(styles.Amber).Render(filled) +
		lipgloss.NewStyle().Foreground(styles.MidGray).Render(empty)

	percent := fmt.Sprintf("%3.0f%%", progress*100)

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(bar + " " + styles.Body.Render(percent))
}

func (v *AnalyzingView) renderPhases(width int) string {
	var lines []string

	for _, phase := range v.phases {
		var icon string
		var labelStyle lipgloss.Style

		if phase.Error != "" {
			icon = styles.Error.Render("✗")
			labelStyle = styles.Error
		} else if phase.Done {
			icon = styles.SuccessText.Render("✓")
			labelStyle = styles.Body
		} else if phase.Active {
			icon = v.spinner.Render()
			labelStyle = styles.Body
		} else {
			icon = styles.Muted.Render("○")
			labelStyle = styles.Muted
		}

		line := icon + " " + labelStyle.Render(phase.Name)

		if phase.Active {
			line += "  " + styles.Muted.Render(fmt.Sprintf("%0.0f%%", phase.Progress*100))
		}

		if phase.Error != "" {
			line += "\n  " + styles.Error.Render(phase.Error)
		}

		lines = append(lines, line)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return styles.Box.Width(width).Render(content)
}

// GetHelpItems returns context-specific help
func (v *AnalyzingView) GetHelpItems() []components.HelpItem {
	return []components.HelpItem{
		{Key: "g", Desc: "play mini game"},
		{Key: "ctrl+c", Desc: "quit"},
	}
}
