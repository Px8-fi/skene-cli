package views

import (
	"fmt"
	"skene-terminal-v2/internal/services/syscheck"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// SysCheckView displays system prerequisite check results
type SysCheckView struct {
	width          int
	height         int
	results        *syscheck.SystemCheckResult
	spinner        *components.Spinner
	checking       bool
	showInstall    bool // Show uv install option
	ideRequestSent bool // Track if IDE request was sent
	ideRequestPath string // Path to the request file
	header         *components.WizardHeader
	buttonGroup    *components.ButtonGroup
}

// NewSysCheckView creates a new system check view
func NewSysCheckView() *SysCheckView {
	return &SysCheckView{
		spinner:     components.NewSpinner(),
		checking:    true,
		showInstall: false,
		header:      components.NewWizardHeader(1, "System Check"),
		buttonGroup: components.NewButtonGroup("Continue", "Install uv", "Ask IDE", "Quit"),
	}
}

// SetSize updates dimensions
func (v *SysCheckView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)
}

// SetResults sets the check results
func (v *SysCheckView) SetResults(results *syscheck.SystemCheckResult) {
	v.results = results
	v.checking = false
	v.showInstall = !results.AllPassed && results.UV.Status == syscheck.StatusFailed
}

// SetIDERequestSent marks that IDE request was sent
func (v *SysCheckView) SetIDERequestSent(filePath string) {
	v.ideRequestSent = true
	v.ideRequestPath = filePath
}

// SetChecking sets the checking state
func (v *SysCheckView) SetChecking(checking bool) {
	v.checking = checking
}

// TickSpinner advances the spinner animation
func (v *SysCheckView) TickSpinner() {
	v.spinner.Tick()
}

// NeedsUVInstall returns true if uv needs to be installed
func (v *SysCheckView) NeedsUVInstall() bool {
	return v.results != nil && v.results.UV.Status == syscheck.StatusFailed
}

// CanProceed returns true if the system is ready to continue
func (v *SysCheckView) CanProceed() bool {
	return v.results != nil && v.results.CanProceed
}

// PythonFailed returns true if Python check failed
func (v *SysCheckView) PythonFailed() bool {
	return v.results != nil && v.results.Python.Status == syscheck.StatusFailed
}

// GetButtonLabel returns current button label
func (v *SysCheckView) GetButtonLabel() string {
	return v.buttonGroup.GetActiveLabel()
}

// HandleLeft handles left key
func (v *SysCheckView) HandleLeft() {
	v.buttonGroup.Previous()
}

// HandleRight handles right key
func (v *SysCheckView) HandleRight() {
	v.buttonGroup.Next()
}

// Render the system check view
func (v *SysCheckView) Render() string {
	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	// Wizard header
	wizHeader := v.header.Render()

	// Check results section
	checksSection := v.renderChecks(sectionWidth)

	// Status message
	statusMsg := v.renderStatus()

	// Buttons (only show when not checking)
	var buttonsSection string
	if !v.checking {
		// Update button states based on results
		if v.results != nil && v.results.CanProceed {
			buttonsSection = v.buttonGroup.Render()
		} else if v.results != nil && v.results.Python.Status == syscheck.StatusFailed {
			// Show Ask IDE and Quit buttons when Python is missing
			ideBtn := components.NewButtonGroup("Ask IDE", "Quit")
			buttonsSection = ideBtn.Render()
		} else if v.results != nil && !v.results.AllPassed {
			// Show Ask IDE button when there are failures
			ideBtn := components.NewButtonGroup("Ask IDE", "Quit")
			buttonsSection = ideBtn.Render()
		} else {
			buttonsSection = v.buttonGroup.Render()
		}
	}

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.WizardHelp())

	// Combine
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		wizHeader,
		"",
		checksSection,
		"",
		statusMsg,
		"",
		buttonsSection,
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

func (v *SysCheckView) renderChecks(width int) string {
	header := styles.SectionHeader.Render("Checking prerequisites...")

	var items []string

	if v.checking {
		items = append(items, v.spinner.SpinnerWithText("Checking Python installation..."))
		items = append(items, styles.Muted.Render("  Waiting..."))
		items = append(items, "")
		items = append(items, v.spinner.SpinnerWithText("Checking uv runtime..."))
		items = append(items, styles.Muted.Render("  Waiting..."))
	} else if v.results != nil {
		// Python check
		items = append(items, renderCheckLine(v.results.Python))
		items = append(items, "")

		// UV check
		items = append(items, renderCheckLine(v.results.UV))
		items = append(items, "")

		// Pip check (fallback info)
		if v.results.UV.Status == syscheck.StatusFailed && v.results.Pip.Status == syscheck.StatusPassed {
			items = append(items, renderCheckLine(v.results.Pip))
			items = append(items, "")
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		lipgloss.JoinVertical(lipgloss.Left, items...),
	)

	return styles.Box.Width(width).Render(content)
}

func renderCheckLine(result syscheck.CheckResult) string {
	var icon string
	var nameStyle lipgloss.Style

	switch result.Status {
	case syscheck.StatusPassed:
		icon = styles.SuccessText.Render("✓")
		nameStyle = styles.Body
	case syscheck.StatusFailed:
		icon = styles.Error.Render("✗")
		nameStyle = styles.Error
	case syscheck.StatusWarning:
		icon = lipgloss.NewStyle().Foreground(styles.Warning).Render("!")
		nameStyle = styles.Body
	default:
		icon = styles.Muted.Render("○")
		nameStyle = styles.Muted
	}

	line := icon + " " + nameStyle.Render(result.Name)
	if result.Message != "" {
		line += "  " + styles.Muted.Render(result.Message)
	}

	if result.Status == syscheck.StatusFailed && result.FixCommand != "" {
		line += "\n  " + styles.Accent.Render("Fix: ") + styles.Muted.Render(result.FixCommand)
	}

	return line
}

func (v *SysCheckView) renderStatus() string {
	if v.checking {
		return styles.Muted.Render("Running system checks...")
	}

	if v.results == nil {
		return ""
	}

	// Show IDE request success message if sent
	if v.ideRequestSent {
		msg := styles.SuccessText.Render("✓ Request sent to IDE! Check Cursor chat or ask: 'help me fix the system check issues'")
		if v.ideRequestPath != "" {
			msg += "\n" + styles.Muted.Render(fmt.Sprintf("Details saved to: %s", v.ideRequestPath))
		}
		return msg
	}

	if v.results.AllPassed {
		return styles.SuccessText.Render("All checks passed! Your system is ready.")
	}

	if v.results.CanProceed {
		return lipgloss.NewStyle().Foreground(styles.Warning).Render(
			"uv is not installed but pip is available. You can continue or install uv for a better experience.")
	}

	if v.results.Python.Status == syscheck.StatusFailed {
		return styles.Error.Render("Python 3.11+ is required. Please install it and try again.")
	}

	return styles.Error.Render("Some requirements are not met. Please fix the issues above.")
}

// GetHelpItems returns context-specific help
func (v *SysCheckView) GetHelpItems() []components.HelpItem {
	if v.checking {
		return []components.HelpItem{
			{Key: "ctrl+c", Desc: "quit"},
		}
	}
	return []components.HelpItem{
		{Key: "enter", Desc: "continue"},
		{Key: "←/→", Desc: "select option"},
		{Key: "ctrl+c", Desc: "quit"},
	}
}
