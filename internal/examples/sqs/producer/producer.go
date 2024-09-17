package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/simodima/squeue"
	sqsexample "github.com/simodima/squeue/internal/examples/sqs"
	"github.com/simodima/squeue/sqs"
)

func init() {
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
}

func main() {
	d, err := sqs.New(
		sqs.WithUrl(os.Getenv("AWS_QUEUE_URL")),
		sqs.WithRegion(os.Getenv("AWS_REGION")),
	)

	if err != nil {
		panic(err)
	}

	pub := squeue.NewProducer(d, "test-simone")
	tick := time.Tick(time.Second * 2)

	for i := 0; ; i++ {
		<-tick
		_ = pub.Enqueue(&sqsexample.MyEvent{Name: fmt.Sprintf("Message #%d", i)})
	}
}
