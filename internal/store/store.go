package store

import (
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/stuttgart-things/homerun2-core-catcher/internal/models"
)

// SortField defines which field to sort by.
type SortField int

const (
	SortByTimestamp SortField = iota
	SortBySeverity
	SortBySystem
	SortByAuthor
	SortByTitle
)

// SortDirection defines sort order.
type SortDirection int

const (
	SortAsc SortDirection = iota
	SortDesc
)

// FilterOptions controls which messages to include.
type FilterOptions struct {
	System   string // exact match (empty = all)
	Severity string // exact match (empty = all)
	Author   string // exact match (empty = all)
	Since    time.Duration // only messages newer than this (0 = all)
	Query    string // free-text search across all fields
}

// ListOptions controls pagination, sorting, and filtering.
type ListOptions struct {
	Offset  int
	Limit   int
	SortBy  SortField
	SortDir SortDirection
	Filter  FilterOptions
}

// ListResult contains paginated results with total count.
type ListResult struct {
	Messages []models.CaughtMessage
	Total    int
}

// MessageStore holds caught messages in memory.
type MessageStore struct {
	mu       sync.RWMutex
	messages []models.CaughtMessage
	byID     map[string]int // objectID -> index
	maxSize  int
}

// New creates a MessageStore with the given capacity.
func New(maxSize int) *MessageStore {
	if maxSize <= 0 {
		maxSize = 10000
	}
	return &MessageStore{
		messages: make([]models.CaughtMessage, 0, maxSize),
		byID:     make(map[string]int, maxSize),
		maxSize:  maxSize,
	}
}

// Add appends a message, evicting the oldest if at capacity.
func (s *MessageStore) Add(msg models.CaughtMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.messages) >= s.maxSize {
		// evict oldest
		old := s.messages[0]
		delete(s.byID, old.ObjectID)
		s.messages = s.messages[1:]
		// rebuild index
		for i, m := range s.messages {
			s.byID[m.ObjectID] = i
		}
	}

	s.byID[msg.ObjectID] = len(s.messages)
	s.messages = append(s.messages, msg)
}

// List returns messages according to the given options, applying filters first.
func (s *MessageStore) List(opts ListOptions) ListResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filtered := s.applyFilters(opts.Filter)
	sortMessages(filtered, opts.SortBy, opts.SortDir)

	total := len(filtered)
	if opts.Offset >= total {
		return ListResult{Total: total}
	}

	end := opts.Offset + opts.Limit
	if opts.Limit <= 0 || end > total {
		end = total
	}

	result := make([]models.CaughtMessage, end-opts.Offset)
	copy(result, filtered[opts.Offset:end])
	return ListResult{Messages: result, Total: total}
}

// Search returns messages where any field contains the query string (case-insensitive).
func (s *MessageStore) Search(query string) []models.CaughtMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	q := strings.ToLower(query)
	var results []models.CaughtMessage

	for _, m := range s.messages {
		if matchesQuery(m, q) {
			results = append(results, m)
		}
	}
	return results
}

// Get returns a message by objectID.
func (s *MessageStore) Get(id string) (models.CaughtMessage, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idx, ok := s.byID[id]
	if !ok {
		return models.CaughtMessage{}, false
	}
	return s.messages[idx], true
}

// Count returns the number of stored messages.
func (s *MessageStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.messages)
}

// DistinctSystems returns all unique system values.
func (s *MessageStore) DistinctSystems() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.distinctField(func(m models.CaughtMessage) string { return m.System })
}

// DistinctSeverities returns all unique severity values.
func (s *MessageStore) DistinctSeverities() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.distinctField(func(m models.CaughtMessage) string { return m.Severity })
}

// DistinctAuthors returns all unique author values.
func (s *MessageStore) DistinctAuthors() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.distinctField(func(m models.CaughtMessage) string { return m.Author })
}

func (s *MessageStore) distinctField(extract func(models.CaughtMessage) string) []string {
	seen := make(map[string]struct{})
	for _, m := range s.messages {
		v := extract(m)
		if v != "" {
			seen[v] = struct{}{}
		}
	}
	result := make([]string, 0, len(seen))
	for v := range seen {
		result = append(result, v)
	}
	slices.Sort(result)
	return result
}

func (s *MessageStore) applyFilters(f FilterOptions) []models.CaughtMessage {
	var cutoff time.Time
	if f.Since > 0 {
		cutoff = time.Now().Add(-f.Since)
	}

	q := strings.ToLower(f.Query)

	result := make([]models.CaughtMessage, 0, len(s.messages))
	for _, m := range s.messages {
		if f.System != "" && m.System != f.System {
			continue
		}
		if f.Severity != "" && m.Severity != f.Severity {
			continue
		}
		if f.Author != "" && m.Author != f.Author {
			continue
		}
		if !cutoff.IsZero() && m.CaughtAt.Before(cutoff) {
			continue
		}
		if q != "" && !matchesQuery(m, q) {
			continue
		}
		result = append(result, m)
	}
	return result
}

func matchesQuery(m models.CaughtMessage, q string) bool {
	fields := []string{
		m.Title, m.Message.Message, m.Severity, m.Author,
		m.System, m.Tags, m.ObjectID, m.AssigneeName,
	}
	for _, f := range fields {
		if strings.Contains(strings.ToLower(f), q) {
			return true
		}
	}
	return false
}

func sortMessages(msgs []models.CaughtMessage, by SortField, dir SortDirection) {
	slices.SortFunc(msgs, func(a, b models.CaughtMessage) int {
		var va, vb string
		switch by {
		case SortBySeverity:
			va, vb = a.Severity, b.Severity
		case SortBySystem:
			va, vb = a.System, b.System
		case SortByAuthor:
			va, vb = a.Author, b.Author
		case SortByTitle:
			va, vb = a.Title, b.Title
		default: // SortByTimestamp
			va, vb = a.CaughtAt.String(), b.CaughtAt.String()
		}
		cmp := strings.Compare(va, vb)
		if dir == SortDesc {
			return -cmp
		}
		return cmp
	})
}
