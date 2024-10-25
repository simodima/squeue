package driver

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func NewMemoryDriver(tick time.Duration) *MemoryDriver {
	return &MemoryDriver{
		l:        &sync.Mutex{},
		messages: make(map[string][]Message),
		tick:     tick,
	}
}

type MemoryDriver struct {
	messages map[string][]Message
	l        *sync.Mutex
	tick     time.Duration
}

func (d *MemoryDriver) Ack(queue, messageID string) error {
	return nil
}

func (d *MemoryDriver) Enqueue(queue string, evt []byte, opts ...func(message any)) error {
	d.l.Lock()
	defer d.l.Unlock()

	if evts, ok := d.messages[queue]; ok {
		d.messages[queue] = append(evts, Message{Body: evt})
	} else {
		d.messages[queue] = []Message{
			{
				Body: evt,
				ID:   fmt.Sprint(rand.Int()),
			},
		}
	}

	return nil
}

func (d *MemoryDriver) Consume(queue string, opts ...func(message any)) (*ConsumerController, error) {
	ctrl := NewConsumerController()

	go func() {
		for {
			select {
			case <-ctrl.Done():
				return
			case <-time.After(d.tick):
				if evt := d.pop(queue); evt != nil {
					ctrl.Send(*evt)
				}
			}
		}
	}()

	return ctrl, nil
}

func (d *MemoryDriver) pop(queue string) *Message {
	d.l.Lock()
	defer d.l.Unlock()

	if _, ok := d.messages[queue]; ok {
		if len(d.messages[queue]) > 0 {
			evt := d.messages[queue][0]
			d.messages[queue] = d.messages[queue][1:]
			return &evt
		}
	}

	return nil
}

func (d *MemoryDriver) Ping() error {
	return nil
}
