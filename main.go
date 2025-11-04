package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type model struct {
	msg string
}

func initialModel() model {
	return model{
		msg: "Hello, Bubble Tea!",
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {

	var styleWelcomeMsg = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("17")).
		Background(lipgloss.Color("211")).
		Padding(1, 1)

	welcomeMsg := styleWelcomeMsg.Render("Welcome to tui-app! 🥳")

	view := ""

	var styleHelp = lipgloss.NewStyle().
		Foreground(lipgloss.Color("247")).
		Italic(true).
		Margin(1, 0)

	help := styleHelp.Render("Ctrl+N: New File . Ctrl+L: List . Esc: back/save . Ctrl+S: Save . Ctrl+Q: Quit")

	return fmt.Sprintf("\n%s\n\n%s\n\n%s", welcomeMsg, view, help)
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
