package syscheck

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// CheckStatus represents the result of a single check
type CheckStatus int

const (
	StatusPending  CheckStatus = iota
	StatusRunning
	StatusPassed
	StatusFailed
	StatusWarning
	StatusSkipped
)

// CheckResult holds the result of a system check
type CheckResult struct {
	Name        string
	Status      CheckStatus
	Message     string
	Detail      string
	FixCommand  string
	FixURL      string
	Version     string
	Required    bool
}

// SystemCheckResult holds all system check results
type SystemCheckResult struct {
	Python    CheckResult
	UV        CheckResult
	Pip       CheckResult
	AllPassed bool
	CanProceed bool
}

// Checker performs system prerequisite checks
type Checker struct {
	results *SystemCheckResult
}

// NewChecker creates a new system checker
func NewChecker() *Checker {
	return &Checker{
		results: &SystemCheckResult{
			Python: CheckResult{
				Name:     "Python",
				Status:   StatusPending,
				Required: true,
			},
			UV: CheckResult{
				Name:     "uv",
				Status:   StatusPending,
				Required: false,
			},
			Pip: CheckResult{
				Name:     "pip",
				Status:   StatusPending,
				Required: false,
			},
		},
	}
}

// GetResults returns the current results
func (c *Checker) GetResults() *SystemCheckResult {
	return c.results
}

// RunAllChecks executes all system checks
func (c *Checker) RunAllChecks() *SystemCheckResult {
	c.CheckPython()
	c.CheckUV()
	c.CheckPip()
	c.evaluateResults()
	return c.results
}

// CheckPython verifies Python 3.11+ is available
func (c *Checker) CheckPython() *CheckResult {
	c.results.Python.Status = StatusRunning

	// Try multiple Python commands in order of preference:
	// 1. python3.12, python3.11 (specific versions)
	// 2. python3 (system default)
	// 3. python (fallback)
	// 4. Check Homebrew installation paths
	commands := []string{"python3.12", "python3.11", "python3", "python"}
	
	// Also check Homebrew paths if on macOS
	if runtime.GOOS == "darwin" {
		homebrewPaths := []string{
			"/opt/homebrew/bin/python3.12",
			"/opt/homebrew/bin/python3.11",
			"/opt/homebrew/bin/python3",
			"/usr/local/bin/python3.12",
			"/usr/local/bin/python3.11",
			"/usr/local/bin/python3",
		}
		commands = append(commands, homebrewPaths...)
	}

	for _, cmd := range commands {
		version, err := getCommandVersion(cmd, "--version")
		if err != nil {
			continue
		}

		// Parse version string like "Python 3.12.1"
		major, minor, patch := parsePythonVersion(version)
		if major >= 3 && minor >= 11 {
			c.results.Python.Status = StatusPassed
			c.results.Python.Version = fmt.Sprintf("%d.%d.%d", major, minor, patch)
			c.results.Python.Message = fmt.Sprintf("Python %d.%d.%d found", major, minor, patch)
			c.results.Python.Detail = fmt.Sprintf("Using: %s", cmd)
			return &c.results.Python
		}

		// If we found Python but wrong version, note it but keep looking for a newer one
		if major >= 3 && minor < 11 {
			// Store the old version info but continue searching
			if c.results.Python.Status == StatusRunning {
				c.results.Python.Version = fmt.Sprintf("%d.%d.%d", major, minor, patch)
			}
			continue
		}
	}

	// If we found an old Python version, report it
	if c.results.Python.Version != "" {
		c.results.Python.Status = StatusFailed
		c.results.Python.Message = fmt.Sprintf("Python %s found (requires 3.11+). Python 3.12 may be installed but not in PATH.", c.results.Python.Version)
		fixCmd := c.getPythonUpgradeCommand()
		fixCmd += "\n\nTo use Python 3.12 from Homebrew:\n"
		fixCmd += "  export PATH=\"/opt/homebrew/bin:$PATH\"\n"
		fixCmd += "  # Or add to ~/.zshrc: echo 'export PATH=\"/opt/homebrew/bin:$PATH\"' >> ~/.zshrc"
		c.results.Python.FixCommand = fixCmd
		c.results.Python.FixURL = "https://python.org/downloads/"
		return &c.results.Python
	}

	c.results.Python.Status = StatusFailed
	c.results.Python.Message = "Python not found in PATH"
	c.results.Python.FixCommand = c.getPythonInstallCommand()
	c.results.Python.FixURL = "https://python.org/downloads/"
	return &c.results.Python
}

// getPythonInstallCommand returns the best command to install Python
func (c *Checker) getPythonInstallCommand() string {
	// Check if uv is available - it can install Python!
	if c.results.UV.Status == StatusPassed {
		return "uv python install 3.12\n(If permission error: mkdir -p ~/.local/share/uv/python first)\n\nOr via Homebrew: brew install python@3.12"
	}
	
	// Check if uv command exists (might not be checked yet)
	if _, err := exec.LookPath("uv"); err == nil {
		return "uv python install 3.12\n(If permission error: mkdir -p ~/.local/share/uv/python first)\n\nOr via Homebrew: brew install python@3.12"
	}
	
	return getInstallPythonCommand()
}

