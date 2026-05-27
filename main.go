package main

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

const (
	insert = iota
	visual
)

const logo = `
█████
█████
`

type layout struct {
	navbarH int
	navbarW int

	footerH int
	footerW int

	browsingH int
	browsingW int

	libraryH int
	libraryW int

	queueH int
	queueW int
}

type model struct {
	dims         layout
	screenWidth  int
	screenHeight int

	navbarPad    int
	navbarBorder int
	navbarMargin int

	footerPad    int
	footerBorder int
	footerMargin int

	browsingPad    int
	browsingBorder int
	browsingMargin int

	libraryPad    int
	libraryBorder int
	libraryMargin int

	queuePad    int
	queueBorder int
	queueMargin int
}

func initialModel() model {

	m := model{
		dims:           layout{},
		screenWidth:    0,
		screenHeight:   0,
		navbarBorder:   1,
		footerBorder:   1,
		browsingBorder: 1,
		libraryBorder:  1,
		queueBorder:    1,
	}

	m.recalc()

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
		m.screenHeight = msg.Height
		m.recalc()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) recalc() {
	screenW := m.screenWidth
	screenH := m.screenHeight

	var navbarH int
	var navbarW int
	navbarPad := m.navbarPad
	navbarBorder := m.navbarBorder
	navbarMargin := m.navbarMargin

	var footerH int
	var footerW int
	footerPad := m.footerPad
	footerBorder := m.footerBorder
	footerMargin := m.footerMargin

	var browsingH int
	var browsingW int
	browsingPad := m.browsingPad
	browsingBorder := m.browsingBorder
	browsingMargin := m.browsingMargin

	var libraryH int
	var libraryW int
	libraryBorder := m.libraryBorder
	libraryMargin := m.libraryMargin

	var queueH int
	var queueW int
	queueBorder := m.queueBorder
	queueMargin := m.queueMargin

	// compute navbar height. It takes priority
	navbarWantHeight := 1 + 2*navbarPad + 2*navbarBorder + 2*navbarMargin
	if screenH > navbarWantHeight {
		navbarH = 1
		// the navbar takes the entire width
		navbarW = screenW - 2*navbarPad - 2*navbarBorder - 2*navbarMargin
	}
	// then the footer
	footerWantHeight := 2 + 2*footerPad + 2*footerBorder + 2*footerMargin
	if screenH > footerWantHeight+navbarWantHeight {
		footerH = 2
		// the footer also takes the entire width
		footerW = screenW - 2*footerPad - 2*footerBorder - 2*footerMargin
	}

	browsingWantHeight := screenH - footerWantHeight - navbarWantHeight - 2*browsingBorder + 2*browsingMargin + 2*browsingPad
	if browsingWantHeight > 0 {
		browsingH = browsingWantHeight
	}

	// middle row: library | browsing | queue. Library takes priority, then queue.
	const browsingMinW = 18
	libraryWantW := 16 + 2*libraryBorder + 2*libraryMargin
	queueWantW := 16 + 2*queueBorder + 2*queueMargin
	browsingBaseW := 2*browsingBorder + 2*browsingMargin

	if browsingH > 0 && screenW >= libraryWantW+browsingBaseW+browsingMinW {
		libraryW = 16
		libraryH = browsingH
	}
	if browsingH > 0 && screenW >= libraryWantW+queueWantW+browsingBaseW+browsingMinW {
		queueW = 16
		queueH = browsingH
	}

	usedW := 0
	if libraryW > 0 {
		usedW += libraryWantW
	}
	if queueW > 0 {
		usedW += queueWantW
	}
	if browsingH > 0 {
		browsingW = screenW - usedW - browsingBaseW
	}

	m.dims.navbarH = navbarH
	m.dims.navbarW = navbarW
	m.dims.footerH = footerH
	m.dims.footerW = footerW
	m.dims.browsingH = browsingH
	m.dims.browsingW = browsingW
	m.dims.libraryH = libraryH
	m.dims.libraryW = libraryW
	m.dims.queueH = queueH
	m.dims.queueW = queueW

}

