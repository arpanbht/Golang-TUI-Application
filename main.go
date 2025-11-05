package main

import (
	"fmt"
	"log"
	"os"

	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var (
	vaultDir string

	cursorStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("110"))
	inputPlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("217"))
	cursorLineStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("189"))
	// bulletPromtStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
)

type model struct {
	newFileInput           textinput.Model
	createFileInputVisible bool
}

func initialModel() model {
	err := os.MkdirAll(vaultDir, 0750)
	if err != nil {
		log.Fatal("Error creating vault directory ⛔", err)
	}

	ti := textinput.New()
	ti.Placeholder = "What would be the file name?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Cursor.Style = cursorStyle
	ti.TextStyle = cursorLineStyle
	ti.PlaceholderStyle = inputPlaceholderStyle
	ti.Prompt = "⚡ "

	return model{
		newFileInput:           ti,
		createFileInputVisible: false,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "ctrl+q":
			return m, tea.Quit

		case "ctrl+n":
			m.createFileInputVisible = true
			return m, nil

		case "enter":

		}
	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {

	var styleWelcomeMsg = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("17")).
		Background(lipgloss.Color("211")).
		Padding(1, 1)

	welcomeMsg := styleWelcomeMsg.Render("Welcome to tui-app! 🥳")

	view := ""

	if m.createFileInputVisible {
		view = m.newFileInput.View()
	}

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

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory ⛔", err)
	}

	vaultDir = fmt.Sprintf("%s/.vault", homeDir)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
