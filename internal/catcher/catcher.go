package catcher

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	homerun "github.com/stuttgart-things/homerun-library/v2"
	"github.com/stuttgart-things/redisqueue"
)

// Catcher defines the interface for message consumption backends.
type Catcher interface {
	Run()
	Shutdown()
	Errors() <-chan error
}

// RedisCatcher consumes messages from Redis Streams and logs them.
type RedisCatcher struct {
	consumer *redisqueue.Consumer
	stream   string
}

// NewRedisCatcher creates a consumer connected to the given Redis stream.
func NewRedisCatcher(rc homerun.RedisConfig, groupName, consumerName string) (*RedisCatcher, error) {
	if consumerName == "" {
		hostname, _ := os.Hostname()
		consumerName = hostname
	}

	addr := fmt.Sprintf("%s:%s", rc.Addr, rc.Port)

	consumer, err := redisqueue.NewConsumerWithOptions(&redisqueue.ConsumerOptions{
		Name:       consumerName,
		GroupName:  groupName,
		BufferSize: 100,
		Concurrency: 10,
		RedisOptions: &redisqueue.RedisOptions{
			Addr:     addr,
			Password: rc.Password,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create redis consumer: %w", err)
	}

	c := &RedisCatcher{
		consumer: consumer,
		stream:   rc.Stream,
	}

	consumer.Register(rc.Stream, c.handleMessage)

	return c, nil
}

// Run starts the consumer. Blocks until Shutdown is called.
func (c *RedisCatcher) Run() {
	c.consumer.Run()
}

// Shutdown gracefully stops the consumer.
func (c *RedisCatcher) Shutdown() {
	c.consumer.Shutdown()
}

// Errors returns the consumer's error channel.
func (c *RedisCatcher) Errors() <-chan error {
	return c.consumer.Errors
}

// handleMessage is called for each message received from the stream.
func (c *RedisCatcher) handleMessage(msg *redisqueue.Message) error {
	// Log the raw stream message
	slog.Info("message caught",
		"stream", msg.Stream,
		"id", msg.ID,
	)

	// The stream entry contains a messageID field referencing the Redis JSON object
	if messageID, ok := msg.Values["messageID"]; ok {
		slog.Info("message reference",
			"messageID", messageID,
			"stream_id", msg.ID,
		)
	}

	// Log all values from the stream entry
	valuesJSON, err := json.Marshal(msg.Values)
	if err != nil {
		slog.Warn("failed to marshal message values", "error", err)
	} else {
		slog.Info("message values", "data", string(valuesJSON))
	}

	return nil
}
