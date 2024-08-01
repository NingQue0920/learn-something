//go:build kafka
// +build kafka

package middleware

import (
	"net"
)

type Kafka struct {
	config kafka.ConnConfig
	conn   net.Conn
}

func (k *Kafka) Initialize() error {
	// init kafka
	var err error
	k.config = kafka.ConnConfig{
		ClientID:        "",
		Topic:           "",
		Partition:       0,
		Broker:          0,
		Rack:            "",
		TransactionalID: "",
	}
	k.conn = kafka.NewConnWith(nil, k.config)
	return err
}
func (k *Kafka) Produce(topic, message string) error {
	_, err := k.conn.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}
