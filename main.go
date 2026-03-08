package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/banner"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/catcher"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/config"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/handlers"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/store"
	"github.com/stuttgart-things/homerun2-core-catcher/internal/tui"

	homerun "github.com/stuttgart-things/homerun-library/v2"
)

//go:embed static/*
var staticFS embed.FS

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
	var msgHandlers []catcher.MessageHandler

	// Log handler is always active (except cli mode where TUI owns the terminal)
	if mode != "cli" {
		msgHandlers = append(msgHandlers, catcher.LogHandler())
	}

	// Store handler for modes that need in-memory storage (cli, web)
	var msgStore *store.MessageStore
	if mode == "cli" || mode == "web" {
		maxMessages, _ := strconv.Atoi(homerun.GetEnv("MAX_MESSAGES", "10000"))
		msgStore = store.New(maxMessages)
		msgHandlers = append(msgHandlers, catcher.StoreHandler(msgStore))
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
		c = catcher.NewFileCatcher(filePath, interval, msgHandlers...)
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
		c, err = catcher.NewRedisCatcher(redisConfig, consumerGroup, consumerName, msgHandlers...)
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

	switch mode {
	case "cli":
		// Run catcher in background, TUI in foreground
		go c.Run()

		if errCh := c.Errors(); errCh != nil {
			go func() {
				for err := range errCh {
					slog.Error("consumer error", "error", err)
				}
			}()
		}

		p := tea.NewProgram(tui.New(msgStore))
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
			os.Exit(1)
		}

		c.Shutdown()

	case "web":
		// Run catcher in background, HTTP server in foreground
		go c.Run()

		if errCh := c.Errors(); errCh != nil {
			go func() {
				for err := range errCh {
					slog.Error("consumer error", "error", err)
				}
			}()
		}

		port := homerun.GetEnv("PORT", "8080")
		buildInfo := handlers.BuildInfo{Version: version, Commit: commit, Date: date}

		mux := http.NewServeMux()
		mux.HandleFunc("/", handlers.MessagesHandler(msgStore))
		mux.HandleFunc("/messages", handlers.MessagesTableHandler(msgStore))
		mux.HandleFunc("/messages/", handlers.MessageDetailHandler(msgStore))
		mux.HandleFunc("/export", handlers.ExportHandler(msgStore))
		mux.HandleFunc("/health", handlers.NewHealthHandler(buildInfo))
		staticSub, _ := fs.Sub(staticFS, "static")
		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

		srv := &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			slog.Info("web server starting", "port", port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("web server error", "error", err)
				os.Exit(1)
			}
		}()

		<-quit
		slog.Info("shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
		c.Shutdown()

		slog.Info("catcher exited gracefully")

	default: // log
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

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
}
