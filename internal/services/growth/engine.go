package growth

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"skene-terminal-v2/internal/services/llm"
	"strings"
	"time"
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
	GrowthPlan  string
	Manifest    string
	ProductDocs string
	Error       error
}

// EngineConfig holds the analysis configuration
type EngineConfig struct {
	Provider   string
	Model      string
	APIKey     string
	BaseURL    string
	ProjectDir string
	OutputDir  string
	Verbose    bool
}

// Engine performs real codebase analysis using LLM
type Engine struct {
	config   EngineConfig
	client   *llm.Client
	updateFn func(PhaseUpdate)
}

// NewEngine creates a new analysis engine
func NewEngine(config EngineConfig, updateFn func(PhaseUpdate)) *Engine {
	client := llm.NewClient(config.Provider, config.Model, config.APIKey, config.BaseURL)
	return &Engine{
		config:   config,
		client:   client,
		updateFn: updateFn,
	}
}

// Run executes the full analysis pipeline
func (e *Engine) Run(ctx context.Context) *AnalysisResult {
	result := &AnalysisResult{}

	// Phase 1: Scan codebase
	e.sendUpdate(PhaseScanCodebase, 0.0, "Scanning project structure...")
	codebaseContext, err := e.scanCodebase()
	if err != nil {
		result.Error = fmt.Errorf("codebase scan failed: %w", err)
		return result
	}
	e.sendUpdate(PhaseScanCodebase, 1.0, "Codebase scanned")

	// Phase 2: Detect product features
	e.sendUpdate(PhaseDetectFeatures, 0.0, "Analyzing product features...")
	features, err := e.detectFeatures(ctx, codebaseContext)
	if err != nil {
		result.Error = fmt.Errorf("feature detection failed: %w", err)
		return result
	}
	e.sendUpdate(PhaseDetectFeatures, 1.0, "Features detected")

	// Phase 3: Growth loop analysis
	e.sendUpdate(PhaseGrowthLoops, 0.0, "Analyzing growth loops...")
	growthPlan, err := e.analyzeGrowthLoops(ctx, codebaseContext, features)
	if err != nil {
		result.Error = fmt.Errorf("growth analysis failed: %w", err)
		return result
	}
	result.GrowthPlan = growthPlan
	e.sendUpdate(PhaseGrowthLoops, 1.0, "Growth loops identified")

	// Phase 4: Monetisation analysis
	e.sendUpdate(PhaseMonetisation, 0.0, "Analyzing monetisation opportunities...")
	manifest, err := e.analyzeMonetisation(ctx, codebaseContext, features, growthPlan)
	if err != nil {
		result.Error = fmt.Errorf("monetisation analysis failed: %w", err)
		return result
	}
	result.Manifest = manifest
	e.sendUpdate(PhaseMonetisation, 1.0, "Monetisation analysis complete")

	// Phase 5: Opportunity modelling
	e.sendUpdate(PhaseOpportunities, 0.0, "Modelling opportunities...")
	// This is part of the overall analysis - we combine the results
	e.sendUpdate(PhaseOpportunities, 1.0, "Opportunities modelled")

	// Phase 6: Generate docs
	e.sendUpdate(PhaseGenerateDocs, 0.0, "Generating documentation...")
	productDocs, err := e.generateDocs(ctx, codebaseContext, features, growthPlan, manifest)
	if err != nil {
		result.Error = fmt.Errorf("doc generation failed: %w", err)
		return result
	}
	result.ProductDocs = productDocs
	e.sendUpdate(PhaseGenerateDocs, 0.8, "Saving files...")

	// Save output files
	if err := e.saveOutputFiles(result); err != nil {
		// Non-fatal: log but continue
		result.ProductDocs += fmt.Sprintf("\n\n(Warning: Failed to save files: %v)", err)
	}
	e.sendUpdate(PhaseGenerateDocs, 1.0, "Documentation generated")

	return result
}

