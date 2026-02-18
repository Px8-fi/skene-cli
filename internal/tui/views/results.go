package views

import (
	"skene/internal/tui/components"
	"skene/internal/tui/styles"
	"strings"

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

// NewResultsView creates a new results view with default placeholder content
func NewResultsView() *ResultsView {
	return NewResultsViewWithContent(
		getResultsGrowthPlanContent(),
		getResultsManifestContent(),
		getResultsProductDocsContent(),
	)
}

// NewResultsViewWithContent creates a new results view with custom content
func NewResultsViewWithContent(growthPlan, manifest, productDocs string) *ResultsView {
	vp := viewport.New(60, 20)

	v := &ResultsView{
		tabs:      []string{"Growth Plan", "Manifest", "Product Docs"},
		activeTab: 0,
		contents:  make(map[string]string),
		viewport:  vp,
		focus:     ResultsFocusTabs,
		header:    components.NewWizardHeader(7, "Analysis Results"),
	}

	// Set content (use provided or fall back to defaults)
	if growthPlan != "" {
		v.contents["Growth Plan"] = growthPlan
	} else {
		v.contents["Growth Plan"] = getResultsGrowthPlanContent()
	}
	if manifest != "" {
		v.contents["Manifest"] = manifest
	} else {
		v.contents["Manifest"] = getResultsManifestContent()
	}
	if productDocs != "" {
		v.contents["Product Docs"] = productDocs
	} else {
		v.contents["Product Docs"] = getResultsProductDocsContent()
	}

	v.viewport.SetContent(v.contents["Growth Plan"])

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

	// Wizard header
	wizHeader := lipgloss.NewStyle().Width(sectionWidth).Render(v.header.Render())

	// Success banner
	banner := styles.SuccessText.Render("Skene Analysis Complete")

	// Tabs
	tabsView := v.renderTabs()

	// Content
	contentBox := v.renderContentBox()

	// Action hint
	actionHint := styles.Accent.Render("Press 'n' for next steps")

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.WizardResultsHelp())

	// Combine
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

func getResultsGrowthPlanContent() string {
	return strings.TrimSpace(`
## CONFIDENTIAL ENGINEERING MEMO

Date: 2026-02-09
Subject: Growth Strategy Analysis
From: Council of Growth Engineers

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

EXECUTIVE SUMMARY

Your codebase has been analyzed for Product-Led Growth
opportunities. Below are the key findings and actionable
recommendations.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1. STRIP TO THE GROWTH CORE

The fundamental growth problem is maximizing the number
of codebases Skene analyzes. Every codebase analyzed is a
potential source of:

  • New Users — owner and collaborators
  • New Codebases — related project recommendations
  • Improved AI — more data to refine analysis

2. SELECTED GROWTH LOOPS

Loop 1: Public Scan Leaderboard
  Priority: HIGH | Impact: 4.2x user acquisition
  - Allow users to make scans public
  - Create competitive leaderboards
  - Gamify code quality improvements

Loop 2: Collaborative Analysis
  Priority: HIGH | Impact: 3.1x retention
  - Team-based scanning features
  - Shared insights and recommendations
  - Cross-project pattern detection

Loop 3: AI-Powered Recommendations
  Priority: MEDIUM | Impact: 2.5x engagement
  - Suggest related codebases to scan
  - Identify similar projects
  - Build a network effect

3. IMPLEMENTATION ROADMAP

Week 1-2: Ship public scan feature
Week 3-4: Build leaderboard infrastructure
Week 5-6: Implement sharing and collaboration
Week 7-8: Launch recommendation engine

4. SUCCESS METRICS

  • Scans per user per week: target 3+
  • Viral coefficient (K-factor): target 1.2
  • Time to first collaboration: <48 hours
  • Cross-project engagement: 35%

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Execute. No meetings.
`)
}

func getResultsManifestContent() string {
	return strings.TrimSpace(`
SKENE GROWTH MANIFEST v2.0
Generated: 2026-02-09

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

TECH STACK DETECTION

  Framework:  detected
  Language:   detected
  Database:   detected
  Auth:       detected
  Deployment: detected

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CURRENT GROWTH FEATURES

  ✓ User authentication flow
  ✓ API key management
  ✓ Configuration system
  ! No social sharing detected
  ! No referral system detected
  ! No usage analytics detected

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

REVENUE LEAKAGE ISSUES

  ⚠ No monetization layer detected
  ⚠ Missing conversion funnels
  ⚠ No upsell triggers found

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

GROWTH OPPORTUNITIES

  1. [HIGH] Social sharing for analysis results
  2. [HIGH] Team collaboration features
  3. [MEDIUM] Public analysis leaderboard
  4. [MEDIUM] Integration marketplace
  5. [LOW] Webhook notifications

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

GENERATED FILES

  ./skene-context/growth-manifest.json
  ./skene-context/growth-template.json
  ./skene-context/product-docs.md
`)
}

func getResultsProductDocsContent() string {
	return strings.TrimSpace(`
PRODUCT DOCUMENTATION
Auto-generated by Skene CLI

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PRODUCT OVERVIEW

  Tagline: [Auto-detected from codebase]
  Target Audience: Developers and teams
  Value Proposition: Growth analysis for code

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

DETECTED FEATURES

  1. Code Analysis Engine
     Scans codebases for growth patterns
     and opportunities.

  2. AI-Powered Insights
     Uses LLM providers to generate growth
     strategies and recommendations.

  3. Configuration Management
     Flexible config system supporting
     multiple providers and models.

  4. Multi-Provider Support
     Works with OpenAI, Anthropic, Gemini,
     local models, and custom endpoints.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

GETTING STARTED

  1. Run 'uvx skene-growth analyze .'
  2. Review growth-manifest.json
  3. Generate growth plan with 'uvx skene-growth plan'
  4. Build implementation with 'uvx skene-growth build'

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

For more information:
  https://github.com/SkeneTechnologies/skene-cli
`)
}
