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

func TestList(t *testing.T) {
	s := New(100)
	for i := 0; i < 10; i++ {
		s.Add(newMsg(fmt.Sprintf("%d", i), fmt.Sprintf("msg-%d", i), "info", "ci"))
	}

	// pagination
	page := s.List(ListOptions{Offset: 2, Limit: 3})
	if len(page) != 3 {
		t.Fatalf("expected 3, got %d", len(page))
	}

	// no limit
	all := s.List(ListOptions{})
	if len(all) != 10 {
		t.Fatalf("expected 10, got %d", len(all))
	}
}

func TestListSorted(t *testing.T) {
	s := New(100)
	s.Add(newMsg("1", "Zebra", "info", "ci"))
	s.Add(newMsg("2", "Alpha", "error", "ci"))
	s.Add(newMsg("3", "Middle", "warning", "ci"))

	result := s.List(ListOptions{SortBy: SortByTitle, SortDir: SortAsc})
	if result[0].Title != "Alpha" || result[2].Title != "Zebra" {
		t.Fatalf("expected sorted asc by title, got %v", []string{result[0].Title, result[1].Title, result[2].Title})
	}

	result = s.List(ListOptions{SortBy: SortByTitle, SortDir: SortDesc})
	if result[0].Title != "Zebra" || result[2].Title != "Alpha" {
		t.Fatalf("expected sorted desc by title")
	}
}
