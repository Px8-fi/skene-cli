package installer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"skene-terminal-v2/internal/services/analyzer"
	"skene-terminal-v2/internal/services/config"
)

// TaskStatus represents the status of an installation task
type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
)

// Task represents an installation task
type Task struct {
	ID       string
	Name     string
	Status   TaskStatus
	Progress float64
	Error    error
	StartAt  time.Time
	EndAt    time.Time
}

// Installer handles the installation process
type Installer struct {
	projectPath string
	analysis    *analyzer.AnalysisResult
	config      *config.Config
	tasks       []*Task
	currentTask int
	logs        []string
}

// NewInstaller creates a new installer
func NewInstaller(projectPath string, analysis *analyzer.AnalysisResult, cfg *config.Config) *Installer {
	return &Installer{
		projectPath: projectPath,
		analysis:    analysis,
		config:      cfg,
		tasks:       make([]*Task, 0),
		logs:        make([]string, 0),
	}
}

// SetupTasks initializes the task list based on analysis
func (i *Installer) SetupTasks() {
	i.tasks = []*Task{
		{ID: "detect", Name: "Detecting project type", Status: TaskPending},
		{ID: "config", Name: "Generating configuration", Status: TaskPending},
		{ID: "manifest", Name: "Creating manifest files", Status: TaskPending},
	}
}

func (i *Installer) removeTask(id string) {
	var filtered []*Task
	for _, t := range i.tasks {
		if t.ID != id {
			filtered = append(filtered, t)
		}
	}
	i.tasks = filtered
}

// GetTasks returns all tasks
func (i *Installer) GetTasks() []*Task {
	return i.tasks
}

// GetCurrentTask returns the current task being executed
func (i *Installer) GetCurrentTask() *Task {
	if i.currentTask >= 0 && i.currentTask < len(i.tasks) {
		return i.tasks[i.currentTask]
	}
	return nil
}

// GetProgress returns overall progress (0.0 to 1.0)
func (i *Installer) GetProgress() float64 {
	if len(i.tasks) == 0 {
		return 0
	}

	completed := 0
	for _, t := range i.tasks {
		if t.Status == TaskCompleted {
			completed++
		}
	}

	return float64(completed) / float64(len(i.tasks))
}

// IsComplete returns true if all tasks are done
func (i *Installer) IsComplete() bool {
	for _, t := range i.tasks {
		if t.Status != TaskCompleted && t.Status != TaskFailed {
			return false
		}
	}
	return true
}

// HasErrors returns true if any task failed
func (i *Installer) HasErrors() bool {
	for _, t := range i.tasks {
		if t.Status == TaskFailed {
			return true
		}
	}
	return false
}

// GetLogs returns installation logs
func (i *Installer) GetLogs() []string {
	return i.logs
}

func (i *Installer) log(msg string) {
	i.logs = append(i.logs, fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), msg))
}

// RunTask executes a single task
func (i *Installer) RunTask(ctx context.Context, taskID string) error {
	var task *Task
	for idx, t := range i.tasks {
		if t.ID == taskID {
			task = t
			i.currentTask = idx
			break
		}
	}

	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.Status = TaskRunning
	task.StartAt = time.Now()
	i.log(fmt.Sprintf("Starting: %s", task.Name))

	var err error
	switch taskID {
	case "detect":
		err = i.runDetect(ctx, task)
	case "config":
		err = i.runConfigGen(ctx, task)
	case "manifest":
		err = i.runManifestGen(ctx, task)
	}

	task.EndAt = time.Now()

	if err != nil {
		task.Status = TaskFailed
		task.Error = err
		i.log(fmt.Sprintf("Failed: %s - %v", task.Name, err))
		return err
	}

	task.Status = TaskCompleted
	task.Progress = 1.0
	i.log(fmt.Sprintf("Completed: %s", task.Name))
	return nil
}

// RunAll executes all tasks sequentially
func (i *Installer) RunAll(ctx context.Context, progressChan chan<- float64) error {
	for _, task := range i.tasks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := i.RunTask(ctx, task.ID); err != nil {
				return err
			}
			if progressChan != nil {
				progressChan <- i.GetProgress()
			}
		}
	}
	return nil
}

