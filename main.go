package main

import (
	"log"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	questions   []string
	width       int
	height      int
	index       int
	answerField textarea.Model
	textState   string
}

func New(questions []string) *model {

	answerField := textarea.New()
	answerField.Focus()
	return &model{
		index:       0,
		questions:   questions,
		answerField: answerField,
		textState:   "",
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.index += 1
			if m.index >= len(m.questions) {
				return m, tea.Quit
			}
			m.answerField.Cursor.SetMode(cursor.CursorStatic)
			return m, nil
		}
	}
	newAns, cmd := m.answerField.Update(msg)
	m.answerField = newAns
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	if m.index >= len(m.questions) {
		return "closing"
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.questions[m.index],
			lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("36")).
				BorderStyle(lipgloss.NormalBorder()).
				Padding(1).
				Width(m.width-2).
				Render(m.answerField.View()),
		),
	)
}

func main() {
	m := New([]string{"hi what is going on", "hello", "it's me"})

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %w", err)
	}
	defer f.Close()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
