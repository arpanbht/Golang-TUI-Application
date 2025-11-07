package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	vaultDir string

	cursorStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("110"))
	inputPlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("217"))
	cursorLineStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("189"))
	textAreaCursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	docStyle              = lipgloss.NewStyle().Margin(1, 2)
	listStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Background(lipgloss.Color("81")).Padding(0, 1)
	// bulletPromtStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	newFileInput           textinput.Model
	createFileInputVisible bool
	currentFile            *os.File
	noteTextArea           textarea.Model
	list                   list.Model
	isListVisible          bool
}

func initialModel() model {
	err := os.MkdirAll(vaultDir, 0750)
	if err != nil {
		log.Fatal("Error creating vault directory ⛔", err)
	}

	// New file input
	ti := textinput.New()
	ti.Placeholder = "What would be the file name?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Cursor.Style = cursorStyle
	ti.TextStyle = cursorLineStyle
	ti.PlaceholderStyle = inputPlaceholderStyle
	ti.Prompt = "⚡ "

	// Note text area
	ta := textarea.New()
	ta.Placeholder = "Write your thoughts 💭..."
	ta.Focus()
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Placeholder = textAreaCursorStyle
	ta.Cursor.Style = cursorStyle

	// List
	noteList := listFiles()
	finalList := list.New(noteList, list.NewDefaultDelegate(), 0, 0)
	finalList.Title = "Your Notes 📚"
	finalList.Styles.Title = listStyle

	return model{
		newFileInput:           ti,
		createFileInputVisible: false,
		noteTextArea:           ta,
		list:                   finalList,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-10)

	case tea.KeyMsg:

		switch msg.String() {

		case "q":
			return m, tea.Quit

		case "ctrl+n":
			m.createFileInputVisible = true
			return m, nil

		case "enter":
			if m.currentFile != nil {
				break
			}

			if m.isListVisible {
				item, ok := m.list.SelectedItem().(item)
				if ok {
					filepath := fmt.Sprintf("%s/%s", vaultDir, item.title)
					content, err := os.ReadFile(filepath)
					if err != nil {
						log.Fatal("Error reading file ⛔", err)
						return m, nil
					}
					m.noteTextArea.SetValue(string(content))

					file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
					if err != nil {
						log.Fatal("Error opening file ⛔", err)
						return m, nil
					}
					m.currentFile = file
					m.isListVisible = false
				}
				return m, nil
			}

			filename := m.newFileInput.Value()
			if filename != "" {
				filepath := fmt.Sprintf("%s/%s.md", vaultDir, filename)

				if _, err := os.Stat(filepath); err == nil {
					return m, nil
				}
				file, err := os.Create(filepath)
				if err != nil {
					log.Fatal("Error creating file ⛔", err)
				}

				m.currentFile = file
				m.createFileInputVisible = false
				m.newFileInput.SetValue("")

			}

			return m, nil

		case "ctrl+s":
			if m.currentFile == nil {
				break
			}

			if err := m.currentFile.Truncate(0); err != nil {
				fmt.Println("Cannot save the file 😥")
				return m, nil
			}

			if _, err := m.currentFile.Seek(0, 0); err != nil {
				fmt.Println("Cannot save the file 🥺")
				return m, nil
			}

			if _, err := m.currentFile.WriteString(m.noteTextArea.Value()); err != nil {
				fmt.Println("Cannot save the file 🥲")
				return m, nil
			}

			if err := m.currentFile.Close(); err != nil {
				fmt.Println("Cannot close the file 😭")
				return m, nil
			}

			m.currentFile = nil
			m.noteTextArea.SetValue("")

			return m, nil

		case "ctrl+l":
			noteList := listFiles()
			m.list.SetItems(noteList)
			m.isListVisible = true
			return m, nil

		case "esc":
			if m.createFileInputVisible {
				m.createFileInputVisible = false
				m.newFileInput.SetValue("")
				return m, nil
			}

			if m.currentFile != nil {
				m.currentFile = nil
				m.noteTextArea.SetValue("")
				return m, nil
			}

			if m.isListVisible {
				if m.list.FilterState() == list.Filtering {
					m.list.ResetFilter()
					return m, nil
				}
				m.isListVisible = false
				return m, nil
			}
			return m, nil

		case "ctrl+d":
			if m.isListVisible {
				item, ok := m.list.SelectedItem().(item)
				if ok {
					filepath := fmt.Sprintf("%s/%s", vaultDir, item.title)
					err := os.Remove(filepath)
					if err != nil {
						log.Fatal("Error deleting file ⛔", err)
						return m, nil
					}
					noteList := listFiles()
					m.list.SetItems(noteList)
				}
			}
			return m, nil
		}
	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
	}

	if m.isListVisible {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {

	var styleWelcomeMsg = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("17")).
		Background(lipgloss.Color("211")).
		Padding(0, 2).
		MarginBottom(1)

	welcomeMsg := styleWelcomeMsg.Render("Totion 🥳")

	view := ""
	if m.createFileInputVisible {
		view = m.newFileInput.View()
	}

	if m.currentFile != nil {
		view = m.noteTextArea.View()
	}

	if m.isListVisible {
		view = m.list.View()
	}

	var styleHelp = lipgloss.NewStyle().
		Foreground(lipgloss.Color("247")).
		Italic(true).
		Margin(1, 0)

	help := styleHelp.Render("Ctrl+N: New File . Ctrl+L: List . Esc: Back . Ctrl+S: Save . Q: Quit")

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

func listFiles() []list.Item {
	items := make([]list.Item, 0)

	entries, err := os.ReadDir(vaultDir)
	if err != nil {
		log.Fatal("Error reading vault directory ❌", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			modTime := info.ModTime().Format("2006-01-02 15:04:05")

			items = append(items, item{
				title: entry.Name(),
				desc:  fmt.Sprintf("Last modified: %s", modTime),
			})
		}
	}

	return items
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
