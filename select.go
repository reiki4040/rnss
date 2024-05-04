package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func NewSelectionCUI(list []string, filter []rune) (*model, error) {
	return &model{
		cursor:       0,
		filter:       filter,
		filteredList: []string{},
		list:         list,
	}, nil
}

type model struct {
	cursor       int
	filter       []rune
	filteredList []string

	list []string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.cursor = -1
			return m, tea.Quit
		case "up":
			m.cursor -= 1
			if m.cursor < 0 {
				m.cursor = len(m.filteredList) - 1
			}
		case "down":
			m.cursor += 1
			if m.cursor >= len(m.filteredList) {
				m.cursor = 0
			}
		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
			}
		case "enter":
			if len(m.filteredList) > 0 {
				return m, tea.Quit
			}
		default:
			m.filter = append(m.filter, []rune(msg.String())...)
		}
	}

	// filter selection with filter ignore upper/lower
	m.filteredList = nil
	filter := strings.ToLower(string(m.filter))
	for _, id := range m.list {
		if strings.Contains(strings.ToLower(id), filter) {
			m.filteredList = append(m.filteredList, id)
		}
	}

	return m, nil
}

func (m model) View() string {
	var sb strings.Builder

	cursorColor := termenv.ColorProfile().Color("#00FFFF")
	cursorLineBackgroundColor := termenv.ColorProfile().Color("#333333")
	cursorLineCharColor := termenv.ColorProfile().Color("#FF00FF")
	filterMatchCharColor := termenv.ColorProfile().Color("#FFD700")

	// show list
	for i, selection := range m.filteredList {
		cursor := " "
		if i == m.cursor {
			cursor = "> "
			sb.WriteString(termenv.String(cursor).Foreground(cursorColor).String())
			sb.WriteString(termenv.String(selection).Background(cursorLineBackgroundColor).Foreground(cursorLineCharColor).String())
		} else {
			sb.WriteString(" ")
			// show selection with highlight if matched filter
			filtered := string(m.filter)
			if len(filtered) > len(selection) {
				filtered = filtered[:len(selection)]
			}
			sb.WriteString(strings.ReplaceAll(selection, filtered, termenv.String(filtered).Foreground(filterMatchCharColor).String()))

		}
		sb.WriteString("\n")
	}

	// show filter input prompt
	sb.WriteString(fmt.Sprintf("\nFilter: %s", string(m.filter)))

	return sb.String()
}
