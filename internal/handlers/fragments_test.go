package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	homerun "github.com/stuttgart-things/homerun-library/v4"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/models"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/store"
)

func TestManifestListsTheMessagesFragment(t *testing.T) {
	rec := httptest.NewRecorder()
	ManifestHandler(BuildInfo{Version: "v1.2.3"}, "/core")(rec, httptest.NewRequest(http.MethodGet, "/ui/manifest.json", nil))

	var m Manifest
	if err := json.NewDecoder(rec.Body).Decode(&m); err != nil {
		t.Fatal(err)
	}
	if m.Version != "v1.2.3" || m.BasePath != "/core" || len(m.Fragments) != 1 {
		t.Fatalf("manifest = %+v", m)
	}
	if f := m.Fragments[0]; f.ID != "messages" || f.URL != "/core/ui/fragment/messages" || f.RefreshSeconds < 2 {
		t.Errorf("fragment = %+v", f)
	}
}

func TestMessagesFragment(t *testing.T) {
	s := store.New(20)
	now := time.Now()
	for i := range 10 {
		s.Add(models.CaughtMessage{
			Message:  homerun.Message{Title: "old", Severity: "info", System: "sys"},
			ObjectID: string(rune('a' + i)),
			CaughtAt: now.Add(-time.Hour),
		})
	}
	s.Add(models.CaughtMessage{
		Message:  homerun.Message{Title: "<b>disk</b> full", Severity: "critical", System: "nas"},
		ObjectID: "new",
		CaughtAt: now,
	})

	rec := httptest.NewRecorder()
	MessagesFragmentHandler(s, "/core")(rec, httptest.NewRequest(http.MethodGet, "/ui/fragment/messages", nil))
	body := rec.Body.String()

	if strings.Contains(body, "<html") || strings.Contains(body, "<script") {
		t.Fatal("a fragment is bare HTML")
	}
	if n := strings.Count(body, `<li class="hr-item">`); n != messagesFragmentLimit {
		t.Errorf("%d items, want %d", n, messagesFragmentLimit)
	}
	for _, want := range []string{
		`<span class="hr-sev hr-sev-error">CRITICAL</span>`,
		"&lt;b&gt;disk&lt;/b&gt; full", // titles come from messages and are escaped
		"11 in memory",
		`href="/core/"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("fragment lacks %s:\n%s", want, body)
		}
	}
	if i := strings.Index(body, "disk"); i < 0 || i > strings.Index(body, "old") {
		t.Error("the newest message is not first")
	}
}

func TestMessagesFragmentEmpty(t *testing.T) {
	rec := httptest.NewRecorder()
	MessagesFragmentHandler(store.New(10), "")(rec, httptest.NewRequest(http.MethodGet, "/ui/fragment/messages", nil))
	if body := rec.Body.String(); !strings.Contains(body, "No messages yet.") || !strings.Contains(body, `href="/"`) {
		t.Errorf("empty fragment = %s", body)
	}
}

func TestHRSeverity(t *testing.T) {
	for in, want := range map[string]string{
		"ERROR": "error", "critical": "error", "warning": "warning",
		"success": "success", "debug": "debug", "info": "info", "": "info", "bogus": "info",
	} {
		if got := hrSeverity(in); got != want {
			t.Errorf("hrSeverity(%q) = %q, want %q", in, got, want)
		}
	}
}
