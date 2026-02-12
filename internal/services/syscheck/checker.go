package syscheck

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
			AllPassed: true,
			CanProceed: true,
		},
	}
}

// GetResults returns the current results
func (c *Checker) GetResults() *SystemCheckResult {
	return c.results
}

// RunAllChecks executes all system checks
func (c *Checker) RunAllChecks() *SystemCheckResult {
	// No checks needed for Rust engine
	c.results.AllPassed = true
	c.results.CanProceed = true
	return c.results
}

// Dummy methods to satisfy interface if needed
func (c *Checker) InstallUV() error { return nil }
func (c *Checker) GetAlternativeInstallCommands() []string { return []string{} }
