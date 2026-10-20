package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stuttgart-things/homerun2-core-catcher/internal/store"
)

const portalTokensLink = `<link rel="stylesheet" href="/static/homerun2-tokens.css">`

func TestThePortalTokensAreLinkedOnlyBehindThePortal(t *testing.T) {
	h := MessagesHandler(store.New(10), BuildInfo{}, "/core")

	alone := httptest.NewRecorder()
	h(alone, httptest.NewRequest(http.MethodGet, "/", http.NoBody))
	if strings.Contains(alone.Body.String(), portalTokensLink) {
		t.Error("the page links the portal's tokens without the portal")
	}

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	req.Header.Set("X-Forwarded-Prefix", "/core")
	behind := httptest.NewRecorder()
	h(behind, req)
	body := behind.Body.String()
	if i := strings.Index(body, portalTokensLink); i < 0 || i < strings.Index(body, "</style>") {
		t.Error("behind the portal the page does not link its tokens after its own copy")
	}
}

func TestThePageSwitchesToTheCompactHeaderInsideThePortal(t *testing.T) {
	rec := httptest.NewRecorder()
	MessagesHandler(store.New(10), BuildInfo{}, "")(rec, httptest.NewRequest(http.MethodGet, "/", http.NoBody))
	head, _, _ := strings.Cut(rec.Body.String(), "</head>")
	script := strings.Index(head, "classList.add('hr-embedded')")
	if script < 0 || script > strings.Index(head, "<style>") {
		t.Error("the hr-embedded script is missing or runs after the styles")
	}
	if !strings.Contains(head, ".hr-embedded .header-bar {") {
		t.Error("no compact header styles")
	}
}