// Task implementations
func (i *Installer) runDetect(ctx context.Context, task *Task) error {
	// Simulate detection (actual detection already done by analyzer)
	for p := 0.0; p <= 1.0; p += 0.2 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			task.Progress = p
			time.Sleep(100 * time.Millisecond)
		}
	}

	i.log(fmt.Sprintf("Detected project type: %s", i.analysis.ProjectType))
	i.log(fmt.Sprintf("Install method: %s", i.analysis.InstallMethod))
	return nil
}

func (i *Installer) runConfigGen(ctx context.Context, task *Task) error {
	task.Progress = 0.3

	// Create config directory
	configDir := filepath.Join(i.projectPath, ".skene")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	task.Progress = 0.6

	// Write main config file
	configPath := filepath.Join(i.projectPath, ".skene.config")
	configContent := fmt.Sprintf(`{
  "provider": "%s",
  "model": "%s",
  "output_dir": "./skene-context",
  "verbose": true
}`, i.config.Provider, i.config.Model)

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	i.log(fmt.Sprintf("Created config: %s", configPath))
	task.Progress = 1.0
	return nil
}

func (i *Installer) runManifestGen(ctx context.Context, task *Task) error {
	task.Progress = 0.3

	// Create output directory
	outputDir := filepath.Join(i.projectPath, "skene-context")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	task.Progress = 0.5

	// Create manifest file
	manifestPath := filepath.Join(outputDir, "manifest.json")
	manifestContent := fmt.Sprintf(`{
  "version": "1.0.0",
  "generated_at": "%s",
  "provider": "%s",
  "model": "%s",
  "files": []
}`, time.Now().Format(time.RFC3339), i.config.Provider, i.config.Model)

	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	i.log(fmt.Sprintf("Created manifest: %s", manifestPath))

	task.Progress = 0.8

	// Create growth plan template
	growthPlanPath := filepath.Join(outputDir, "growth-plan.md")
	growthPlanContent := `# Growth Plan

## Overview
This document outlines the growth strategy generated by Skene.

## Objectives
- [ ] Define key metrics
- [ ] Identify growth opportunities
- [ ] Implement tracking

## Next Steps
Run ` + "`skene analyze`" + ` to generate detailed recommendations.
`

	if err := os.WriteFile(growthPlanPath, []byte(growthPlanContent), 0644); err != nil {
		return fmt.Errorf("failed to write growth plan: %w", err)
	}

	i.log(fmt.Sprintf("Created growth plan: %s", growthPlanPath))
	task.Progress = 1.0
	return nil
}

func (i *Installer) executeCommand(ctx context.Context, cmdStr string) error {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Dir = i.projectPath
	cmd.Stdout = nil // Suppress output
	cmd.Stderr = nil

	return cmd.Run()
}

// InstallError represents a detailed installation error
type InstallError struct {
	Code       string
	Message    string
	Suggestion string
	Retryable  bool
}

func (e *InstallError) Error() string {
	return e.Message
}

// Common installation errors
var (
	ErrPythonNotFound = &InstallError{
		Code:       "PYTHON_NOT_FOUND",
		Message:    "Python is not installed",
		Suggestion: "Install Python 3.8+ from python.org",
		Retryable:  false,
	}

	ErrPipFailed = &InstallError{
		Code:       "PIP_FAILED",
		Message:    "pip installation failed",
		Suggestion: "Try running 'pip install --upgrade pip' first",
		Retryable:  true,
	}

	ErrPermissionDenied = &InstallError{
		Code:       "PERMISSION_DENIED",
		Message:    "Permission denied",
		Suggestion: "Check write permissions or try with sudo",
		Retryable:  true,
	}

	ErrNetworkFailed = &InstallError{
		Code:       "NETWORK_FAILED",
		Message:    "Network connection failed",
		Suggestion: "Check your internet connection and try again",
		Retryable:  true,
	}
)
