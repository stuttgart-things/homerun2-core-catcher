package store

import (
	"strings"
	"sync"

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

// ListOptions controls pagination and sorting.
type ListOptions struct {
	Offset    int
	Limit     int
	SortBy    SortField
	SortDir   SortDirection
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

// List returns messages according to the given options.
func (s *MessageStore) List(opts ListOptions) []models.CaughtMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sorted := s.sorted(opts.SortBy, opts.SortDir)

	if opts.Offset >= len(sorted) {
		return nil
	}

	end := opts.Offset + opts.Limit
	if opts.Limit <= 0 || end > len(sorted) {
		end = len(sorted)
	}

	result := make([]models.CaughtMessage, end-opts.Offset)
	copy(result, sorted[opts.Offset:end])
	return result
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

func (s *MessageStore) sorted(by SortField, dir SortDirection) []models.CaughtMessage {
	cp := make([]models.CaughtMessage, len(s.messages))
	copy(cp, s.messages)

	less := func(i, j int) bool {
		var a, b string
		switch by {
		case SortBySeverity:
			a, b = cp[i].Severity, cp[j].Severity
		case SortBySystem:
			a, b = cp[i].System, cp[j].System
		case SortByAuthor:
			a, b = cp[i].Author, cp[j].Author
		case SortByTitle:
			a, b = cp[i].Title, cp[j].Title
		default: // SortByTimestamp
			a, b = cp[i].CaughtAt.String(), cp[j].CaughtAt.String()
		}
		if dir == SortDesc {
			return a > b
		}
		return a < b
	}

	// simple insertion sort (fine for bounded store)
	for i := 1; i < len(cp); i++ {
		for j := i; j > 0 && less(j, j-1); j-- {
			cp[j], cp[j-1] = cp[j-1], cp[j]
		}
	}
	return cp
}
