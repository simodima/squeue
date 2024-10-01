package driver

func NewConsumerController() *ConsumerController {
	return &ConsumerController{
		data: make(chan Message),
		done: make(chan struct{}),
	}
}

type ConsumerController struct {
	data chan Message
	done chan struct{}
}

func (c *ConsumerController) Send(m Message) {
	c.data <- m
}

func (c *ConsumerController) Data() <-chan Message {
	return c.data
}

func (c *ConsumerController) Done() chan struct{} {
	return c.done
}

func (c *ConsumerController) Stop() {
	c.done <- struct{}{}
	close(c.data)
}

type Driver interface {
	Enqueue(queue string, data []byte, opts ...func(message any)) error
	Consume(queue string, opts ...func(message any)) (*ConsumerController, error)
	Ack(queue string, messageID string) error
}
