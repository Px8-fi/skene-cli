package views

import (
	"fmt"
	"skene/internal/constants"
	"skene/internal/services/syscheck"
	"skene/internal/tui/components"
	"skene/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// SysCheckView displays system prerequisite check results
type SysCheckView struct {
	width          int
	height         int
	results        *syscheck.SystemCheckResult
	spinner        *components.Spinner
	checking       bool
	ideRequestSent bool
	ideRequestPath string
	header         *components.WizardHeader
	buttonGroup    *components.ButtonGroup
}

// NewSysCheckView creates a new system check view
func NewSysCheckView() *SysCheckView {
	return &SysCheckView{
		spinner:  components.NewSpinner(),
		checking: true,
		header:   components.NewTitleHeader(constants.StepNameSysCheck),
		buttonGroup: components.NewButtonGroup(constants.ButtonContinue, constants.ButtonQuit),
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

	wizHeader := lipgloss.NewStyle().Width(sectionWidth).Render(v.header.Render())
	checksSection := v.renderChecks(sectionWidth)
	statusMsg := v.renderStatus()

	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.FooterHelp([]components.HelpItem{
			{Key: constants.HelpKeyEnter, Desc: constants.HelpDescContinue},
			{Key: constants.HelpKeyEsc, Desc: constants.HelpDescBack},
			{Key: constants.HelpKeyCtrlC, Desc: constants.HelpDescQuit},
		}))

	block := lipgloss.JoinVertical(
		lipgloss.Left,
		wizHeader,
		"",
		checksSection,
		"",
		statusMsg,
	)

	content := lipgloss.NewStyle().PaddingTop(2).Render(block)

	centered := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Top,
		content,
	)

	return centered + "\n" + footer
}

func (v *SysCheckView) renderChecks(width int) string {
	header := styles.SectionHeader.Render(constants.SysCheckHeader)

	var items []string

	if v.checking {
		items = append(items, v.spinner.SpinnerWithText(constants.SysCheckSettingUp))
		items = append(items, styles.Muted.Render("  "+constants.SysCheckFirstRun))
	} else if v.results != nil {
		items = append(items, renderCheckLine(v.results.UV))
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
		return styles.Muted.Render(constants.SysCheckRunning)
	}

	if v.results == nil {
		return ""
	}

	if v.ideRequestSent {
		msg := styles.SuccessText.Render("✓ " + constants.SysCheckIDESent)
		if v.ideRequestPath != "" {
			msg += "\n" + styles.Muted.Render(fmt.Sprintf("Details saved to: %s", v.ideRequestPath))
		}
		return msg
	}

	if v.results.AllPassed {
		return styles.SuccessText.Render(constants.SysCheckAllPassed)
	}

	return styles.Error.Render(constants.SysCheckFailed)
}

// GetHelpItems returns context-specific help
func (v *SysCheckView) GetHelpItems() []components.HelpItem {
	if v.checking {
		return []components.HelpItem{
			{Key: constants.HelpKeyCtrlC, Desc: constants.HelpDescQuit},
		}
	}
	return []components.HelpItem{
		{Key: constants.HelpKeyEnter, Desc: constants.HelpDescContinue},
		{Key: constants.HelpKeyLeftRight, Desc: constants.HelpDescSelectOption},
		{Key: constants.HelpKeyCtrlC, Desc: constants.HelpDescQuit},
	}
}
