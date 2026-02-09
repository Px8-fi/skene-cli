package tui

import (
	"fmt"
	"time"

	"skene-terminal-v2/internal/game"
	"skene-terminal-v2/internal/services/config"
	"skene-terminal-v2/internal/services/syscheck"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/views"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
)

// ═══════════════════════════════════════════════════════════════════
// WIZARD STATE MACHINE
// ═══════════════════════════════════════════════════════════════════

// AppState represents the current wizard step
type AppState int

const (
	StateWelcome        AppState = iota // Welcome screen
	StateSysCheck                       // System prerequisite checks
	StateInstallMethod                  // Select uvx vs pip
	StateInstalling                     // Installation progress
	StateProviderSelect                 // AI provider selection
	StateModelSelect                    // Model selection for chosen provider
	StateAuth                           // Skene magic link authentication
	StateAPIKey                         // Manual API key entry
	StateLocalModel                     // Local model detection (Ollama/LM Studio)
	StateProjectDir                     // Project directory selection
	StateAnalysisConfig                 // Analysis configuration
	StateAnalyzing                      // Analysis progress
	StateResults                        // Results dashboard
	StateNextSteps                      // Next steps after analysis
	StateError                          // Error display
	StateGame                           // Mini game during wait
)

// ═══════════════════════════════════════════════════════════════════
// MESSAGES
// ═══════════════════════════════════════════════════════════════════

// TickMsg is sent on each animation frame
type TickMsg time.Time

// CountdownMsg is sent during auth countdown
type CountdownMsg int

// SysCheckDoneMsg is sent when system checks complete
type SysCheckDoneMsg struct {
	Results *syscheck.SystemCheckResult
}

// InstallDoneMsg is sent when installation completes
type InstallDoneMsg struct {
	Error error
}

// AnalysisDoneMsg is sent when analysis completes
type AnalysisDoneMsg struct {
	Error error
}

// UVInstallDoneMsg is sent when uv installation completes
type UVInstallDoneMsg struct {
	Error error
}

// LocalModelDetectMsg is sent with local model detection results
type LocalModelDetectMsg struct {
	Models []string
	Error  error
}

// ═══════════════════════════════════════════════════════════════════
// APP MODEL
// ═══════════════════════════════════════════════════════════════════

// App is the main Bubble Tea application model implementing the wizard
type App struct {
	// Core state
	state     AppState
	prevState AppState
	width     int
	height    int
	time      float64

	// Services
	configMgr    *config.Manager
	sysChecker   *syscheck.Checker
	sysCheckDone bool

	// Selected configuration
	selectedProvider *config.Provider
	selectedModel    *config.Model
	installMethod    string // "uvx" or "pip"

	// Views
	welcomeView        *views.WelcomeView
	sysCheckView       *views.SysCheckView
	installMethodView  *views.InstallMethodView
	installingView     *views.InstallingView
	providerView       *views.ProviderView
	modelView          *views.ModelView
	authView           *views.AuthView
	apiKeyView         *views.APIKeyView
	localModelView     *views.LocalModelView
	projectDirView     *views.ProjectDirView
	analysisConfigView *views.AnalysisConfigView
	analyzingView      *views.AnalyzingView
	resultsView        *views.ResultsView
	nextStepsView      *views.NextStepsView
	errorView          *views.ErrorView

	// Help overlay
	helpOverlay *components.HelpOverlay
	showHelp    bool

	// Game
	game *game.Game

	// Timing
	installStartTime  time.Time
	analysisStartTime time.Time

	// Auth state
	authCountdown int

	// Error state
	currentError *views.ErrorInfo
}

// ═══════════════════════════════════════════════════════════════════
// INITIALIZATION
// ═══════════════════════════════════════════════════════════════════

// NewApp creates a new wizard application
func NewApp() *App {
	configMgr := config.NewManager(".")
	configMgr.LoadConfig()

	// Set default values if not present
	if configMgr.Config.OutputDir == "" {
		configMgr.Config.OutputDir = "./skene-context"
	}

	app := &App{
		state:        StateWelcome,
		configMgr:    configMgr,
		sysChecker:   syscheck.NewChecker(),
		welcomeView:  views.NewWelcomeView(),
		providerView: views.NewProviderView(),
		helpOverlay:  components.NewHelpOverlay(),
	}

	return app
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, tick())
	cmds = append(cmds, textinput.Blink)
	// Initialize welcome animation
	if a.welcomeView != nil {
		animCmd := a.welcomeView.InitAnimation()
		if animCmd != nil {
			cmds = append(cmds, animCmd)
		}
	}
	return tea.Batch(cmds...)
}

