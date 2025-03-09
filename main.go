package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// Define application screens
type screen int

const (
	mainMenu screen = iota
	environment
	desktop
	packages
	install
	help
)

// Run Installation Script
func runInstallScript() {
	fmt.Println("installing tools ....")
	cmd := exec.Command("/bin/bash", "./scripts/install.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running install script:", err)
	}
}

// Menu Items
var menuItems = []string{"Environment", "Desktop", "Packages", "Install", "Help", "Quit"}

// Environment Options (checkboxes)
type envItem struct {
	name    string
	enabled bool
}

var environmentItems = []envItem{
	{"Hacking", false},
	{"Development", false},
	{"Experiment", false},
}

// Desktop Options (radio selection)
type desktopItem struct {
	name     string
	selected bool
}

var desktopItems = []desktopItem{
	{"KDE", false},
	{"GNOME", false},
	{"Hyprland", false},
}

// Model holds the state of the application
type model struct {
	currentScreen screen
	cursor        int // Cursor for main menu
	envCursor     int // Cursor for environment menu
	desktopCursor int // Cursor for desktop menu
}

// Init initializes the program
func (m model) Init() tea.Cmd {
	return nil
}

// Get Packages based on selections
func getPackages() []string {
	var selectedEnv []string
	for _, env := range environmentItems {
		if env.enabled {
			selectedEnv = append(selectedEnv, env.name)
		}
	}

	var selectedDesktop string
	for _, desk := range desktopItems {
		if desk.selected {
			selectedDesktop = desk.name
			break
		}
	}

	// Base packages for all setups
	packages := []string{"Base System", "Essential Tools"}

	// Add environment-specific packages
	for _, env := range selectedEnv {
		switch env {
		case "Hacking":
			packages = append(packages, "Metasploit", "Wireshark", "Nmap", "Burp Suite")

		case "Development":
			packages = append(packages, "VS Code", "Git", "Docker", "Node.js", "Go")
		case "Experiment":

			packages = append(packages, "VirtualBox", "QEMU", "Pentesting Labs")
		}
	}

	// Add desktop-specific packages
	switch selectedDesktop {
	case "KDE":
		packages = append(packages, "Plasma", "Dolphin", "Konsole")
	case "GNOME":
		packages = append(packages, "Gnome Shell", "Nautilus", "Gnome Terminal")

	case "Hyprland":

		packages = append(packages, "Hyprland", "Waybar", "Alacritty")
	}

	return packages
}

// Update function to handle key events
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			if m.currentScreen == mainMenu && m.cursor > 0 {
				m.cursor--
			} else if m.currentScreen == environment && m.envCursor > 0 {
				m.envCursor--
			} else if m.currentScreen == desktop && m.desktopCursor > 0 {
				m.desktopCursor--
			}
		case "k":

			if m.currentScreen == mainMenu && m.cursor < len(menuItems)-1 {
				m.cursor++
			} else if m.currentScreen == environment && m.envCursor < len(environmentItems)-1 {
				m.envCursor++
			} else if m.currentScreen == desktop && m.desktopCursor < len(desktopItems)-1 {
				m.desktopCursor++
			}
		case "enter":
			if m.currentScreen == mainMenu {
				switch m.cursor {
				case 0:
					m.currentScreen = environment
				case 1:
					m.currentScreen = desktop
				case 2:
					m.currentScreen = packages
				case 3:
					m.currentScreen = install
				case 4:
					m.currentScreen = help
				case 5:
					return m, tea.Quit
				}
			}
		case " ": // Toggle checkboxes in environment
			if m.currentScreen == environment {
				environmentItems[m.envCursor].enabled = !environmentItems[m.envCursor].enabled
			}
		case "d": // Select desktop environment (radio selection)
			if m.currentScreen == desktop {
				for i := range desktopItems {
					desktopItems[i].selected = false
				}
				desktopItems[m.desktopCursor].selected = true
			}
		case "b": // Go back to the previous screen
			m.currentScreen = mainMenu
		case "i":
			if m.currentScreen == install {
				runInstallScript()
				return m, tea.Quit
			}
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View function to render different screens

func (m model) View() string {
	switch m.currentScreen {
	case mainMenu:
		s := "Main Menu\nUse ↑/↓ to navigate, Enter to select:\n"
		for i, item := range menuItems {
			cursor := "  "
			if i == m.cursor {
				cursor = "> "
			}
			s += fmt.Sprintf("%s%s\n", cursor, item)
		}
		return s

	case environment:
		s := "Environment Selection\nUse ↑/↓ to navigate, Space to toggle, b: Back\n\n"
		for i, env := range environmentItems {
			cursor := "  "
			if i == m.envCursor {
				cursor = "> "
			}
			checkbox := "[ ]"
			if env.enabled {
				checkbox = "[x]"
			}
			s += fmt.Sprintf("%s%s %s\n", cursor, checkbox, env.name)
		}
		return s

	case desktop:
		s := "Desktop Selection\nUse ↑/↓ to navigate, d: Select, b: Back\n\n"
		for i, desk := range desktopItems {
			cursor := "  "

			if i == m.desktopCursor {
				cursor = "> "
			}

			radio := "( )"
			if desk.selected {
				radio = "(x)"
			}
			s += fmt.Sprintf("%s%s %s\n", cursor, radio, desk.name)

		}
		return s

	case packages:
		s := "Selected Packages:\nUse b: Back\n\n"
		for _, pkg := range getPackages() {
			s += fmt.Sprintf("- %s\n", pkg)
		}
		return s

	case install:
		return "Installation Setup\nPress 'i' to install, 'b' to go back.\n"

	case help:
		return "Help Screen\nUse ↑/↓ to navigate, Enter to select, q to quit, b to go back."

	}
	return "Unknown screen"
}

// Main function
func main() {
	p := tea.NewProgram(model{currentScreen: mainMenu})
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
