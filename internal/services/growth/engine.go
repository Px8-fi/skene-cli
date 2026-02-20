package growth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"skene/internal/constants"
	"skene/internal/services/uvresolver"
)

// AnalysisPhase represents a phase of the analysis
type AnalysisPhase int

const (
	PhaseScanCodebase AnalysisPhase = iota
	PhaseDetectFeatures
	PhaseGrowthLoops
	PhaseMonetisation
	PhaseOpportunities
	PhaseGenerateDocs
)

// PhaseUpdate is sent during analysis to update progress
type PhaseUpdate struct {
	Phase    AnalysisPhase
	Progress float64
	Message  string
}

// AnalysisResult holds the complete analysis output
type AnalysisResult struct {
	GrowthPlan     string
	Manifest       string
	GrowthTemplate string
	Error          error
}

// EngineConfig holds the configuration passed to uvx commands
type EngineConfig struct {
	Provider   string
	Model      string
	APIKey     string
	BaseURL    string
	ProjectDir string
	OutputDir  string
	UseGrowth bool
}

// Engine spawns uvx commands to run Skene libraries in the selected repository
type Engine struct {
	config   EngineConfig
	updateFn func(PhaseUpdate)
}

// NewEngine creates a new engine that delegates to uvx
func NewEngine(config EngineConfig, updateFn func(PhaseUpdate)) *Engine {
	return &Engine{
		config:   config,
		updateFn: updateFn,
	}
}

// Run executes the analysis by spawning uvx skene-growth analyze
func (e *Engine) Run(ctx context.Context) *AnalysisResult {
	result := &AnalysisResult{}

	e.sendUpdate(PhaseScanCodebase, 0.0, "Starting analysis via uvx skene-growth...")

	args := []string{constants.GrowthPackageName, "analyze", "."}
	args = append(args, e.buildCommonFlags()...)

	if err := e.runUVX(ctx, args); err != nil {
		result.Error = fmt.Errorf("analysis failed: %w", err)
		return result
	}

	e.sendUpdate(PhaseGenerateDocs, 1.0, "Analysis complete")

	outputDir := e.resolveOutputDir()
	result.GrowthPlan = loadFileContent(filepath.Join(outputDir, constants.GrowthPlanFile))
	result.Manifest = loadFileContent(filepath.Join(outputDir, constants.GrowthManifestFile))
	result.GrowthTemplate = loadFileContent(filepath.Join(outputDir, constants.GrowthTemplateFile))

	return result
}

// GeneratePlan spawns uvx skene-growth plan
func (e *Engine) GeneratePlan() *AnalysisResult {
	result := &AnalysisResult{}

	args := []string{constants.GrowthPackageName, "plan"}
	args = append(args, e.buildCommonFlags()...)

	if err := e.runUVX(context.Background(), args); err != nil {
		result.Error = fmt.Errorf("plan generation failed: %w", err)
		return result
	}

	outputDir := e.resolveOutputDir()
	result.GrowthPlan = loadFileContent(filepath.Join(outputDir, constants.GrowthPlanFile))
	return result
}

// GenerateBuild spawns uvx skene-growth build
func (e *Engine) GenerateBuild() *AnalysisResult {
	result := &AnalysisResult{}

	args := []string{constants.GrowthPackageName, "build"}
	args = append(args, e.buildCommonFlags()...)

	if err := e.runUVX(context.Background(), args); err != nil {
		result.Error = fmt.Errorf("build generation failed: %w", err)
		return result
	}

	outputDir := e.resolveOutputDir()
	result.GrowthPlan = loadFileContent(filepath.Join(outputDir, constants.ImplementationPromptFile))
	return result
}

// ValidateManifest spawns uvx skene-growth validate
func (e *Engine) ValidateManifest() *AnalysisResult {
	result := &AnalysisResult{}

	manifestPath := filepath.Join(e.resolveOutputDir(), constants.GrowthManifestFile)
	args := []string{constants.GrowthPackageName, "validate", manifestPath}

	if err := e.runUVX(context.Background(), args); err != nil {
		result.Error = fmt.Errorf("validation failed: %w", err)
		return result
	}

	return result
}

// runUVX spawns a uvx command in the project directory and streams output.
// It auto-provisions uv if not already installed.
func (e *Engine) runUVX(ctx context.Context, args []string) error {
	uvxPath, err := uvresolver.Resolve()
	if err != nil {
		return fmt.Errorf("failed to locate uvx: %w", err)
	}

	cmd := exec.CommandContext(ctx, uvxPath, args...)
	cmd.Dir = e.config.ProjectDir

	cmd.Env = append(os.Environ(), e.buildEnvVars()...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start uvx: %w", err)
	}

	var lastLines []string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		e.sendUpdate(PhaseDetectFeatures, 0.5, line)
		lastLines = append(lastLines, line)
		if len(lastLines) > 10 {
			lastLines = lastLines[1:]
		}
	}

	if err := cmd.Wait(); err != nil {
		tail := strings.Join(lastLines, "\n")
		if tail != "" {
			return fmt.Errorf("uvx command failed:\n%s", tail)
		}
		return fmt.Errorf("uvx command failed: %w", err)
	}

	return nil
}

func (e *Engine) buildCommonFlags() []string {
	var flags []string
	if e.config.Provider != "" {
		flags = append(flags, "--provider", e.config.Provider)
	}
	if e.config.Model != "" {
		flags = append(flags, "--model", e.config.Model)
	}
	if e.config.APIKey != "" {
		flags = append(flags, "--api-key", e.config.APIKey)
	}
	if e.config.BaseURL != "" {
		flags = append(flags, "--base-url", e.config.BaseURL)
	}
	return flags
}

func (e *Engine) buildEnvVars() []string {
	var envs []string
	if e.config.APIKey != "" {
		envs = append(envs, "SKENE_API_KEY="+e.config.APIKey)
	}
	if e.config.Provider != "" {
		envs = append(envs, "SKENE_PROVIDER="+e.config.Provider)
	}
	if e.config.Model != "" {
		envs = append(envs, "SKENE_MODEL="+e.config.Model)
	}
	if e.config.BaseURL != "" {
		envs = append(envs, "SKENE_BASE_URL="+e.config.BaseURL)
	}
	return envs
}

func (e *Engine) resolveOutputDir() string {
	if e.config.OutputDir != "" {
		if filepath.IsAbs(e.config.OutputDir) {
			return e.config.OutputDir
		}
		return filepath.Join(e.config.ProjectDir, e.config.OutputDir)
	}
	return filepath.Join(e.config.ProjectDir, constants.OutputDirName)
}

func (e *Engine) sendUpdate(phase AnalysisPhase, progress float64, message string) {
	if e.updateFn != nil {
		e.updateFn(PhaseUpdate{
			Phase:    phase,
			Progress: progress,
			Message:  message,
		})
	}
}

func loadFileContent(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
