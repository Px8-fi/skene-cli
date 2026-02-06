package views

import (
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// ErrorSeverity represents error severity level
type ErrorSeverity int

const (
	SeverityWarning ErrorSeverity = iota
	SeverityError
	SeverityCritical
)

// ErrorInfo contains error details
type ErrorInfo struct {
	Code       string
	Title      string
	Message    string
	Suggestion string
	Severity   ErrorSeverity
	Retryable  bool
}

// ErrorView displays errors with suggested fixes
type ErrorView struct {
	width       int
	height      int
	error       *ErrorInfo
	buttonGroup *components.ButtonGroup
	logs        []string
	showLogs    bool
}

// NewErrorView creates a new error view
func NewErrorView(err *ErrorInfo) *ErrorView {
	var buttons *components.ButtonGroup
	if err.Retryable {
		buttons = components.NewButtonGroup("Retry", "View Logs", "Quit")
	} else {
		buttons = components.NewButtonGroup("View Logs", "Quit")
	}

	return &ErrorView{
		error:       err,
		buttonGroup: buttons,
		logs:        make([]string, 0),
	}
}

// SetSize updates dimensions
func (v *ErrorView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// SetError updates the error to display
func (v *ErrorView) SetError(err *ErrorInfo) {
	v.error = err
}

// SetLogs sets the logs to display
func (v *ErrorView) SetLogs(logs []string) {
	v.logs = logs
}

// ToggleLogs toggles log view
func (v *ErrorView) ToggleLogs() {
	v.showLogs = !v.showLogs
}

// HandleLeft moves button focus left
func (v *ErrorView) HandleLeft() {
	v.buttonGroup.Previous()
}

// HandleRight moves button focus right
func (v *ErrorView) HandleRight() {
	v.buttonGroup.Next()
}

// GetSelectedButton returns selected button
func (v *ErrorView) GetSelectedButton() string {
	return v.buttonGroup.GetActiveLabel()
}

// Render the error view
func (v *ErrorView) Render() string {
	if v.showLogs {
		return v.renderLogs()
	}

	sectionWidth := 70

	// Error icon based on severity
	var icon string
	var titleStyle lipgloss.Style
	switch v.error.Severity {
	case SeverityWarning:
		icon = "⚠️"
		titleStyle = lipgloss.NewStyle().Foreground(styles.Warning)
	case SeverityError:
		icon = "❌"
		titleStyle = styles.Error
	case SeverityCritical:
		icon = "🚨"
		titleStyle = styles.Error.Bold(true)
	}

	// Title
	title := titleStyle.Render(icon + " " + v.error.Title)

	// Error code
	code := styles.Muted.Render("[" + v.error.Code + "]")

	// Message
	message := styles.Body.Render(v.error.Message)

	// Suggestion box
	suggestionHeader := styles.SectionHeader.Render("Suggested Fix")
	suggestion := styles.SuccessText.Render("→ " + v.error.Suggestion)

	suggestionBox := styles.Box.
		Width(sectionWidth - 8).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			suggestionHeader,
			"",
			suggestion,
		))

	// Buttons
	buttons := v.buttonGroup.Render()

	// Combine
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		code,
		"",
		message,
		"",
		suggestionBox,
		"",
		buttons,
	)

	box := styles.Box.
		Width(sectionWidth).
		Render(content)

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.FooterHelp([]components.HelpItem{
			{Key: "←/→", Desc: "select"},
			{Key: "enter", Desc: "confirm"},
			{Key: "l", Desc: "toggle logs"},
			{Key: "q", Desc: "quit"},
		}))

	// Center
	centered := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)

	return centered + "\n" + footer
}

func (v *ErrorView) renderLogs() string {
	sectionWidth := v.width - 20
	if sectionWidth > 100 {
		sectionWidth = 100
	}

	title := styles.SectionHeader.Render("Installation Logs")

	// Build log content
	var logLines []string
	maxLines := v.height - 15
	if maxLines < 10 {
		maxLines = 10
	}

	startIdx := 0
	if len(v.logs) > maxLines {
		startIdx = len(v.logs) - maxLines
	}

	for i := startIdx; i < len(v.logs); i++ {
		logLines = append(logLines, styles.Muted.Render(v.logs[i]))
	}

	if len(logLines) == 0 {
		logLines = append(logLines, styles.Muted.Render("No logs available"))
	}

	logContent := lipgloss.JoinVertical(lipgloss.Left, logLines...)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		logContent,
	)

	box := styles.Box.
		Width(sectionWidth).
		Height(v.height - 10).
		Render(content)

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.FooterHelp([]components.HelpItem{
			{Key: "l", Desc: "close logs"},
			{Key: "q", Desc: "quit"},
		}))

	centered := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)

	return centered + "\n" + footer
}

// GetHelpItems returns context-specific help
func (v *ErrorView) GetHelpItems() []components.HelpItem {
	return []components.HelpItem{
		{Key: "←/→", Desc: "select option"},
		{Key: "enter", Desc: "confirm"},
		{Key: "l", Desc: "view/hide logs"},
		{Key: "q", Desc: "quit"},
	}
}

// Common errors
var (
	ErrPythonNotFound = &ErrorInfo{
		Code:       "PYTHON_NOT_FOUND",
		Title:      "Python Not Found",
		Message:    "Python is required but was not found in your PATH.",
		Suggestion: "Install Python 3.8+ from python.org or your package manager.",
		Severity:   SeverityError,
		Retryable:  false,
	}

	ErrPipFailed = &ErrorInfo{
		Code:       "PIP_FAILED",
		Title:      "Package Installation Failed",
		Message:    "pip failed to install skene-growth package.",
		Suggestion: "Run 'pip install --upgrade pip' and try again.",
		Severity:   SeverityError,
		Retryable:  true,
	}

	ErrNetworkFailed = &ErrorInfo{
		Code:       "NETWORK_ERROR",
		Title:      "Network Connection Failed",
		Message:    "Could not connect to package registry.",
		Suggestion: "Check your internet connection and try again.",
		Severity:   SeverityWarning,
		Retryable:  true,
	}

	ErrPermissionDenied = &ErrorInfo{
		Code:       "PERMISSION_DENIED",
		Title:      "Permission Denied",
		Message:    "Insufficient permissions to write to target directory.",
		Suggestion: "Check directory permissions or run with elevated privileges.",
		Severity:   SeverityError,
		Retryable:  true,
	}

	ErrInvalidAPIKey = &ErrorInfo{
		Code:       "INVALID_API_KEY",
		Title:      "Invalid API Key",
		Message:    "The provided API key was rejected by the provider.",
		Suggestion: "Double-check your API key and ensure it has the required permissions.",
		Severity:   SeverityError,
		Retryable:  true,
	}
)
