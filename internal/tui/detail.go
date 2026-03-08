package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/models"
)

type detailModel struct {
	message *models.CaughtMessage
}

func newDetail() detailModel {
	return detailModel{}
}

func (d detailModel) Render(width int) string {
	if d.message == nil {
		return ""
	}

	m := d.message
	var b strings.Builder

	fields := []struct {
		label string
		value string
	}{
		{"Object ID", m.ObjectID},
		{"Stream ID", m.StreamID},
		{"Title", m.Title},
		{"Message", m.Message.Message},
		{"Severity", m.Severity},
		{"Author", m.Author},
		{"System", m.System},
		{"Timestamp", m.Timestamp},
		{"Caught At", m.CaughtAt.Format("2006-01-02 15:04:05")},
		{"Tags", m.Tags},
		{"Assignee", m.AssigneeName},
		{"Assignee Addr", m.AssigneeAddress},
		{"Artifacts", m.Artifacts},
		{"URL", m.Url},
	}

	for _, f := range fields {
		if f.value == "" {
			continue
		}
		label := detailLabelStyle.Render(fmt.Sprintf("  %-15s", f.label))
		value := detailValueStyle.Render(f.value)
		b.WriteString(label + " " + value + "\n")
	}

	return b.String()
}

var (
	detailLabelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00CC66")).
		Bold(true)

	detailValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))
)
