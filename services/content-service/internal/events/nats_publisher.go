package events

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(natsURL string) (*Publisher, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	return &Publisher{conn: conn}, nil
}

func (p *Publisher) Publish(subject string, data any) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return
	}

	if err := p.conn.Publish(subject, payload); err != nil {
		log.Printf("failed to publish event: %v", err)
		return
	}

	log.Printf("[EVENT PUBLISHED] subject=%s payload=%s", subject, string(payload))
}

func (p *Publisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
