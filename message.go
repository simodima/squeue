package squeue

// Message represent a queue message that can be consumed from any queue
//
// type myMessage struct {
//     content string
// }
// func (e *myMessage) MarshalJSON() ([]byte, error) { ... }
// func (e *myMessage) UnmarshalJSON(data []byte) error { ... }

// sub := squeue.NewConsumer[*myMessage](d)
// for m := range sub.Consume() {
//     do something with m
// }
type Message[T any] struct {
	Content T
	ID      string
	Error   error
}