// ═══════════════════════════════════════════════════════════════════
// UPDATE
// ═══════════════════════════════════════════════════════════════════

// Update handles messages and updates state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global: ctrl+c always quits
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}

		// Help toggle
		if msg.String() == "?" && a.state != StateAPIKey && a.state != StateProjectDir {
			a.showHelp = !a.showHelp
			return a, nil
		}

		// Close help on any key
		if a.showHelp && msg.String() != "?" {
			a.showHelp = false
			return a, nil
		}

		// State-specific key handling
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

		// Update welcome animation
		if a.state == StateWelcome {
			a.welcomeView.SetTime(a.time)
		}

		// Tick spinners for active views
		if a.state == StateSysCheck && a.sysCheckView != nil {
			a.sysCheckView.TickSpinner()
		}
		if a.state == StateInstalling && a.installingView != nil {
			a.installingView.TickSpinner()
			elapsed := time.Since(a.installStartTime).Seconds()
			a.installingView.SetElapsedTime(elapsed)
			a.simulateInstallProgress()
		}
		if a.state == StateAnalyzing && a.analyzingView != nil {
			a.analyzingView.TickSpinner()
			elapsed := time.Since(a.analysisStartTime).Seconds()
			a.analyzingView.SetElapsedTime(elapsed)
			a.simulateAnalysisProgress()
		}
		if a.state == StateAuth && a.authView != nil {
			a.authView.TickSpinner()
		}
		if a.state == StateAPIKey && a.apiKeyView != nil {
			a.apiKeyView.TickSpinner()
		}
		if a.state == StateLocalModel && a.localModelView != nil {
			a.localModelView.TickSpinner()
		}

		// Update game if active
		if a.state == StateGame && a.game != nil {
			a.game.Update()
		}

		cmds = append(cmds, tick())

	case CountdownMsg:
		a.authCountdown = int(msg)
		if a.authCountdown <= 0 {
			// Open browser and wait
			if a.authView != nil {
				browser.OpenURL(a.authView.GetAuthURL())
				a.authView.SetAuthState(views.AuthStateBrowserOpen)
			}
			// Move to API key view for manual fallback after browser opens
			a.transitionToAPIKey()
		} else if a.authView != nil {
			a.authView.SetCountdown(a.authCountdown)
			cmds = append(cmds, countdown(a.authCountdown-1))
		}

	case SysCheckDoneMsg:
		if a.sysCheckView != nil {
			a.sysCheckView.SetResults(msg.Results)
			a.sysCheckDone = true
		}

	case InstallDoneMsg:
		if msg.Error != nil {
			a.showError(&views.ErrorInfo{
				Code:       "INSTALL_FAILED",
				Title:      "Installation Failed",
				Message:    msg.Error.Error(),
				Suggestion: "Check the logs and try again.",
				Severity:   views.SeverityError,
				Retryable:  true,
			})
		} else {
			// Move to provider selection
			a.state = StateProviderSelect
			a.providerView.SetSize(a.width, a.height)
		}

	case AnalysisDoneMsg:
		if msg.Error != nil {
			a.showError(views.ErrAnalysisFailed)
		} else {
			a.state = StateResults
			a.resultsView = views.NewResultsView()
			a.resultsView.SetSize(a.width, a.height)
		}

	case UVInstallDoneMsg:
		if msg.Error != nil {
			a.showError(views.ErrUVInstallFailed)
		} else {
			// Re-run system checks
			a.startSysCheck()
		}

	case LocalModelDetectMsg:
		if a.localModelView != nil {
			if msg.Error != nil {
				a.localModelView.SetError(msg.Error.Error())
			} else {
				a.localModelView.SetModels(msg.Models)
			}
		}

	case game.GameTickMsg:
		if a.state == StateGame && a.game != nil {
			a.game.Update()
			cmds = append(cmds, game.GameTickCmd())
		}

	default:
		// Forward messages to welcome animation
		if a.state == StateWelcome && a.welcomeView != nil {
			animCmd := a.welcomeView.UpdateAnimation(msg)
			if animCmd != nil {
				cmds = append(cmds, animCmd)
			}
		}
	}

	return a, tea.Batch(cmds...)
}