// getPythonUpgradeCommand returns the best command to upgrade Python
func (c *Checker) getPythonUpgradeCommand() string {
	// Check if uv is available - it can install Python!
	if c.results.UV.Status == StatusPassed {
		return "uv python install 3.12\n(If permission error: mkdir -p ~/.local/share/uv/python first)\n\nOr via Homebrew: brew upgrade python"
	}
	
	// Check if uv command exists (might not be checked yet)
	if _, err := exec.LookPath("uv"); err == nil {
		return "uv python install 3.12\n(If permission error: mkdir -p ~/.local/share/uv/python first)\n\nOr via Homebrew: brew upgrade python"
	}
	
	return getUpgradePythonCommand()
}

// CheckUV verifies uv is installed
func (c *Checker) CheckUV() *CheckResult {
	c.results.UV.Status = StatusRunning

	version, err := getCommandVersion("uv", "--version")
	if err != nil {
		c.results.UV.Status = StatusFailed
		c.results.UV.Message = "uv not found"
		c.results.UV.FixCommand = "curl -LsSf https://astral.sh/uv/install.sh | sh"
		c.results.UV.FixURL = "https://docs.astral.sh/uv/"
		return &c.results.UV
	}

	// Parse version like "uv 0.5.1"
	c.results.UV.Status = StatusPassed
	c.results.UV.Version = strings.TrimPrefix(strings.TrimSpace(version), "uv ")
	c.results.UV.Message = fmt.Sprintf("uv %s found", c.results.UV.Version)
	return &c.results.UV
}

// CheckPip verifies pip is available
func (c *Checker) CheckPip() *CheckResult {
	c.results.Pip.Status = StatusRunning

	for _, cmd := range []string{"pip3", "pip"} {
		version, err := getCommandVersion(cmd, "--version")
		if err != nil {
			continue
		}

		// Parse "pip 24.0 from ..."
		parts := strings.Fields(version)
		if len(parts) >= 2 {
			c.results.Pip.Status = StatusPassed
			c.results.Pip.Version = parts[1]
			c.results.Pip.Message = fmt.Sprintf("pip %s found", parts[1])
			return &c.results.Pip
		}
	}

	c.results.Pip.Status = StatusFailed
	c.results.Pip.Message = "pip not found"
	c.results.Pip.FixCommand = "python3 -m ensurepip --upgrade"
	return &c.results.Pip
}

// evaluateResults determines overall status
func (c *Checker) evaluateResults() {
	pythonOK := c.results.Python.Status == StatusPassed
	uvOK := c.results.UV.Status == StatusPassed

	// All passed if Python and at least one package manager (uv or pip)
	c.results.AllPassed = pythonOK && uvOK
	c.results.CanProceed = pythonOK && (uvOK || c.results.Pip.Status == StatusPassed)
}

// CanInstallUV returns true if we can attempt to install uv
func (c *Checker) CanInstallUV() bool {
	return c.results.UV.Status == StatusFailed
}

// GetUVInstallCommand returns the platform-specific uv install command
func (c *Checker) GetUVInstallCommand() string {
	if runtime.GOOS == "windows" {
		return "powershell -ExecutionPolicy ByPass -c \"irm https://astral.sh/uv/install.ps1 | iex\""
	}
	return "curl -LsSf https://astral.sh/uv/install.sh | sh"
}

// InstallUVError represents a uv installation error with alternatives
type InstallUVError struct {
	OriginalError error
	Output        string
	Alternatives  []string
}

func (e *InstallUVError) Error() string {
	msg := fmt.Sprintf("failed to install uv: %v", e.OriginalError)
	if e.Output != "" {
		msg += "\nOutput: " + e.Output
	}
	return msg
}

// InstallUV attempts to install uv
func (c *Checker) InstallUV() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-ExecutionPolicy", "ByPass", "-c",
			"irm https://astral.sh/uv/install.ps1 | iex")
	} else {
		cmd = exec.Command("sh", "-c", "curl -LsSf https://astral.sh/uv/install.sh | sh")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		alternatives := c.getAlternativeInstallMethods(outputStr)
		return &InstallUVError{
			OriginalError: err,
			Output:        outputStr,
			Alternatives:   alternatives,
		}
	}

	// Re-check uv after install
	c.CheckUV()
	return nil
}

