package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"pem-parser/internal/app"
	"pem-parser/internal/port/ui"
	"time"
)

var debug bool

func main() {
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.Parse()

	logLevel := &slog.LevelVar{}
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	if debug {
		logLevel.Set(slog.LevelDebug)
		logger.Debug("Debug logging enabled")
	}

	slog.SetDefault(logger)

	application := app.NewApplication()

	server, err := ui.NewServer(logger, application)
	if err != nil {
		logger.Error("failed to create server", slog.Any("error", err))
	}

	errCh := make(chan error, 1)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		err := server.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to listen and serve: %w", err)
		}

		close(errCh)
	}()

	logger.Info("server running on Port 8080")

	select {
	// Wait until we receive SIGINT (ctrl+c on cli)
	case <-ctx.Done():
		logger.Info("server shutting down")
		break
	case err := <-errCh:
		logger.Error("failed to start server", slog.Any("error", err))
	}

	sCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err = server.Stop(sCtx); err != nil {
		logger.Error("failed to stop server", slog.Any("error", err))
	}
}
