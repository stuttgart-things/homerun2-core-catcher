package catcher

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	homerun "github.com/stuttgart-things/homerun-library/v2"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/models"
)

// FileCatcher reads messages from a JSON file and replays them.
type FileCatcher struct {
	path     string
	interval time.Duration
	handlers []MessageHandler
	quit     chan struct{}
	done     chan struct{}
	once     sync.Once
}

// NewFileCatcher creates a FileCatcher that replays messages from a JSON file.
func NewFileCatcher(path string, interval time.Duration, handlers ...MessageHandler) *FileCatcher {
	if interval <= 0 {
		interval = time.Second
	}
	return &FileCatcher{
		path:     path,
		interval: interval,
		handlers: handlers,
		quit:     make(chan struct{}),
		done:     make(chan struct{}),
	}
}

// Run reads messages from the file and dispatches them to handlers. Blocks until Shutdown.
func (f *FileCatcher) Run() {
	defer close(f.done)

	messages, err := f.loadMessages()
	if err != nil {
		slog.Error("failed to load messages from file", "path", f.path, "error", err)
		return
	}

	slog.Info("file catcher loaded messages", "count", len(messages), "path", f.path)

	for i, msg := range messages {
		select {
		case <-f.quit:
			return
		default:
		}

		caught := models.CaughtMessage{
			Message:  msg,
			ObjectID: fmt.Sprintf("file-%d-%s", i, msg.System),
			StreamID: "file",
			CaughtAt: time.Now(),
		}

		for _, h := range f.handlers {
			h(caught)
		}

		if i < len(messages)-1 {
			select {
			case <-f.quit:
				return
			case <-time.After(f.interval):
			}
		}
	}

	// Keep running until shutdown
	<-f.quit
}

// Shutdown stops the file catcher.
func (f *FileCatcher) Shutdown() {
	f.once.Do(func() {
		close(f.quit)
	})
	<-f.done
}

// Errors returns a nil channel (FileCatcher doesn't produce errors this way).
func (f *FileCatcher) Errors() <-chan error {
	return nil
}

func (f *FileCatcher) loadMessages() ([]homerun.Message, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var messages []homerun.Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return messages, nil
}
