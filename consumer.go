package squeue

import (
	"context"
	"encoding/json"

	"github.com/toretto460/squeue/driver"
)

//go:generate mockgen -source=driver/driver.go -package=squeue_test -destination=driver_test.go

// NewConsumer
func NewConsumer[T json.Unmarshaler](d driver.Driver) Consumer[T] {
	return Consumer[T]{d}
}

type Consumer[T json.Unmarshaler] struct {
	driver driver.Driver
}

// TODO: document Consume
func (p *Consumer[T]) Consume(ctx context.Context, topic string, opts ...func(message any)) (chan Message[T], error) {
	messages, err := p.driver.Consume(ctx, topic, opts...)
	if err != nil {
		return nil, wrapErr(err, ErrDriver, nil)
	}

	outMsg := make(chan Message[T])

	go func() {
		for message := range messages {
			if message.Error != nil {
				outMsg <- Message[T]{
					Error: wrapErr(message.Error, ErrDriver, nil),
				}
				continue
			}

			var content T
			err := json.Unmarshal(message.Body, &content)
			if err != nil {
				outMsg <- Message[T]{
					ID:    message.ID,
					Error: wrapErr(err, ErrUnmarshal, message.Body),
				}
				continue
			}

			outMsg <- Message[T]{
				Content: content,
				ID:      message.ID,
			}
		}
		close(outMsg)
	}()

	return outMsg, nil
}

// TODO: document Ack
func (p *Consumer[T]) Ack(queue string, m Message[T]) error {
	return p.driver.Ack(queue, m.ID)
}
