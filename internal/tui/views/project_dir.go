package views

import (
	"os"
	"path/filepath"
	"skene-terminal-v2/internal/tui/components"
	"skene-terminal-v2/internal/tui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// ProjectDirView handles project directory selection
type ProjectDirView struct {
	width       int
	height      int
	textInput   textinput.Model
	buttonGroup *components.ButtonGroup
	inputFocus  bool
	currentDir  string
	isValid     bool
	validMsg    string
	warningMsg  string
	header      *components.WizardHeader
	browsing         bool
	dirBrowser       *components.DirBrowser
	browseButtons    *components.ButtonGroup
	browseFocusList  bool // true = focus on dir listing, false = focus on buttons
}

// NewProjectDirView creates a new project directory view
func NewProjectDirView() *ProjectDirView {
	cwd, _ := os.Getwd()

	ti := textinput.New()
	ti.Placeholder = cwd
	ti.SetValue(cwd)
	ti.CharLimit = 256
	ti.Width = 50
	ti.Focus()

	bg := components.NewButtonGroup("Use Current", "Browse", "Continue")
	// Input has focus initially, so no button should be highlighted
	bg.SetActiveIndex(-1)

	return &ProjectDirView{
		textInput:   ti,
		buttonGroup: bg,
		inputFocus:  true,
		currentDir:  cwd,
		isValid:     true,
		header:      components.NewWizardHeader(5, "Project Directory"),
	}
}

// SetSize updates dimensions
func (v *ProjectDirView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(width)
}

// Update handles text input updates
func (v *ProjectDirView) Update(msg interface{}) {
	if v.inputFocus {
		v.textInput, _ = v.textInput.Update(msg)
		v.validatePath()
	}
}

// HandleTab toggles between input and buttons
func (v *ProjectDirView) HandleTab() {
	v.inputFocus = !v.inputFocus
	if v.inputFocus {
		v.textInput.Focus()
		// Deactivate buttons so none appear highlighted while typing
		v.buttonGroup.SetActiveIndex(-1)
	} else {
		v.textInput.Blur()
		// Activate the first button when moving focus to buttons
		v.buttonGroup.SetActiveIndex(0)
	}
}

// HandleLeft handles left key in buttons
func (v *ProjectDirView) HandleLeft() {
	if !v.inputFocus {
		v.buttonGroup.Previous()
	}
}

// HandleRight handles right key in buttons
func (v *ProjectDirView) HandleRight() {
	if !v.inputFocus {
		v.buttonGroup.Next()
	}
}

// IsInputFocused returns if the input is focused
func (v *ProjectDirView) IsInputFocused() bool {
	return v.inputFocus
}

// GetButtonLabel returns the selected button label
func (v *ProjectDirView) GetButtonLabel() string {
	return v.buttonGroup.GetActiveLabel()
}

// UseCurrentDir sets the current working directory
func (v *ProjectDirView) UseCurrentDir() {
	cwd, _ := os.Getwd()
	v.textInput.SetValue(cwd)
	v.currentDir = cwd
	v.validatePath()
}

// StartBrowsing activates the directory browser from the current path
func (v *ProjectDirView) StartBrowsing() {
	startPath := v.GetProjectDir()
	v.dirBrowser = components.NewDirBrowser(startPath)
	// Give the browser a reasonable height based on available space
	browserHeight := v.height - 16
	if browserHeight < 6 {
		browserHeight = 6
	}
	if browserHeight > 18 {
		browserHeight = 18
	}
	v.dirBrowser.SetHeight(browserHeight)
	v.browseButtons = components.NewButtonGroup("Select This Directory", "Cancel")
	// Deactivate buttons initially -- focus starts on the listing
	v.browseButtons.SetActiveIndex(-1)
	v.browseFocusList = true
	v.browsing = true
	v.textInput.Blur()
}

// StopBrowsing exits the directory browser without selecting
func (v *ProjectDirView) StopBrowsing() {
	v.browsing = false
	v.dirBrowser = nil
	v.browseButtons = nil
	v.inputFocus = false // return focus to buttons
}

// IsBrowsing returns true if the directory browser is active
func (v *ProjectDirView) IsBrowsing() bool {
	return v.browsing
}

// BrowseFocusOnList returns true if the dir listing has focus (not buttons)
func (v *ProjectDirView) BrowseFocusOnList() bool {
	return v.browseFocusList
}

// BrowseConfirm selects the current browsed directory and exits browsing
func (v *ProjectDirView) BrowseConfirm() {
	if v.dirBrowser == nil {
		return
	}
	// Always use the directory we're currently inside
	selectedPath := v.dirBrowser.CurrentPath()
	v.textInput.SetValue(selectedPath)
	v.currentDir = selectedPath
	v.validatePath()
	v.browsing = false
	v.dirBrowser = nil
	v.browseButtons = nil
	v.inputFocus = false
}

// GetBrowseButtonLabel returns the active browse button label
func (v *ProjectDirView) GetBrowseButtonLabel() string {
	if v.browseButtons == nil {
		return ""
	}
	return v.browseButtons.GetActiveLabel()
}

// HandleBrowseTab toggles focus between dir listing and buttons
func (v *ProjectDirView) HandleBrowseTab() {
	v.browseFocusList = !v.browseFocusList
	if v.browseFocusList {
		// Deactivate buttons when moving focus to listing
		v.browseButtons.SetActiveIndex(-1)
	} else {
		// Activate first button when moving focus to buttons
		v.browseButtons.SetActiveIndex(0)
	}
}

// HandleBrowseLeft handles left key in browse button area
func (v *ProjectDirView) HandleBrowseLeft() {
	if !v.browseFocusList && v.browseButtons != nil {
		v.browseButtons.Previous()
	}
}

// HandleBrowseRight handles right key in browse button area
func (v *ProjectDirView) HandleBrowseRight() {
	if !v.browseFocusList && v.browseButtons != nil {
		v.browseButtons.Next()
	}
}

// HandleBrowseKey handles key input for the directory listing
func (v *ProjectDirView) HandleBrowseKey(key string) {
	if v.dirBrowser == nil {
		return
	}
	switch key {
	case "up", "k":
		v.dirBrowser.CursorUp()
	case "down", "j":
		v.dirBrowser.CursorDown()
	case "enter":
		// If it's a directory, navigate into it
		if v.dirBrowser.SelectedIsDir() {
			v.dirBrowser.Enter()
		}
	case "backspace":
		v.dirBrowser.GoUp()
	case ".":
		v.dirBrowser.ToggleHidden()
	}
}

// GetProjectDir returns the entered/selected directory
func (v *ProjectDirView) GetProjectDir() string {
	val := v.textInput.Value()
	if val == "" {
		return v.currentDir
	}
	// Expand ~ to home dir
	if len(val) > 0 && val[0] == '~' {
		home, _ := os.UserHomeDir()
		val = filepath.Join(home, val[1:])
	}
	return val
}

// IsValid returns if the path is valid
func (v *ProjectDirView) IsValid() bool {
	return v.isValid
}

// HasWarning returns true if there's a non-blocking warning
func (v *ProjectDirView) HasWarning() bool {
	return v.warningMsg != ""
}

func (v *ProjectDirView) validatePath() {
	path := v.GetProjectDir()

	info, err := os.Stat(path)
	if err != nil {
		v.isValid = false
		v.validMsg = "Directory not found"
		v.warningMsg = ""
		return
	}

	if !info.IsDir() {
		v.isValid = false
		v.validMsg = "Path is not a directory"
		v.warningMsg = ""
		return
	}

	v.isValid = true
	v.validMsg = ""

	// Check for common project indicators
	hasProject := false
	projectMarkers := []string{
		"package.json", "pyproject.toml", "requirements.txt",
		"go.mod", "Cargo.toml", "pom.xml", "build.gradle",
		".git", "Makefile",
	}
	for _, marker := range projectMarkers {
		if _, err := os.Stat(filepath.Join(path, marker)); err == nil {
			hasProject = true
			break
		}
	}

	if !hasProject {
		v.warningMsg = "No recognizable project structure detected. Analysis may be limited."
	} else {
		v.warningMsg = ""
	}
}

// Render the project directory view
func (v *ProjectDirView) Render() string {
	sectionWidth := v.width - 20
	if sectionWidth < 60 {
		sectionWidth = 60
	}
	if sectionWidth > 80 {
		sectionWidth = 80
	}

	// Wizard header
	wizHeader := lipgloss.NewStyle().Width(sectionWidth).Render(v.header.Render())

	if v.browsing && v.dirBrowser != nil {
		// Render the directory browser
		browserSection := v.dirBrowser.Render(sectionWidth)

		// Buttons (uses the real browseButtons group which tracks focus)
		browseBtns := lipgloss.NewStyle().
			Width(sectionWidth).
			Align(lipgloss.Center).
			Render(v.browseButtons.Render())

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			wizHeader,
			"",
			browserSection,
			"",
			browseBtns,
		)

		padded := lipgloss.NewStyle().PaddingTop(2).Render(content)

		centered := lipgloss.Place(
			v.width,
			v.height-3,
			lipgloss.Center,
			lipgloss.Top,
			padded,
		)

		return centered
	}

	// Directory selection section
	dirSection := v.renderDirSection(sectionWidth)

	// Buttons
	buttons := lipgloss.NewStyle().
		Width(sectionWidth).
		Align(lipgloss.Center).
		Render(v.buttonGroup.Render())

	// Footer
	footer := lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Render(components.WizardInputHelp())

	// Combine
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		wizHeader,
		"",
		dirSection,
		"",
		buttons,
	)

	padded := lipgloss.NewStyle().PaddingTop(2).Render(content)

	centered := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Top,
		padded,
	)

	return centered + "\n" + footer
}

func (v *ProjectDirView) renderDirSection(width int) string {
	header := styles.SectionHeader.Render("Select project to analyze")
	subtitle := styles.Muted.Render("Enter the path to your project's root directory")

	// Directory input
	dirLabel := styles.Label.Render("Directory:")
	inputField := v.textInput.View()

	// Validation status
	var validationLine string
	if !v.isValid && v.validMsg != "" {
		validationLine = styles.Error.Render("✗ " + v.validMsg)
	} else if v.warningMsg != "" {
		validationLine = lipgloss.NewStyle().Foreground(styles.Warning).Render("! " + v.warningMsg)
	} else if v.isValid {
		validationLine = styles.SuccessText.Render("✓ Valid project directory")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		subtitle,
		"",
		dirLabel,
		inputField,
		"",
		validationLine,
	)

	return styles.Box.Width(width).Render(content)
}

// GetHelpItems returns context-specific help
func (v *ProjectDirView) GetHelpItems() []components.HelpItem {
	return []components.HelpItem{
		{Key: "enter", Desc: "confirm"},
		{Key: "tab", Desc: "switch focus"},
		{Key: "esc", Desc: "go back"},
		{Key: "ctrl+c", Desc: "quit"},
	}
}
