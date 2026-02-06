# Skene Terminal v2

A beautiful, production-ready terminal installer and configuration tool for [skene-growth](https://github.com/skene-ai/skene-growth). Built with Go using the Bubble Tea framework.

<p align="center">
  <img src="designs/Step%201.png" width="400" alt="Welcome Screen" />
  <img src="designs/Step%203.png" width="400" alt="Provider Selection" />
</p>

## Features

- 🎨 **Beautiful Terminal UI** - Rich, animated interface with retro terminal aesthetics
- 🔄 **Interactive Configuration** - Step-by-step guided setup
- 🔐 **Secure API Key Entry** - Masked input with validation
- 📊 **Progress Tracking** - Animated progress bars with elapsed time
- 🎮 **Easter Egg Game** - Space shooter mini-game during loading
- 📖 **Multi-tab Dashboard** - View growth plans, manifests, and contribution guides
- ⌨️ **Fully Keyboard Navigable** - No mouse required
- 🌐 **Cross-platform** - Works on macOS, Linux, and Windows

## Prerequisites

- Go 1.22 or later
- Python 3.8+ (for skene-growth installation)

## Quick Start

### Using Make (Recommended)

```bash
# Clone the repository
git clone https://github.com/skene-ai/skene-terminal-v2
cd skene-terminal-v2

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

### Application Flow

1. **Welcome Screen** - Press `ENTER` to begin
2. **Configuration Review** - View existing config, choose to edit or proceed
3. **Provider Selection** - Choose your LLM provider (OpenAI, Gemini, Anthropic, etc.)
4. **Model Selection** - Select the specific model for your provider
5. **API Key Entry** - Enter your API key securely
6. **Installation** - Watch progress as skene-growth is configured
7. **Dashboard** - Review your setup and next steps

### Keyboard Controls

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate up/down |
| `←/→` or `h/l` | Navigate left/right |
| `Enter` | Confirm/Select |
| `Esc` | Go back |
| `Tab` | Switch focus area |
| `?` | Toggle help overlay |
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
│       └── main.go           # Application entry point
├── internal/
│   ├── tui/
│   │   ├── app.go            # Main application model
│   │   ├── styles/
│   │   │   └── styles.go     # Lip Gloss styling system
│   │   ├── components/
│   │   │   ├── button.go     # Button components
│   │   │   ├── help.go       # Help overlay
│   │   │   ├── logo.go       # ASCII logo animations
│   │   │   └── progress.go   # Progress bars
│   │   └── views/
│   │       ├── intro.go      # Welcome screen
│   │       ├── config.go     # Configuration review
│   │       ├── provider.go   # Provider selection
│   │       ├── model.go      # Model selection
│   │       ├── auth.go       # Auth simulation
│   │       ├── apikey.go     # API key entry
│   │       ├── generating.go # Progress view
│   │       ├── dashboard.go  # Final dashboard
│   │       └── error.go      # Error handling
│   ├── services/
│   │   ├── analyzer/
│   │   │   └── analyzer.go   # Project analysis
│   │   ├── config/
│   │   │   └── manager.go    # Configuration management
│   │   └── installer/
│   │       └── installer.go  # Installation engine
│   └── game/
│       └── shooter.go        # Space shooter game
├── designs/                   # Design reference images
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Architecture

### State Management

The application uses a finite state machine pattern:

```
Intro → Config → Provider → Model → APIKey → Generating → Dashboard
                    ↓
                  Auth (for Skene provider)
```

### Views

Each view implements a consistent interface:
- `SetSize(width, height int)` - Handle terminal resize
- `Render() string` - Return the view content
- `GetHelpItems() []HelpItem` - Context-specific help

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
  "output_dir": "./skene-context",
  "verbose": true
}
```

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

# Output will be in build/
# - skene-linux-amd64
# - skene-darwin-amd64
# - skene-darwin-arm64
# - skene-windows-amd64.exe
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Support

- 📖 [Documentation](https://docs.skene.ai)
- 🐛 [Issues](https://github.com/skene-ai/skene-terminal-v2/issues)
- 💬 [Discord](https://discord.gg/skene)
