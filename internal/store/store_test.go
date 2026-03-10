package store

import (
	"fmt"
	"testing"
	"time"

	homerun "github.com/stuttgart-things/homerun-library/v2"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/models"
)

func newMsg(id, title, severity, system string) models.CaughtMessage {
	return models.CaughtMessage{
		Message: homerun.Message{
			Title:    title,
			Severity: severity,
			System:   system,
			Author:   "test",
		},
		ObjectID: id,
		StreamID: "messages",
		CaughtAt: time.Now(),
	}
}

func TestAddAndCount(t *testing.T) {
	s := New(100)
	s.Add(newMsg("1", "first", "info", "ci"))
	s.Add(newMsg("2", "second", "error", "ci"))

	if s.Count() != 2 {
		t.Fatalf("expected 2, got %d", s.Count())
	}
}

func TestGet(t *testing.T) {
	s := New(100)
	s.Add(newMsg("abc", "test msg", "info", "ci"))

	m, ok := s.Get("abc")
	if !ok {
		t.Fatal("expected to find message")
	}
	if m.Title != "test msg" {
		t.Fatalf("expected 'test msg', got %q", m.Title)
	}

	_, ok = s.Get("nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestEviction(t *testing.T) {
	s := New(3)
	s.Add(newMsg("1", "first", "info", "a"))
	s.Add(newMsg("2", "second", "info", "b"))
	s.Add(newMsg("3", "third", "info", "c"))
	s.Add(newMsg("4", "fourth", "info", "d"))

	if s.Count() != 3 {
		t.Fatalf("expected 3 after eviction, got %d", s.Count())
	}

	_, ok := s.Get("1")
	if ok {
		t.Fatal("oldest message should have been evicted")
	}

	m, ok := s.Get("4")
	if !ok || m.Title != "fourth" {
		t.Fatal("newest message should exist")
	}
}

func TestSearch(t *testing.T) {
	s := New(100)
	s.Add(newMsg("1", "Deploy complete", "success", "argocd"))
	s.Add(newMsg("2", "Build failed", "error", "jenkins"))
	s.Add(newMsg("3", "Config update", "info", "argocd"))

	results := s.Search("argocd")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	results = s.Search("FAILED")
	if len(results) != 1 {
		t.Fatalf("expected 1 result for case-insensitive search, got %d", len(results))
	}
}

func newMsgWithAuthor(id, title, severity, system, author string, caughtAt time.Time) models.CaughtMessage {
	return models.CaughtMessage{
		Message: homerun.Message{
			Title:    title,
			Severity: severity,
			System:   system,
			Author:   author,
		},
		ObjectID: id,
		StreamID: "messages",
		CaughtAt: caughtAt,
	}
}

func TestList(t *testing.T) {
	s := New(100)
	for i := 0; i < 10; i++ {
		s.Add(newMsg(fmt.Sprintf("%d", i), fmt.Sprintf("msg-%d", i), "info", "ci"))
	}

	// pagination
	result := s.List(ListOptions{Offset: 2, Limit: 3})
	if len(result.Messages) != 3 {
		t.Fatalf("expected 3, got %d", len(result.Messages))
	}
	if result.Total != 10 {
		t.Fatalf("expected total 10, got %d", result.Total)
	}

	// no limit
	all := s.List(ListOptions{})
	if len(all.Messages) != 10 {
		t.Fatalf("expected 10, got %d", len(all.Messages))
	}
}

func TestListSorted(t *testing.T) {
	s := New(100)
	s.Add(newMsg("1", "Zebra", "info", "ci"))
	s.Add(newMsg("2", "Alpha", "error", "ci"))
	s.Add(newMsg("3", "Middle", "warning", "ci"))

	result := s.List(ListOptions{SortBy: SortByTitle, SortDir: SortAsc})
	if result.Messages[0].Title != "Alpha" || result.Messages[2].Title != "Zebra" {
		t.Fatalf("expected sorted asc by title, got %v", []string{result.Messages[0].Title, result.Messages[1].Title, result.Messages[2].Title})
	}

	result = s.List(ListOptions{SortBy: SortByTitle, SortDir: SortDesc})
	if result.Messages[0].Title != "Zebra" || result.Messages[2].Title != "Alpha" {
		t.Fatalf("expected sorted desc by title")
	}
}

func TestFilterBySystem(t *testing.T) {
	s := New(100)
	s.Add(newMsg("1", "msg1", "info", "movie-scripts"))
	s.Add(newMsg("2", "msg2", "error", "labul-staging"))
	s.Add(newMsg("3", "msg3", "info", "movie-scripts"))

	result := s.List(ListOptions{Filter: FilterOptions{System: "movie-scripts"}})
	if result.Total != 2 {
		t.Fatalf("expected 2 filtered by system, got %d", result.Total)
	}
	for _, m := range result.Messages {
		if m.System != "movie-scripts" {
			t.Fatalf("expected system movie-scripts, got %q", m.System)
		}
	}
}

func TestFilterBySeverity(t *testing.T) {
	s := New(100)
	s.Add(newMsg("1", "msg1", "info", "ci"))
	s.Add(newMsg("2", "msg2", "error", "ci"))
	s.Add(newMsg("3", "msg3", "error", "ci"))
	s.Add(newMsg("4", "msg4", "warning", "ci"))

	result := s.List(ListOptions{Filter: FilterOptions{Severity: "error"}})
	if result.Total != 2 {
		t.Fatalf("expected 2 errors, got %d", result.Total)
	}
}

func TestFilterByAuthor(t *testing.T) {
	s := New(100)
	now := time.Now()
	s.Add(newMsgWithAuthor("1", "msg1", "info", "ci", "k8s-pitcher", now))
	s.Add(newMsgWithAuthor("2", "msg2", "info", "ci", "omni-pitcher", now))
	s.Add(newMsgWithAuthor("3", "msg3", "info", "ci", "k8s-pitcher", now))

	result := s.List(ListOptions{Filter: FilterOptions{Author: "k8s-pitcher"}})
	if result.Total != 2 {
		t.Fatalf("expected 2 by author, got %d", result.Total)
	}
}

func TestFilterByTime(t *testing.T) {
	s := New(100)
	now := time.Now()
	s.Add(newMsgWithAuthor("1", "old", "info", "ci", "test", now.Add(-2*time.Hour)))
	s.Add(newMsgWithAuthor("2", "recent", "info", "ci", "test", now.Add(-30*time.Minute)))
	s.Add(newMsgWithAuthor("3", "newest", "info", "ci", "test", now.Add(-5*time.Minute)))

	result := s.List(ListOptions{Filter: FilterOptions{Since: time.Hour}})
	if result.Total != 2 {
		t.Fatalf("expected 2 within last hour, got %d", result.Total)
	}
}

func TestFilterCombined(t *testing.T) {
	s := New(100)
	now := time.Now()
	s.Add(newMsgWithAuthor("1", "msg1", "error", "movie-scripts", "k8s-pitcher", now.Add(-30*time.Minute)))
	s.Add(newMsgWithAuthor("2", "msg2", "info", "movie-scripts", "k8s-pitcher", now.Add(-30*time.Minute)))
	s.Add(newMsgWithAuthor("3", "msg3", "error", "labul", "k8s-pitcher", now.Add(-30*time.Minute)))
	s.Add(newMsgWithAuthor("4", "msg4", "error", "movie-scripts", "omni", now.Add(-30*time.Minute)))
	s.Add(newMsgWithAuthor("5", "msg5", "error", "movie-scripts", "k8s-pitcher", now.Add(-2*time.Hour)))

	result := s.List(ListOptions{
		Filter: FilterOptions{
			System:   "movie-scripts",
			Severity: "error",
			Author:   "k8s-pitcher",
			Since:    time.Hour,
		},
	})
	if result.Total != 1 {
		t.Fatalf("expected 1 combined filter match, got %d", result.Total)
	}
	if result.Messages[0].ObjectID != "1" {
		t.Fatalf("expected message 1, got %q", result.Messages[0].ObjectID)
	}
}

func TestFilterWithQuery(t *testing.T) {
	s := New(100)
	s.Add(newMsg("1", "Deploy complete", "success", "movie-scripts"))
	s.Add(newMsg("2", "Deploy failed", "error", "movie-scripts"))
	s.Add(newMsg("3", "Build ok", "success", "movie-scripts"))

	result := s.List(ListOptions{Filter: FilterOptions{Query: "deploy", Severity: "error"}})
	if result.Total != 1 {
		t.Fatalf("expected 1 result for query+severity filter, got %d", result.Total)
	}
}

func TestDistinctSystems(t *testing.T) {
	s := New(100)
	s.Add(newMsg("1", "msg1", "info", "movie-scripts"))
	s.Add(newMsg("2", "msg2", "info", "labul"))
	s.Add(newMsg("3", "msg3", "info", "movie-scripts"))

	systems := s.DistinctSystems()
	if len(systems) != 2 {
		t.Fatalf("expected 2 systems, got %d", len(systems))
	}
	if systems[0] != "labul" || systems[1] != "movie-scripts" {
		t.Fatalf("expected [labul, movie-scripts], got %v", systems)
	}
}

func TestDistinctSeverities(t *testing.T) {
	s := New(100)
	s.Add(newMsg("1", "msg1", "info", "ci"))
	s.Add(newMsg("2", "msg2", "error", "ci"))
	s.Add(newMsg("3", "msg3", "info", "ci"))

	sev := s.DistinctSeverities()
	if len(sev) != 2 {
		t.Fatalf("expected 2 severities, got %d: %v", len(sev), sev)
	}
}

func TestFilterWithPagination(t *testing.T) {
	s := New(100)
	for i := 0; i < 10; i++ {
		s.Add(newMsg(fmt.Sprintf("%d", i), fmt.Sprintf("msg-%d", i), "error", "ci"))
	}
	for i := 10; i < 20; i++ {
		s.Add(newMsg(fmt.Sprintf("%d", i), fmt.Sprintf("msg-%d", i), "info", "ci"))
	}

	result := s.List(ListOptions{
		Offset: 0,
		Limit:  5,
		Filter: FilterOptions{Severity: "error"},
	})
	if result.Total != 10 {
		t.Fatalf("expected total 10 errors, got %d", result.Total)
	}
	if len(result.Messages) != 5 {
		t.Fatalf("expected 5 on page, got %d", len(result.Messages))
	}
}
