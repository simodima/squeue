package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/toretto460/squeue"
	"github.com/toretto460/squeue/driver"
)

type myMessage struct {
	name string
}

func (e *myMessage) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	e.name = raw["name"].(string)

	return nil
}

func (e *myMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"name": e.name,
	})
}

func cancelOnSignal(fn func(), signals ...os.Signal) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, signals...)

	s := <-sigch

	for _, sig := range signals {
		if s.String() == sig.String() {
			log.Printf("Signal %s intercepted", s)
			fn()
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go cancelOnSignal(cancel, syscall.SIGINT, syscall.SIGTERM)

	d := driver.NewMemoryDriver(time.Microsecond)

	queue := squeue.NewProducer(d)
	consumer := squeue.NewConsumer[*myMessage](d)

	events, err := consumer.Consume(ctx, "queue.test")
	if err != nil {
		panic(err)
	}

	_ = queue.Enqueue("queue.test", &myMessage{"foo"})
	_ = queue.Enqueue("queue.test", &myMessage{"bar"})
	_ = queue.Enqueue("queue.test", &myMessage{"baz"})

	// Consumer gorouting
	go func() {
		for evt := range events {
			log.Print("Received ", evt.Content)
		}
	}()

	<-ctx.Done()
}
