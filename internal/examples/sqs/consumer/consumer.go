package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/toretto460/squeue"
	sqsexample "github.com/toretto460/squeue/internal/examples/sqs"
	"github.com/toretto460/squeue/sqs"
)

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d, err := sqs.New(
		sqs.WithUrl(os.Getenv("AWS_QUEUE_URL")),
		sqs.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		panic(err)
	}

	sub := squeue.NewConsumer[*sqsexample.MyEvent](d)
	q := "test-simone"

	messages, err := sub.Consume(ctx, q)
	if err != nil {
		panic(err)
	}

	log.Print("Waiting for consuming")

	for m := range messages {
		go func(message squeue.Message[*sqsexample.MyEvent]) {
			if message.Error != nil {
				log.Printf("Received an error %+v", message.Error)
				return
			}

			log.Printf("Received %+v", message)
			if message.Content.Name != "user" {
				sub.Ack(q, message)
			}
		}(m)
	}
}