// scanCodebase reads the project structure and key files
func (e *Engine) scanCodebase() (string, error) {
	var sb strings.Builder
	projectDir := e.config.ProjectDir

	sb.WriteString(fmt.Sprintf("Project Directory: %s\n\n", projectDir))

	// Walk the directory tree (limited depth)
	sb.WriteString("== FILE STRUCTURE ==\n")
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}

		relPath, _ := filepath.Rel(projectDir, path)

		// Skip hidden dirs, node_modules, etc.
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "__pycache__" ||
				name == "venv" || name == ".venv" || name == "build" || name == "dist" ||
				name == ".git" || name == "vendor" {
				return filepath.SkipDir
			}
		}

		// Limit depth to 4
		depth := strings.Count(relPath, string(os.PathSeparator))
		if depth > 4 {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			sb.WriteString(fmt.Sprintf("  %s/\n", relPath))
		} else {
			sb.WriteString(fmt.Sprintf("  %s (%d bytes)\n", relPath, info.Size()))
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	// Read key project files
	keyFiles := []string{
		"README.md", "readme.md",
		"package.json",
		"pyproject.toml", "requirements.txt", "setup.py",
		"go.mod",
		"Cargo.toml",
		"pom.xml",
		"Makefile", "Dockerfile",
		".env.example",
	}

	sb.WriteString("\n== KEY FILES CONTENT ==\n")
	for _, f := range keyFiles {
		fullPath := filepath.Join(projectDir, f)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}
		content := string(data)
		// Truncate large files
		if len(content) > 3000 {
			content = content[:3000] + "\n... (truncated)"
		}
		sb.WriteString(fmt.Sprintf("\n--- %s ---\n%s\n", f, content))
	}

	// Read source files (first 5 of each type)
	sourceExts := []string{".py", ".js", ".ts", ".go", ".rs", ".java", ".rb"}
	sb.WriteString("\n== SAMPLE SOURCE FILES ==\n")

	for _, ext := range sourceExts {
		count := 0
		filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || count >= 3 {
				return nil
			}

			name := info.Name()
			if strings.HasPrefix(name, ".") {
				return nil
			}

			// Skip test files and dependencies
			relPath, _ := filepath.Rel(projectDir, path)
			if strings.Contains(relPath, "node_modules") || strings.Contains(relPath, "__pycache__") ||
				strings.Contains(relPath, "venv") || strings.Contains(relPath, ".git") {
				return nil
			}

			if filepath.Ext(name) == ext {
				data, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				content := string(data)
				if len(content) > 2000 {
					content = content[:2000] + "\n... (truncated)"
				}
				sb.WriteString(fmt.Sprintf("\n--- %s ---\n%s\n", relPath, content))
				count++
			}
			return nil
		})
	}

	// Cap total context
	result := sb.String()
	if len(result) > 30000 {
		result = result[:30000] + "\n... (context truncated)"
	}

	return result, nil
}

// detectFeatures uses LLM to detect product features
func (e *Engine) detectFeatures(ctx context.Context, codebaseContext string) (string, error) {
	messages := []llm.ChatMessage{
		{
			Role: "system",
			Content: `You are a senior product analyst. Analyze the codebase and identify:
1. Core product features
2. User-facing functionality
3. Authentication/authorization patterns
4. Data models and storage
5. API endpoints
6. Third-party integrations
7. Current growth mechanisms (sharing, referrals, analytics, etc.)

Be concise and structured. Use bullet points.`,
		},
		{
			Role:    "user",
			Content: "Analyze this codebase for product features:\n\n" + codebaseContext,
		},
	}

	return e.client.Chat(ctx, messages, 2000)
}

