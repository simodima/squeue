package squeue

import (
	"encoding/json"

	"github.com/simodima/squeue/driver"
)

func NewProducer(d driver.Driver) Producer {
	return Producer{
		driver: d,
	}
}

type Producer struct {
	driver driver.Driver
}

// Enqueue sends a message to the given queue
// any provided options will be sent to the
// underlying driver
func (q *Producer) Enqueue(queue string, message json.Marshaler, opts ...func(message any)) error {
	data, err := json.Marshal(message)
	if err != nil {
		return wrapErr(err, ErrMarshal, nil)
	}

	return q.driver.Enqueue(queue, data, opts...)
}
