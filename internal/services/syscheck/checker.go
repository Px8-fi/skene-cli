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

	// Try python3 first, then python
	for _, cmd := range []string{"python3", "python"} {
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

		if major >= 3 {
			c.results.Python.Status = StatusFailed
			c.results.Python.Version = fmt.Sprintf("%d.%d.%d", major, minor, patch)
			c.results.Python.Message = fmt.Sprintf("Python %d.%d.%d found (requires 3.11+)", major, minor, patch)
			c.results.Python.FixCommand = getUpgradePythonCommand()
			c.results.Python.FixURL = "https://python.org/downloads/"
			return &c.results.Python
		}
	}

	c.results.Python.Status = StatusFailed
	c.results.Python.Message = "Python not found in PATH"
	c.results.Python.FixCommand = getInstallPythonCommand()
	c.results.Python.FixURL = "https://python.org/downloads/"
	return &c.results.Python
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
		return fmt.Errorf("failed to install uv: %w\nOutput: %s", err, string(output))
	}

	// Re-check uv after install
	c.CheckUV()
	return nil
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
	switch runtime.GOOS {
	case "darwin":
		return "brew upgrade python"
	case "linux":
		return "sudo apt install python3.12"
	default:
		return "Download from https://python.org/downloads/"
	}
}
