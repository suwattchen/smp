package nats

import "github.com/nats-io/nats.go"

// JSMessage adapts *nats.Msg to the Message interface expected by the ack handler.
type JSMessage struct {
	msg *nats.Msg
}

// NewJSMessage wraps an incoming JetStream message.
func NewJSMessage(msg *nats.Msg) JSMessage {
	return JSMessage{msg: msg}
}

func (m JSMessage) Ack() error {
	return m.msg.Ack()
}

func (m JSMessage) Metadata() (*nats.MsgMetadata, error) {
	return m.msg.Metadata()
}

func (m JSMessage) Subject() string { return m.msg.Subject }
func (m JSMessage) Reply() string   { return m.msg.Reply }
func (m JSMessage) Data() []byte    { return m.msg.Data }
