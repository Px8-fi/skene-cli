package tui

import (
	"os"
	"time"

	"skene-terminal-v2/internal/game"
	"skene-terminal-v2/internal/services/analyzer"
	"skene-terminal-v2/internal/services/config"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/views"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
)

// AppState represents the current application state
type AppState int

const (
	StateIntro AppState = iota
	StateAnalyzing
	StateConfig
	StateProviderSelect
	StateModelSelect
	StateAuth
	StateAPIKey
	StateGenerating
	StateDashboard
	StateError
	StateGame
)

// TickMsg is sent on each animation frame
type TickMsg time.Time

// CountdownMsg is sent during auth countdown
type CountdownMsg int

// ProgressMsg is sent during generation
type ProgressMsg struct {
	TaskIndex int
	Progress  float64
}

// GenerationDoneMsg is sent when generation completes
type GenerationDoneMsg struct{}

// App is the main Bubble Tea application model
type App struct {
	// Core state
	state     AppState
	prevState AppState
	width     int
	height    int
	time      float64

	// Services
	configMgr *config.Manager
	analyzer  *analyzer.Analyzer
	analysis  *analyzer.AnalysisResult

	// Selected configuration
	selectedProvider *config.Provider
	selectedModel    *config.Model

	// Views
	introView      *views.IntroView
	configView     *views.ConfigView
	providerView   *views.ProviderView
	modelView      *views.ModelView
	authView       *views.AuthView
	apiKeyView     *views.APIKeyView
	generatingView *views.GeneratingView
	dashboardView  *views.DashboardView
	errorView      *views.ErrorView

	// Help overlay
	helpOverlay *components.HelpOverlay
	showHelp    bool

	// Game
	game *game.Game

	// Generation state
	generationStartTime time.Time
	generationTasks     []views.GeneratingTask

	// Auth state
	authCountdown int

	// Error state
	currentError *views.ErrorInfo
}

// NewApp creates a new application
func NewApp() *App {
	cwd, _ := os.Getwd()

	configMgr := config.NewManager(cwd)
	configMgr.LoadConfig()

	// Set default values if not present
	if configMgr.Config.Provider == "" {
		configMgr.Config.Provider = "gemini"
	}
	if configMgr.Config.Model == "" {
		configMgr.Config.Model = "gemini-3-flash-preview"
	}
	if configMgr.Config.OutputDir == "" {
		configMgr.Config.OutputDir = "./skene-context"
	}

	app := &App{
		state:        StateIntro,
		configMgr:    configMgr,
		introView:    views.NewIntroView(),
		configView:   views.NewConfigView(configMgr),
		providerView: views.NewProviderView(),
		helpOverlay:  components.NewHelpOverlay(),
	}

	return app
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		tick(),
		textinput.Blink,
	)
}

// Update handles messages and updates state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global key handlers
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}

		// Help toggle
		if msg.String() == "?" {
			a.showHelp = !a.showHelp
			return a, nil
		}

		// If help is shown, close it on any key
		if a.showHelp && msg.String() != "?" {
			a.showHelp = false
			return a, nil
		}

		// State-specific handling
		cmd := a.handleKeyPress(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.updateViewSizes()

	case TickMsg:
		a.time += 0.05

		// Update views that need time
		if a.state == StateIntro {
			a.introView.SetTime(a.time)
		}

		// Update generating view
		if a.state == StateGenerating {
			elapsed := time.Since(a.generationStartTime).Seconds()
			a.generatingView.SetElapsedTime(elapsed)

			// Simulate progress
			for i := range a.generationTasks {
				if !a.generationTasks[i].Done {
					a.generationTasks[i].Progress += 0.02
					if a.generationTasks[i].Progress >= 1.0 {
						a.generationTasks[i].Progress = 1.0
						a.generationTasks[i].Done = true
					}
					a.generatingView.UpdateTask(i, a.generationTasks[i].Progress)
				}
			}

			// Check if all done
			if a.generatingView.AllTasksDone() {
				cmds = append(cmds, func() tea.Msg {
					return GenerationDoneMsg{}
				})
			}
		}

		// Update game if active
		if a.state == StateGame && a.game != nil {
			a.game.Update()
		}

		cmds = append(cmds, tick())

	case CountdownMsg:
		a.authCountdown = int(msg)
		if a.authCountdown <= 0 {
			// Open browser and move to API key view
			browser.OpenURL(a.authView.GetAuthURL())
			a.state = StateAPIKey
			a.apiKeyView = views.NewAPIKeyView(a.selectedProvider, a.selectedModel)
			a.apiKeyView.SetSize(a.width, a.height)
		} else {
			a.authView.SetCountdown(a.authCountdown)
			cmds = append(cmds, countdown(a.authCountdown-1))
		}

	case GenerationDoneMsg:
		a.state = StateDashboard
		if a.dashboardView == nil {
			a.dashboardView = views.NewDashboardView()
			a.dashboardView.SetSize(a.width, a.height)
		}

	case game.GameTickMsg:
		if a.state == StateGame && a.game != nil {
			a.game.Update()
			cmds = append(cmds, game.GameTickCmd())
		}
	}

	return a, tea.Batch(cmds...)
}

