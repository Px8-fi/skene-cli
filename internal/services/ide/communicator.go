package ide

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"skene-terminal-v2/internal/services/syscheck"
	"strings"
	"time"
)

// IDEIssue represents an issue that can be sent to the IDE
type IDEIssue struct {
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	FailedChecks []FailedCheck         `json:"failed_checks"`
	Timestamp   string                 `json:"timestamp"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// FailedCheck represents a failed system check
type FailedCheck struct {
	Name       string `json:"name"`
	Message    string `json:"message"`
	FixCommand string `json:"fix_command"`
	FixURL     string `json:"fix_url,omitempty"`
	Required   bool   `json:"required"`
}

// Communicator handles communication with the IDE
type Communicator struct {
	workspacePath string
}

// NewCommunicator creates a new IDE communicator
func NewCommunicator(workspacePath string) *Communicator {
	return &Communicator{
		workspacePath: workspacePath,
	}
}

// SendSystemCheckIssues sends failed system check results to the IDE
func (c *Communicator) SendSystemCheckIssues(results *syscheck.SystemCheckResult) error {
	var failedChecks []FailedCheck

	// Collect failed checks
	if results.Python.Status == syscheck.StatusFailed {
		// Add alternative installation methods for Python
		alternatives := []string{
			"brew install python@3.12",
			"brew upgrade python",
		}
		// If uv is available, add it as first option but note potential permission issues
		if results.UV.Status == syscheck.StatusPassed {
			alternatives = append([]string{
				"uv python install 3.12",
				"(If permission error: mkdir -p ~/.local/share/uv/python first)",
			}, alternatives...)
		}
		
		fixCmd := results.Python.FixCommand
		if len(alternatives) > 0 {
			fixCmd += "\n\nAlternative methods:\n" + strings.Join(alternatives, "\n")
		}
		
		failedChecks = append(failedChecks, FailedCheck{
			Name:       results.Python.Name,
			Message:    results.Python.Message,
			FixCommand: fixCmd,
			FixURL:     results.Python.FixURL,
			Required:   results.Python.Required,
		})
	}

	if results.UV.Status == syscheck.StatusFailed {
		// Add alternative installation methods
		alternatives := []string{
			"brew install uv",
			"mkdir -p ~/.local/bin && curl -LsSf https://astral.sh/uv/install.sh | sh",
			"curl -LsSf https://astral.sh/uv/install.sh | UV_INSTALL_DIR=~/bin sh",
		}
		fixCmd := results.UV.FixCommand
		fixCmd += "\n\nAlternative methods:\n" + strings.Join(alternatives, "\n")
		
		failedChecks = append(failedChecks, FailedCheck{
			Name:       results.UV.Name,
			Message:    results.UV.Message,
			FixCommand: fixCmd,
			FixURL:     results.UV.FixURL,
			Required:   results.UV.Required,
		})
	}

	if results.Pip.Status == syscheck.StatusFailed {
		failedChecks = append(failedChecks, FailedCheck{
			Name:       results.Pip.Name,
			Message:    results.Pip.Message,
			FixCommand: results.Pip.FixCommand,
			Required:   results.Pip.Required,
		})
	}

	if len(failedChecks) == 0 {
		return fmt.Errorf("no failed checks to send")
	}

	// Determine issue title and description
	title := "System Prerequisites Check Failed"
	description := "The system check found issues that need to be resolved before proceeding."
	if results.Python.Status == syscheck.StatusFailed {
		title = "Python Version Issue"
		description = "Python 3.11+ is required but not found or outdated."
	}

	issue := IDEIssue{
		Type:         "system_check",
		Title:        title,
		Description:  description,
		FailedChecks: failedChecks,
		Timestamp:    time.Now().Format(time.RFC3339),
		Context: map[string]interface{}{
			"can_proceed": results.CanProceed,
			"all_passed":  results.AllPassed,
		},
	}

	return c.writeIssueToFile(issue)
}

// writeIssueToFile writes the issue to a file that the IDE can pick up
func (c *Communicator) writeIssueToFile(issue IDEIssue) error {
	// Try multiple locations for IDE communication
	jsonLocation := filepath.Join(c.workspacePath, ".cursor", "skene-request.json")
	mdLocation := filepath.Join(c.workspacePath, ".cursor", "skene-request.md")
	
	// Ensure .cursor directory exists
	dir := filepath.Dir(jsonLocation)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create .cursor directory: %w", err)
	}

	// Write JSON file
	data, err := json.MarshalIndent(issue, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(jsonLocation, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	// Write Markdown file for easier reading
	mdContent := c.generateMarkdown(issue)
	if err := os.WriteFile(mdLocation, []byte(mdContent), 0644); err != nil {
		// Don't fail if markdown write fails, JSON is more important
	}

	// Success - also write to stderr in a format Cursor can detect
	fmt.Fprintf(os.Stderr, "\n[CURSOR_REQUEST] File written to: %s\n", jsonLocation)
	fmt.Fprintf(os.Stderr, "[CURSOR_REQUEST] Markdown: %s\n", mdLocation)
	fmt.Fprintf(os.Stderr, "[CURSOR_REQUEST] Type: %s\n", issue.Type)
	fmt.Fprintf(os.Stderr, "[CURSOR_REQUEST] Title: %s\n", issue.Title)
	
	return nil
}

// generateMarkdown creates a markdown representation of the issue
func (c *Communicator) generateMarkdown(issue IDEIssue) string {
	var md strings.Builder
	
	md.WriteString(fmt.Sprintf("# %s\n\n", issue.Title))
	md.WriteString(fmt.Sprintf("%s\n\n", issue.Description))
	md.WriteString(fmt.Sprintf("**Type:** `%s`  \n", issue.Type))
	md.WriteString(fmt.Sprintf("**Timestamp:** %s\n\n", issue.Timestamp))
	
	md.WriteString("## Failed Checks\n\n")
	for i, check := range issue.FailedChecks {
		md.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, check.Name))
		md.WriteString("- **Status:** ❌ Failed\n")
		md.WriteString(fmt.Sprintf("- **Message:** %s\n", check.Message))
		md.WriteString(fmt.Sprintf("- **Required:** %v\n", check.Required))
		
		if check.FixCommand != "" {
			// Split fix command by newlines to handle alternatives better
			commands := strings.Split(check.FixCommand, "\n")
			md.WriteString("- **Fix Commands:**\n")
			for _, cmd := range commands {
				cmd = strings.TrimSpace(cmd)
				if cmd != "" {
					if strings.HasPrefix(cmd, "Alternative") {
						md.WriteString(fmt.Sprintf("  **%s**\n", cmd))
					} else {
						md.WriteString(fmt.Sprintf("  ```bash\n  %s\n  ```\n", cmd))
					}
				}
			}
		}
		
		if check.FixURL != "" {
			md.WriteString(fmt.Sprintf("- **More Info:** %s\n", check.FixURL))
		}
		
		md.WriteString("\n")
	}
	
	md.WriteString("## Context\n\n")
	md.WriteString(fmt.Sprintf("- **Can Proceed:** %v\n", issue.Context["can_proceed"]))
	md.WriteString(fmt.Sprintf("- **All Passed:** %v\n", issue.Context["all_passed"]))
	
	md.WriteString("\n---\n\n")
	md.WriteString("*This file was automatically generated by skene-cli. You can ask Cursor IDE to help resolve these issues.*\n")
	
	return md.String()
}

// GetRequestFilePath returns the path where the request file was written
func (c *Communicator) GetRequestFilePath() string {
	locations := []string{
		filepath.Join(c.workspacePath, ".cursor", "skene-request.json"),
		filepath.Join(c.workspacePath, ".skene", "ide-request.json"),
		filepath.Join(os.TempDir(), "skene-ide-request.json"),
	}

	for _, location := range locations {
		if _, err := os.Stat(location); err == nil {
			return location
		}
	}

	return ""
}