// ═══════════════════════════════════════════════════════════════════
// KEY HANDLERS
// ═══════════════════════════════════════════════════════════════════

func (a *App) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	switch a.state {
	case StateWelcome:
		return a.handleWelcomeKeys(key)
	case StateSysCheck:
		return a.handleSysCheckKeys(key)
	case StateInstallMethod:
		return a.handleInstallMethodKeys(key)
	case StateInstalling:
		return a.handleInstallingKeys(key)
	case StateProviderSelect:
		return a.handleProviderKeys(msg)
	case StateModelSelect:
		return a.handleModelKeys(msg)
	case StateAuth:
		return a.handleAuthKeys(key)
	case StateAPIKey:
		return a.handleAPIKeyKeys(msg)
	case StateLocalModel:
		return a.handleLocalModelKeys(key)
	case StateProjectDir:
		return a.handleProjectDirKeys(msg)
	case StateAnalysisConfig:
		return a.handleAnalysisConfigKeys(key)
	case StateAnalyzing:
		return a.handleAnalyzingKeys(key)
	case StateResults:
		return a.handleResultsKeys(key)
	case StateNextSteps:
		return a.handleNextStepsKeys(key)
	case StateError:
		return a.handleErrorKeys(key)
	case StateGame:
		return a.handleGameKeys(msg)
	}

	return nil
}

func (a *App) handleWelcomeKeys(key string) tea.Cmd {
	switch key {
	case "enter":
		a.state = StateSysCheck
		a.sysCheckView = views.NewSysCheckView()
		a.sysCheckView.SetSize(a.width, a.height)
		return a.startSysCheckCmd()
	case "q":
		return tea.Quit
	}
	return nil
}

func (a *App) handleSysCheckKeys(key string) tea.Cmd {
	if !a.sysCheckDone {
		return nil // Ignore keys during check
	}

	switch key {
	case "enter":
		btn := a.sysCheckView.GetButtonLabel()
		switch btn {
		case "Continue":
			if a.sysCheckView.CanProceed() {
				a.transitionToInstallMethod()
			}
		case "Install uv":
			return a.startUVInstallCmd()
		case "Quit":
			return tea.Quit
		}
	case "left", "h":
		a.sysCheckView.HandleLeft()
	case "right", "l":
		a.sysCheckView.HandleRight()
	case "q":
		return tea.Quit
	}
	return nil
}

func (a *App) handleInstallMethodKeys(key string) tea.Cmd {
	switch key {
	case "up", "k":
		a.installMethodView.HandleUp()
	case "down", "j":
		a.installMethodView.HandleDown()
	case "enter":
		method := a.installMethodView.GetSelectedMethod()
		// Check if uvx is available when selected
		sysResults := a.sysChecker.GetResults()
		if method == "uvx" && sysResults.UV.Status != syscheck.StatusPassed {
			// Can't use uvx, show warning
			return nil
		}
		a.installMethod = method
		a.startInstalling()
	case "esc":
		a.state = StateSysCheck
	case "q":
		return tea.Quit
	}
	return nil
}

func (a *App) handleInstallingKeys(key string) tea.Cmd {
	switch key {
	case "g":
		a.prevState = a.state
		a.state = StateGame
		if a.game == nil {
			a.game = game.NewGame(60, 20)
		} else {
			a.game.Restart()
		}
		a.game.SetSize(60, 20)
		return game.GameTickCmd()
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
				a.state = StateInstallMethod
			} else {
				return a.selectProvider()
			}
		} else {
			return a.selectProvider()
		}
	case "esc":
		a.state = StateInstallMethod
	case "q":
		return tea.Quit
	}
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

func (a *App) handleAuthKeys(key string) tea.Cmd {
	switch key {
	case "m":
		// Skip to manual entry
		if a.authView != nil {
			a.authView.ShowFallback()
		}
	case "enter":
		if a.authView != nil && a.authView.IsFallbackShown() {
			a.transitionToAPIKey()
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
				if a.apiKeyView.GetBaseURL() != "" {
					a.configMgr.SetBaseURL(a.apiKeyView.GetBaseURL())
				}
				a.transitionToProjectDir()
			}
		case "tab":
			a.apiKeyView.HandleTab()
		case "esc":
			a.navigateBackFromAPIKey()
		default:
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
				a.navigateBackFromAPIKey()
			} else {
				if a.apiKeyView.Validate() {
					a.configMgr.SetAPIKey(a.apiKeyView.GetAPIKey())
					if a.apiKeyView.GetBaseURL() != "" {
						a.configMgr.SetBaseURL(a.apiKeyView.GetBaseURL())
					}
					a.transitionToProjectDir()
				}
			}
		case "tab":
			a.apiKeyView.HandleTab()
		case "esc":
			a.navigateBackFromAPIKey()
		}
	}
	return nil
}

