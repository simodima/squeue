package driver

import "context"

type Driver interface {
	Enqueue(queue string, data []byte) error
	Consume(ctx context.Context, topic string) (chan Message, error)
	Ack(queue string, messageID string) error
}