// handleKeyPress handles state-specific key presses
func (a *App) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	switch a.state {
	case StateIntro:
		if key == "enter" {
			a.state = StateConfig
			return nil
		}
		if key == "q" {
			return tea.Quit
		}

	case StateConfig:
		return a.handleConfigKeys(msg)

	case StateProviderSelect:
		return a.handleProviderKeys(msg)

	case StateModelSelect:
		return a.handleModelKeys(msg)

	case StateAuth:
		return a.handleAuthKeys(msg)

	case StateAPIKey:
		return a.handleAPIKeyKeys(msg)

	case StateGenerating:
		if key == "g" {
			// Start game
			a.prevState = a.state
			a.state = StateGame
			if a.game == nil {
				a.game = game.NewGame(60, 20)
			} else {
				a.game.Restart()
			}
			a.game.SetSize(60, 20)
			return game.GameTickCmd()
		}
		if key == "q" {
			return tea.Quit
		}

	case StateDashboard:
		return a.handleDashboardKeys(msg)

	case StateError:
		return a.handleErrorKeys(msg)

	case StateGame:
		return a.handleGameKeys(msg)
	}

	return nil
}

func (a *App) handleConfigKeys(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	switch key {
	case "left", "h":
		a.configView.HandleLeft()
	case "right", "l":
		a.configView.HandleRight()
	case "enter":
		if a.configView.GetSelectedButton() == "Yes" {
			a.state = StateProviderSelect
		} else {
			// Skip to generating with current config
			a.startGenerating()
		}
	case "y":
		a.state = StateProviderSelect
	case "n":
		a.startGenerating()
	case "q":
		return tea.Quit
	}

	return nil
}

func (a *App) handleProviderKeys(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	switch key {
	case "up", "k":
		a.providerView.HandleUp()
	case "down", "j":
		a.providerView.HandleDown()
	case "left", "h":
		a.providerView.HandleLeft()
	case "right", "l":
		a.providerView.HandleRight()
	case "enter":
		if a.providerView.IsButtonFocused() {
			if a.providerView.GetButtonLabel() == "Back" {
				a.state = StateConfig
			} else {
				a.selectProvider()
			}
		} else {
			a.selectProvider()
		}
	case "esc":
		a.state = StateConfig
	case "q":
		return tea.Quit
	}

	return nil
}

func (a *App) selectProvider() tea.Cmd {
	provider := a.providerView.GetSelectedProvider()
	if provider == nil {
		return nil
	}

	a.selectedProvider = provider
	a.configMgr.SetProvider(provider.ID)

	// If Skene provider, go to auth flow
	if provider.ID == "skene" {
		a.authView = views.NewAuthView(provider)
		a.authView.SetSize(a.width, a.height)
		a.authCountdown = 3
		a.state = StateAuth
		return countdown(3)
	}

	// Otherwise go to model selection
	a.modelView = views.NewModelView(provider)
	a.modelView.SetSize(a.width, a.height)
	a.state = StateModelSelect

	return nil
}

