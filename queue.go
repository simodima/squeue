package squeue

import (
	"encoding/json"

	"github.com/toretto460/squeue/driver"
)

func NewQueue(d driver.Driver, q string) Queue {
	return Queue{
		driver:     d,
		identifier: q,
	}
}

type Queue struct {
	driver     driver.Driver
	identifier string
}

func (q *Queue) Enqueue(message json.Marshaler) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return q.driver.Enqueue(q.identifier, data)
}
