package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stuttgart-things/homerun2-core-catcher/internal/store"
)

// ExportHandler exports messages as JSON or CSV.
func ExportHandler(s *store.MessageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		format := r.URL.Query().Get("format")
		if format == "" {
			format = "json"
		}

		messages := s.List(store.ListOptions{
			SortBy:  store.SortByTimestamp,
			SortDir: store.SortDesc,
		})

		switch format {
		case "csv":
			w.Header().Set("Content-Type", "text/csv")
			w.Header().Set("Content-Disposition", "attachment; filename=messages.csv")

			writer := csv.NewWriter(w)
			writer.Write([]string{"ObjectID", "Title", "Message", "Severity", "Author", "System", "Timestamp", "Tags"})

			for _, m := range messages {
				writer.Write([]string{
					m.ObjectID,
					m.Title,
					m.Message.Message,
					m.Severity,
					m.Author,
					m.System,
					m.Timestamp,
					m.Tags,
				})
			}
			writer.Flush()

		default:
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Disposition", "attachment; filename=messages.json")
			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			if err := enc.Encode(messages); err != nil {
				http.Error(w, fmt.Sprintf("export error: %v", err), http.StatusInternalServerError)
			}
		}
	}
}