func (a *App) handleModelKeys(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	switch key {
	case "up", "k":
		a.modelView.HandleUp()
	case "down", "j":
		a.modelView.HandleDown()
	case "left", "h":
		a.modelView.HandleLeft()
	case "right", "l":
		a.modelView.HandleRight()
	case "enter":
		if a.modelView.IsButtonFocused() {
			if a.modelView.GetButtonLabel() == "Back" {
				a.state = StateProviderSelect
			} else {
				a.selectModel()
			}
		} else {
			a.selectModel()
		}
	case "esc":
		a.state = StateProviderSelect
	case "q":
		return tea.Quit
	}

	return nil
}

func (a *App) selectModel() {
	model := a.modelView.GetSelectedModel()
	if model == nil {
		return
	}

	a.selectedModel = model
	a.configMgr.SetModel(model.ID)

	// Go to API key entry
	a.apiKeyView = views.NewAPIKeyView(a.selectedProvider, a.selectedModel)
	a.apiKeyView.SetSize(a.width, a.height)
	a.state = StateAPIKey
}

func (a *App) handleAuthKeys(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	switch key {
	case "m":
		// Manual entry
		a.authView.ShowFallback()
	case "enter":
		if a.authView.IsFallbackShown() {
			a.apiKeyView = views.NewAPIKeyView(a.selectedProvider, a.selectedModel)
			a.apiKeyView.SetSize(a.width, a.height)
			a.state = StateAPIKey
		}
	case "esc":
		a.state = StateProviderSelect
	case "q":
		return tea.Quit
	}

	return nil
}

func (a *App) handleAPIKeyKeys(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	if a.apiKeyView.IsInputFocused() {
		switch key {
		case "enter":
			if a.apiKeyView.Validate() {
				a.configMgr.SetAPIKey(a.apiKeyView.GetAPIKey())
				a.startGenerating()
			}
		case "tab":
			a.apiKeyView.HandleTab()
		case "esc":
			if a.selectedProvider != nil && a.selectedProvider.ID == "skene" {
				a.state = StateAuth
			} else {
				a.state = StateModelSelect
			}
		default:
			// Pass to text input
			a.apiKeyView.Update(msg)
		}
	} else {
		switch key {
		case "left", "h":
			a.apiKeyView.HandleLeft()
		case "right", "l":
			a.apiKeyView.HandleRight()
		case "enter":
			if a.apiKeyView.GetButtonLabel() == "Back" {
				if a.selectedProvider != nil && a.selectedProvider.ID == "skene" {
					a.state = StateAuth
				} else {
					a.state = StateModelSelect
				}
			} else {
				if a.apiKeyView.Validate() {
					a.configMgr.SetAPIKey(a.apiKeyView.GetAPIKey())
					a.startGenerating()
				}
			}
		case "tab":
			a.apiKeyView.HandleTab()
		case "esc":
			if a.selectedProvider != nil && a.selectedProvider.ID == "skene" {
				a.state = StateAuth
			} else {
				a.state = StateModelSelect
			}
		case "q":
			return tea.Quit
		}
	}

	return nil
}

func (a *App) handleDashboardKeys(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	switch key {
	case "left", "h":
		a.dashboardView.HandleLeft()
	case "right", "l":
		a.dashboardView.HandleRight()
	case "up", "k":
		a.dashboardView.HandleUp()
	case "down", "j":
		a.dashboardView.HandleDown()
	case "tab":
		a.dashboardView.HandleTab()
	case "q":
		return tea.Quit
	}

	return nil
}

func (a *App) handleErrorKeys(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	switch key {
	case "left", "h":
		a.errorView.HandleLeft()
	case "right":
		a.errorView.HandleRight()
	case "enter":
		btn := a.errorView.GetSelectedButton()
		switch btn {
		case "Retry":
			a.startGenerating()
		case "View Logs":
			a.errorView.ToggleLogs()
		case "Quit":
			return tea.Quit
		}
	case "l", "L":
		a.errorView.ToggleLogs()
	case "q":
		return tea.Quit
	}

	return nil
}

