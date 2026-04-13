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

	if !sqsexample.CheckEnvVariables("AWS_REGION", "AWS_QUEUE_URL") {
		log.Fatal(`Please set the env variables
		AWS_REGION=eu-central-1
		AWS_QUEUE_URL=https://sqs.eu-central-1.amazonaws.com/...
		Credentials are resolved via the AWS SDK default chain (env vars, ~/.aws/credentials, Pod Identity, etc.)
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
	tick := time.Tick(time.Millisecond)

	for i := 0; ; i++ {
		<-tick
		_ = pub.Enqueue(&sqsexample.MyEvent{Name: fmt.Sprintf("Message #%d", i)})
	}
}
