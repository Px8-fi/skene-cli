package main

import (
	"fmt"
	"os"

	"skene-terminal-v2/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create the application
	app := tui.NewApp()

	// Create the program with alt screen
	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Set program reference for background task communication
	app.SetProgram(p)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