// getAlternativeInstallMethods returns alternative ways to install uv based on the error
func (c *Checker) getAlternativeInstallMethods(errorOutput string) []string {
	alternatives := []string{}
	
	// Check for permission denied errors
	if strings.Contains(errorOutput, "Permission denied") || 
	   strings.Contains(errorOutput, "Operation not permitted") ||
	   strings.Contains(errorOutput, "mkdir") {
		alternatives = append(alternatives, 
			"Create directory first: mkdir -p ~/.local/bin",
			"Install via Homebrew: brew install uv",
			"Use custom install path: curl -LsSf https://astral.sh/uv/install.sh | UV_INSTALL_DIR=~/bin sh",
			"Install to /usr/local/bin (requires sudo): curl -LsSf https://astral.sh/uv/install.sh | UV_INSTALL_DIR=/usr/local/bin sh",
		)
	}
	
	// Check for network errors
	if strings.Contains(errorOutput, "curl") || strings.Contains(errorOutput, "network") {
		alternatives = append(alternatives,
			"Check your internet connection",
			"Try installing via Homebrew: brew install uv",
		)
	}
	
	// Default alternatives if none matched
	if len(alternatives) == 0 {
		alternatives = append(alternatives,
			"Install via Homebrew: brew install uv",
			"Create ~/.local/bin directory first: mkdir -p ~/.local/bin",
			"Use custom install path: curl -LsSf https://astral.sh/uv/install.sh | UV_INSTALL_DIR=~/bin sh",
		)
	}
	
	return alternatives
}

// GetAlternativeInstallCommands returns formatted alternative installation commands
func (c *Checker) GetAlternativeInstallCommands() []string {
	return []string{
		"brew install uv",
		"mkdir -p ~/.local/bin && curl -LsSf https://astral.sh/uv/install.sh | sh",
		"curl -LsSf https://astral.sh/uv/install.sh | UV_INSTALL_DIR=~/bin sh",
	}
}

// Helper functions

func getCommandVersion(command string, args ...string) (string, error) {
	path, err := exec.LookPath(command)
	if err != nil {
		return "", fmt.Errorf("command not found: %s", command)
	}

	cmd := exec.Command(path, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run %s: %w", command, err)
	}

	return strings.TrimSpace(string(output)), nil
}

func parsePythonVersion(versionStr string) (major, minor, patch int) {
	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(versionStr)
	if len(matches) < 4 {
		return 0, 0, 0
	}

	major, _ = strconv.Atoi(matches[1])
	minor, _ = strconv.Atoi(matches[2])
	patch, _ = strconv.Atoi(matches[3])
	return
}

func getInstallPythonCommand() string {
	// Check if uv is available - it can install Python!
	if _, err := exec.LookPath("uv"); err == nil {
		return "uv python install 3.12"
	}
	
	switch runtime.GOOS {
	case "darwin":
		return "brew install python@3.12"
	case "linux":
		return "sudo apt install python3.12"
	default:
		return "Download from https://python.org/downloads/"
	}
}

func getUpgradePythonCommand() string {
	// Check if uv is available - it can install Python!
	if _, err := exec.LookPath("uv"); err == nil {
		return "uv python install 3.12"
	}
	
	switch runtime.GOOS {
	case "darwin":
		return "brew upgrade python"
	case "linux":
		return "sudo apt install python3.12"
	default:
		return "Download from https://python.org/downloads/"
	}
}

// InstallPythonViaUVError represents a Python installation error via uv with alternatives
type InstallPythonViaUVError struct {
	OriginalError error
	Output        string
	Alternatives  []string
}

func (e *InstallPythonViaUVError) Error() string {
	msg := fmt.Sprintf("failed to install Python via uv: %v", e.OriginalError)
	if e.Output != "" {
		msg += "\nOutput: " + e.Output
	}
	return msg
}

// InstallPythonViaUV attempts to install Python using uv
func (c *Checker) InstallPythonViaUV(version string) error {
	if version == "" {
		version = "3.12"
	}
	
	cmd := exec.Command("uv", "python", "install", version)
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		alternatives := c.getAlternativePythonInstallMethods(outputStr)
		return &InstallPythonViaUVError{
			OriginalError: err,
			Output:        outputStr,
			Alternatives:  alternatives,
		}
	}
	
	// Re-check Python after install
	c.CheckPython()
	return nil
}

// getAlternativePythonInstallMethods returns alternative ways to install Python based on the error
func (c *Checker) getAlternativePythonInstallMethods(errorOutput string) []string {
	alternatives := []string{}
	
	// Check for permission denied errors
	if strings.Contains(errorOutput, "Permission denied") || 
	   strings.Contains(errorOutput, "Operation not permitted") ||
	   strings.Contains(errorOutput, "failed to create directory") {
		alternatives = append(alternatives,
			"Create directory first: mkdir -p ~/.local/share/uv/python",
			"Install via Homebrew: brew install python@3.12",
			"Install via Homebrew (upgrade): brew upgrade python",
		)
	}
	
	// Check for network errors
	if strings.Contains(errorOutput, "network") || strings.Contains(errorOutput, "download") {
		alternatives = append(alternatives,
			"Check your internet connection",
			"Try installing via Homebrew: brew install python@3.12",
		)
	}
	
	// Default alternatives if none matched
	if len(alternatives) == 0 {
		alternatives = append(alternatives,
			"Install via Homebrew: brew install python@3.12",
			"Create ~/.local/share/uv/python directory first: mkdir -p ~/.local/share/uv/python",
		)
	}
	
	return alternatives
}

// CanInstallPythonViaUV returns true if uv is available and can install Python
func (c *Checker) CanInstallPythonViaUV() bool {
	return c.results.UV.Status == StatusPassed && c.results.Python.Status == StatusFailed
}
