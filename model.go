package main

import (
	"bufio"
	"os"
)

type (
	focusArea                int
	installDoneWithOutputMsg string
)

const (
	focusLeft focusArea = iota
	focusRight
)

type model struct {
	cursor        int
	rightCursor   int
	focus         focusArea
	width, height int
	items         []string
	content       map[string][]string
	selectedRight map[string]map[int]bool // category → index → selected
	notification  string
}
type clearNotificationMsg struct{}

func initialModel() model {
	// Load Dev packages by default
	pkgList, err := readLinesFromFile("./config/pkglist.dev.txt")
	if err != nil {
		pkgList = []string{"[error reading file]"} // fallback
	}

	content := map[string][]string{
		"Environment": {"Dev", "Hacking"},
		"OS":          {"Kali", "Arch", "Ubuntu"},
		"Packages":    pkgList,
	}

	// Initialize selectedRight map
	selected := make(map[string]map[int]bool)

	// Select "Dev" in Environment
	selected["Environment"] = make(map[int]bool)
	selected["Environment"][0] = true // 0 = "Dev"

	// Select all packages by default
	selected["Packages"] = make(map[int]bool)
	for i := range pkgList {
		selected["Packages"][i] = true
	}
	// Optional: You could initialize OS too
	selected["OS"] = make(map[int]bool)

	return model{
		cursor:        0,
		rightCursor:   0,
		focus:         focusLeft,
		items:         []string{"Environment", "OS", "Packages"},
		content:       content,
		selectedRight: selected,
	}
}

func readLinesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}
