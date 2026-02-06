package views

import (
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// DashboardFocus represents which element is focused
type DashboardFocus int

const (
	FocusTabs DashboardFocus = iota
	FocusContent
)

// DashboardView shows the final result with tabs
type DashboardView struct {
	width         int
	height        int
	tabs          []string
	activeTab     int
	contents      map[string]string
	viewport      viewport.Model
	focus         DashboardFocus
}

// NewDashboardView creates a new dashboard view
func NewDashboardView() *DashboardView {
	vp := viewport.New(60, 20)

	v := &DashboardView{
		tabs:      []string{"Growth plan", "Manifest", "Contribute"},
		activeTab: 0,
		contents:  make(map[string]string),
		viewport:  vp,
		focus:     FocusTabs,
	}

	// Set default content
	v.contents["Growth plan"] = getGrowthPlanContent()
	v.contents["Manifest"] = getManifestContent()
	v.contents["Contribute"] = getContributeContent()

	v.viewport.SetContent(v.contents["Growth plan"])

	return v
}

// SetSize updates dimensions
func (v *DashboardView) SetSize(width, height int) {
	v.width = width
	v.height = height

	// Update viewport size
	vpWidth := width - 10
	if vpWidth < 40 {
		vpWidth = 40
	}
	if vpWidth > 100 {
		vpWidth = 100
	}

	vpHeight := height - 12
	if vpHeight < 10 {
		vpHeight = 10
	}

	v.viewport.Width = vpWidth
	v.viewport.Height = vpHeight
}

// HandleLeft moves tab left
func (v *DashboardView) HandleLeft() {
	if v.focus == FocusTabs && v.activeTab > 0 {
		v.activeTab--
		v.updateContent()
	}
}

// HandleRight moves tab right
func (v *DashboardView) HandleRight() {
	if v.focus == FocusTabs && v.activeTab < len(v.tabs)-1 {
		v.activeTab++
		v.updateContent()
	}
}

// HandleUp scrolls content up
func (v *DashboardView) HandleUp() {
	if v.focus == FocusContent {
		v.viewport.LineUp(3)
	}
}

// HandleDown scrolls content down
func (v *DashboardView) HandleDown() {
	if v.focus == FocusContent {
		v.viewport.LineDown(3)
	}
}

// HandleTab cycles focus
func (v *DashboardView) HandleTab() {
	if v.focus == FocusTabs {
		v.focus = FocusContent
	} else {
		v.focus = FocusTabs
	}
}

func (v *DashboardView) updateContent() {
	tabName := v.tabs[v.activeTab]
	if content, ok := v.contents[tabName]; ok {
		v.viewport.SetContent(content)
		v.viewport.GotoTop()
	}
}

// GetViewport returns the viewport for updates
func (v *DashboardView) GetViewport() *viewport.Model {
	return &v.viewport
}

// Render the dashboard view
func (v *DashboardView) Render() string {
	// Render tabs
	tabsView := v.renderTabs()

	// Content area with border
	contentBox := v.renderContentBox()

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.DashboardHelp())

	// Combine
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		tabsView,
		contentBox,
	)

	// Add padding
	padded := lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)

	// Fill remaining height
	mainContent := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Top,
		padded,
	)

	return mainContent + "\n" + footer
}

