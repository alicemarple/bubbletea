package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Colors and highlight function
var (
	highlight = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#eba0ac")).
			Render

	focusedBorder   = lipgloss.Color("#fab387")
	unfocusedBorder = lipgloss.Color("#89b4fa")
)

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
