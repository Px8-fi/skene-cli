package growth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// RustEngine wraps the Rust skene-engine binary
type RustEngine struct {
	config   EngineConfig
	updateFn func(PhaseUpdate)
	binPath  string
}

// NewRustEngine creates a new Rust-based analysis engine
func NewRustEngine(config EngineConfig, updateFn func(PhaseUpdate)) (*RustEngine, error) {
	// Locate the binary (in order of preference):
	// 1. SKENE_ENGINE_PATH env var
	// 2. Same directory as the skene executable (for build/ and installed binaries)
	// 3. Current directory and build/
	// 4. Dev path: engine/target/release/
	// 5. PATH
	binName := "skene-engine"
	if os.Getenv("GOOS") == "windows" {
		binName += ".exe"
	}

	cwd, _ := os.Getwd()
	var candidates []string

	// 1. Explicit path from env
	if p := os.Getenv("SKENE_ENGINE_PATH"); p != "" {
		candidates = append(candidates, p)
	}

	// 2. Same directory as skene binary (most reliable)
	if execPath, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(execPath), binName))
	}

	// 3. Current directory and build/
	candidates = append(candidates,
		filepath.Join(cwd, binName),
		filepath.Join(cwd, "build", binName),
	)

	// 4. Dev path
	candidates = append(candidates,
		filepath.Join(cwd, "engine", "target", "release", binName),
	)

	// 5. PATH lookup
	if path, err := exec.LookPath(binName); err == nil {
		candidates = append(candidates, path)
	}

	var binPath string
	for _, p := range candidates {
		if p == "" {
			continue
		}
		// Resolve to absolute path so exec always gets a full path (avoids "not found in $PATH")
		abs, err := filepath.Abs(p)
		if err != nil {
			continue
		}
		if _, err := os.Stat(abs); err == nil {
			binPath = abs
			break
		}
	}

	if binPath == "" {
		return nil, fmt.Errorf("skene-engine binary not found: place it next to the skene executable or set SKENE_ENGINE_PATH")
	}

	return &RustEngine{
		config:   config,
		updateFn: updateFn,
		binPath:  binPath,
	}, nil
}

// EngineInput matches the Rust protocol Input struct
type EngineInput struct {
	Command        string   `json:"command"`
	Provider       string   `json:"provider"`
	Model          string   `json:"model"`
	APIKey         string   `json:"api_key"`
	BaseURL        *string  `json:"base_url,omitempty"`
	ProjectDir     string   `json:"project_dir"`
	OutputDir      string   `json:"output_dir"`
	ProductDocs    bool     `json:"product_docs"`
	ExcludeFolders []string `json:"exclude_folders"`
	Debug          bool     `json:"debug"`
	ManifestPath   *string  `json:"manifest_path,omitempty"`
	TemplatePath   *string  `json:"template_path,omitempty"`
	Onboarding     *bool    `json:"onboarding,omitempty"`
}

// EngineOutput matches the Rust protocol Output enum structure
type EngineOutput struct {
	Type         string  `json:"type"`
	Phase        string  `json:"phase,omitempty"`
	Step         int     `json:"step,omitempty"`
	TotalSteps   int     `json:"total_steps,omitempty"`
	Progress     float64 `json:"progress,omitempty"`
	Message      string  `json:"message,omitempty"`
	ManifestPath *string `json:"manifest_path,omitempty"`
	TemplatePath *string `json:"template_path,omitempty"`
	DocsPath     *string `json:"docs_path,omitempty"`
	PlanPath     *string `json:"plan_path,omitempty"`
	Code         *string `json:"code,omitempty"`
}

// Run executes the analysis using the Rust binary
func (e *RustEngine) Run() *AnalysisResult {
	result := &AnalysisResult{}
	
	input := EngineInput{
		Command:        "analyze",
		Provider:       e.config.Provider,
		Model:          e.config.Model,
		APIKey:         e.config.APIKey,
		ProjectDir:     e.config.ProjectDir,
		OutputDir:      e.config.OutputDir,
		ProductDocs:    false, // Can be exposed in config if needed
		ExcludeFolders: []string{}, // Can be exposed
		Debug:          e.config.Verbose,
	}
	
	if e.config.BaseURL != "" {
		input.BaseURL = &e.config.BaseURL
	}

	if err := e.execute(input, result); err != nil {
		result.Error = err
	}
	
	return result
}

