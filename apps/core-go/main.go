package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "time"

    "github.com/nats-io/nats.go"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

type server struct {
    js     nats.JetStreamContext
    stream string
    port   string
}

func main() {
    zerolog.TimeFieldFormat = time.RFC3339
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    natsURL := getEnv("NATS_URL", "nats://nats:4222")
    stream := getEnv("CORE_NATS_STREAM", "events")
    port := getEnv("CORE_HTTP_PORT", "8080")

    nc, err := nats.Connect(natsURL)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to connect to NATS")
    }
    defer nc.Drain()

    js, err := nc.JetStream()
    if err != nil {
        log.Fatal().Err(err).Msg("failed to init JetStream context")
    }

    s := &server{js: js, stream: stream, port: port}

    if err := s.ensureStream(); err != nil {
        log.Fatal().Err(err).Msg("failed to ensure stream")
    }

    if err := s.subscribe(); err != nil {
        log.Fatal().Err(err).Msg("failed to subscribe")
    }

    go func() {
        if err := s.serveHTTP(); err != nil && err != http.ErrServerClosed {
            log.Fatal().Err(err).Msg("http server error")
        }
    }()

    <-ctx.Done()
    log.Info().Msg("shutting down core-go")
}

func (s *server) ensureStream() error {
    _, err := s.js.AddStream(&nats.StreamConfig{
        Name:     s.stream,
        Subjects: []string{"events.*"},
    })
    if err != nil && err != nats.ErrStreamNameAlreadyInUse {
        return err
    }
    return nil
}

func (s *server) subscribe() error {
    _, err := s.js.Subscribe("events.>", func(msg *nats.Msg) {
        log.Info().Str("subject", msg.Subject).Msg("received event")
        if err := msg.Ack(); err != nil {
            log.Error().Err(err).Str("subject", msg.Subject).Msg("ack failed")
        }
    })
    return err
}

func (s *server) serveHTTP() error {
    http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
        _, _ = fmt.Fprintln(w, "ok")
    })
    http.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
        _, _ = fmt.Fprintln(w, "requests_total 1")
    })

    addr := ":" + s.port
    log.Info().Str("addr", addr).Msg("starting http server")
    return http.ListenAndServe(addr, nil)
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
