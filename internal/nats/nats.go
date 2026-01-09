package nats

import (
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
)

// Константы для описания операций
const (
	opConn            = "nats.Connect"
	opCreateJetStream = "nats.CreateJetStreamContext"
	opCreateStream    = "nats.CreateStream"
	opPubMessage      = "nats.PublishMessage"
)

// Nats описывает брокер сообщений
type Nats struct {
	log  *slog.Logger
	Conn *nats.Conn
	js   nats.JetStreamContext
}

// NewNats конструктор для Nats
func NewNats(log *slog.Logger, url string) (*Nats, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opConn))
		return nil, fmt.Errorf("%s: %w", opConn, err)
	}
	log.Info("connected to "+url, slog.String("operation", opConn))

	js, err := nc.JetStream()
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opCreateJetStream))
		return nil, fmt.Errorf("%s: %w", opCreateJetStream, err)
	}
	log.Info("create jetstream context", slog.String("operation", opCreateJetStream))

	return &Nats{
		log:  log,
		Conn: nc,
		js:   js,
	}, nil
}

// CreateStream создаёт поток и топики
func (n *Nats) CreateStream(name string, subjects []string) (*nats.StreamInfo, error) {
	res, err := n.js.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: subjects,
	})
	if err != nil {
		n.log.Error("error", err.Error(), slog.String("operation", opCreateStream))
		return nil, fmt.Errorf("%s: %w", opCreateStream, err)
	}

	n.log.Info("create stream "+name, slog.String("operation", opCreateStream))

	return res, nil
}

// PublishMessage публикует сообщений в соответствующий топик
func (n *Nats) PublishMessage(topic string, data []byte) error {
	_, err := n.js.Publish(topic, data)
	if err != nil {
		n.log.Error("error", err.Error(), slog.String("operation", opPubMessage))
		return fmt.Errorf("%s: %w", opPubMessage, err)
	}

	n.log.Info("pub to "+topic, slog.String("operation", opPubMessage))

	return nil
}
