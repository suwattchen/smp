package nats

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/nats-io/nats.go"
)

type fakeMessage struct {
	meta       *nats.MsgMetadata
	metaErr    error
	ackErr     error
	acked      bool
	subjectStr string
}

func (f *fakeMessage) Ack() error {
	if f.ackErr != nil {
		return f.ackErr
	}
	f.acked = true
	return nil
}

func (f *fakeMessage) Metadata() (*nats.MsgMetadata, error) {
	return f.meta, f.metaErr
}

func (f *fakeMessage) Subject() string { return f.subjectStr }
func (f *fakeMessage) Reply() string   { return "" }
func (f *fakeMessage) Data() []byte    { return []byte("demo") }

func TestHandleMessageAck_JetStream(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	msg := &fakeMessage{meta: &nats.MsgMetadata{Stream: "events", Sequence: nats.SequencePair{Stream: 12}}, subjectStr: "events.created"}

	if err := HandleMessageAck(logger, msg); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !msg.acked {
		t.Fatalf("expected message to be acked")
	}
}

func TestHandleMessageAck_PlainNATS(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	msg := &fakeMessage{metaErr: errors.New("message lacks jetstream context"), subjectStr: "events.created"}

	if err := HandleMessageAck(logger, msg); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if msg.acked {
		t.Fatalf("expected message not to be acked for plain NATS")
	}
}

func TestHandleMessageAck_NilMessage(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if err := HandleMessageAck(logger, nil); err == nil {
		t.Fatalf("expected nil message error, got %v", err)
	}
}

func TestHandleMessageAck_AckFailure(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	msg := &fakeMessage{meta: &nats.MsgMetadata{Stream: "events", Sequence: nats.SequencePair{Stream: 1}}, subjectStr: "events.created", ackErr: errors.New("boom")}

	if err := HandleMessageAck(logger, msg); err == nil {
		t.Fatalf("expected error when ack fails")
	}
}
