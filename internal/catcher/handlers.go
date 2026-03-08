package catcher

import (
	"log/slog"

	"github.com/stuttgart-things/homerun2-core-catcher/internal/models"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/store"
)

// LogHandler returns a MessageHandler that logs messages with severity-aware levels.
func LogHandler() MessageHandler {
	return func(msg models.CaughtMessage) {
		level := severityToLevel(msg.Severity)

		slog.Log(nil, level, "message caught",
			"objectId", msg.ObjectID,
			"streamId", msg.StreamID,
			"title", msg.Title,
			"message", msg.Message.Message,
			"severity", msg.Severity,
			"author", msg.Author,
			"system", msg.System,
			"timestamp", msg.Timestamp,
			"tags", msg.Tags,
		)
	}
}

// StoreHandler returns a MessageHandler that stores messages in the given store.
func StoreHandler(s *store.MessageStore) MessageHandler {
	return func(msg models.CaughtMessage) {
		s.Add(msg)
	}
}

func severityToLevel(severity string) slog.Level {
	switch severity {
	case "error":
		return slog.LevelError
	case "warning":
		return slog.LevelWarn
	case "debug":
		return slog.LevelDebug
	default: // info, success, etc.
		return slog.LevelInfo
	}
}
