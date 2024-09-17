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

// TODO: document Enqueue
func (q *Producer) Enqueue(queue string, message json.Marshaler, opts ...func(message any)) error {
	data, err := json.Marshal(message)
	if err != nil {
		return wrapErr(err, ErrMarshal, nil)
	}

	return q.driver.Enqueue(queue, data, opts...)
}
