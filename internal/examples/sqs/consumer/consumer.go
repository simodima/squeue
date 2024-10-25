package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/simodima/squeue"
	sqsexample "github.com/simodima/squeue/internal/examples/sqs"
	"github.com/simodima/squeue/sqs"
)

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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if !sqsexample.CheckEnvVariables("AWS_PROFILE", "AWS_SHARED_CREDENTIALS_FILE", "AWS_REGION", "AWS_QUEUE_URL") {
		log.Fatal(`Please set the env variables
		AWS_PROFILE=user-dev-admin
		AWS_SHARED_CREDENTIALS_FILE=/Users/{name.lastname}/.aws/credentials
		AWS_REGION=eu-central-1
		AWS_QUEUE_URL=https://sqs.eu-central-1.amazonaws.com/...
		`)
	}

	d, err := sqs.New(
		sqs.WithUrl(os.Getenv("AWS_QUEUE_URL")),
		sqs.WithRegion(os.Getenv("AWS_REGION")),
		sqs.AutoTestConnection(),
	)
	if err != nil {
		panic(err)
	}

	q := "test-simone"
	sub := squeue.NewConsumer[*sqsexample.MyEvent](d, q)
	go cancelOnSignal(func() {
		log.Println("Stopping consumer...")
		sub.Stop()
		log.Println("Stopped")
	}, syscall.SIGINT, syscall.SIGTERM)

	messages, err := sub.Consume(
		sqs.WithConsumeMaxNumberOfMessages(10),
		sqs.WithConsumeWaitTimeSeconds(2),
	)
	if err != nil {
		panic(err)
	}

	// Ping the sqs to check connectivity with AWS
	go func() {
		for {
			<-time.After(time.Second * 10)
			if err := sub.Ping(); err != nil {
				log.Println("Ping failed: " + err.Error())
			} else {
				log.Println("Ping OK")
			}
		}
	}()

	log.Print("Waiting for consuming")
	wg := &sync.WaitGroup{}

	for m := range messages {
		wg.Add(1)
		go func(message squeue.Message[*sqsexample.MyEvent]) {
			defer wg.Done()
			if message.Error != nil {
				log.Printf("Received message with an error %+v", message.Error)
				return
			}

			log.Printf("Received %s", message.Content.Name)
			if err := sub.Ack(message); err != nil {
				log.Print("Failed sending ack ", err)
			}
		}(m)
	}

	wg.Wait()
}
