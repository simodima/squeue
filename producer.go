package squeue

import (
	"encoding/json"

	"github.com/simodima/squeue/driver"
)

func NewProducer(d driver.Driver, queue string) Producer {
	return Producer{
		driver: d,
		queue:  queue,
	}
}

type Producer struct {
	driver driver.Driver
	queue  string
}

// Enqueue sends a message to the underlying queue,
// any provided options will be sent to the driver.
func (q *Producer) Enqueue(message json.Marshaler, opts ...func(message any)) error {
	data, err := json.Marshal(message)
	if err != nil {
		return wrapErr(err, ErrMarshal, nil)
	}

	return q.driver.Enqueue(q.queue, data, opts...)
}