// GeneratePlan generates a growth plan
func (e *RustEngine) GeneratePlan(manifestPath string, onboarding bool) *AnalysisResult {
	result := &AnalysisResult{}
	
	input := EngineInput{
		Command:      "plan",
		Provider:     e.config.Provider,
		Model:        e.config.Model,
		APIKey:       e.config.APIKey,
		ProjectDir:   e.config.ProjectDir,
		OutputDir:    e.config.OutputDir,
		ManifestPath: &manifestPath,
		Onboarding:   &onboarding,
		Debug:        e.config.Verbose,
	}
	
	if e.config.BaseURL != "" {
		input.BaseURL = &e.config.BaseURL
	}

	if err := e.execute(input, result); err != nil {
		result.Error = err
	}
	
	return result
}

// GenerateBuild generates an implementation prompt from the manifest
func (e *RustEngine) GenerateBuild() *AnalysisResult {
	result := &AnalysisResult{}

	manifestPath := filepath.Join(e.config.OutputDir, "growth-manifest.json")

	input := EngineInput{
		Command:      "build",
		Provider:     e.config.Provider,
		Model:        e.config.Model,
		APIKey:       e.config.APIKey,
		ProjectDir:   e.config.ProjectDir,
		OutputDir:    e.config.OutputDir,
		ManifestPath: &manifestPath,
		Debug:        e.config.Verbose,
	}

	if e.config.BaseURL != "" {
		input.BaseURL = &e.config.BaseURL
	}

	if err := e.execute(input, result); err != nil {
		result.Error = err
	}

	return result
}

// CheckStatus checks growth loop implementation status
func (e *RustEngine) CheckStatus() *AnalysisResult {
	result := &AnalysisResult{}

	manifestPath := filepath.Join(e.config.OutputDir, "growth-manifest.json")

	input := EngineInput{
		Command:      "status",
		Provider:     e.config.Provider,
		Model:        e.config.Model,
		APIKey:       e.config.APIKey,
		ProjectDir:   e.config.ProjectDir,
		OutputDir:    e.config.OutputDir,
		ManifestPath: &manifestPath,
		Debug:        e.config.Verbose,
	}

	if e.config.BaseURL != "" {
		input.BaseURL = &e.config.BaseURL
	}

	if err := e.execute(input, result); err != nil {
		result.Error = err
	}

	return result
}

func (e *RustEngine) execute(input EngineInput, result *AnalysisResult) error {
	cmd := exec.Command(e.binPath)
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start skene-engine: %w", err)
	}
	
	// Send input
	go func() {
		defer stdin.Close()
		json.NewEncoder(stdin).Encode(input)
	}()
	
	// Read stdout
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		var output EngineOutput
		if err := json.Unmarshal([]byte(line), &output); err != nil {
			// Not JSON, maybe log?
			continue
		}
		
		switch output.Type {
		case "progress":
			// Map Rust phases to Go AnalysisPhase
			var phase AnalysisPhase
			switch output.Phase {
			case "tech_stack":
				phase = PhaseScanCodebase
			case "growth_features":
				phase = PhaseDetectFeatures
			case "revenue_leakage":
				phase = PhaseGrowthLoops // Mapping roughly
			case "industry":
				phase = PhaseMonetisation
			case "manifest":
				phase = PhaseGenerateDocs
			default:
				phase = PhaseOpportunities
			}
			
			if e.updateFn != nil {
				e.updateFn(PhaseUpdate{
					Phase:    phase,
					Progress: output.Progress,
					Message:  output.Message,
				})
			}
			
		case "result":
			// Read result files if needed, but Rust already wrote them.
			// Just populate result struct paths/content if needed.
			if output.ManifestPath != nil {
				if content, err := os.ReadFile(*output.ManifestPath); err == nil {
					result.Manifest = string(content)
				}
			}
			if output.PlanPath != nil {
				if content, err := os.ReadFile(*output.PlanPath); err == nil {
					result.GrowthPlan = string(content)
				}
			}
			if output.DocsPath != nil {
				if content, err := os.ReadFile(*output.DocsPath); err == nil {
					result.ProductDocs = string(content)
				}
			}
			
		case "error":
			return fmt.Errorf("%s", output.Message)
		}
	}
	
	if err := cmd.Wait(); err != nil {
		// Read stderr for details
		errBytes, _ := io.ReadAll(stderr)
		return fmt.Errorf("skene-engine failed: %w\n%s", err, string(errBytes))
	}
	
	return nil
}
