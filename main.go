package main

import (
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/stuttgart-things/homerun2-core-catcher/internal/banner"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/catcher"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/config"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/store"

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

	mode := homerun.GetEnv("CATCHER_MODE", "log")
	backend := homerun.GetEnv("CATCHER_BACKEND", "redis")

	// Build message handlers based on mode
	var handlers []catcher.MessageHandler

	// Log handler is always active
	handlers = append(handlers, catcher.LogHandler())

	// Store handler for modes that need in-memory storage (cli, web)
	var msgStore *store.MessageStore
	if mode == "cli" || mode == "web" {
		maxMessages, _ := strconv.Atoi(homerun.GetEnv("MAX_MESSAGES", "10000"))
		msgStore = store.New(maxMessages)
		handlers = append(handlers, catcher.StoreHandler(msgStore))
		_ = msgStore // will be used by TUI/web in future
	}

	// Create catcher backend
	var c catcher.Catcher

	switch backend {
	case "file":
		filePath := homerun.GetEnv("CATCHER_FILE_PATH", "messages.json")
		intervalStr := homerun.GetEnv("CATCHER_FILE_INTERVAL", "1s")
		interval, err := time.ParseDuration(intervalStr)
		if err != nil {
			interval = time.Second
		}
		c = catcher.NewFileCatcher(filePath, interval, handlers...)
		slog.Info("catcher configured",
			"backend", "file",
			"path", filePath,
			"interval", interval,
			"mode", mode,
		)
	default:
		redisConfig := config.LoadRedisConfig()
		consumerGroup := homerun.GetEnv("CONSUMER_GROUP", "homerun2-core-catcher")
		consumerName := homerun.GetEnv("CONSUMER_NAME", "")

		var err error
		c, err = catcher.NewRedisCatcher(redisConfig, consumerGroup, consumerName, handlers...)
		if err != nil {
			slog.Error("failed to create catcher", "error", err)
			os.Exit(1)
		}
		slog.Info("catcher configured",
			"backend", "redis",
			"redis_addr", redisConfig.Addr,
			"redis_port", redisConfig.Port,
			"stream", redisConfig.Stream,
			"consumer_group", consumerGroup,
			"mode", mode,
		)
	}

	// Handle shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Collect errors
	if errCh := c.Errors(); errCh != nil {
		go func() {
			for err := range errCh {
				slog.Error("consumer error", "error", err)
			}
		}()
	}

	go func() {
		<-quit
		slog.Info("shutting down catcher")
		c.Shutdown()
	}()

	slog.Info("catcher running, waiting for messages...")
	c.Run()

	slog.Info("catcher exited gracefully")
}