func (a *App) handleLocalModelKeys(key string) tea.Cmd {
	if a.localModelView == nil {
		return nil
	}

	switch key {
	case "up", "k":
		a.localModelView.HandleUp()
	case "down", "j":
		a.localModelView.HandleDown()
	case "enter":
		if a.localModelView.IsFound() {
			model := a.localModelView.GetSelectedModel()
			a.configMgr.SetModel(model)
			a.configMgr.SetBaseURL(a.localModelView.GetBaseURL())
			a.transitionToProjectDir()
		}
	case "r":
		// Retry detection
		return a.detectLocalModels()
	case "esc":
		a.state = StateProviderSelect
	case "q":
		return tea.Quit
	}
	return nil
}

func (a *App) handleProjectDirKeys(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	if a.projectDirView.IsInputFocused() {
		switch key {
		case "enter":
			if a.projectDirView.IsValid() {
				a.configMgr.SetProjectDir(a.projectDirView.GetProjectDir())
				a.transitionToAnalysisConfig()
			}
		case "tab":
			a.projectDirView.HandleTab()
		case "esc":
			a.navigateBackFromProjectDir()
		default:
			a.projectDirView.Update(msg)
		}
	} else {
		switch key {
		case "left", "h":
			a.projectDirView.HandleLeft()
		case "right", "l":
			a.projectDirView.HandleRight()
		case "enter":
			btn := a.projectDirView.GetButtonLabel()
			switch btn {
			case "Use Current":
				a.projectDirView.UseCurrentDir()
			case "Continue":
				if a.projectDirView.IsValid() {
					a.configMgr.SetProjectDir(a.projectDirView.GetProjectDir())
					a.transitionToAnalysisConfig()
				}
			}
		case "tab":
			a.projectDirView.HandleTab()
		case "esc":
			a.navigateBackFromProjectDir()
		}
	}
	return nil
}

func (a *App) handleAnalysisConfigKeys(key string) tea.Cmd {
	switch key {
	case "left", "h":
		a.analysisConfigView.HandleLeft()
	case "right", "l":
		a.analysisConfigView.HandleRight()
	case "up", "k":
		a.analysisConfigView.HandleUp()
	case "down", "j":
		a.analysisConfigView.HandleDown()
	case " ":
		a.analysisConfigView.HandleSpace()
	case "enter":
		if a.analysisConfigView.IsDefaultMode() {
			btn := a.analysisConfigView.GetButtonLabel()
			if btn == "Yes" {
				a.startAnalysis()
			} else {
				a.analysisConfigView.SetCustomMode()
			}
		} else {
			a.startAnalysis()
		}
	case "y":
		if a.analysisConfigView.IsDefaultMode() {
			a.startAnalysis()
		}
	case "n":
		if a.analysisConfigView.IsDefaultMode() {
			a.analysisConfigView.SetCustomMode()
		}
	case "esc":
		a.state = StateProjectDir
	case "q":
		return tea.Quit
	}
	return nil
}

func (a *App) handleAnalyzingKeys(key string) tea.Cmd {
	switch key {
	case "g":
		a.prevState = a.state
		a.state = StateGame
		if a.game == nil {
			a.game = game.NewGame(60, 20)
		} else {
			a.game.Restart()
		}
		a.game.SetSize(60, 20)
		return game.GameTickCmd()
	case "q":
		return tea.Quit
	}
	return nil
}

func (a *App) handleResultsKeys(key string) tea.Cmd {
	switch key {
	case "left", "h":
		a.resultsView.HandleLeft()
	case "right", "l":
		a.resultsView.HandleRight()
	case "up", "k":
		a.resultsView.HandleUp()
	case "down", "j":
		a.resultsView.HandleDown()
	case "tab":
		a.resultsView.HandleTab()
	case "n":
		a.state = StateNextSteps
		a.nextStepsView = views.NewNextStepsView()
		a.nextStepsView.SetSize(a.width, a.height)
	case "q":
		return tea.Quit
	}
	return nil
}

