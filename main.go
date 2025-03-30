package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Width(100). // Set width
			Height(10).                                                         // Set height for vertical centering
			Align(lipgloss.Center)
	progressStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("72"))
)

type tickMsg struct{}

type Item struct {
	name   string
	bought bool
}

type model struct {
	cursor     int
	items      []Item
	selected   []string
	installing bool
	progress   int
}

func initialModel() model {
	return model{
		items: []Item{
			{"Neovim", false},
			{"Bat", false},
			{"Fzf", false},
			{"Ripgrep", false},
			{"Eza", false},
			{"Lazygit", false},
		},
	}
}

// Init is the initial setup function
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages (user input)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case " ": // Toggle checkbox
			m.items[m.cursor].bought = !m.items[m.cursor].bought
		case "enter":
			m.selected = nil
			for _, item := range m.items {
				if item.bought {
					m.selected = append(m.selected, item.name)
				}
			}
			if len(m.selected) > 0 {
				m.installing = true
				m.progress = 0
				return m, installPackages()
			}
			return m, tea.Quit
		}
	case tickMsg:
		if m.installing {
			if m.progress < 100 {
				m.progress += 10
				return m, tea.Tick(time.Second/2, func(time.Time) tea.Msg {
					return tickMsg{}
				})
			} else {
				m.installing = false
			}
		}
	}
	return m, nil
}

func installPackages() tea.Cmd {
	return tea.Tick(time.Second/2, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

// View renders the UI
func (m model) View() string {
	var output string
	for i, item := range m.items {
		cursor := " " // Cursor indicator
		if i == m.cursor {
			cursor = ">" // Highlight current selection
		}

		checkbox := "[ ]" // Checkbox
		if item.bought {
			checkbox = "[x]" // Marked as bought
		}

		output += fmt.Sprintf("%s %s %s\n", cursorStyle.Render(cursor), keywordStyle.Render(checkbox), keywordStyle.Render(item.name))
	}

	if m.installing {
		output += "\nInstalling selected packages:\n"
		output += progressStyle.Render(fmt.Sprintf("[%s] %d%%", string(repeat('#', m.progress/10)), m.progress)) + "\n"
	} else if len(m.selected) > 0 {
		output += "\nSelected items:\n"
		for _, item := range m.selected {
			output += fmt.Sprintf("- %s\n", keywordStyle.Render(item))
		}
	} else {
		output += "\nNo items selected.\n"
	}
	output += "\n\n\n\n" + helpStyle.Render("  j: down • k: up • space: toggle • enter: install • q: exit\n")

	return output
}

func repeat(char rune, count int) string {
	return string(make([]rune, count, count))
}

func main() {
	m := initialModel()
	if _, err := tea.NewProgram(&m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
