package analyzer

import (
	"os"
	"path/filepath"
)

// ProjectType represents the detected project type
type ProjectType string

const (
	ProjectTypePython  ProjectType = "python"
	ProjectTypeNode    ProjectType = "node"
	ProjectTypeGo      ProjectType = "go"
	ProjectTypeUnknown ProjectType = "unknown"
)

// InstallMethod represents how to install dependencies
type InstallMethod string

const (
	InstallMethodPip    InstallMethod = "pip"
	InstallMethodPoetry InstallMethod = "poetry"
	InstallMethodPipenv InstallMethod = "pipenv"
	InstallMethodUV     InstallMethod = "uv"
)

// AnalysisResult contains project analysis results
type AnalysisResult struct {
	ProjectPath   string
	ProjectType   ProjectType
	InstallMethod InstallMethod
	HasVenv       bool
	VenvPath      string
	PythonVersion string
	HasConfig     bool
	ConfigPath    string
	Dependencies  []string
	Errors        []AnalysisError
}

// AnalysisError represents an error during analysis
type AnalysisError struct {
	Code    string
	Message string
	Fix     string
}

// Analyzer analyzes project structure
type Analyzer struct {
	projectPath string
	result      *AnalysisResult
}

// NewAnalyzer creates a new analyzer
func NewAnalyzer(projectPath string) *Analyzer {
	return &Analyzer{
		projectPath: projectPath,
		result: &AnalysisResult{
			ProjectPath: projectPath,
			Errors:      make([]AnalysisError, 0),
		},
	}
}

// Analyze performs full project analysis
func (a *Analyzer) Analyze() *AnalysisResult {
	a.detectProjectType()
	a.detectInstallMethod()
	a.detectVirtualEnv()
	a.detectExistingConfig()
	a.validateEnvironment()

	return a.result
}

// detectProjectType determines what kind of project this is
func (a *Analyzer) detectProjectType() {
	// Check for Python project markers
	pythonMarkers := []string{
		"pyproject.toml",
		"requirements.txt",
		"setup.py",
		"setup.cfg",
		"Pipfile",
	}

	for _, marker := range pythonMarkers {
		if a.fileExists(marker) {
			a.result.ProjectType = ProjectTypePython
			return
		}
	}

	// Check for Node.js
	if a.fileExists("package.json") {
		a.result.ProjectType = ProjectTypeNode
		return
	}

	// Check for Go
	if a.fileExists("go.mod") {
		a.result.ProjectType = ProjectTypeGo
		return
	}

	a.result.ProjectType = ProjectTypeUnknown
}

// detectInstallMethod determines how to install Python packages
func (a *Analyzer) detectInstallMethod() {
	if a.result.ProjectType != ProjectTypePython {
		return
	}

	// Priority: poetry > uv > pipenv > pip
	if a.fileExists("poetry.lock") || a.containsKey("pyproject.toml", "[tool.poetry]") {
		a.result.InstallMethod = InstallMethodPoetry
		return
	}

	if a.fileExists("uv.lock") {
		a.result.InstallMethod = InstallMethodUV
		return
	}

	if a.fileExists("Pipfile") || a.fileExists("Pipfile.lock") {
		a.result.InstallMethod = InstallMethodPipenv
		return
	}

	// Default to pip
	a.result.InstallMethod = InstallMethodPip
}

// detectVirtualEnv checks for existing virtual environments
func (a *Analyzer) detectVirtualEnv() {
	if a.result.ProjectType != ProjectTypePython {
		return
	}

	venvPaths := []string{
		".venv",
		"venv",
		".env",
		"env",
	}

	for _, vp := range venvPaths {
		fullPath := filepath.Join(a.projectPath, vp)
		if a.dirExists(fullPath) {
			// Check for activate script
			activatePath := filepath.Join(fullPath, "bin", "activate")
			if a.fileExists(activatePath) || a.fileExists(filepath.Join(fullPath, "Scripts", "activate.bat")) {
				a.result.HasVenv = true
				a.result.VenvPath = fullPath
				return
			}
		}
	}
}

// detectExistingConfig checks for skene config
func (a *Analyzer) detectExistingConfig() {
	configPaths := []string{
		".skene.config",
		"skene.config.json",
		".skene/config.json",
	}

	for _, cp := range configPaths {
		fullPath := filepath.Join(a.projectPath, cp)
		if a.fileExists(fullPath) {
			a.result.HasConfig = true
			a.result.ConfigPath = fullPath
			return
		}
	}
}

// validateEnvironment checks for required tools
func (a *Analyzer) validateEnvironment() {
	if a.result.ProjectType == ProjectTypePython {
		// Check for Python
		if !a.commandExists("python3") && !a.commandExists("python") {
			a.result.Errors = append(a.result.Errors, AnalysisError{
				Code:    "PYTHON_NOT_FOUND",
				Message: "Python is not installed or not in PATH",
				Fix:     "Install Python 3.8+ from python.org or use your package manager",
			})
		}

		// Check for pip based on install method
		switch a.result.InstallMethod {
		case InstallMethodPip:
			if !a.commandExists("pip3") && !a.commandExists("pip") {
				a.result.Errors = append(a.result.Errors, AnalysisError{
					Code:    "PIP_NOT_FOUND",
					Message: "pip is not installed",
					Fix:     "Run: python -m ensurepip --upgrade",
				})
			}
		case InstallMethodPoetry:
			if !a.commandExists("poetry") {
				a.result.Errors = append(a.result.Errors, AnalysisError{
					Code:    "POETRY_NOT_FOUND",
					Message: "poetry is not installed",
					Fix:     "Run: curl -sSL https://install.python-poetry.org | python3 -",
				})
			}
		case InstallMethodUV:
			if !a.commandExists("uv") {
				a.result.Errors = append(a.result.Errors, AnalysisError{
					Code:    "UV_NOT_FOUND",
					Message: "uv is not installed",
					Fix:     "Run: curl -LsSf https://astral.sh/uv/install.sh | sh",
				})
			}
		}
	}
}

// Helper methods
func (a *Analyzer) fileExists(name string) bool {
	path := filepath.Join(a.projectPath, name)
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func (a *Analyzer) dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func (a *Analyzer) containsKey(filename, key string) bool {
	path := filepath.Join(a.projectPath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return containsString(string(data), key)
}

func (a *Analyzer) commandExists(cmd string) bool {
	// Simple check - in production, use exec.LookPath
	paths := []string{
		"/usr/bin/" + cmd,
		"/usr/local/bin/" + cmd,
		"/opt/homebrew/bin/" + cmd,
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

func containsString(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (haystack == needle || len(haystack) > 0 && containsString(haystack[1:], needle) || haystack[:len(needle)] == needle)
}

// GetInstallCommand returns the command to install skene-growth
func (a *Analyzer) GetInstallCommand() string {
	switch a.result.InstallMethod {
	case InstallMethodPoetry:
		return "poetry add skene-growth"
	case InstallMethodPipenv:
		return "pipenv install skene-growth"
	case InstallMethodUV:
		return "uv pip install skene-growth"
	default:
		return "pip install skene-growth"
	}
}

// GetVenvCreateCommand returns command to create venv
func (a *Analyzer) GetVenvCreateCommand() string {
	switch a.result.InstallMethod {
	case InstallMethodPoetry:
		return "poetry install"
	case InstallMethodPipenv:
		return "pipenv install"
	case InstallMethodUV:
		return "uv venv"
	default:
		return "python -m venv .venv"
	}
}
