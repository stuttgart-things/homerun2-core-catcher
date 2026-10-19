package handlers

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	homerun "github.com/stuttgart-things/homerun-library/v4"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/models"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/store"
)

// localURL matches every same-origin URL the templates hand to the browser.
var localURL = regexp.MustCompile(`(?:href|src|hx-get)="(/[^"]*)"`)

func TestTemplatesCarryBasePath(t *testing.T) {
	s := store.New(10)
	s.Add(models.CaughtMessage{
		Message:  homerun.Message{Title: "t", Severity: "info", System: "sys"},
		ObjectID: "abc",
		CaughtAt: time.Now(),
	})

	for _, tc := range []struct {
		name    string
		handler http.HandlerFunc
	}{
		{"page", MessagesHandler(s, BuildInfo{}, "/core")},
		{"table", MessagesTableHandler(s, BuildInfo{}, "/core")},
	} {
		rec := httptest.NewRecorder()
		tc.handler(rec, httptest.NewRequest(http.MethodGet, "/", nil))
		body := rec.Body.String()

		urls := localURL.FindAllStringSubmatch(body, -1)
		if len(urls) == 0 {
			t.Fatalf("%s: no local URLs rendered", tc.name)
		}
		for _, m := range urls {
			if !strings.HasPrefix(m[1], "/core/") {
				t.Errorf("%s: URL %q lacks the base path", tc.name, m[1])
			}
		}
		if !strings.Contains(body, `hx-get="/core/messages/abc"`) {
			t.Errorf("%s: row detail URL missing or unprefixed", tc.name)
		}
	}
}
