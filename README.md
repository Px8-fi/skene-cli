# Skene Terminal v2

A beautiful, production-ready terminal installer and configuration tool for [skene-growth](https://github.com/SkeneTechnologies/skene-growth). Built with Go using the Bubble Tea framework.

<p align="center">
  <img src="designs/Step%201.png" width="400" alt="Welcome Screen" />
  <img src="designs/Step%203.png" width="400" alt="Provider Selection" />
</p>

## Features

- 🧙 **Wizard-Guided Flow** - Step-by-step installation and configuration wizard
- 🎨 **Beautiful Terminal UI** - Rich, animated interface with retro terminal aesthetics
- ✅ **System Checks** - Automatic verification of Python 3.11+ and uv runtime
- 🔄 **Installation Methods** - Choose between quick run (uvx) or full install (pip)
- 🔐 **Multiple Auth Flows** - Skene magic link, API key entry, and local model detection
- 🤖 **AI Provider Support** - OpenAI, Anthropic, Gemini, Skene, Ollama, LM Studio, and generic endpoints
- 📊 **Progress Tracking** - Animated progress bars with elapsed time and phase tracking
- 🎮 **Easter Egg Game** - Space shooter mini-game during loading
- 📖 **Tabbed Results Dashboard** - View growth plans, manifests, and product docs
- ⌨️ **Fully Keyboard Navigable** - No mouse required
- 🌐 **Cross-platform** - Works on macOS, Linux, and Windows
- 🛡️ **Error Handling** - Robust error handling with retry mechanisms and detailed error messages

## Prerequisites

- Go 1.22 or later (for building from source)
- Python 3.11+ (required for skene-growth)
- `uv` runtime (will be installed automatically if missing)

## Installation

### Quick Install (Recommended)

Install `skene` with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/Px8-fi/skene-cli/Rust-impelementation/install.sh | bash
```

This will:
- Automatically detect your platform (macOS, Linux, or Windows)
- Download the appropriate binary from the latest GitHub release (if available)
- Or build from source if releases aren't available
- Install to `/usr/local/bin/skene`
- Make it executable and verify installation

### Alternative: Clone and Install

If you prefer to clone the repository first:

```bash
# Clone the repository
git clone https://github.com/Px8-fi/skene-cli
cd skene-cli

# Run the install script
./install.sh
```

The install script will automatically detect you're in the repository directory and use the local `build/skene` binary if available.

This will:
- Automatically detect your platform (macOS, Linux, or Windows)
- Download the appropriate binary from the latest GitHub release (if public)
- Or build from source if releases aren't available
- Install to `/usr/local/bin/skene`
- Make it executable and verify installation

### Manual Installation

If you prefer to download the script first:

```bash
# Download the install script
curl -fsSL https://raw.githubusercontent.com/Px8-fi/skene-cli/Rust-impelementation/install.sh -o install.sh

# Make it executable
chmod +x install.sh

# Run the installer
./install.sh
```

### Custom Installation Location

Install to a custom directory (no sudo required):

```bash
INSTALL_DIR=~/bin ./install.sh
```

Or set it as an environment variable:

```bash
export INSTALL_DIR=~/bin
./install.sh
```

### Install Specific Version

Install a specific release version:

```bash
VERSION=v1.0.0 curl -fsSL https://raw.githubusercontent.com/Px8-fi/skene-cli/Rust-impelementation/install.sh | bash
```

### Local Development Installation

If you're working on the project locally:

```bash
# Clone the repository
git clone https://github.com/Px8-fi/skene-cli
cd skene-cli

# Build the project
make build

# Install using the script (uses local build)
./install.sh

# Or force local build
USE_LOCAL=true ./install.sh
```

The install script will automatically detect you're in the repository directory and use the local `build/skene` binary if available.

### Build from Source

For developers who want to build from source:

```bash
# Clone the repository
git clone https://github.com/Px8-fi/skene-cli
cd skene-cli

# Install Go dependencies
make install

# Build for your platform
make build

# Or build for all platforms
make build-all

# Install manually
sudo cp build/skene /usr/local/bin/
```

### Verify Installation

After installation, verify it works:

```bash
skene --version
# or
skene --help
```

If the command is not found, make sure `/usr/local/bin` is in your PATH:

```bash
# Add to ~/.zshrc (macOS/Linux with zsh)
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Add to ~/.bashrc (Linux with bash)
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

## Quick Start

### Using Make (Recommended)

```bash
# Clone the repository
git clone https://github.com/Px8-fi/skene-cli
cd skene-cli

# Install dependencies and build
make install
make build

# Run the application
make run
```

### Using Go directly

```bash
# Install dependencies
go mod download
go mod tidy

# Run directly
go run ./cmd/skene

# Or build and run
go build -o skene ./cmd/skene
./skene
```

## Usage

### Wizard Flow

The application guides you through a complete installation and analysis workflow:

1. **Welcome Screen** - Animated welcome with ASCII art
2. **System Checks** - Verifies Python 3.11+ and uv runtime (installs uv if needed)
3. **Install Method** - Choose between:
   - `uvx` (quick run, ephemeral environment)
   - `pip` (full installation)
4. **Installing** - Progress tracking for skene-growth installation
5. **AI Provider Selection** - Choose from:
   - Skene (with magic link authentication)
   - OpenAI
   - Anthropic (Claude)
   - Gemini
   - Local models (Ollama, LM Studio)
   - Other OpenAI-compatible APIs
6. **Model Selection** - Select the specific model for your provider
7. **Authentication**:
   - **Skene**: Magic link flow with browser redirect (falls back to manual API key)
   - **Other providers**: Manual API key entry with validation
   - **Local models**: Automatic detection and model selection
8. **Project Directory** - Select or enter the project directory to analyze
9. **Analysis Configuration** - Configure analysis settings or use recommended defaults
10. **Analyzing** - Multi-phase analysis progress (scanning, feature detection, growth analysis)
11. **Results Dashboard** - Tabbed view of:
    - Growth Plan
    - Growth Manifest
    - Product Documentation
12. **Next Steps** - Choose to generate roadmap, validate manifest, re-run analysis, or exit

### Keyboard Controls

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate up/down |
| `←/→` or `h/l` | Navigate left/right |
| `Enter` | Confirm/Select |
| `Esc` | Go back to previous step |
| `Tab` | Switch focus area (in multi-input views) |
| `Space` | Toggle checkboxes/options |
| `?` | Toggle help overlay |
| `L` | Toggle error logs (in error view) |
| `q` | Quit |
| `g` | Play mini-game (during loading) |

### Mini-Game Controls (Space Shooter)

| Key | Action |
|-----|--------|
| `←/→` or `a/d` | Move ship |
| `Space` | Shoot |
| `p` | Pause |
| `r` | Restart (when game over) |
| `Esc` | Exit game |

## Project Structure

```
skene-terminal-v2/
├── cmd/
│   └── skene/
│       └── main.go                    # Application entry point
├── internal/
│   ├── tui/
│   │   ├── app.go                     # Main wizard state machine
│   │   ├── styles/
│   │   │   └── styles.go              # Lip Gloss styling system
│   │   ├── components/
│   │   │   ├── button.go              # Button components
│   │   │   ├── button_group.go        # Button group component
│   │   │   ├── help.go                # Help overlay
│   │   │   ├── logo.go                # ASCII logo animations
│   │   │   ├── progress.go            # Progress bars
│   │   │   ├── spinner.go             # Loading spinners
│   │   │   └── wizard_header.go       # Wizard step progress header
│   │   └── views/
│   │       ├── welcome.go              # Welcome screen
│   │       ├── syscheck.go            # System checks view
│   │       ├── install_method.go      # Install method selection
│   │       ├── installing.go          # Installation progress
│   │       ├── provider.go            # Provider selection
│   │       ├── model.go               # Model selection
│   │       ├── auth.go                # Skene magic link auth
│   │       ├── apikey.go              # API key entry
│   │       ├── local_model.go         # Local model detection
│   │       ├── project_dir.go         # Project directory selection
│   │       ├── analysis_config.go     # Analysis configuration
│   │       ├── analyzing.go           # Analysis progress
│   │       ├── results.go             # Results dashboard
│   │       ├── next_steps.go          # Next steps menu
│   │       └── error.go               # Error handling
│   ├── services/
│   │   ├── analyzer/
│   │   │   └── analyzer.go            # Project analysis
│   │   ├── config/
│   │   │   └── manager.go             # Configuration management
│   │   ├── installer/
│   │   │   └── installer.go           # Installation engine
│   │   └── syscheck/
│   │       └── checker.go             # System prerequisite checks
│   └── game/
│       └── shooter.go                 # Space shooter game
├── designs/                            # Design reference images
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Architecture

### State Management

The application uses a finite state machine pattern implementing a complete wizard flow:

```
Welcome → SystemCheck → InstallMethod → Installing → Provider → Model → Auth/APIKey → ProjectDir → AnalysisConfig → Analyzing → Results → NextSteps
                                                                         ↓
                                                                    LocalModel (if local provider)
```

### Wizard States

- **Welcome**: Initial animated welcome screen
- **SystemCheck**: Verifies Python 3.11+ and uv runtime
- **InstallMethod**: Choose between uvx and pip installation
- **Installing**: Track installation progress with task status
- **Provider**: Select AI provider (Skene, OpenAI, Anthropic, Gemini, Local, Generic)
- **Model**: Select model for chosen provider
- **Auth**: Skene magic link authentication flow
- **APIKey**: Manual API key entry with validation
- **LocalModel**: Detect and select local models (Ollama/LM Studio)
- **ProjectDir**: Select project directory for analysis
- **AnalysisConfig**: Configure analysis settings
- **Analyzing**: Multi-phase analysis progress tracking
- **Results**: Tabbed dashboard with growth plan, manifest, and docs
- **NextSteps**: Post-analysis action menu
- **Error**: Error display with retry options

### Views

Each view implements a consistent interface:
- `SetSize(width, height int)` - Handle terminal resize
- `Render() string` - Return the view content
- `GetHelpItems() []HelpItem` - Context-specific help

Views are organized by wizard step, with each step handling its own state, validation, and user interaction. The main `app.go` orchestrates state transitions and delegates to view-specific handlers.

### Styling

Uses Lip Gloss with a warm, retro terminal aesthetic:
- **Primary**: `#EDC29C` (Cream/Amber)
- **Background**: Dark (`#1A1A1A`)
- **Accent**: Subtle gold highlights
- **Typography**: Monospace terminal fonts

## Configuration

Configuration is stored in:
- **Project**: `.skene.config` in the current directory
- **User**: `~/.config/skene/config`

Example configuration:
```json
{
  "provider": "gemini",
  "model": "gemini-3-flash-preview",
  "api_key": "your-api-key",
  "base_url": "https://api.example.com/v1",
  "output_dir": "./skene-context",
  "verbose": true
}
```

### Supported Providers

- **Skene**: `skene` (magic link auth or API key)
- **OpenAI**: `openai` (requires API key)
- **Anthropic**: `anthropic` or `claude` (requires API key)
- **Gemini**: `gemini` (requires API key)
- **Ollama**: `ollama` (local, no API key needed)
- **LM Studio**: `lmstudio` (local, no API key needed)
- **Generic**: `generic` (OpenAI-compatible endpoint, requires `base_url`)

For generic providers, set `base_url` to your endpoint URL (e.g., `http://localhost:8000/v1` for local servers).

## Development

### Live Reload

```bash
# Install air for live reload
go install github.com/cosmtrek/air@latest

# Run with live reload
make dev
```

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

### Formatting

```bash
make fmt
```

## Building for Distribution

```bash
# Build for all platforms
make build-all

# Zip files
make release

# Output will be in build/
# - skene-linux-amd64
# - skene-darwin-amd64
# - skene-darwin-arm64
# - skene-windows-amd64.exe
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components (textinput, viewport, etc.)
- [pkg/browser](https://github.com/pkg/browser) - Browser opening for magic link auth

## Error Handling

The wizard includes robust error handling:

- **System Check Failures**: Clear messages with installation instructions
- **Installation Errors**: Retry options and detailed error logs
- **Authentication Failures**: Validation hints and retry mechanisms
- **Analysis Errors**: Phase-specific error messages with recovery options
- **Network Errors**: Retry suggestions and fallback options

All errors display in a dedicated error view with:
- Clear error messages
- Actionable suggestions
- Retry/View Logs/Quit options
- Expandable error logs (press `L`)

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Support

- 📖 [Documentation](https://docs.skene.ai)
- 🐛 [Issues](https://github.com/Px8-fi/skene-cli/issues)
- 💬 [Discord](https://discord.gg/skene)
