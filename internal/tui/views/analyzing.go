package views

import (
	"skene/internal/constants"
	"skene/internal/tui/components"
	"skene/internal/tui/styles"

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
	header      *components.WizardHeader
	spinner     *components.Spinner
	terminal    *components.TerminalOutput
	failed      bool
	done        bool
	failMessage string
	currentIdx int
}

// NewAnalyzingView creates a new analysis progress view
func NewAnalyzingView() *AnalyzingView {
	phases := []AnalysisPhase{
		{Name: constants.PhaseScanCodebase, Active: true},
		{Name: constants.PhaseDetectFeatures, Active: false},
		{Name: constants.PhaseGrowthLoops, Active: false},
		{Name: constants.PhaseMonetisation, Active: false},
		{Name: constants.PhaseOpportunities, Active: false},
		{Name: constants.PhaseGenerateDocs, Active: false},
	}

	return &AnalyzingView{
		phases:   phases,
		header:   components.NewTitleHeader(constants.StepNameAnalyzing),
		spinner:  components.NewSpinner(),
		terminal: components.NewTerminalOutput(14, 300),
	}
}

// NewCommandView creates a view for running a generic command with terminal output
func NewCommandView(title string) *AnalyzingView {
	return &AnalyzingView{
		phases:   []AnalysisPhase{},
		header:   components.NewTitleHeader(title),
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
	v.terminal.AddLine("✓ " + constants.AnalyzingDone)
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
	wizHeader := lipgloss.NewStyle().Width(sectionWidth).Render(v.header.Render())

	// Current phase status
	var statusLine string
	if v.failed {
		statusLine = styles.Error.Render("✗ " + constants.AnalyzingFailed)
		if v.failMessage != "" {
			statusLine += "\n" + styles.Muted.Render("  "+v.failMessage)
		}
	} else if v.done {
		statusLine = styles.SuccessText.Render("✓ " + constants.AnalyzingComplete)
	} else if len(v.phases) > 0 && v.AllPhasesDone() {
		statusLine = styles.SuccessText.Render("✓ " + constants.AnalyzingComplete)
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
			statusLine = v.spinner.Render() + " " + styles.Body.Render(constants.AnalyzingRunning)
		}
	}

	// Terminal output
	termOutput := v.terminal.Render(sectionWidth)

	// Footer
	var footerContent string
	if v.done || v.failed {
		footerContent = components.FooterHelp([]components.HelpItem{
			{Key: constants.HelpKeyEsc, Desc: constants.HelpDescGoBack},
			{Key: constants.HelpKeyCtrlC, Desc: constants.HelpDescQuit},
		})
	} else {
		footerContent = components.FooterHelp([]components.HelpItem{
			{Key: constants.HelpKeyEsc, Desc: constants.HelpDescCancel},
			{Key: constants.HelpKeyG, Desc: constants.HelpDescPlayMiniGame},
			{Key: constants.HelpKeyCtrlC, Desc: constants.HelpDescQuit},
		})
	}
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(footerContent)

	// Combine
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		wizHeader,
		"",
		statusLine,
		"",
		termOutput,
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

// GetHelpItems returns context-specific help
func (v *AnalyzingView) GetHelpItems() []components.HelpItem {
	if v.done || v.failed {
		return []components.HelpItem{
			{Key: constants.HelpKeyEsc, Desc: constants.HelpDescGoBack},
			{Key: constants.HelpKeyCtrlC, Desc: constants.HelpDescQuit},
		}
	}
	return []components.HelpItem{
		{Key: constants.HelpKeyEsc, Desc: constants.HelpDescCancel},
		{Key: constants.HelpKeyG, Desc: constants.HelpDescPlayMiniGame},
		{Key: constants.HelpKeyCtrlC, Desc: constants.HelpDescQuit},
	}
}
