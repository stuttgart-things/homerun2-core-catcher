package tui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/models"
)

type tableModel struct {
	rows   []models.CaughtMessage
	cursor int
}

func newTable() tableModel {
	return tableModel{}
}

func (t tableModel) Render(width, maxRows int) string {
	if len(t.rows) == 0 {
		return emptyStyle.Render("  no messages yet — waiting for stream...")
	}

	// Column widths (proportional)
	colTitle := width * 30 / 100
	colSev := 10
	colSystem := width * 18 / 100
	colAuthor := width * 15 / 100
	colTime := 20
	colRemaining := width - colTitle - colSev - colSystem - colAuthor - colTime - 6
	if colRemaining > 0 {
		colTitle += colRemaining
	}

	var b strings.Builder

	// Header row
	hdr := fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s  %-*s",
		colTitle, "TITLE",
		colSev, "SEVERITY",
		colSystem, "SYSTEM",
		colAuthor, "AUTHOR",
		colTime, "TIMESTAMP",
	)
	b.WriteString(tableHeaderStyle.Width(width).Render(hdr) + "\n")

	// Data rows
	visible := t.rows
	if maxRows > 0 && len(visible) > maxRows {
		visible = visible[:maxRows]
	}

	for i, msg := range visible {
		title := truncate(msg.Title, colTitle)
		sev := truncate(msg.Severity, colSev)
		sys := truncate(msg.System, colSystem)
		author := truncate(msg.Author, colAuthor)
		ts := truncate(msg.CaughtAt.Format(time.DateTime), colTime)

		row := fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s  %-*s",
			colTitle, title,
			colSev, sev,
			colSystem, sys,
			colAuthor, author,
			colTime, ts,
		)

		style := rowStyle
		if i == t.cursor {
			style = selectedRowStyle
		}

		// Color severity
		if i == t.cursor {
			b.WriteString(style.Width(width).Render(row) + "\n")
		} else {
			styled := applySeverityColor(row, msg.Severity, style)
			b.WriteString(styled + "\n")
		}
	}

	return b.String()
}

func truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func applySeverityColor(row, severity string, base lipgloss.Style) string {
	switch severity {
	case "error":
		return base.Foreground(lipgloss.Color("#FF4444")).Render(row)
	case "warning":
		return base.Foreground(lipgloss.Color("#FFAA00")).Render(row)
	case "success":
		return base.Foreground(lipgloss.Color("#00CC66")).Render(row)
	case "debug":
		return base.Foreground(lipgloss.Color("#888888")).Render(row)
	default:
		return base.Render(row)
	}
}

var (
	tableHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00CC66")).
		Background(lipgloss.Color("#1A1A1A"))

	rowStyle = lipgloss.NewStyle()

	selectedRowStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#004422")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	emptyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Italic(true)
)
