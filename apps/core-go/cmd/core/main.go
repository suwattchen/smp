package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	snats "github.com/smp/core-go/internal/nats"

	"github.com/nats-io/nats.go"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	natsURL := getenv("NATS_URL", nats.DefaultURL)
	subject := getenv("NATS_SUBJECT", "events.demo")

	// start http health server
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	srv := &http.Server{Addr: ":8080", Handler: httpMux}
	go func() {
		logger.Info("starting http server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server exited", "err", err)
		}
	}()

	go startNATS(logger, natsURL, subject)

	<-ctx.Done()
	logger.Info("shutting down core service")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to shutdown http server", "err", err)
	}
}

func startNATS(logger *slog.Logger, url, subject string) {
	nc, err := nats.Connect(url, nats.MaxReconnects(-1))
	if err != nil {
		logger.Error("failed to connect to nats", "url", url, "err", err)
		return
	}
	logger.Info("connected to nats", "url", url)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		logger.Error("jetstream context error", "err", err)
		return
	}

	_, err = js.Subscribe(subject, func(msg *nats.Msg) {
		wrapped := snats.NewJSMessage(msg)
		if err := snats.HandleMessageAck(logger, wrapped); err != nil {
			logger.Error("failed to handle message", "err", err)
			return
		}

		logger.Info("processed message", "subject", wrapped.Subject(), "bytes", len(wrapped.Data()))
	}, nats.ManualAck())

	if err != nil {
		logger.Error("failed to subscribe", "subject", subject, "err", err)
		return
	}

	select {}
}

func getenv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
