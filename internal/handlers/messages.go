package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/stuttgart-things/homerun2-core-catcher/internal/store"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/web"
)

var templates *template.Template

func init() {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"seq": func(n int) []int {
			s := make([]int, n)
			for i := range s {
				s[i] = i
			}
			return s
		},
		"severityClass": func(s string) string {
			switch s {
			case "error":
				return "severity-error"
			case "warning":
				return "severity-warning"
			case "success":
				return "severity-success"
			case "debug":
				return "severity-debug"
			default:
				return "severity-info"
			}
		},
	}
	templates = template.Must(template.New("").Funcs(funcMap).ParseFS(web.TemplateFS, "templates/*.html"))
}

// MessagesHandler serves the main dashboard page.
func MessagesHandler(s *store.MessageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := buildTableData(s, r)
		templates.ExecuteTemplate(w, "index.html", data)
	}
}

// MessagesTableHandler serves the HTMX partial for the message table.
func MessagesTableHandler(s *store.MessageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := buildTableData(s, r)
		templates.ExecuteTemplate(w, "table.html", data)
	}
}

// MessageDetailHandler serves the HTMX partial for a single message detail.
func MessageDetailHandler(s *store.MessageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/messages/")
		if id == "" {
			http.Error(w, "missing message ID", http.StatusBadRequest)
			return
		}

		msg, ok := s.Get(id)
		if !ok {
			http.Error(w, "message not found", http.StatusNotFound)
			return
		}

		templates.ExecuteTemplate(w, "detail.html", msg)
	}
}

type tableData struct {
	Messages   interface{}
	Query      string
	SortBy     string
	SortDir    string
	Page       int
	PageSize   int
	Total      int
	TotalAll   int
	TotalPages int
	HasPrev    bool
	HasNext    bool
	// Filter values
	FilterSystem   string
	FilterSeverity string
	FilterAuthor   string
	FilterTime     string
	// Distinct values for dropdowns
	Systems    []string
	Severities []string
	Authors    []string
}

// parseSinceDuration converts a time filter string to a duration.
func parseSinceDuration(s string) time.Duration {
	switch s {
	case "1h":
		return time.Hour
	case "6h":
		return 6 * time.Hour
	case "24h":
		return 24 * time.Hour
	case "7d":
		return 7 * 24 * time.Hour
	default:
		return 0
	}
}

func buildTableData(s *store.MessageStore, r *http.Request) tableData {
	q := r.URL.Query()

	query := q.Get("q")
	sortBy := q.Get("sort")
	sortDir := q.Get("dir")
	pageStr := q.Get("page")
	pageSizeStr := q.Get("size")

	// Filter params
	filterSystem := q.Get("system")
	filterSeverity := q.Get("severity")
	filterAuthor := q.Get("author")
	filterTime := q.Get("time")

	page, _ := strconv.Atoi(pageStr)
	if page < 0 {
		page = 0
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize <= 0 {
		pageSize = 25
	}

	if sortDir == "" {
		sortDir = "desc"
	}

	var sortField store.SortField
	switch sortBy {
	case "severity":
		sortField = store.SortBySeverity
	case "system":
		sortField = store.SortBySystem
	case "author":
		sortField = store.SortByAuthor
	case "title":
		sortField = store.SortByTitle
	default:
		sortField = store.SortByTimestamp
		sortBy = "timestamp"
	}

	dir := store.SortDesc
	if sortDir == "asc" {
		dir = store.SortAsc
	}

	result := s.List(store.ListOptions{
		Offset:  page * pageSize,
		Limit:   pageSize,
		SortBy:  sortField,
		SortDir: dir,
		Filter: store.FilterOptions{
			System:   filterSystem,
			Severity: filterSeverity,
			Author:   filterAuthor,
			Since:    parseSinceDuration(filterTime),
			Query:    query,
		},
	})

	total := result.Total
	totalPages := 1
	if total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}

	return tableData{
		Messages:       result.Messages,
		Query:          query,
		SortBy:         sortBy,
		SortDir:        sortDir,
		Page:           page,
		PageSize:       pageSize,
		Total:          total,
		TotalAll:       s.Count(),
		TotalPages:     totalPages,
		HasPrev:        page > 0,
		HasNext:        page < totalPages-1,
		FilterSystem:   filterSystem,
		FilterSeverity: filterSeverity,
		FilterAuthor:   filterAuthor,
		FilterTime:     filterTime,
		Systems:        s.DistinctSystems(),
		Severities:     s.DistinctSeverities(),
		Authors:        s.DistinctAuthors(),
	}
}
