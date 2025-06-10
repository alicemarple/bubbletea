package main

import (
	"fmt"
	"os"
	"strings"

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

type focusArea int

const (
	focusLeft focusArea = iota
	focusRight
)

type model struct {
	cursor  int
	focus   focusArea
	width   int
	height  int
	items   []string
	content map[string][]string
}

func initialModel() model {
	return model{
		cursor: 0,
		focus:  focusLeft,
		items:  []string{"Environment", "OS", "Packages"},
		content: map[string][]string{
			"Environment": {"Dev", "Hacking"},
			"OS":          {"Kali", "Arch", "Ubuntu"},
			"Packages":    {"Neovim", "Tmux", "Zsh"},
		},
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.focus == focusLeft && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.focus == focusLeft && m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "h":
			m.focus = focusLeft
		case "l":
			m.focus = focusRight
		}
	}
	return m, nil
}

func (m model) View() string {
	leftStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(unfocusedBorder).
		Width(m.width*15/100).
		Height(m.height*30/100).
		Margin(2, 1)

	rightStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(unfocusedBorder).
		Width(m.width*60/100).
		Height(m.height*80/100).
		Margin(2, 0)
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")).
		Align(lipgloss.Center).
		Width(m.width).
		MarginTop(1)

	if m.focus == focusLeft {
		leftStyle = leftStyle.BorderForeground(focusedBorder)
	} else {
		rightStyle = rightStyle.BorderForeground(focusedBorder)
	}

	left := m.renderLeft()
	right := m.renderRight()
	help := helpStyle.Render("  space: toggle select • i: install • j: up • k: down • q: exit\n")

	up := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(left),
		rightStyle.Render(right),
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		up,
		help,
	)
}

func (m model) renderLeft() string {
	var b strings.Builder
	for i, item := range m.items {
		if i == m.cursor {
			fmt.Fprintf(&b, "%s\n", highlight("> "+item))
		} else {
			fmt.Fprintf(&b, "  %s\n", item)
		}
	}
	return b.String()
}

func (m model) renderRight() string {
	selected := m.items[m.cursor]
	return strings.Join(m.content[selected], "\n")
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
