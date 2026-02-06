package components

import (
	"skene-terminal-v2/internal/tui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HelpItem represents a single help entry
type HelpItem struct {
	Key  string
	Desc string
}

// HelpOverlay renders a help panel overlay
type HelpOverlay struct {
	Items   []HelpItem
	Title   string
	Visible bool
}

// NewHelpOverlay creates a new help overlay
func NewHelpOverlay() *HelpOverlay {
	return &HelpOverlay{
		Items: []HelpItem{
			{Key: "?", Desc: "toggle help"},
			{Key: "q", Desc: "quit"},
		},
		Title:   "Help",
		Visible: false,
	}
}

// Toggle visibility
func (h *HelpOverlay) Toggle() {
	h.Visible = !h.Visible
}

// SetItems updates help items
func (h *HelpOverlay) SetItems(items []HelpItem) {
	h.Items = items
}

// Render the help overlay
func (h *HelpOverlay) Render(width, height int) string {
	if !h.Visible {
		return ""
	}

	// Build help content
	var lines []string
	lines = append(lines, styles.SectionHeader.Render(h.Title))
	lines = append(lines, "")

	for _, item := range h.Items {
		key := styles.HelpKey.Render(item.Key)
		desc := styles.HelpDesc.Render(item.Desc)
		lines = append(lines, key+"  "+desc)
	}

	content := strings.Join(lines, "\n")

	// Style the box
	box := styles.Box.
		Width(40).
		Render(content)

	// Center in screen
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}

// FooterHelp renders inline footer help
func FooterHelp(items []HelpItem) string {
	var parts []string
	for _, item := range items {
		part := styles.HelpKey.Render(item.Key) + " " + styles.HelpDesc.Render(item.Desc)
		parts = append(parts, part)
	}
	return strings.Join(parts, styles.HelpSeparator.String())
}

// DefaultFooterHelp returns common footer help
func DefaultFooterHelp() string {
	return FooterHelp([]HelpItem{
		{Key: "?", Desc: "help"},
		{Key: "q", Desc: "quit"},
	})
}

// IntroHelp returns intro screen help
func IntroHelp() string {
	return FooterHelp([]HelpItem{
		{Key: "?", Desc: "help"},
		{Key: "q", Desc: "quit"},
	})
}

// ConfigHelp returns config screen help
func ConfigHelp() string {
	return FooterHelp([]HelpItem{
		{Key: "↑/↓", Desc: "navigate"},
		{Key: "enter", Desc: "select"},
		{Key: "?", Desc: "help"},
		{Key: "q", Desc: "quit"},
	})
}

// NavHelp returns navigation help
func NavHelp() string {
	return FooterHelp([]HelpItem{
		{Key: "←/→", Desc: "navigate"},
		{Key: "enter", Desc: "confirm"},
		{Key: "esc", Desc: "back"},
		{Key: "q", Desc: "quit"},
	})
}

// InputHelp returns input screen help
func InputHelp() string {
	return FooterHelp([]HelpItem{
		{Key: "enter", Desc: "submit"},
		{Key: "esc", Desc: "back"},
		{Key: "q", Desc: "quit"},
	})
}

// LoadingHelp returns loading screen help
func LoadingHelp() string {
	return FooterHelp([]HelpItem{
		{Key: "g", Desc: "play game"},
		{Key: "q", Desc: "quit"},
	})
}

// DashboardHelp returns dashboard help
func DashboardHelp() string {
	return FooterHelp([]HelpItem{
		{Key: "←/→", Desc: "tabs"},
		{Key: "↑/↓", Desc: "scroll"},
		{Key: "tab", Desc: "focus"},
		{Key: "q", Desc: "quit"},
	})
}
