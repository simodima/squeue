package squeue

import (
	"encoding/json"

	"github.com/toretto460/squeue/driver"
)

func NewProducer(d driver.Driver) Producer {
	return Producer{
		driver: d,
	}
}

type Producer struct {
	driver driver.Driver
}

func (q *Producer) Enqueue(queue string, message json.Marshaler, opts ...func(message any)) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return q.driver.Enqueue(queue, data, opts...)
}
