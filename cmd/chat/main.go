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
	"syscall"
	"time"

	"github.com/nebiros/jobsity-chat-assignment/internal/rest"
	"github.com/nebiros/jobsity-chat-assignment/pkg/db"
	"github.com/nebiros/jobsity-chat-assignment/pkg/httpext"
	"github.com/nebiros/jobsity-chat-assignment/pkg/session"
)

var (
	address    = flag.String("address", ":3000", "Server address")
	dbFilePath = flag.String("dbFilePath", "./data/db.sqlite", "DB file path")
	sessionKey = flag.String("sessionKey", "", "Session key")
	debug      = flag.Bool("debug", false, "Set debug mode")
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
}

func main() {
	flag.Parse()

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags]\n", os.Args[0])
		flag.PrintDefaults()
	}

	if len(os.Args) < 1 {
		flag.Usage()
		os.Exit(0)
	}

	if *sessionKey == "" {
		slog.Error("you must specify a session key")

		flag.Usage()
		os.Exit(1)
	}

	if *debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}

	if err := run(); err != nil {
		slog.Error("unable to run application", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL}...)
	defer cancel()

	dbClient, err := db.NewClient(*dbFilePath)
	if err != nil {
		return err
	}

	mux, err := rest.MakeRoutes(rest.WithDBClient(dbClient), rest.WithSessionStore(session.NewCookieStore(*sessionKey)))
	if err != nil {
		return err
	}

	s := httpext.NewServer(*address, mux)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := dbClient.Close(); err != nil {
			slog.Error("unable to close DB connection", slog.Any("error", err))
		}

		if err := s.Shutdown(ctx); err != nil {
			slog.Error("unable to shutdown server", slog.Any("error", err))
		}
	}()

	_, _ = fmt.Fprintf(os.Stdout, "=====> start receiving at '%s'\n", *address)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server listen and serve errored", slog.Any("error", err))
		return err
	}

	return nil
}
