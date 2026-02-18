package views

import (
	"skene/internal/tui/components"
	"skene/internal/tui/styles"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// ResultsFocus represents which element is focused
type ResultsFocus int

const (
	ResultsFocusTabs ResultsFocus = iota
	ResultsFocusContent
)

// ResultsView shows the analysis results in a tabbed dashboard
type ResultsView struct {
	width     int
	height    int
	tabs      []string
	activeTab int
	contents  map[string]string
	viewport  viewport.Model
	focus     ResultsFocus
	header    *components.WizardHeader
}

// NewResultsView creates an empty results view (no tabs until content is provided)
func NewResultsView() *ResultsView {
	return NewResultsViewWithContent("", "", "")
}

// NewResultsViewWithContent creates a results view showing only the tabs
// for which skene-growth produced actual content.
func NewResultsViewWithContent(growthPlan, manifest, productDocs string) *ResultsView {
	vp := viewport.New(60, 20)

	v := &ResultsView{
		activeTab: 0,
		contents:  make(map[string]string),
		viewport:  vp,
		focus:     ResultsFocusTabs,
		header:    components.NewWizardHeader(7, "Analysis Results"),
	}

	if growthPlan != "" {
		v.tabs = append(v.tabs, "Growth Plan")
		v.contents["Growth Plan"] = growthPlan
	}
	if manifest != "" {
		v.tabs = append(v.tabs, "Manifest")
		v.contents["Manifest"] = manifest
	}
	if productDocs != "" {
		v.tabs = append(v.tabs, "Product Docs")
		v.contents["Product Docs"] = productDocs
	}

	if len(v.tabs) == 0 {
		v.tabs = append(v.tabs, "Results")
		v.contents["Results"] = "No analysis output found in skene-context/.\nRun the analysis first."
	}

	v.viewport.SetContent(v.contents[v.tabs[0]])

	return v
}

// SetSize updates dimensions
func (v *ResultsView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)

	vpWidth := width - 10
	if vpWidth < 40 {
		vpWidth = 40
	}
	if vpWidth > 100 {
		vpWidth = 100
	}

	vpHeight := height - 16
	if vpHeight < 10 {
		vpHeight = 10
	}

	v.viewport.Width = vpWidth
	v.viewport.Height = vpHeight
}

// HandleLeft moves tab left
func (v *ResultsView) HandleLeft() {
	if v.focus == ResultsFocusTabs && v.activeTab > 0 {
		v.activeTab--
		v.updateContent()
	}
}

// HandleRight moves tab right
func (v *ResultsView) HandleRight() {
	if v.focus == ResultsFocusTabs && v.activeTab < len(v.tabs)-1 {
		v.activeTab++
		v.updateContent()
	}
}

// HandleUp scrolls content up
func (v *ResultsView) HandleUp() {
	if v.focus == ResultsFocusContent {
		v.viewport.LineUp(3)
	}
}

// HandleDown scrolls content down
func (v *ResultsView) HandleDown() {
	if v.focus == ResultsFocusContent {
		v.viewport.LineDown(3)
	}
}

// HandleTab cycles focus
func (v *ResultsView) HandleTab() {
	if v.focus == ResultsFocusTabs {
		v.focus = ResultsFocusContent
	} else {
		v.focus = ResultsFocusTabs
	}
}

func (v *ResultsView) updateContent() {
	tabName := v.tabs[v.activeTab]
	if content, ok := v.contents[tabName]; ok {
		v.viewport.SetContent(content)
		v.viewport.GotoTop()
	}
}

// Render the results view
func (v *ResultsView) Render() string {
	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	wizHeader := lipgloss.NewStyle().Width(sectionWidth).Render(v.header.Render())

	banner := styles.SuccessText.Render("Skene Analysis Complete")

	tabsView := v.renderTabs()

	contentBox := v.renderContentBox()

	actionHint := styles.Accent.Render("Press 'n' for next steps")

	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.WizardResultsHelp())

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		wizHeader,
		"",
		banner,
		"",
		tabsView,
		contentBox,
		"",
		actionHint,
	)

	mainContent := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.NewStyle().Padding(1, 2).Render(content),
	)

	return mainContent + "\n" + footer
}

func (v *ResultsView) renderTabs() string {
	var tabs []string
	for i, tab := range v.tabs {
		var style lipgloss.Style
		if i == v.activeTab {
			style = styles.TabActive
		} else {
			style = styles.TabInactive
		}
		tabs = append(tabs, style.Render(tab))
	}
	return lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...)
}

func (v *ResultsView) renderContentBox() string {
	borderColor := styles.MidGray
	if v.focus == ResultsFocusContent {
		borderColor = styles.Cream
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(v.viewport.Width + 6)

	return boxStyle.Render(v.viewport.View())
}

// GetHelpItems returns context-specific help
func (v *ResultsView) GetHelpItems() []components.HelpItem {
	if v.focus == ResultsFocusTabs {
		return []components.HelpItem{
			{Key: "←/→", Desc: "switch tabs"},
			{Key: "tab", Desc: "focus content"},
			{Key: "n", Desc: "next steps"},
			{Key: "ctrl+c", Desc: "quit"},
		}
	}
	return []components.HelpItem{
		{Key: "↑/↓", Desc: "scroll"},
		{Key: "tab", Desc: "focus tabs"},
		{Key: "n", Desc: "next steps"},
		{Key: "ctrl+c", Desc: "quit"},
	}
}
