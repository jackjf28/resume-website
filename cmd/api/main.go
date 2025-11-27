package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/jackjf28/resume-website/server"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	if len(args) > 1 {
		fmt.Fprintf(w, "%s\n", strings.Join(args, ","))
	}

	//Logging
	logger := slog.New(slog.NewTextHandler(w, nil))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	s := server.NewServer(ctx)

	httpServer := &http.Server{
		Addr:    ":4000",
		Handler: s,
	}

	fmt.Fprintf(w, fmt.Sprintf("listening on %s.\n", httpServer.Addr))
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(w, "error listening on server: %v\n", err)
		}
	}()

	<-ctx.Done()
	fmt.Fprintf(w, "\nshutting down server...\n")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(w, "error during shutdown: %v\n", err)
		return err
	}

	fmt.Fprintf(w, "server shut down.\n")
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
