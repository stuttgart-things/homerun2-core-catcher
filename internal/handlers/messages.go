package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

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
	Messages  interface{}
	Query     string
	SortBy    string
	SortDir   string
	Page      int
	PageSize  int
	Total     int
	TotalPages int
	HasPrev   bool
	HasNext   bool
}

func buildTableData(s *store.MessageStore, r *http.Request) tableData {
	q := r.URL.Query()

	query := q.Get("q")
	sortBy := q.Get("sort")
	sortDir := q.Get("dir")
	pageStr := q.Get("page")
	pageSizeStr := q.Get("size")

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

	var messages interface{}
	var total int

	if query != "" {
		all := s.Search(query)
		total = len(all)
		start := page * pageSize
		if start >= total {
			messages = nil
		} else {
			end := start + pageSize
			if end > total {
				end = total
			}
			slice := all[start:end]
			messages = slice
		}
	} else {
		total = s.Count()
		messages = s.List(store.ListOptions{
			Offset:  page * pageSize,
			Limit:   pageSize,
			SortBy:  sortField,
			SortDir: dir,
		})
	}

	totalPages := 1
	if total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}

	// Toggle sort direction for column header links
	nextDir := "asc"
	if sortDir == "asc" {
		nextDir = "desc"
	}
	_ = nextDir

	return tableData{
		Messages:   messages,
		Query:      query,
		SortBy:     sortBy,
		SortDir:    fmt.Sprintf("%s", sortDir),
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasPrev:    page > 0,
		HasNext:    page < totalPages-1,
	}
}