func (m model) View() string {
	var rows []string

	if m.dims.navbarH > 0 {
		// navbar
		rows = append(rows, lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("23")).
			BorderStyle(lipgloss.NormalBorder()).
			Width(m.dims.navbarW).
			Height(m.dims.navbarH).
			Render("navbar"))
	}

	var middle []string

	if m.dims.libraryH > 0 {
		// library content
		logoStr := strings.TrimSpace(logo)
		logoW := lipgloss.Width(logoStr)
		logoH := lipgloss.Height(logoStr)
		gap := 1
		maxTitleW := m.dims.libraryW - logoW - gap
		if maxTitleW < 1 {
			maxTitleW = 1
		}
		maxRows := m.dims.libraryH / logoH

		albums := []struct {
			title, artist string
		}{
			{"Album title 1", "Artist 1"},
			{"Album title 2", "Artist 2"},
			{"Album title 3", "Artist 3"},
			{"Album title 4", "Artist 4"},
			{"Album title 5", "Artist 5"},
			{"Album title 6", "Artist 6"},
			{"Album title 7", "Artist 7"},
			{"Album title 8", "Artist 8"},
			{"Album title 9", "Artist 9"},
			{"Album title 10", "Artist 10"},
			{"Album title 11", "Artist 11"},
			{"Album title 12", "Artist 12"},
			{"Album title 13", "Artist 13"},
			{"Album title 14", "Artist 14"},
			{"Album title 15", "Artist 15"},
			{"Album title 16", "Artist 16"},
			{"Album title 17", "Artist 17"},
			{"Album title 18", "Artist 18"},
			{"Album title 19", "Artist 19"},
			{"Album title 20", "Artist 20"},
			{"Album title 21", "Artist 21"},
		}
		if len(albums) > maxRows {
			albums = albums[:maxRows]
		}

		artistStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

		var libRows []string
		for _, a := range albums {
			title := truncate.StringWithTail(a.title, uint(maxTitleW), "…")
			artist := truncate.StringWithTail(a.artist, uint(maxTitleW), "…")
			info := lipgloss.JoinVertical(lipgloss.Left, title, artistStyle.Render(artist))
			row := lipgloss.JoinHorizontal(lipgloss.Center, logoStr, " ", info)
			libRows = append(libRows, row)
		}
		libContent := lipgloss.JoinVertical(lipgloss.Left, libRows...)

		// library
		middle = append(middle, lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("33")).
			BorderStyle(lipgloss.NormalBorder()).
			Width(m.dims.libraryW).
			Height(m.dims.libraryH).
			Render(libContent))
	}
	if m.dims.browsingH > 0 {
		//browsing
		middle = append(middle, lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("27")).
			BorderStyle(lipgloss.NormalBorder()).
			Width(m.dims.browsingW).
			Height(m.dims.browsingH).
			Render("browsing"))
	}
	if m.dims.queueH > 0 {
		// queue
		middle = append(middle, lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("45")).
			BorderStyle(lipgloss.NormalBorder()).
			Width(m.dims.queueW).
			Height(m.dims.queueH).
			Render("queue"))
	}
	if len(middle) > 0 {
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, middle...))
	}

	if m.dims.footerH > 0 {
		logoStr := strings.TrimSpace(logo)

		footerContent := lipgloss.JoinHorizontal(
			lipgloss.Center,
			logoStr,
			" ",
			lipgloss.JoinVertical(
				lipgloss.Left,
				"Song title",
				lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Song artist"),
			),
		)

		rows = append(rows, lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("36")).
			BorderStyle(lipgloss.NormalBorder()).
			Height(m.dims.footerH).
			Width(m.dims.footerW).
			Render(footerContent))
	}

	return lipgloss.Place(
		m.screenWidth,
		m.screenHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, rows...),
	)
}

func main() {
	m := initialModel()

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