func (a *App) handleNextStepsKeys(key string) tea.Cmd {
	switch key {
	case "up", "k":
		a.nextStepsView.HandleUp()
	case "down", "j":
		a.nextStepsView.HandleDown()
	case "enter":
		action := a.nextStepsView.GetSelectedAction()
		if action == nil {
			return nil
		}
		switch action.ID {
		case "exit":
			return tea.Quit
		case "rerun":
			a.startAnalysis()
		case "config":
			a.state = StateProviderSelect
		case "plan", "validate", "open":
			// For now, show the command
			return nil
		}
	case "esc":
		a.state = StateResults
	case "q":
		return tea.Quit
	}
	return nil
}

func (a *App) handleErrorKeys(key string) tea.Cmd {
	switch key {
	case "left", "h":
		a.errorView.HandleLeft()
	case "right", "l":
		a.errorView.HandleRight()
	case "enter":
		btn := a.errorView.GetSelectedButton()
		switch btn {
		case "Retry":
			// Go back to previous state and retry
			a.state = a.prevState
		case "View Logs":
			a.errorView.ToggleLogs()
		case "Quit":
			return tea.Quit
		}
	case "L":
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

// ═══════════════════════════════════════════════════════════════════
// STATE TRANSITIONS
// ═══════════════════════════════════════════════════════════════════

func (a *App) selectProvider() tea.Cmd {
	provider := a.providerView.GetSelectedProvider()
	if provider == nil {
		return nil
	}

	a.selectedProvider = provider
	a.configMgr.SetProvider(provider.ID)

	// Branch based on provider type
	if provider.ID == "skene" {
		// Skene: go to magic link auth
		a.authView = views.NewAuthView(provider)
		a.authView.SetSize(a.width, a.height)
		a.authCountdown = 3
		a.state = StateAuth
		return countdown(3)
	}

	if provider.IsLocal {
		// Local model: detect runtime
		a.localModelView = views.NewLocalModelView(provider.ID)
		a.localModelView.SetSize(a.width, a.height)
		a.state = StateLocalModel
		return a.detectLocalModels()
	}

	// Regular providers: go to model selection
	a.modelView = views.NewModelView(provider)
	a.modelView.SetSize(a.width, a.height)
	a.state = StateModelSelect
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
	a.transitionToAPIKey()
}

func (a *App) transitionToAPIKey() {
	a.apiKeyView = views.NewAPIKeyView(a.selectedProvider, a.selectedModel)
	a.apiKeyView.SetSize(a.width, a.height)
	a.state = StateAPIKey
}

func (a *App) transitionToInstallMethod() {
	sysResults := a.sysChecker.GetResults()
	uvAvailable := sysResults.UV.Status == syscheck.StatusPassed
	a.installMethodView = views.NewInstallMethodView(uvAvailable)
	a.installMethodView.SetSize(a.width, a.height)
	a.state = StateInstallMethod
}

func (a *App) transitionToProjectDir() {
	a.projectDirView = views.NewProjectDirView()
	a.projectDirView.SetSize(a.width, a.height)
	a.state = StateProjectDir
}

func (a *App) transitionToAnalysisConfig() {
	providerName := ""
	modelName := ""
	if a.selectedProvider != nil {
		providerName = a.selectedProvider.Name
	}
	if a.selectedModel != nil {
		modelName = a.selectedModel.Name
	}
	projectDir := a.configMgr.Config.ProjectDir
	if projectDir == "" {
		projectDir = "."
	}

	a.analysisConfigView = views.NewAnalysisConfigView(providerName, modelName, projectDir)
	a.analysisConfigView.SetSize(a.width, a.height)
	a.state = StateAnalysisConfig
}

func (a *App) navigateBackFromAPIKey() {
	if a.selectedProvider != nil {
		if a.selectedProvider.ID == "skene" {
			a.state = StateAuth
		} else if a.selectedProvider.IsGeneric {
			a.state = StateProviderSelect
		} else {
			a.state = StateModelSelect
		}
	} else {
		a.state = StateProviderSelect
	}
}

func (a *App) navigateBackFromProjectDir() {
	if a.selectedProvider != nil && a.selectedProvider.IsLocal {
		a.state = StateLocalModel
	} else {
		a.state = StateAPIKey
	}
}

// ═══════════════════════════════════════════════════════════════════
// ASYNC OPERATIONS
// ═══════════════════════════════════════════════════════════════════

func (a *App) startSysCheck() {
	a.sysCheckView = views.NewSysCheckView()
	a.sysCheckView.SetSize(a.width, a.height)
	a.sysCheckDone = false
	a.state = StateSysCheck
}

func (a *App) startSysCheckCmd() tea.Cmd {
	checker := a.sysChecker
	return func() tea.Msg {
		results := checker.RunAllChecks()
		return SysCheckDoneMsg{Results: results}
	}
}

func (a *App) startUVInstallCmd() tea.Cmd {
	checker := a.sysChecker
	return func() tea.Msg {
		err := checker.InstallUV()
		return UVInstallDoneMsg{Error: err}
	}
}

func (a *App) startInstalling() {
	a.installingView = views.NewInstallingView(a.installMethod)
	a.installingView.SetSize(a.width, a.height)
	a.installStartTime = time.Now()
	a.state = StateInstalling
}

func (a *App) startAnalysis() {
	a.analyzingView = views.NewAnalyzingView()
	a.analyzingView.SetSize(a.width, a.height)
	a.analysisStartTime = time.Now()
	a.state = StateAnalyzing
}

func (a *App) detectLocalModels() tea.Cmd {
	providerID := ""
	if a.selectedProvider != nil {
		providerID = a.selectedProvider.ID
	}

	return func() tea.Msg {
		// Simulate detection with some default models
		time.Sleep(500 * time.Millisecond)

		var models []string
		switch providerID {
		case "ollama":
			models = []string{"llama3.3", "mistral", "codellama", "deepseek-r1"}
		case "lmstudio":
			models = []string{"Currently loaded model"}
		}

		if len(models) > 0 {
			return LocalModelDetectMsg{Models: models}
		}
		return LocalModelDetectMsg{
			Error: fmt.Errorf("could not connect to local model server"),
		}
	}
}

func (a *App) showError(err *views.ErrorInfo) {
	a.prevState = a.state
	a.currentError = err
	a.errorView = views.NewErrorView(err)
	a.errorView.SetSize(a.width, a.height)
	a.state = StateError
}

// ═══════════════════════════════════════════════════════════════════
// PROGRESS SIMULATION
// ═══════════════════════════════════════════════════════════════════

func (a *App) simulateInstallProgress() {
	if a.installingView == nil || a.installingView.AllTasksDone() {
		if a.installingView != nil && a.installingView.AllTasksDone() {
			// Transition to provider selection
			a.state = StateProviderSelect
			a.providerView.SetSize(a.width, a.height)
		}
		return
	}

	elapsed := time.Since(a.installStartTime).Seconds()

	// Simulate tasks completing over time
	taskCount := 3
	if a.installMethod == "pip" {
		taskCount = 4
	}

	taskDuration := 2.0 // seconds per task
	for i := 0; i < taskCount; i++ {
		taskStart := float64(i) * taskDuration
		taskEnd := taskStart + taskDuration

		if elapsed > taskStart && elapsed <= taskEnd {
			progress := (elapsed - taskStart) / taskDuration
			if progress > 1.0 {
				progress = 1.0
			}
			a.installingView.UpdateTask(i, progress)
		} else if elapsed > taskEnd {
			a.installingView.UpdateTask(i, 1.0)
		}
	}
}

func (a *App) simulateAnalysisProgress() {
	if a.analyzingView == nil || a.analyzingView.AllPhasesDone() {
		if a.analyzingView != nil && a.analyzingView.AllPhasesDone() {
			// Transition to results
			a.state = StateResults
			a.resultsView = views.NewResultsView()
			a.resultsView.SetSize(a.width, a.height)
		}
		return
	}

	elapsed := time.Since(a.analysisStartTime).Seconds()

	// Simulate 6 phases completing over ~12 seconds
	phaseDuration := 2.0
	for i := 0; i < 6; i++ {
		phaseStart := float64(i) * phaseDuration
		phaseEnd := phaseStart + phaseDuration

		if elapsed > phaseStart && elapsed <= phaseEnd {
			progress := (elapsed - phaseStart) / phaseDuration
			if progress > 1.0 {
				progress = 1.0
			}
			a.analyzingView.UpdatePhase(i, progress)
		} else if elapsed > phaseEnd {
			a.analyzingView.UpdatePhase(i, 1.0)
		}
	}
}

// ═══════════════════════════════════════════════════════════════════
// VIEW SIZING
// ═══════════════════════════════════════════════════════════════════

func (a *App) updateViewSizes() {
	if a.welcomeView != nil {
		a.welcomeView.SetSize(a.width, a.height)
	}
	if a.sysCheckView != nil {
		a.sysCheckView.SetSize(a.width, a.height)
	}
	if a.installMethodView != nil {
		a.installMethodView.SetSize(a.width, a.height)
	}
	if a.installingView != nil {
		a.installingView.SetSize(a.width, a.height)
	}
	if a.providerView != nil {
		a.providerView.SetSize(a.width, a.height)
	}
	if a.modelView != nil {
		a.modelView.SetSize(a.width, a.height)
	}
	if a.authView != nil {
		a.authView.SetSize(a.width, a.height)
	}
	if a.apiKeyView != nil {
		a.apiKeyView.SetSize(a.width, a.height)
	}
	if a.localModelView != nil {
		a.localModelView.SetSize(a.width, a.height)
	}
	if a.projectDirView != nil {
		a.projectDirView.SetSize(a.width, a.height)
	}
	if a.analysisConfigView != nil {
		a.analysisConfigView.SetSize(a.width, a.height)
	}
	if a.analyzingView != nil {
		a.analyzingView.SetSize(a.width, a.height)
	}
	if a.resultsView != nil {
		a.resultsView.SetSize(a.width, a.height)
	}
	if a.nextStepsView != nil {
		a.nextStepsView.SetSize(a.width, a.height)
	}
	if a.errorView != nil {
		a.errorView.SetSize(a.width, a.height)
	}
	if a.game != nil {
		a.game.SetSize(60, 20)
	}
}

// ═══════════════════════════════════════════════════════════════════
// VIEW RENDERING
// ═══════════════════════════════════════════════════════════════════

// View renders the current wizard step
func (a *App) View() string {
	var content string

	switch a.state {
	case StateWelcome:
		content = a.welcomeView.Render()
	case StateSysCheck:
		if a.sysCheckView != nil {
			content = a.sysCheckView.Render()
		}
	case StateInstallMethod:
		if a.installMethodView != nil {
			content = a.installMethodView.Render()
		}
	case StateInstalling:
		if a.installingView != nil {
			content = a.installingView.Render()
		}
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
	case StateLocalModel:
		if a.localModelView != nil {
			content = a.localModelView.Render()
		}
	case StateProjectDir:
		if a.projectDirView != nil {
			content = a.projectDirView.Render()
		}
	case StateAnalysisConfig:
		if a.analysisConfigView != nil {
			content = a.analysisConfigView.Render()
		}
	case StateAnalyzing:
		if a.analyzingView != nil {
			content = a.analyzingView.Render()
		}
	case StateResults:
		if a.resultsView != nil {
			content = a.resultsView.Render()
		}
	case StateNextSteps:
		if a.nextStepsView != nil {
			content = a.nextStepsView.Render()
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
	case StateWelcome:
		return a.welcomeView.GetHelpItems()
	case StateSysCheck:
		if a.sysCheckView != nil {
			return a.sysCheckView.GetHelpItems()
		}
	case StateInstallMethod:
		if a.installMethodView != nil {
			return a.installMethodView.GetHelpItems()
		}
	case StateInstalling:
		if a.installingView != nil {
			return a.installingView.GetHelpItems()
		}
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
	case StateLocalModel:
		if a.localModelView != nil {
			return a.localModelView.GetHelpItems()
		}
	case StateProjectDir:
		if a.projectDirView != nil {
			return a.projectDirView.GetHelpItems()
		}
	case StateAnalysisConfig:
		if a.analysisConfigView != nil {
			return a.analysisConfigView.GetHelpItems()
		}
	case StateAnalyzing:
		if a.analyzingView != nil {
			return a.analyzingView.GetHelpItems()
		}
	case StateResults:
		if a.resultsView != nil {
			return a.resultsView.GetHelpItems()
		}
	case StateNextSteps:
		if a.nextStepsView != nil {
			return a.nextStepsView.GetHelpItems()
		}
	case StateError:
		if a.errorView != nil {
			return a.errorView.GetHelpItems()
		}
	}

	return components.NewHelpOverlay().Items
}

// ═══════════════════════════════════════════════════════════════════
// HELPERS
// ═══════════════════════════════════════════════════════════════════

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
