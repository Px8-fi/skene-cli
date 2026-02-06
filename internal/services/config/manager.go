package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the skene-growth configuration
type Config struct {
	Provider   string `json:"provider"`
	Model      string `json:"model"`
	APIKey     string `json:"api_key"`
	OutputDir  string `json:"output_dir"`
	Verbose    bool   `json:"verbose"`
	ProjectDir string `json:"project_dir"`
}

// Manager handles configuration file operations
type Manager struct {
	ProjectConfigPath string
	UserConfigPath    string
	Config            *Config
}

// NewManager creates a new config manager
func NewManager(projectDir string) *Manager {
	homeDir, _ := os.UserHomeDir()

	return &Manager{
		ProjectConfigPath: filepath.Join(projectDir, ".skene.config"),
		UserConfigPath:    filepath.Join(homeDir, ".config", "skene", "config"),
		Config: &Config{
			OutputDir: "./skene-context",
			Verbose:   true,
		},
	}
}

// ConfigStatus represents config file status
type ConfigStatus struct {
	Type   string
	Path   string
	Found  bool
	Config *Config
}

// CheckConfigs checks for existing configuration files
func (m *Manager) CheckConfigs() []ConfigStatus {
	statuses := []ConfigStatus{
		{
			Type:  "Project",
			Path:  m.ProjectConfigPath,
			Found: fileExists(m.ProjectConfigPath),
		},
		{
			Type:  "User",
			Path:  m.UserConfigPath,
			Found: fileExists(m.UserConfigPath),
		},
	}

	// Load existing configs
	for i, status := range statuses {
		if status.Found {
			config, err := m.loadConfigFile(status.Path)
			if err == nil {
				statuses[i].Config = config
			}
		}
	}

	return statuses
}

// LoadConfig loads configuration from files (project takes precedence)
func (m *Manager) LoadConfig() error {
	// Try project config first
	if fileExists(m.ProjectConfigPath) {
		config, err := m.loadConfigFile(m.ProjectConfigPath)
		if err == nil {
			m.Config = config
			return nil
		}
	}

	// Fall back to user config
	if fileExists(m.UserConfigPath) {
		config, err := m.loadConfigFile(m.UserConfigPath)
		if err == nil {
			m.Config = config
			return nil
		}
	}

	// No config found, use defaults
	return nil
}

func (m *Manager) loadConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves configuration to project config file
func (m *Manager) SaveConfig() error {
	// Ensure directory exists
	dir := filepath.Dir(m.ProjectConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(m.Config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(m.ProjectConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// SaveUserConfig saves configuration to user config file
func (m *Manager) SaveUserConfig() error {
	// Ensure directory exists
	dir := filepath.Dir(m.UserConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(m.Config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(m.UserConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// SetProvider sets the LLM provider
func (m *Manager) SetProvider(provider string) {
	m.Config.Provider = provider
}

// SetModel sets the model name
func (m *Manager) SetModel(model string) {
	m.Config.Model = model
}

// SetAPIKey sets the API key
func (m *Manager) SetAPIKey(key string) {
	m.Config.APIKey = key
}

// SetProjectDir sets the project directory
func (m *Manager) SetProjectDir(dir string) {
	m.Config.ProjectDir = dir
}

// GetMaskedAPIKey returns masked API key for display
func (m *Manager) GetMaskedAPIKey() string {
	if len(m.Config.APIKey) <= 8 {
		return "****"
	}
	return m.Config.APIKey[:4] + ".." + m.Config.APIKey[len(m.Config.APIKey)-4:]
}

// HasValidConfig checks if config has minimum required values
func (m *Manager) HasValidConfig() bool {
	return m.Config.Provider != "" && m.Config.Model != "" && m.Config.APIKey != ""
}

// GetShortenedPath returns a shortened path for display
func GetShortenedPath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	return "..." + path[len(path)-maxLen+3:]
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Provider represents an LLM provider with its models
type Provider struct {
	ID          string
	Name        string
	Description string
	Models      []Model
	RequiresKey bool
	AuthURL     string // For browser-based auth
}

// Model represents an LLM model
type Model struct {
	ID          string
	Name        string
	Description string
}

// GetProviders returns all available providers
func GetProviders() []Provider {
	return []Provider{
		{
			ID:          "skene",
			Name:        "Skene",
			Description: "LLM & Growth context",
			RequiresKey: true,
			AuthURL:     "https://www.skene.ai/login?=retrieve-api-key",
			Models: []Model{
				{ID: "skene-growth-v1", Name: "skene-growth-v1", Description: "Growth analysis model"},
			},
		},
		{
			ID:          "openai",
			Name:        "openai",
			Description: "LLM Provider",
			RequiresKey: true,
			Models: []Model{
				{ID: "gpt-4o", Name: "gpt-4o", Description: "Most capable GPT-4 model"},
				{ID: "gpt-4-turbo", Name: "gpt-4-turbo", Description: "Fast GPT-4 model"},
				{ID: "gpt-3.5-turbo", Name: "gpt-3.5-turbo", Description: "Fast and affordable"},
			},
		},
		{
			ID:          "gemini",
			Name:        "gemini",
			Description: "LLM Provider",
			RequiresKey: true,
			Models: []Model{
				{ID: "gemini-3-flash-preview", Name: "gemini-3-flash-preview", Description: "Fast Gemini model"},
				{ID: "gemini-3-pro-preview", Name: "gemini-3-pro-preview", Description: "Advanced Gemini model"},
				{ID: "gemini-2.5-flash", Name: "gemini-2.5-flash", Description: "Balanced performance"},
			},
		},
		{
			ID:          "anthropic",
			Name:        "anthropic",
			Description: "LLM Provider",
			RequiresKey: true,
			Models: []Model{
				{ID: "claude-3-opus", Name: "claude-3-opus", Description: "Most capable Claude"},
				{ID: "claude-3-sonnet", Name: "claude-3-sonnet", Description: "Balanced performance"},
				{ID: "claude-3-haiku", Name: "claude-3-haiku", Description: "Fast and efficient"},
			},
		},
		{
			ID:          "meta",
			Name:        "meta",
			Description: "LLM Provider",
			RequiresKey: true,
			Models: []Model{
				{ID: "llama-3-70b", Name: "llama-3-70b", Description: "Large Llama model"},
				{ID: "llama-3-8b", Name: "llama-3-8b", Description: "Efficient Llama model"},
			},
		},
		{
			ID:          "cohere",
			Name:        "cohere",
			Description: "LLM Provider",
			RequiresKey: true,
			Models: []Model{
				{ID: "command-r-plus", Name: "command-r-plus", Description: "Enterprise model"},
				{ID: "command-r", Name: "command-r", Description: "RAG optimized"},
			},
		},
		{
			ID:          "huggingface",
			Name:        "huggingface",
			Description: "LLM Provider",
			RequiresKey: true,
			Models: []Model{
				{ID: "mistral-7b", Name: "mistral-7b", Description: "Open source model"},
				{ID: "mixtral-8x7b", Name: "mixtral-8x7b", Description: "MoE architecture"},
			},
		},
	}
}

// GetProviderByID returns a provider by ID
func GetProviderByID(id string) *Provider {
	providers := GetProviders()
	for _, p := range providers {
		if p.ID == id {
			return &p
		}
	}
	return nil
}
