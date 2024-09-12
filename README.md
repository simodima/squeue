<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./docs/squeue-white.png">
  <source media="(prefers-color-scheme: light)"  srcset="./docs/squeue-black.png">
  <img alt="S-queue mascotte" src="./docs/squeue-black.png">
</picture>

**s-queue** is a library to interact with Amazon SQS queues with 100% type safety, powered by golang generics.

---

[![codecov](https://codecov.io/github/simodima/squeue/graph/badge.svg?token=DW7C57P2VW)](https://codecov.io/github/simodima/squeue)
[![Go Report Card](https://goreportcard.com/badge/github.com/toretto460/squeue)](https://goreportcard.com/report/github.com/toretto460/squeue)
[![Go Reference](https://pkg.go.dev/badge/github.com/simodima/squeue.svg)](https://pkg.go.dev/github.com/simodima/squeue)

The library provides the following features:

### Driver based implementation
The interface exposed by the library abstract the internal driver implementation, that allows to plug various different drivers depending of the development needs and the environment.

The currently implemented drivers are :
- In Memory
- Amazon SQS

You can easly use the *Amazon SQS* driver in any staging/production environment and use the *In Memory* one for a lightweight testing environment.

### Type Safety
The library uses an opinionated JSON format to serialize/deserialize the messages. To be 100% typesafe you just need to provide a message type implementing the golang `json.Marshaler` and `json.Unmarshaler` interfaces; in that way the raw data serialization/deserialization is in your direct control.

### Channel of messages 
Consuming message is as easy as looping a channel of messages

```golang
    messages, _ := sub.Consume(ctx, "queue")
    for m := range messages {
        go func(message squeue.Message[*my.Message]) {
        // Do the magic...
        }(m)
    }
```

<details>

<summary>Examples</summary>

For a more clear documentation look at the `internal/examples/` directory

**In Memory driver**
```bash
go run internal/examples/memory/main.go
```

**Amazon SQS driver**
```bash
go run internal/examples/sqs/consumer/consumer.go

## in a different shell ↓
go run internal/examples/sqs/producer/producer.go
```

</details>
