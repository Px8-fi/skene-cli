package components

import (
	"fmt"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// WizardStep represents a step in the wizard
type WizardStep struct {
	Number int
	Name   string
}

// WizardSteps defines the complete wizard flow
var WizardSteps = []WizardStep{
	{Number: 1, Name: "System Check"},
	{Number: 2, Name: "Install Method"},
	{Number: 3, Name: "AI Provider"},
	{Number: 4, Name: "Authentication"},
	{Number: 5, Name: "Project Setup"},
	{Number: 6, Name: "Analysis"},
	{Number: 7, Name: "Results"},
}

// WizardHeader renders the wizard progress header
type WizardHeader struct {
	CurrentStep int
	TotalSteps  int
	StepName    string
	Width       int
}

// NewWizardHeader creates a new wizard header
func NewWizardHeader(currentStep int, stepName string) *WizardHeader {
	return &WizardHeader{
		CurrentStep: currentStep,
		TotalSteps:  len(WizardSteps),
		StepName:    stepName,
		Width:       80,
	}
}

// SetWidth sets the header width
func (h *WizardHeader) SetWidth(width int) {
	h.Width = width
}

// SetStep updates the current step
func (h *WizardHeader) SetStep(step int, name string) {
	h.CurrentStep = step
	h.StepName = name
}

// Render the wizard header
func (h *WizardHeader) Render() string {
	headerWidth := h.Width - 4
	if headerWidth < 60 {
		headerWidth = 60
	}
	if headerWidth > 100 {
		headerWidth = 100
	}

	// Step counter
	stepCounter := styles.Muted.Render(fmt.Sprintf("Step %d of %d", h.CurrentStep, h.TotalSteps))

	// Step name
	stepName := styles.Title.Render(h.StepName)

	// Progress dots
	dots := renderWizardDots(h.CurrentStep, h.TotalSteps)

	// Left side: step info
	leftSide := lipgloss.JoinVertical(
		lipgloss.Left,
		stepCounter,
		stepName,
	)

	// Layout
	topBar := lipgloss.JoinHorizontal(
		lipgloss.Center,
		leftSide,
	)

	// Separator line
	sep := styles.Divider(headerWidth)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		topBar,
		dots,
		sep,
	)

	return lipgloss.NewStyle().
		Width(headerWidth).
		Padding(0, 2).
		Render(content)
}

// RenderCompact renders a compact single-line header
func (h *WizardHeader) RenderCompact() string {
	dots := renderWizardDots(h.CurrentStep, h.TotalSteps)
	stepInfo := styles.Muted.Render(fmt.Sprintf("Step %d/%d", h.CurrentStep, h.TotalSteps))
	stepName := styles.Accent.Render(h.StepName)

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		dots,
		"  ",
		stepInfo,
		"  ",
		stepName,
	)
}

func renderWizardDots(current, total int) string {
	var dots string
	for i := 1; i <= total; i++ {
		if i < current {
			dots += styles.SuccessText.Render("●")
		} else if i == current {
			dots += styles.Accent.Render("●")
		} else {
			dots += styles.Muted.Render("○")
		}
		if i < total {
			if i < current {
				dots += styles.SuccessText.Render("─")
			} else {
				dots += styles.Muted.Render("─")
			}
		}
	}
	return dots
}
