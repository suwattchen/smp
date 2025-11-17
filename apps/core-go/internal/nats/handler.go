package nats

import (
	"errors"
	"log/slog"

	"github.com/nats-io/nats.go"
)

// Message wraps the subset of nats.Msg we rely on so it can be tested easily.
type Message interface {
	Ack() error
	Metadata() (*nats.MsgMetadata, error)
	Subject() string
	Reply() string
	Data() []byte
}

// HandleMessageAck ensures JetStream messages are acknowledged, while plain NATS messages are left untouched.
func HandleMessageAck(logger *slog.Logger, msg Message) error {
	if msg == nil {
		return errors.New("nil message received")
	}

	meta, err := msg.Metadata()
	if err != nil || meta == nil {
		logger.Info("received plain NATS message, skipping ack", "subject", msg.Subject())
		return nil
	}

	if err := msg.Ack(); err != nil {
		logger.Error("failed to ack JetStream message", "subject", msg.Subject(), "err", err)
		return err
	}

	logger.Info("acknowledged JetStream message", "subject", msg.Subject(), "stream", meta.Stream, "sequence", meta.Sequence.Stream)
	return nil
}
