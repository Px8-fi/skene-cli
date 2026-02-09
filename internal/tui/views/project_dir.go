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

	return &ProjectDirView{
		textInput:   ti,
		buttonGroup: components.NewButtonGroup("Use Current", "Browse", "Continue"),
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
	} else {
		v.textInput.Blur()
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
	wizHeader := v.header.Render()

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
		lipgloss.Center,
		wizHeader,
		"",
		dirSection,
		"",
		buttons,
	)

	centered := lipgloss.Place(
		v.width,
		v.height-3,
		lipgloss.Center,
		lipgloss.Center,
		content,
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
