package squeue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/toretto460/squeue/driver"
)

// NewConsumer
func NewConsumer[T json.Unmarshaler](d driver.Driver) Consumer[T] {
	return Consumer[T]{d}
}

type Consumer[T any] struct {
	driver driver.Driver
}

func (p *Consumer[T]) Consume(ctx context.Context, topic string) (chan Message[T], error) {
	messages, err := p.driver.Consume(ctx, topic)
	if err != nil {
		return nil, fmt.Errorf("error consuming messages: %w", err)
	}

	outMsg := make(chan Message[T])

	go func() {
		for message := range messages {
			if message.Error != nil {
				outMsg <- Message[T]{
					Error: message.Error,
				}
				continue
			}

			var content T
			err := json.Unmarshal(message.Body, &content)
			if err != nil {
				outMsg <- Message[T]{
					Content: content,
					ID:      message.ID,
					Error:   err,
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

func (p *Consumer[T]) Ack(queue string, m Message[T]) error {
	return p.driver.Ack(queue, m.ID)
}
