// Copyright the original author or authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"syscall"
	"time"

	"l7e.io/vanity"
	"l7e.io/vanity/cmd/vanity/server/interceptors"
	"l7e.io/vanity/pkg/toml"
	"l7e.io/yama"
	"m4o.io/gslog"
	"m4o.io/gslog/stdout"
)

const (
	// readHeaderTimeout is the time that the server waits for the headers of a
	// request. It stops a client that sends the headers slowly.
	readHeaderTimeout = 10 * time.Second

	// shutdownTimeout is the time that the server gets to close itself after a
	// signal to stop.
	shutdownTimeout = 2 * time.Second
)

func main() {
	handler := stdout.NewHandler(os.Stdout, gslog.WithLogLevelFromEnvVar("LOG_LEVEL"))
	slog.SetDefault(slog.New(handler))

	ctx := context.Background()

	// The library sends its request failures and each served import path
	// through one logger, at one severity. Debug keeps the served paths out
	// of the error stream.
	vanity.SetLogger(vanity.LoggerFunc(func(format string, v ...interface{}) {
		slog.DebugContext(ctx, fmt.Sprintf(format, v...))
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	slog.InfoContext(ctx, "Port configured", slog.String("addr", addr))

	file := os.Getenv("TOML_FILE")
	if file == "" {
		file = "/var/vanity/config.toml"
	}
	slog.InfoContext(ctx, "TOML file location", slog.String("file", file))

	table := os.Getenv("TOML_TABLE")
	if table == "" {
		table = "vanity"
	}
	slog.InfoContext(ctx, "TOML table name", slog.String("table", table))

	be, err := toml.NewTOMLBackend(toml.InTable(table), toml.FromFile(file))
	if err != nil {
		slog.ErrorContext(ctx, "Unable to read the vanity entries", slog.Any("error", err))
		os.Exit(1)
	}

	consumer := vanity.ConsumerFunc(func(ctx context.Context, importPath, vcs, vcsPath string) {
		slog.InfoContext(ctx, "Vanity entry",
			slog.String("importPath", importPath),
			slog.String("vcs", vcs),
			slog.String("vcsPath", vcsPath))
	})

	err = be.List(ctx, consumer)
	if err != nil {
		slog.ErrorContext(ctx, "Unable to list the vanity entries", slog.Any("error", err))
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/", interceptors.WrapHandler(vanity.NewVanityHandler(be)))

	s := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	watcher := yama.NewWatcher(
		yama.WatchingSignals(syscall.SIGINT, syscall.SIGTERM),
		yama.WithTimeout(shutdownTimeout),
		yama.WithClosers(s))

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		slog.ErrorContext(ctx, "Unable to serve", slog.Any("error", err))
		_ = watcher.Close()
	}

	if err := watcher.Wait(); err != nil {
		slog.WarnContext(ctx, "Shutdown error", slog.Any("error", err))
	}
}
