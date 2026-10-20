package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/stuttgart-things/homerun2-core-catcher/internal/models"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/store"
)

// Micro-frontend fragments for homerun2-portal. The portal reads
// <base>/ui/manifest.json to learn which fragments a service offers and loads
// each <base>/ui/fragment/<id> with htmx into its overview. A fragment is bare
// HTML styled only by the portal's hr-* classes (static/homerun2.css in
// homerun2-portal); every URL in it carries the base path, because it is
// served on the portal's origin.

const (
	messagesFragmentLimit   = 8
	messagesFragmentRefresh = 5
)

// Fragment is one entry of the manifest.
type Fragment struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	RefreshSeconds int    `json:"refreshSeconds"`
}

// Manifest lists the fragments the service offers to the portal.
type Manifest struct {
	Name      string     `json:"name"`
	Title     string     `json:"title"`
	Version   string     `json:"version"`
	BasePath  string     `json:"basePath"`
	Fragments []Fragment `json:"fragments"`
}

// ManifestHandler serves GET /ui/manifest.json.
func ManifestHandler(info BuildInfo, basePath string) http.HandlerFunc {
	m := Manifest{
		Name:     "homerun2-core-catcher",
		Title:    "Core Catcher",
		Version:  info.Version,
		BasePath: basePath,
		Fragments: []Fragment{{
			ID:             "messages",
			Title:          "Latest messages",
			URL:            basePath + "/ui/fragment/messages",
			RefreshSeconds: messagesFragmentRefresh,
		}},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(m)
	}
}

// html/template escapes every field: titles and systems come from messages,
// and the portal shows the fragment to every viewer.
var messagesFragment = template.Must(template.New("messages").Funcs(template.FuncMap{
	"hrSeverity": hrSeverity,
	"upper":      strings.ToUpper,
}).Parse(`<div class="hr-frag"><ul class="hr-list">
{{- range .Messages}}<li class="hr-item"><span class="hr-time">{{.CaughtAt.Format "15:04:05"}}</span><span class="hr-sev hr-sev-{{hrSeverity .Severity}}">{{upper .Severity}}</span><span class="hr-system">{{.System}}</span><span class="hr-title">{{.Title}}</span></li>
{{- else}}<li class="hr-item hr-muted">No messages yet.</li>{{end -}}
</ul><span class="hr-muted">{{.Total}} in memory</span><a class="hr-frag-more" href="{{.BasePath}}/">Open dashboard →</a></div>`))

// MessagesFragmentHandler serves GET /ui/fragment/messages: the newest
// messages as an hr-list.
func MessagesFragmentHandler(s *store.MessageStore, basePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := s.List(store.ListOptions{
			Limit:   messagesFragmentLimit,
			SortBy:  store.SortByTimestamp,
			SortDir: store.SortDesc,
		})
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = messagesFragment.Execute(w, struct {
			Messages []models.CaughtMessage
			Total    int
			BasePath string
		}{result.Messages, s.Count(), basePath})
	}
}

// hrSeverity maps a severity to the portal's hr-sev-* classes.
func hrSeverity(severity string) string {
	switch s := strings.ToLower(severity); s {
	case "error", "critical":
		return "error"
	case "warning", "success", "debug":
		return s
	default:
		return "info"
	}
}
