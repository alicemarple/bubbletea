package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

	notificationStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")).
		Width(m.width).
		Align(lipgloss.Center).
		MarginTop(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")).
		Align(lipgloss.Center).
		Width(m.width).
		MarginTop(1)

	// -- Render left and right focus
	if m.focus == focusLeft {
		leftStyle = leftStyle.BorderForeground(focusedBorder)
	} else {
		rightStyle = rightStyle.BorderForeground(focusedBorder)
	}

	left := m.renderLeft()
	right := m.renderRight()

	// -- Bottom block (for progress or errors)
	// -- Notification block (only success/info)
	var notification string
	if m.notification != "" {
		notification = notificationStyle.Render(m.notification)
	}

	help := helpStyle.Render("  space: toggle select • i: install • s: save • j: up • k: down • q: exit")

	up := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(left),
		rightStyle.Render(right),
	)
	return lipgloss.JoinVertical(
		lipgloss.Top,
		up,
		notification,
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
	lines := m.content[selected]
	selectedMap := m.selectedRight[selected]

	var b strings.Builder
	for i, line := range lines {
		cursor := "  "
		check := "[ ]"
		if i == m.rightCursor {
			cursor = highlight("> ")
		}
		if selectedMap[i] {
			check = "[x]"
		}
		fmt.Fprintf(&b, "%s%s %s\n", cursor, check, line)
	}
	return b.String()
}

func saveSelectionToFile(m model) error {
	// Get selected environment (first one marked true)
	env := "unknown"
	for i, v := range m.selectedRight["Environment"] {
		if v {
			env = strings.ToLower(m.content["Environment"][i])
			break
		}
	}

	// Get selected OS
	selectedOS := "none"
	for i, v := range m.selectedRight["OS"] {
		if v {
			selectedOS = strings.ToLower(m.content["OS"][i])
			break
		}
	}

	// Get selected Packages
	var selectedPackages []string
	for i, v := range m.selectedRight["Packages"] {
		if v && i < len(m.content["Packages"]) {
			selectedPackages = append(selectedPackages, m.content["Packages"][i])
		}
	}

	// Build output string
	var b strings.Builder
	b.WriteString("environment : " + env + "\n")
	b.WriteString("os : " + selectedOS + "\n")
	b.WriteString("packages :\n")
	for _, pkg := range selectedPackages {
		b.WriteString(pkg + "\n")
	}

	// Write to file
	filename := fmt.Sprintf("./download.%s.txt", env)
	return os.WriteFile(filename, []byte(b.String()), 0644)
}

func clearNotificationAfter(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return clearNotificationMsg{}
	})
}