func (v *DashboardView) renderTabs() string {
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

func (v *DashboardView) renderContentBox() string {
	// Border style based on focus
	borderColor := styles.MidGray
	if v.focus == FocusContent {
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
func (v *DashboardView) GetHelpItems() []components.HelpItem {
	if v.focus == FocusTabs {
		return []components.HelpItem{
			{Key: "←/→", Desc: "switch tabs"},
			{Key: "tab", Desc: "focus content"},
			{Key: "q", Desc: "quit"},
		}
	}
	return []components.HelpItem{
		{Key: "↑/↓", Desc: "scroll"},
		{Key: "tab", Desc: "focus tabs"},
		{Key: "q", Desc: "quit"},
	}
}

// Content generators
func getGrowthPlanContent() string {
	return strings.TrimSpace(`
## CONFIDENTIAL ENGINEERING MEMO

**Date:** 2026-01-30

**Subject:** Project skene-growth: Dominance through AI-Powered Virality

**To:** CEO

**From:** Council of Growth Engineers

**Executive Summary:**

Skene-growth, as it stands, is a glorified feature. Its potential
is being squandered on linear improvements like "better onboarding."
The real game is leveraging the AI core to engineer a viral loop
where every scan *automatically* generates more scans, more users,
and a rapidly expanding moat. We will achieve this by turning code
analysis into a collaborative, competitive, and inherently viral
activity powered by AI.

### 1. Strip to the Growth Core

The fundamental growth problem is *not* user onboarding. It's
maximizing the number of codebases Skene analyzes. Every codebase
analyzed is a potential source of:

*   **New Users:** The owner and collaborators of the code.
*   **New Codebases:** Recommendations for related projects.
*   **Improved AI:** More data to refine the analysis and
    recommendations.

### 2. The Viral Mechanics

**Public Scans & Leaderboards**
- Allow users to make scans public
- Create competitive leaderboards
- Gamify code quality improvements

**Collaborative Analysis**
- Team-based scanning features
- Shared insights and recommendations
- Cross-project pattern detection

**AI-Powered Recommendations**
- Suggest related codebases to scan
- Identify similar projects
- Build a network effect

### 3. Implementation Priorities

1. Ship the public scan feature immediately
2. Build the leaderboard infrastructure
3. Implement sharing and collaboration
4. Launch the recommendation engine
5. Monitor and optimize viral metrics

### 4. Success Metrics

- Scans per user per week
- Viral coefficient (K-factor)
- Time to first collaboration
- Cross-project engagement rate

Execute. No meetings.
`)
}

func getManifestContent() string {
	return strings.TrimSpace(`
## Skene-Growth Manifest

**Version:** 1.0.0
**Generated:** 2026-01-30

### Configuration

| Key | Value |
|-----|-------|
| Provider | gemini |
| Model | gemini-3-flash-preview |
| Output | ./skene-context |
| Verbose | true |

### Generated Files

- .skene.config
- skene-context/manifest.json
- skene-context/growth-plan.md
- skene-context/analysis.json

### Installation Details

**Package:** skene-growth
**Method:** pip install
**Dependencies:** 12 packages

### Next Steps

1. Run 'skene analyze' to generate insights
2. Review growth-plan.md for recommendations
3. Configure your preferred LLM provider
4. Start tracking growth metrics

### Support

- Documentation: https://docs.skene.ai
- Issues: https://github.com/skene-ai/skene-growth/issues
- Discord: https://discord.gg/skene
`)
}

func getContributeContent() string {
	return strings.TrimSpace(`
## Contributing to Skene-Growth

We welcome contributions from the community!

### Ways to Contribute

**Code Contributions**
- Bug fixes and improvements
- New features and integrations
- Performance optimizations
- Documentation updates

**Non-Code Contributions**
- Bug reports and feature requests
- Documentation improvements
- Community support
- Translations

### Getting Started

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write or update tests
5. Submit a pull request

### Development Setup

` + "```bash" + `
git clone https://github.com/skene-ai/skene-growth
cd skene-growth
python -m venv .venv
source .venv/bin/activate
pip install -e ".[dev]"
pytest
` + "```" + `

### Code Style

- Follow PEP 8 guidelines
- Use type hints where possible
- Write docstrings for public APIs
- Keep functions focused and small

### Community

- Be respectful and inclusive
- Help newcomers get started
- Share knowledge and ideas
- Celebrate contributions

### Contact

- GitHub: @skene-ai
- Twitter: @skene_ai
- Email: contribute@skene.ai

Thank you for helping make Skene-Growth better!
`)
}
