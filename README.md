# Skene CLI

A terminal interface for the [Skene Growth](https://github.com/SkeneTechnologies/skene-growth) ecosystem. Guides you through codebase analysis, growth plan generation, and implementation -- all from the terminal. Built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

<p align="center">
  <img src="designs/Step%201.png" width="400" alt="Welcome Screen" />
  <img src="designs/Step%203.png" width="400" alt="Provider Selection" />
</p>

## What It Does

Skene CLI is the interactive front-end for three Skene packages that work together:

| Package | Purpose |
|---------|---------|
| [**Skene Growth**](https://github.com/SkeneTechnologies/skene-growth) | Tech stack detection, growth feature discovery, revenue leakage analysis, growth plan generation |
| [**Skene Skills**](https://github.com/SkeneTechnologies/skene-skills) | PLG analysis skills for Claude Code -- analyze, generate manifests and templates |
| [**Skene Cookbook**](https://github.com/SkeneTechnologies/skene-cookbook) | 700+ AI skills for PLG, marketing, security, DevEx, and more |

By default all three packages are used during analysis. You can choose a subset in the advanced configuration screen.

## Features

- **Step-by-step setup** -- provider, model, authentication, and project selection
- **Multiple AI providers** -- Skene, OpenAI, Anthropic, Gemini, Ollama, LM Studio, or any OpenAI-compatible endpoint
- **Authentication flows** -- Skene magic link, API key entry, and local model auto-detection
- **Skene package selection** -- choose which Skene packages to include (Growth, Skills, Cookbook)
- **Existing analysis detection** -- detects previous `skene-context/` output and offers to view or rerun
- **Multi-phase analysis** -- animated progress with live terminal output
- **Tabbed results dashboard** -- view growth plans, manifests, and product docs
- **Next steps menu** -- run `plan`, `build`, `status`, or re-analyze directly from the CLI
- **Self-contained** -- everything runs through the bundled Rust engine, no Python or uv required
- **Error handling** -- severity-based error display with suggestions, retry, and go back
- **Cross-platform** -- macOS, Linux, and Windows
- **Mini-game** -- space shooter while you wait

## Prerequisites

- Go 1.22+ (for building from source)
- Rust toolchain (for building the analysis engine from source; a local copy is installed automatically by `make build` if not present)

## Installation

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/SkeneTechnologies/skene-cli/main/install.sh | bash
```

This detects your platform, downloads the correct binary (or builds from source), and installs to `/usr/local/bin/skene`.

### Clone and Install

```bash
git clone https://github.com/SkeneTechnologies/skene-cli
cd skene-cli
./install.sh
```

### Build from Source

```bash
git clone https://github.com/SkeneTechnologies/skene-cli
cd skene-cli
make install   # download Go dependencies
make build     # build binary + Rust engine
make run       # run the application
```

### Custom Install Location

```bash
INSTALL_DIR=~/bin ./install.sh
```

### Specific Version

```bash
VERSION=v0.2.0 ./install.sh
```

### Verify

```bash
skene --version
```

If the command is not found, add the install directory to your PATH:

```bash
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

## Usage

Run `skene` and follow the prompts. The flow is:

```
Welcome
  -> AI Provider selection (Skene, OpenAI, Anthropic, Gemini, Ollama, LM Studio, Generic)
    -> Model selection
      -> Authentication (magic link, API key, or local model detection)
        -> Project directory (detects existing analysis if present)
          -> Analysis configuration (default or custom package/option selection)
            -> Running analysis (live progress + terminal output)
              -> Results dashboard (Growth Plan | Manifest | Product Docs)
                -> Next steps
```

### Existing Analysis Detection

When you select a project directory that already contains a `skene-context/` folder, you are given two options:

- **View Analysis** -- loads the existing results into the dashboard
- **Rerun Analysis** -- proceeds to configuration and runs a fresh analysis

### Analysis Configuration

The default configuration runs all three Skene packages. Selecting "No" on the default question opens the advanced screen where you can:

- Toggle individual packages (Skene Growth, Skene Skills, Skene Cookbook)
- Enable/disable product docs generation
- Set business type
- Toggle verbose output

### Next Steps

After analysis completes, the next steps menu offers:

| Action | Description |
|--------|-------------|
| Generate Growth Plan | Prioritised growth plan with implementation roadmap |
| Build Implementation Prompt | Ready-to-use prompt for Cursor, Claude, or other AI tools |
| Check Loop Status | Verify which growth loop requirements are implemented |
| Re-run Analysis | Analyse the codebase again |
| Open Generated Files | Open `./skene-context/` in your file manager |
| Change Configuration | Return to provider selection |
| Exit | Close Skene CLI |

All commands run through the bundled Rust engine -- no external tools required.

### Keyboard Controls

| Key | Action |
|-----|--------|
| `Up/Down` or `j/k` | Navigate |
| `Left/Right` or `h/l` | Navigate / switch tabs |
| `Enter` | Confirm / select |
| `Esc` | Go back |
| `Tab` | Switch focus area |
| `Space` | Toggle checkbox / option |
| `?` | Toggle help overlay |
| `g` | Play mini-game (during loading) |
| `Ctrl+C` | Quit |

## Configuration

Configuration files are checked in order (first found wins):

1. **Project** -- `.skene.config` in the project directory
2. **User** -- `~/.config/skene/config`

Example `.skene.config`:

```json
{
  "provider": "gemini",
  "model": "gemini-3-flash-preview",
  "api_key": "your-api-key",
  "base_url": "",
  "output_dir": "./skene-context",
  "verbose": true,
  "use_growth": true,
  "use_skills": true,
  "use_cookbook": true
}
```

### Supported Providers

| Provider | ID | Auth | Notes |
|----------|----|------|-------|
| Skene | `skene` | Magic link or API key | Recommended -- built-in growth model |
| OpenAI | `openai` | API key | GPT-4o, GPT-4 Turbo, GPT-3.5 Turbo |
| Anthropic | `anthropic` | API key | Claude Sonnet 4.5, Opus, Haiku |
| Gemini | `gemini` | API key | Gemini 3 Flash, 3 Pro, 2.5 Flash |
| Ollama | `ollama` | None (local) | Llama 3.3, Mistral, CodeLlama, DeepSeek R1 |
| LM Studio | `lmstudio` | None (local) | Uses currently loaded model |
| Generic | `generic` | API key + base URL | Any OpenAI-compatible endpoint |

## Project Structure

```
skene-cli/
├── cmd/skene/
│   └── main.go                          # Entry point
├── internal/
│   ├── tui/
│   │   ├── app.go                       # State machine and message handling
│   │   ├── styles/
│   │   │   └── styles.go                # Colour palette and Lip Gloss styles
│   │   ├── components/
│   │   │   ├── button.go                # Button and ButtonGroup
│   │   │   ├── dir_browser.go           # Directory browser
│   │   │   ├── help.go                  # Help overlay and footer
│   │   │   ├── logo.go                  # ASCII logo
│   │   │   ├── ascii_motion_placeholder.go
│   │   │   ├── progress.go              # Progress bar
│   │   │   ├── terminal_output.go       # Scrollable terminal output
│   │   │   └── wizard_header.go         # Step progress header
│   │   └── views/
│   │       ├── welcome.go               # Welcome screen
│   │       ├── syscheck.go              # System prerequisite checks
│   │       ├── install_method.go        # uvx vs pip selection
│   │       ├── installing.go            # Installation progress
│   │       ├── provider.go              # AI provider selection
│   │       ├── model.go                 # Model selection
│   │       ├── auth.go                  # Skene magic link auth
│   │       ├── apikey.go                # API key entry
│   │       ├── local_model.go           # Local model detection
│   │       ├── project_dir.go           # Project directory + existing analysis detection
│   │       ├── analysis_config.go       # Package selection and options
│   │       ├── analyzing.go             # Live analysis progress
│   │       ├── results.go               # Tabbed results dashboard
│   │       ├── next_steps.go            # Post-analysis actions
│   │       └── error.go                 # Error display
│   ├── services/
│   │   ├── analyzer/analyzer.go         # Project analysis
│   │   ├── auth/callback.go             # OAuth callback server
│   │   ├── config/manager.go            # Config file management
│   │   ├── growth/
│   │   │   ├── engine.go                # Analysis engine
│   │   │   └── engine_rust.go           # Rust engine wrapper
│   │   ├── ide/communicator.go          # IDE integration
│   │   ├── installer/installer.go       # Package installer
│   │   ├── llm/client.go               # LLM API client
│   │   └── syscheck/checker.go          # System checks
│   └── game/
│       └── shooter.go                   # Space shooter mini-game
├── engine/                              # Rust analysis engine source
├── designs/                             # Design reference images
├── Makefile
├── go.mod
└── go.sum
```

## Architecture

### Rust Engine

All analysis, planning, building, and status checks run through `skene-engine`, a Rust binary that communicates with the Go TUI via JSON over stdin/stdout. The engine handles:

- **analyze** -- codebase scanning, feature detection, growth loop analysis, manifest generation
- **plan** -- generates a prioritised growth plan from the manifest
- **build** -- generates a copy-paste implementation prompt for AI coding tools
- **status** -- checks which growth features have been implemented

The Go side (`internal/services/growth/engine_rust.go`) locates the binary automatically and streams progress updates back to the UI in real time. No Python, uv, or external package managers are needed at runtime.

### State Machine

The application is a finite state machine driven by Bubble Tea messages:

```
Welcome -> Provider -> Model -> Auth/APIKey -> ProjectDir -> AnalysisConfig -> Analyzing -> Results -> NextSteps
                                    |                |
                                LocalModel      (existing analysis?)
                                                 /          \
                                          View Results    Rerun -> AnalysisConfig
```

Every view in `internal/tui/views/` implements a consistent interface:

```go
SetSize(width, height int)
Render() string
GetHelpItems() []HelpItem
```

The main `app.go` orchestrates state transitions and delegates key handling to view-specific methods.

### Styling

Lip Gloss with a warm, retro terminal palette:

- **Cream/Amber** `#EDC29C` -- primary accent
- **Dark background** `#1A1A1A`
- **Success** `#7CB374`, **Warning** `#E6B450`, **Error** `#F05D5E`

Light terminal backgrounds are detected automatically and colours are adjusted for contrast.

### Error Handling

Errors are categorised by severity (Warning, Error, Critical) and displayed in a consistent view with:

- Severity indicator and error code
- Clear message and suggested fix
- Retry / Go Back / Quit buttons
- Press `Esc` to go back to the previous screen

## Development

```bash
# Live reload (requires air)
make dev

# Run tests
make test

# Lint
make lint

# Format
make fmt

# Build for all platforms
make build-all

# Package releases
make release
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) -- TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) -- Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) -- UI components (textinput, viewport)
- [pkg/browser](https://github.com/pkg/browser) -- Browser opening for magic link auth
- [termenv](https://github.com/muesli/termenv) -- Terminal capability detection

## Related Projects

- [skene-growth](https://github.com/SkeneTechnologies/skene-growth) -- PLG analysis toolkit (CLI + MCP server)
- [skene-skills](https://github.com/SkeneTechnologies/skene-skills) -- Claude Code plugin for PLG analysis
- [skene-cookbook](https://github.com/SkeneTechnologies/skene-cookbook) -- 700+ AI skills for Claude and Cursor

## License

MIT
