package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMount(t *testing.T) {
	app := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("app " + r.URL.Path))
	})
	health := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("health"))
	})

	cases := []struct {
		basePath, path string
		code           int
		body, location string
	}{
		{"", "/", 200, "app /", ""},
		{"", "/messages", 200, "app /messages", ""},
		{"", "/health", 200, "health", ""},
		{"/core", "/core/", 200, "app /", ""},
		{"/core", "/core/messages/abc", 200, "app /messages/abc", ""},
		{"/core", "/core/static/htmx.min.js", 200, "app /static/htmx.min.js", ""},
		{"/core", "/health", 200, "health", ""},
		{"/core", "/core/health", 200, "app /health", ""},
		{"/core", "/core", 307, "", "/core/"},
		{"/core", "/", 307, "", "/core/"},
		{"/core", "/messages", 404, "", ""},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		Mount(tc.basePath, app, health).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if rec.Code != tc.code {
			t.Errorf("base %q GET %s: code %d, want %d", tc.basePath, tc.path, rec.Code, tc.code)
			continue
		}
		if tc.body != "" && rec.Body.String() != tc.body {
			t.Errorf("base %q GET %s: body %q, want %q", tc.basePath, tc.path, rec.Body.String(), tc.body)
		}
		if loc := rec.Header().Get("Location"); tc.location != "" && loc != tc.location {
			t.Errorf("base %q GET %s: Location %q, want %q", tc.basePath, tc.path, loc, tc.location)
		}
	}
}
