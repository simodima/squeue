package driver

type Message struct {
	Body  []byte
	ID    string
	Error error
}
