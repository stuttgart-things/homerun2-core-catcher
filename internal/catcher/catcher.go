package catcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
	homerun "github.com/stuttgart-things/homerun-library/v2"
	"github.com/stuttgart-things/redisqueue"
)

// Catcher defines the interface for message consumption backends.
type Catcher interface {
	Run()
	Shutdown()
	Errors() <-chan error
}

// RedisCatcher consumes messages from Redis Streams and resolves full payloads from Redis JSON.
type RedisCatcher struct {
	consumer    *redisqueue.Consumer
	redisClient *redis.Client
	stream      string
}

// NewRedisCatcher creates a consumer connected to the given Redis stream.
func NewRedisCatcher(rc homerun.RedisConfig, groupName, consumerName string) (*RedisCatcher, error) {
	if consumerName == "" {
		hostname, _ := os.Hostname()
		consumerName = hostname
	}

	addr := fmt.Sprintf("%s:%s", rc.Addr, rc.Port)

	consumer, err := redisqueue.NewConsumerWithOptions(&redisqueue.ConsumerOptions{
		Name:        consumerName,
		GroupName:   groupName,
		BufferSize:  100,
		Concurrency: 10,
		RedisOptions: &redisqueue.RedisOptions{
			Addr:     addr,
			Password: rc.Password,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create redis consumer: %w", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: rc.Password,
	})

	c := &RedisCatcher{
		consumer:    consumer,
		redisClient: redisClient,
		stream:      rc.Stream,
	}

	consumer.Register(rc.Stream, c.handleMessage)

	return c, nil
}

// Run starts the consumer. Blocks until Shutdown is called.
func (c *RedisCatcher) Run() {
	c.consumer.Run()
}

// Shutdown gracefully stops the consumer and closes the Redis client.
func (c *RedisCatcher) Shutdown() {
	c.consumer.Shutdown()
	if c.redisClient != nil {
		c.redisClient.Close()
	}
}

// Errors returns the consumer's error channel.
func (c *RedisCatcher) Errors() <-chan error {
	return c.consumer.Errors
}

// handleMessage is called for each message received from the stream.
func (c *RedisCatcher) handleMessage(msg *redisqueue.Message) error {
	slog.Info("message caught",
		"stream", msg.Stream,
		"id", msg.ID,
	)

	messageID, ok := msg.Values["messageID"]
	if !ok {
		slog.Warn("stream entry missing messageID field", "id", msg.ID)
		return nil
	}

	messageIDStr, ok := messageID.(string)
	if !ok {
		slog.Warn("messageID is not a string", "id", msg.ID, "messageID", messageID)
		return nil
	}

	slog.Info("message reference",
		"messageID", messageIDStr,
		"stream_id", msg.ID,
	)

	// Resolve full message payload from Redis JSON
	payload, err := c.resolveMessage(messageIDStr)
	if err != nil {
		slog.Warn("failed to resolve message from Redis JSON",
			"messageID", messageIDStr,
			"error", err,
		)
		return nil
	}

	slog.Info("message payload",
		"messageID", messageIDStr,
		"title", payload.Title,
		"message", payload.Message,
		"severity", payload.Severity,
		"author", payload.Author,
		"system", payload.System,
		"timestamp", payload.Timestamp,
		"tags", payload.Tags,
	)

	return nil
}

// resolveMessage fetches the full message payload from Redis JSON using JSON.GET.
func (c *RedisCatcher) resolveMessage(messageID string) (*homerun.Message, error) {
	ctx := context.Background()

	result, err := c.redisClient.Do(ctx, "JSON.GET", messageID, ".").Text()
	if err != nil {
		return nil, fmt.Errorf("JSON.GET %s: %w", messageID, err)
	}

	var msg homerun.Message
	if err := json.Unmarshal([]byte(result), &msg); err != nil {
		return nil, fmt.Errorf("unmarshal message: %w", err)
	}

	return &msg, nil
}
