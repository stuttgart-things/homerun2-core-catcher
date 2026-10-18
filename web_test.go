package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	homerun "github.com/stuttgart-things/homerun-library/v4"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/handlers"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/models"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/store"
)

func testStore() *store.MessageStore {
	s := store.New(10)
	s.Add(models.CaughtMessage{
		Message:  homerun.Message{Title: "disk full", Severity: "error", System: "node-03", Author: "prometheus"},
		ObjectID: "obj-1",
		StreamID: "messages",
		CaughtAt: time.Now(),
	})
	return s
}

func get(t *testing.T, h http.Handler, path string) (int, string, http.Header) {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	body, _ := io.ReadAll(rec.Body)
	return rec.Code, string(body), rec.Header()
}

// Same-origin URLs in the pages: an attribute value starting with a single "/".
var rootURL = regexp.MustCompile(`(?:href|src|hx-get|hx-post|action)="(/[^/"][^"]*|/)"`)

func unprefixed(body, basePath string) []string {
	var out []string
	for _, m := range rootURL.FindAllStringSubmatch(body, -1) {
		if !strings.HasPrefix(m[1], basePath+"/") {
			out = append(out, m[1])
		}
	}
	return out
}

func TestWebUnderBasePath(t *testing.T) {
	h := webHandler(testStore(), handlers.BuildInfo{Version: "test"}, "/core")

	code, page, _ := get(t, h, "/core/")
	if code != http.StatusOK || !strings.Contains(page, "disk full") {
		t.Fatalf("GET /core/ = %d", code)
	}
	for _, want := range []string{
		`src="/core/static/htmx.min.js"`, `href="/core/export?format=csv"`,
		`hx-get="/core/messages"`, `hx-get="/core/messages/obj-1"`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("page lacks %s", want)
		}
	}
	if u := unprefixed(page, "/core"); len(u) > 0 {
		t.Errorf("page links outside /core/: %v", u)
	}

	_, table, _ := get(t, h, "/core/messages?sort=title")
	if !strings.Contains(table, `hx-get="/core/messages?sort=`) {
		t.Error("table partial does not carry the prefix")
	}
	if u := unprefixed(table, "/core"); len(u) > 0 {
		t.Errorf("table links outside /core/: %v", u)
	}

	for path, want := range map[string]string{
		"/core/messages/obj-1":     "disk full",
		"/core/export?format=json": "disk full",
		"/core/static/htmx.min.js": "htmx",
		"/core/health":             `"status":"healthy"`,
		"/health":                  `"status":"healthy"`, // probes and the portal, at the root
	} {
		if code, body, _ := get(t, h, path); code != http.StatusOK || !strings.Contains(body, want) {
			t.Errorf("GET %s = %d, body lacks %q", path, code, want)
		}
	}

	if code, _, hdr := get(t, h, "/"); code/100 != 3 || hdr.Get("Location") != "/core/" {
		t.Errorf("GET / = %d → %q, want a redirect to /core/", code, hdr.Get("Location"))
	}
	if code, _, _ := get(t, h, "/messages"); code != http.StatusNotFound {
		t.Errorf("GET /messages without the prefix = %d, want 404", code)
	}
}

func TestWebWithoutBasePathIsUnchanged(t *testing.T) {
	h := webHandler(testStore(), handlers.BuildInfo{}, "")
	code, page, _ := get(t, h, "/")
	if code != http.StatusOK {
		t.Fatalf("GET / = %d", code)
	}
	if !strings.Contains(page, `src="/static/htmx.min.js"`) || !strings.Contains(page, `hx-get="/messages"`) {
		t.Error("root URLs changed without a base path")
	}
	if code, _, _ := get(t, h, "/health"); code != http.StatusOK {
		t.Errorf("GET /health = %d", code)
	}
}
