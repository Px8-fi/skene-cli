package components

import (
	"skene/internal/tui/styles"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// TerminalOutput displays scrolling terminal/process output in a box
type TerminalOutput struct {
	lines      []string
	maxLines   int // max lines to keep in buffer
	width      int
	height     int // visible lines
	mu         sync.Mutex
}

// NewTerminalOutput creates a new terminal output display
func NewTerminalOutput(visibleLines, maxBuffer int) *TerminalOutput {
	if maxBuffer < visibleLines {
		maxBuffer = visibleLines * 3
	}
	return &TerminalOutput{
		lines:    make([]string, 0),
		maxLines: maxBuffer,
		height:   visibleLines,
	}
}

// SetSize updates the display dimensions
func (t *TerminalOutput) SetSize(width, height int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.width = width
	t.height = height
}

// AddLine appends a line of output
func (t *TerminalOutput) AddLine(line string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Split on newlines in case multiple lines come at once
	newLines := strings.Split(line, "\n")
	for _, l := range newLines {
		// Strip trailing carriage return
		l = strings.TrimRight(l, "\r")
		t.lines = append(t.lines, l)
	}

	// Trim buffer if too large
	if len(t.lines) > t.maxLines {
		t.lines = t.lines[len(t.lines)-t.maxLines:]
	}
}

// AddOutput appends raw output that may contain multiple lines
func (t *TerminalOutput) AddOutput(output string) {
	if output == "" {
		return
	}
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	for _, line := range lines {
		t.AddLine(line)
	}
}

// Clear resets the output
func (t *TerminalOutput) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lines = make([]string, 0)
}

// LineCount returns the number of lines
func (t *TerminalOutput) LineCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.lines)
}

// Render the terminal output box
func (t *TerminalOutput) Render(width int) string {
	t.mu.Lock()
	defer t.mu.Unlock()

	if width < 20 {
		width = 20
	}

	// Content width inside the box (border 2 + padding 4)
	contentWidth := width - 6
	if contentWidth < 10 {
		contentWidth = 10
	}

	// Get the last N visible lines (auto-scroll to bottom)
	visibleCount := t.height
	if visibleCount <= 0 {
		visibleCount = 8
	}

	startIdx := 0
	if len(t.lines) > visibleCount {
		startIdx = len(t.lines) - visibleCount
	}

	var displayLines []string

	defaultStyle := lipgloss.NewStyle().
		Foreground(styles.Cream).
		MaxWidth(contentWidth)
	errorStyle := lipgloss.NewStyle().
		Foreground(styles.Coral).
		MaxWidth(contentWidth)
	successStyle := lipgloss.NewStyle().
		Foreground(styles.Success).
		MaxWidth(contentWidth)
	warningStyle := lipgloss.NewStyle().
		Foreground(styles.Warning).
		MaxWidth(contentWidth)

	for i := startIdx; i < len(t.lines); i++ {
		line := t.lines[i]
		if len(line) > contentWidth {
			line = line[:contentWidth-1] + "~"
		}

		upper := strings.ToUpper(line)
		var styled string
		if strings.Contains(upper, "ERROR") || strings.Contains(upper, "FAILED") ||
			strings.Contains(upper, "TRACEBACK") || strings.Contains(upper, "EXCEPTION") {
			styled = errorStyle.Render(line)
		} else if strings.Contains(line, "✓") || strings.Contains(upper, "SUCCESS") ||
			strings.Contains(upper, "COMPLETE") || strings.Contains(upper, "DONE") {
			styled = successStyle.Render(line)
		} else if strings.Contains(upper, "WARNING") || strings.Contains(upper, "WARN") {
			styled = warningStyle.Render(line)
		} else {
			styled = defaultStyle.Render(line)
		}
		displayLines = append(displayLines, styled)
	}

	// Pad with empty lines if not enough output yet
	for len(displayLines) < visibleCount {
		displayLines = append(displayLines, "")
	}

	content := strings.Join(displayLines, "\n")

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.MidGray).
		Padding(0, 1).
		Width(width - 2) // account for border

	return boxStyle.Render(content)
}
