package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/models"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/store"
)

const (
	viewList   = iota
	viewDetail
)

// refreshMsg triggers a UI refresh.
type refreshMsg struct{}

// Model is the top-level Bubble Tea model.
type Model struct {
	store      *store.MessageStore
	view       int
	table      tableModel
	detail     detailModel
	search     searchModel
	width      int
	height     int
	sortField  store.SortField
	sortDir    store.SortDirection
	page       int
	pageSize   int
	searching  bool
}

// New creates a new TUI model.
func New(s *store.MessageStore) Model {
	return Model{
		store:    s,
		view:     viewList,
		table:    newTable(),
		detail:   newDetail(),
		search:   newSearch(),
		sortField: store.SortByTimestamp,
		sortDir:  store.SortDesc,
		page:     0,
		pageSize: 20,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return refreshMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case refreshMsg:
		m.table.rows = m.fetchRows()
		return m, tea.Tick(time.Second, func(time.Time) tea.Msg {
			return refreshMsg{}
		})

	case tea.KeyMsg:
		// Global keys
		if !m.searching {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}

		switch m.view {
		case viewList:
			return m.updateList(msg)
		case viewDetail:
			return m.updateDetail(msg)
		}
	}

	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.searching {
		switch msg.String() {
		case "enter":
			m.searching = false
			m.page = 0
			m.table.rows = m.fetchRows()
			return m, nil
		case "esc":
			m.searching = false
			m.search.query = ""
			m.page = 0
			m.table.rows = m.fetchRows()
			return m, nil
		default:
			m.search, _ = m.search.Update(msg)
			m.page = 0
			m.table.rows = m.fetchRows()
			return m, nil
		}
	}

	switch msg.String() {
	case "/":
		m.searching = true
		m.search.Focus()
		return m, nil
	case "s":
		m.sortField = (m.sortField + 1) % 5
		m.table.rows = m.fetchRows()
		return m, nil
	case "S":
		if m.sortDir == store.SortAsc {
			m.sortDir = store.SortDesc
		} else {
			m.sortDir = store.SortAsc
		}
		m.table.rows = m.fetchRows()
		return m, nil
	case "j", "down":
		if m.table.cursor < len(m.table.rows)-1 {
			m.table.cursor++
		}
		return m, nil
	case "k", "up":
		if m.table.cursor > 0 {
			m.table.cursor--
		}
		return m, nil
	case "n", "right":
		totalPages := m.totalPages()
		if m.page < totalPages-1 {
			m.page++
			m.table.cursor = 0
			m.table.rows = m.fetchRows()
		}
		return m, nil
	case "p", "left":
		if m.page > 0 {
			m.page--
			m.table.cursor = 0
			m.table.rows = m.fetchRows()
		}
		return m, nil
	case "enter":
		if len(m.table.rows) > 0 && m.table.cursor < len(m.table.rows) {
			m.detail.message = &m.table.rows[m.table.cursor]
			m.view = viewDetail
		}
		return m, nil
	}

	return m, nil
}

func (m Model) updateDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "backspace", "q":
		m.view = viewList
		return m, nil
	}
	return m, nil
}

func (m Model) View() tea.View {
	var content string
	if m.width == 0 {
		content = "loading..."
	} else {
		switch m.view {
		case viewDetail:
			content = m.renderDetail()
		default:
			content = m.renderList()
		}
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m Model) renderList() string {
	var b strings.Builder

	// Header
	header := headerStyle.Width(m.width).Render(
		fmt.Sprintf("  HOMERUN² core-catcher    [/] search  [s] sort  [S] dir  [n/p] page  [enter] detail  [q] quit"),
	)
	b.WriteString(header + "\n")

	// Search bar
	if m.searching {
		b.WriteString(m.search.View() + "\n")
	} else if m.search.query != "" {
		b.WriteString(searchActiveStyle.Render(fmt.Sprintf("  filter: %s", m.search.query)) + "\n")
	}

	// Sort info
	sortNames := []string{"timestamp", "severity", "system", "author", "title"}
	dirName := "↓"
	if m.sortDir == store.SortAsc {
		dirName = "↑"
	}
	b.WriteString(sortInfoStyle.Render(fmt.Sprintf("  sort: %s %s", sortNames[m.sortField], dirName)) + "\n")

	// Table
	tableHeight := m.height - 7
	if m.searching || m.search.query != "" {
		tableHeight--
	}
	b.WriteString(m.table.Render(m.width, tableHeight))

	// Footer
	total := m.store.Count()
	totalPages := m.totalPages()
	footer := footerStyle.Width(m.width).Render(
		fmt.Sprintf("  %d messages | page %d/%d | stream: messages", total, m.page+1, max(totalPages, 1)),
	)
	b.WriteString("\n" + footer)

	return b.String()
}

func (m Model) renderDetail() string {
	var b strings.Builder

	header := headerStyle.Width(m.width).Render("  Message Detail    [esc/backspace] back")
	b.WriteString(header + "\n\n")

	if m.detail.message != nil {
		b.WriteString(m.detail.Render(m.width))
	}

	return b.String()
}

func (m Model) fetchRows() []models.CaughtMessage {
	if m.search.query != "" {
		all := m.store.Search(m.search.query)
		start := m.page * m.pageSize
		if start >= len(all) {
			return nil
		}
		end := start + m.pageSize
		if end > len(all) {
			end = len(all)
		}
		return all[start:end]
	}
	return m.store.List(store.ListOptions{
		Offset:  m.page * m.pageSize,
		Limit:   m.pageSize,
		SortBy:  m.sortField,
		SortDir: m.sortDir,
	})
}

func (m Model) totalPages() int {
	total := m.store.Count()
	if total == 0 {
		return 1
	}
	return (total + m.pageSize - 1) / m.pageSize
}

// Styles
var (
	headerStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#00CC66")).
		Foreground(lipgloss.Color("#000000")).
		Bold(true)

	footerStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(lipgloss.Color("#CCCCCC"))

	sortInfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	searchActiveStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFAA00"))
)
