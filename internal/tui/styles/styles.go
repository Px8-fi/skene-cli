package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette - warm, retro terminal aesthetic
var (
	// Primary colors
	Cream     = lipgloss.Color("#EDC29C")
	Sand      = lipgloss.Color("#C9A97A")
	Charcoal  = lipgloss.Color("#1A1A1A")
	DarkGray  = lipgloss.Color("#2D2D2D")
	MidGray   = lipgloss.Color("#4A4A4A")
	LightGray = lipgloss.Color("#6A6A6A")
	White     = lipgloss.Color("#FFFFFF")

	// Accent colors
	Amber   = lipgloss.Color("#EDC29C")
	Rust    = lipgloss.Color("#D4A574")
	Coral   = lipgloss.Color("#F05D5E")
	Success = lipgloss.Color("#7CB374")
	Warning = lipgloss.Color("#E6B450")

	// Game colors
	GameCyan    = lipgloss.Color("#4CC9F0")
	GameMagenta = lipgloss.Color("#F72585")
	GameYellow  = lipgloss.Color("#FFD93D")
)

// Text styles
var (
	Title = lipgloss.NewStyle().
		Foreground(Cream).
		Bold(true)

	Subtitle = lipgloss.NewStyle().
			Foreground(Sand)

	Body = lipgloss.NewStyle().
		Foreground(White)

	Muted = lipgloss.NewStyle().
		Foreground(MidGray)

	Accent = lipgloss.NewStyle().
		Foreground(Amber)

	Error = lipgloss.NewStyle().
		Foreground(Coral)

	SuccessText = lipgloss.NewStyle().
			Foreground(Success)

	Label = lipgloss.NewStyle().
		Foreground(Cream).
		Bold(false)

	Value = lipgloss.NewStyle().
		Foreground(White)
)

// Layout styles
var (
	// Box with border
	Box = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(MidGray).
		Padding(1, 2)

	// Active box with highlighted border
	BoxActive = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Cream).
			Padding(1, 2)

	// Rounded box
	RoundedBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(MidGray).
			Padding(1, 2)

	// Section header
	SectionHeader = lipgloss.NewStyle().
			Foreground(Cream).
			Bold(false).
			MarginBottom(1)
)

// Button styles
var (
	Button = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(MidGray).
		Foreground(White).
		Padding(0, 3)

	ButtonActive = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Cream).
			Foreground(Charcoal).
			Background(Cream).
			Padding(0, 3)

	ButtonMuted = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(MidGray).
			Foreground(MidGray).
			Padding(0, 3)
)

// Tab styles
var (
	TabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	TabInactive = lipgloss.NewStyle().
			Border(TabBorder, true).
			BorderForeground(MidGray).
			Foreground(MidGray).
			Padding(0, 2)

	TabActive = lipgloss.NewStyle().
			Border(TabBorder, true).
			BorderForeground(Cream).
			Foreground(White).
			Padding(0, 2)
)

// List styles
var (
	ListItem = lipgloss.NewStyle().
			Foreground(White).
			PaddingLeft(2)

	ListItemSelected = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(Amber).
				Foreground(Amber).
				PaddingLeft(1)

	ListItemDimmed = lipgloss.NewStyle().
			Foreground(MidGray).
			PaddingLeft(2)

	ListDescription = lipgloss.NewStyle().
			Foreground(MidGray).
			PaddingLeft(2)

	ListDescriptionSelected = lipgloss.NewStyle().
				Foreground(Sand).
				PaddingLeft(1)
)

// Progress bar colors
var (
	ProgressFilled = Amber
	ProgressEmpty  = MidGray
)

// Help styles
var (
	HelpKey = lipgloss.NewStyle().
		Foreground(Cream).
		Bold(true)

	HelpDesc = lipgloss.NewStyle().
			Foreground(MidGray)

	HelpSeparator = lipgloss.NewStyle().
			Foreground(MidGray).
			SetString(" • ")
)

// ASCII art style
var (
	ASCII = lipgloss.NewStyle().
		Foreground(Cream)

	ASCIIAnimated = lipgloss.NewStyle().
			Foreground(Amber)
)

// Table styles
var (
	TableHeader = lipgloss.NewStyle().
			Foreground(MidGray).
			Bold(false)

	TableRow = lipgloss.NewStyle().
			Foreground(White)

	TableSeparator = lipgloss.NewStyle().
			Foreground(MidGray)
)

// Modal styles
var (
	Modal = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(MidGray).
		Padding(1, 2).
		Align(lipgloss.Center)

	ModalTitle = lipgloss.NewStyle().
			Foreground(White).
			Bold(false).
			MarginBottom(1)
)

// Footer help bar
var (
	FooterHelp = lipgloss.NewStyle().
			Foreground(MidGray).
			MarginTop(1)
)

// Spinner style
var (
	Spinner = lipgloss.NewStyle().
		Foreground(Amber)
)

// Center helper
func Center(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center)
}

// PlaceCenter centers content in given dimensions
func PlaceCenter(width, height int, content string) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// Divider creates a horizontal line
func Divider(width int) string {
	return lipgloss.NewStyle().
		Foreground(MidGray).
		Render(repeatString("─", width))
}

func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// PageTitle creates a centered page title
func PageTitle(title string, width int) string {
	return lipgloss.NewStyle().
		Foreground(Cream).
		Width(width).
		Align(lipgloss.Center).
		Render(title)
}
