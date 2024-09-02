package main

import (
	"context"
	"encoding/json"
	"log"
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

func main() {
	wait := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
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

	log.Print("Waiting for consuming")
	time.Sleep(time.Second)

	go func() {
		defer func() {
			log.Print("Consumed all messages")
		}()
		for evt := range events {
			log.Print("MSG CONTENT RECEIVED\t", evt.Content)
		}
		log.Print("closed chan")
		wait <- struct{}{}
	}()

	// let the consumer goroutine consume the messages
	time.Sleep(time.Second)
	log.Print("cancelling...")
	cancel()
	time.Sleep(time.Second * 10)
	<-wait
}
