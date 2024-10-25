package squeue

import (
	"encoding/json"

	"github.com/simodima/squeue/driver"
)

//go:generate mockgen -source=driver/driver.go -package=squeue_test -destination=driver_test.go

// NewConsumer creates a new consumer for the given T type of messages
func NewConsumer[T json.Unmarshaler](d driver.Driver, queue string) Consumer[T] {
	return Consumer[T]{
		driver: d,
		queue:  queue,
	}
}

type Consumer[T json.Unmarshaler] struct {
	driver     driver.Driver
	queue      string
	controller *driver.ConsumerController
}

// Consume retrieves messages from the given queue.
// Any provided options will be sent to the underlying driver.
// The messages are indefinetely consumed from the queue and
// sent to the chan Message[T].
func (p *Consumer[T]) Consume(opts ...func(message any)) (chan Message[T], error) {
	ctrl, err := p.driver.Consume(p.queue, opts...)
	if err != nil {
		return nil, wrapErr(err, ErrDriver, nil)
	}
	p.controller = ctrl

	outMsg := make(chan Message[T])

	go func() {
		for message := range ctrl.Data() {
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

// Ping runs a ping on the underlying driver
func (p Consumer[T]) Ping() error {
	return p.driver.Ping()
}

func (p *Consumer[T]) Stop() {
	if p.controller != nil {
		p.controller.Stop()
	}
}

// Ack explicitly acknowldge the message handling.
// It can be implemented as No Operation for some drivers.
func (p *Consumer[T]) Ack(m Message[T]) error {
	return p.driver.Ack(p.queue, m.ID)
}