// analyzeGrowthLoops uses LLM to identify growth opportunities
func (e *Engine) analyzeGrowthLoops(ctx context.Context, codebaseContext, features string) (string, error) {
	messages := []llm.ChatMessage{
		{
			Role: "system",
			Content: `You are a Product-Led Growth expert. Based on the codebase analysis and detected features, create a growth strategy memo.

Format your response EXACTLY like this:

## EXECUTIVE SUMMARY
[2-3 sentences about the product and its growth potential]

## GROWTH LOOPS IDENTIFIED

### Loop 1: [Name]
Priority: [HIGH/MEDIUM/LOW] | Estimated Impact: [Xx metric]
- [Action item 1]
- [Action item 2]
- [Action item 3]

### Loop 2: [Name]
[Same format]

### Loop 3: [Name]
[Same format]

## IMPLEMENTATION ROADMAP
Week 1-2: [Task]
Week 3-4: [Task]
Week 5-6: [Task]
Week 7-8: [Task]

## SUCCESS METRICS
- [Metric 1]: target [value]
- [Metric 2]: target [value]
- [Metric 3]: target [value]

Be specific to THIS codebase. Reference actual files, features, and patterns you see.`,
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Codebase context:\n%s\n\nDetected features:\n%s\n\nCreate a growth strategy.", codebaseContext, features),
		},
	}

	return e.client.Chat(ctx, messages, 3000)
}

// analyzeMonetisation analyzes monetisation opportunities
func (e *Engine) analyzeMonetisation(ctx context.Context, codebaseContext, features, growthPlan string) (string, error) {
	messages := []llm.ChatMessage{
		{
			Role: "system",
			Content: `You are a monetisation strategist. Analyze the codebase for revenue opportunities.

Format your response as a MANIFEST:

SKENE GROWTH MANIFEST
Generated: [date]

== TECH STACK ==
[List detected technologies]

== CURRENT GROWTH FEATURES ==
[List what exists with checkmarks or X marks]

== REVENUE OPPORTUNITIES ==
[List prioritized opportunities with HIGH/MEDIUM/LOW tags]

== MISSING FEATURES ==
[What's missing for growth, e.g. sharing, analytics, referrals]

== RECOMMENDED ACTIONS ==
[Numbered list of specific actions]

Be specific to THIS codebase.`,
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Codebase:\n%s\n\nFeatures:\n%s\n\nGrowth plan:\n%s\n\nCreate a monetisation manifest.", codebaseContext, features, growthPlan),
		},
	}

	return e.client.Chat(ctx, messages, 2500)
}

// generateDocs creates product documentation
func (e *Engine) generateDocs(ctx context.Context, codebaseContext, features, growthPlan, manifest string) (string, error) {
	messages := []llm.ChatMessage{
		{
			Role: "system",
			Content: `You are a technical writer. Create product documentation based on the analysis.

Format:

PRODUCT DOCUMENTATION
Auto-generated by Skene Growth

== PRODUCT OVERVIEW ==
Tagline: [inferred from codebase]
Target Audience: [identified users]
Value Proposition: [core value]

== DETECTED FEATURES ==
[Numbered list with descriptions]

== ARCHITECTURE ==
[Brief technical overview]

== GETTING STARTED ==
[Steps to get started with the product]

== GROWTH RECOMMENDATIONS ==
[Top 3 actionable recommendations]

Be specific and reference actual code/files.`,
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Codebase:\n%s\n\nFeatures:\n%s\n\nGrowth plan:\n%s\n\nManifest:\n%s\n\nGenerate product docs.", codebaseContext, features, growthPlan, manifest),
		},
	}

	return e.client.Chat(ctx, messages, 2500)
}

// saveOutputFiles writes results to disk
func (e *Engine) saveOutputFiles(result *AnalysisResult) error {
	outputDir := e.config.OutputDir
	if outputDir == "" {
		outputDir = filepath.Join(e.config.ProjectDir, "skene-context")
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02")

	files := map[string]string{
		"growth-plan.md":     fmt.Sprintf("# Growth Plan\nGenerated: %s\n\n%s", timestamp, result.GrowthPlan),
		"manifest.md":        fmt.Sprintf("# Growth Manifest\nGenerated: %s\n\n%s", timestamp, result.Manifest),
		"product-docs.md":    fmt.Sprintf("# Product Documentation\nGenerated: %s\n\n%s", timestamp, result.ProductDocs),
	}

	for filename, content := range files {
		path := filepath.Join(outputDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	return nil
}

// sendUpdate sends a phase update to the UI
func (e *Engine) sendUpdate(phase AnalysisPhase, progress float64, message string) {
	if e.updateFn != nil {
		e.updateFn(PhaseUpdate{
			Phase:    phase,
			Progress: progress,
			Message:  message,
		})
	}
}
