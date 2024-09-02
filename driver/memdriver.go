package driver

import (
	"context"
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

func (d *MemoryDriver) Enqueue(queue string, evt []byte) error {
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

func (d *MemoryDriver) Consume(ctx context.Context, topic string) (chan Message, error) {
	results := make(chan Message)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(results)
				return
			case <-time.After(d.tick):
				if evt := d.pop(topic); evt != nil {
					results <- *evt
				}
			}
		}
	}()

	return results, nil
}

func (d *MemoryDriver) pop(topic string) *Message {
	d.l.Lock()
	defer d.l.Unlock()

	if _, ok := d.messages[topic]; ok {
		if len(d.messages[topic]) > 0 {
			evt := d.messages[topic][0]
			d.messages[topic] = d.messages[topic][1:]
			return &evt
		}
	}

	return nil
}
