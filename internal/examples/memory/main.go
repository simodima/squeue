package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/simodima/squeue"
	"github.com/simodima/squeue/driver"
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

	d := driver.NewMemoryDriver(time.Microsecond)

	queue := squeue.NewProducer(d, "queue.test")
	consumer := squeue.NewConsumer[*myMessage](d, "queue.test")

	messages, err := consumer.Consume()
	if err != nil {
		panic(err)
	}

	go cancelOnSignal(func() {
		consumer.Stop()
	}, syscall.SIGINT, syscall.SIGTERM)

	_ = queue.Enqueue(&myMessage{"foo"})
	_ = queue.Enqueue(&myMessage{"bar"})
	_ = queue.Enqueue(&myMessage{"baz"})

	for evt := range messages {
		log.Print("Received ", evt.Content)
	}
}
