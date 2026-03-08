package main

import (
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/stuttgart-things/homerun2-core-catcher/internal/banner"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/catcher"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/config"

	homerun "github.com/stuttgart-things/homerun-library/v2"
)

// Build-time variables set via ldflags
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	banner.Show()
	config.SetupLogging()

	slog.Info("starting homerun2-core-catcher",
		"version", version,
		"commit", commit,
		"date", date,
		"go", runtime.Version(),
	)

	redisConfig := config.LoadRedisConfig()
	consumerGroup := homerun.GetEnv("CONSUMER_GROUP", "homerun2-core-catcher")
	consumerName := homerun.GetEnv("CONSUMER_NAME", "")

	c, err := catcher.NewRedisCatcher(redisConfig, consumerGroup, consumerName)
	if err != nil {
		slog.Error("failed to create catcher", "error", err)
		os.Exit(1)
	}

	slog.Info("catcher configured",
		"redis_addr", redisConfig.Addr,
		"redis_port", redisConfig.Port,
		"stream", redisConfig.Stream,
		"consumer_group", consumerGroup,
	)

	// Handle shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Run consumer in goroutine, collect errors
	go func() {
		for err := range c.Errors() {
			slog.Error("consumer error", "error", err)
		}
	}()

	go func() {
		<-quit
		slog.Info("shutting down catcher")
		c.Shutdown()
	}()

	slog.Info("catcher running, waiting for messages...")
	c.Run()

	slog.Info("catcher exited gracefully")
}
