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

// AnalyzingView shows analysis progress with live terminal output
type AnalyzingView struct {
	width       int
	height      int
	phases      []AnalysisPhase
	elapsedTime float64
	header      *components.WizardHeader
	spinner     *components.Spinner
	terminal    *components.TerminalOutput
	failed      bool
	done        bool
	failMessage string
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
		phases:   phases,
		header:   components.NewWizardHeader(6, "Running Analysis"),
		spinner:  components.NewSpinner(),
		terminal: components.NewTerminalOutput(14, 300),
	}
}

// SetSize updates dimensions
func (v *AnalyzingView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)
	// Adjust terminal visible lines based on available height
	termHeight := height - 18
	if termHeight < 6 {
		termHeight = 6
	}
	if termHeight > 22 {
		termHeight = 22
	}
	v.terminal.SetSize(width, termHeight)
}

// SetElapsedTime updates elapsed time
func (v *AnalyzingView) SetElapsedTime(t float64) {
	v.elapsedTime = t
}

// TickSpinner advances spinner animation
func (v *AnalyzingView) TickSpinner() {
	v.spinner.Tick()
}

// UpdatePhase updates a phase's progress and logs the message to terminal
func (v *AnalyzingView) UpdatePhase(index int, progress float64, message string) {
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
	// Log the message to terminal output
	if message != "" {
		v.terminal.AddLine(message)
	}
}

// SetDone marks the command as successfully completed
func (v *AnalyzingView) SetDone() {
	v.done = true
	v.terminal.AddLine("✓ Done")
}

// SetCommandFailed marks the view as failed with the error visible in terminal
func (v *AnalyzingView) SetCommandFailed(errMsg string) {
	v.failed = true
	v.failMessage = errMsg
	if errMsg != "" {
		v.terminal.AddLine("")
		v.terminal.AddLine("ERROR: " + errMsg)
	}
}

// IsDone returns true if the command completed (success or failure)
func (v *AnalyzingView) IsDone() bool {
	return v.done || v.failed
}

// SetPhaseError marks a phase as failed
func (v *AnalyzingView) SetPhaseError(index int, errMsg string) {
	if index >= 0 && index < len(v.phases) {
		v.phases[index].Error = errMsg
		v.phases[index].Active = false
		v.failed = true
	}
	if errMsg != "" {
		v.terminal.AddLine("ERROR: " + errMsg)
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

	// Current phase status
	var statusLine string
	if v.failed {
		statusLine = styles.Error.Render("✗ Failed")
		if v.failMessage != "" {
			statusLine += "\n" + styles.Muted.Render("  "+v.failMessage)
		}
	} else if v.done {
		statusLine = styles.SuccessText.Render("✓ Complete")
	} else if v.AllPhasesDone() {
		statusLine = styles.SuccessText.Render("✓ Analysis complete")
	} else {
		currentPhase := ""
		for _, p := range v.phases {
			if p.Active {
				currentPhase = p.Name
				break
			}
		}
		if currentPhase != "" {
			statusLine = v.spinner.Render() + " " + styles.Body.Render(currentPhase)
		} else {
			statusLine = v.spinner.Render() + " " + styles.Body.Render("Running...")
		}
	}

	// Overall progress bar
	overallBar := v.renderOverallProgress(sectionWidth)

	// Terminal output
	termOutput := v.terminal.Render(sectionWidth)

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
		statusLine,
		"",
		overallBar,
		"",
		termOutput,
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

// GetHelpItems returns context-specific help
func (v *AnalyzingView) GetHelpItems() []components.HelpItem {
	if v.done || v.failed {
		return []components.HelpItem{
			{Key: "esc", Desc: "go back"},
			{Key: "ctrl+c", Desc: "quit"},
		}
	}
	return []components.HelpItem{
		{Key: "g", Desc: "play mini game"},
		{Key: "ctrl+c", Desc: "quit"},
	}
}
