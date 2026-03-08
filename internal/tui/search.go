package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type searchModel struct {
	query   string
	focused bool
}

func newSearch() searchModel {
	return searchModel{}
}

func (s *searchModel) Focus() {
	s.focused = true
}

func (s searchModel) Update(msg tea.KeyMsg) (searchModel, tea.Cmd) {
	switch msg.String() {
	case "backspace":
		if len(s.query) > 0 {
			s.query = s.query[:len(s.query)-1]
		}
	default:
		if len(msg.String()) == 1 {
			s.query += msg.String()
		}
	}
	return s, nil
}

func (s searchModel) View() string {
	return searchStyle.Render(fmt.Sprintf("  search: %s█", s.query))
}

var searchStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFAA00")).
	Bold(true)
