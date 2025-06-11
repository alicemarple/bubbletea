package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return nil
}

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
			} else if m.focus == focusRight && m.rightCursor > 0 {
				m.rightCursor--
			}

		case "down", "j":
			if m.focus == focusLeft && m.cursor < len(m.items)-1 {
				m.cursor++
			} else if m.focus == focusRight {
				selected := m.items[m.cursor]
				if m.rightCursor < len(m.content[selected])-1 {
					m.rightCursor++
				}
			}

		case "h":
			m.focus = focusLeft

		case "l":
			m.focus = focusRight

		case "s":
			m.notification = ""
			if err := saveSelectionToFile(m); err != nil {
				m.notification = "  Failed to save file"
				return m, clearNotificationAfter(5 * time.Second)
			}
			m.notification = "  Selection saved successfully!"
			return m, clearNotificationAfter(5 * time.Second)

		case " ":
			if m.focus == focusRight {
				category := m.items[m.cursor]

				if category == "Environment" {
					// Deselect all environment options
					for i := range m.content["Environment"] {
						m.selectedRight["Environment"][i] = false
					}

					// Select current
					m.selectedRight["Environment"][m.rightCursor] = true
					selectedEnv := m.content["Environment"][m.rightCursor]

					// Load package list
					var pkgFile string
					switch selectedEnv {
					case "Dev":
						pkgFile = "./config/pkglist.dev.txt"
					case "Hacking":
						pkgFile = "./config/pkglist.cy.txt"
					}

					if pkgFile != "" {
						if pkgList, err := readLinesFromFile(pkgFile); err == nil {
							m.content["Packages"] = pkgList
							m.selectedRight["Packages"] = make(map[int]bool)
							for i := range pkgList {
								m.selectedRight["Packages"][i] = true
							}
						} else {
							m.content["Packages"] = []string{"[error reading file]"}
						}
					}
				} else {
					m.selectedRight[category][m.rightCursor] = !m.selectedRight[category][m.rightCursor]
				}
			}
		}

	case clearNotificationMsg:
		m.notification = ""
		return m, nil
	}

	return m, nil
}
