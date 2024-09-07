package driver

import "context"

type Driver interface {
	Enqueue(queue string, data []byte, opts ...func(message any)) error
	Consume(ctx context.Context, topic string, opts ...func(message any)) (chan Message, error)
	Ack(queue string, messageID string) error
}
