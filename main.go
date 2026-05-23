package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	insert = iota
	visual
)

const logo = `
████████████
████████████
████████████
████████████
████████████
`

type model struct {
	questions   []string
	width       int
	height      int
	index       int
	answerField textarea.Model
	textState   string
	mode        int
}

func New(questions []string) *model {

	answerField := textarea.New()
	answerField.Focus()
	answerField.Placeholder = "Once upon a time ..."

	return &model{
		index:       0,
		questions:   questions,
		answerField: answerField,
		mode:        visual,
		textState:   "",
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mode = visual
		case "ctrl+c":
			return m, tea.Quit
		default:
			switch m.mode {
			case visual:
				switch msg.String() {
				case "i":
					m.mode = insert
				case "j":
					m.answerField, cmd = m.answerField.Update(tea.KeyMsg{Type: tea.KeyDown})
					cmds = append(cmds, cmd)
				case "k":
					m.answerField, cmd = m.answerField.Update(tea.KeyMsg{Type: tea.KeyUp})
					cmds = append(cmds, cmd)
				case "h":
					m.answerField, cmd = m.answerField.Update(tea.KeyMsg{Type: tea.KeyLeft})
					cmds = append(cmds, cmd)
				case "l":
					m.answerField, cmd = m.answerField.Update(tea.KeyMsg{Type: tea.KeyRight})
					cmds = append(cmds, cmd)
				}
			case insert:
				m.answerField, cmd = m.answerField.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	if !m.answerField.Focused() {
		cmd = m.answerField.Focus()
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	if m.index >= len(m.questions) {
		return "closing"
	}

	logoW := lipgloss.Width(logo)
	// chrome budget: outer border+padding (4) + logo PaddingRight (2) + rightSide border (2) + outer Width offset (2) = 10
	textW := m.width - logoW - 10
	if textW < 0 {
		textW = 0
	}

	rightSide := lipgloss.NewStyle().
		Width(textW). // only the text side is sized; it can wrap safely
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("21")).
		Render("Hi this is a longer string to use")

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("25")).
				Padding(1).
				BorderStyle(lipgloss.RoundedBorder()).
				Width(m.width-2).
				Align(lipgloss.Center).
				Render(
					lipgloss.JoinHorizontal(
						lipgloss.Center,
						lipgloss.NewStyle().
							PaddingRight(2).
							Foreground(lipgloss.Color("36")).
							Render(logo),
						rightSide,
					)),
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