func (a *App) handleGameKeys(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	switch key {
	case "left", "a":
		a.game.MoveLeft()
	case "right", "d":
		a.game.MoveRight()
	case " ":
		a.game.Shoot()
	case "p":
		a.game.TogglePause()
	case "r":
		if a.game.IsGameOver() {
			a.game.Restart()
		}
	case "esc", "q":
		a.state = a.prevState
	}

	return nil
}

func (a *App) startGenerating() {
	a.generatingView = views.NewGeneratingView()
	a.generatingView.SetSize(a.width, a.height)
	a.generationStartTime = time.Now()
	a.generationTasks = []views.GeneratingTask{
		{Name: "Generating config", Progress: 0, Done: false},
		{Name: "Building prompt context", Progress: 0, Done: false},
	}
	a.generatingView.SetTasks(a.generationTasks)
	a.state = StateGenerating
}

func (a *App) updateViewSizes() {
	a.introView.SetSize(a.width, a.height)
	a.configView.SetSize(a.width, a.height)
	a.providerView.SetSize(a.width, a.height)

	if a.modelView != nil {
		a.modelView.SetSize(a.width, a.height)
	}
	if a.authView != nil {
		a.authView.SetSize(a.width, a.height)
	}
	if a.apiKeyView != nil {
		a.apiKeyView.SetSize(a.width, a.height)
	}
	if a.generatingView != nil {
		a.generatingView.SetSize(a.width, a.height)
	}
	if a.dashboardView != nil {
		a.dashboardView.SetSize(a.width, a.height)
	}
	if a.errorView != nil {
		a.errorView.SetSize(a.width, a.height)
	}
	if a.game != nil {
		a.game.SetSize(60, 20)
	}
}

// View renders the current view
func (a *App) View() string {
	var content string

	switch a.state {
	case StateIntro:
		content = a.introView.Render()
	case StateConfig:
		content = a.configView.Render()
	case StateProviderSelect:
		content = a.providerView.Render()
	case StateModelSelect:
		if a.modelView != nil {
			content = a.modelView.Render()
		}
	case StateAuth:
		if a.authView != nil {
			content = a.authView.Render()
		}
	case StateAPIKey:
		if a.apiKeyView != nil {
			content = a.apiKeyView.Render()
		}
	case StateGenerating:
		if a.generatingView != nil {
			content = a.generatingView.Render()
		}
	case StateDashboard:
		if a.dashboardView != nil {
			content = a.dashboardView.Render()
		}
	case StateError:
		if a.errorView != nil {
			content = a.errorView.Render()
		}
	case StateGame:
		if a.game != nil {
			content = lipgloss.Place(
				a.width,
				a.height,
				lipgloss.Center,
				lipgloss.Center,
				a.game.Render(),
			)
		}
	}

	// Overlay help if visible
	if a.showHelp {
		helpItems := a.getCurrentHelpItems()
		a.helpOverlay.SetItems(helpItems)
		overlay := a.helpOverlay.Render(a.width, a.height)
		if overlay != "" {
			content = overlay
		}
	}

	return content
}

func (a *App) getCurrentHelpItems() []components.HelpItem {
	switch a.state {
	case StateIntro:
		return a.introView.GetHelpItems()
	case StateConfig:
		return a.configView.GetHelpItems()
	case StateProviderSelect:
		return a.providerView.GetHelpItems()
	case StateModelSelect:
		if a.modelView != nil {
			return a.modelView.GetHelpItems()
		}
	case StateAuth:
		if a.authView != nil {
			return a.authView.GetHelpItems()
		}
	case StateAPIKey:
		if a.apiKeyView != nil {
			return a.apiKeyView.GetHelpItems()
		}
	case StateGenerating:
		if a.generatingView != nil {
			return a.generatingView.GetHelpItems()
		}
	case StateDashboard:
		if a.dashboardView != nil {
			return a.dashboardView.GetHelpItems()
		}
	case StateError:
		if a.errorView != nil {
			return a.errorView.GetHelpItems()
		}
	}

	return components.NewHelpOverlay().Items
}

// Helper functions
func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func countdown(seconds int) tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return CountdownMsg(seconds)
	})
}
