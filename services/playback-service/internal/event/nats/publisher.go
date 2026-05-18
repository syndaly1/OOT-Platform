package nats

import (
	natsgo "github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *natsgo.Conn
}

func NewPublisher(conn *natsgo.Conn) *Publisher {
	return &Publisher{conn: conn}
}

func (p *Publisher) Publish(subject string, data []byte) error {
	return p.conn.Publish(subject, data)
}
