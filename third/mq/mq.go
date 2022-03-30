package mq

import "github.com/nsqio/go-nsq"

// MQ message platform
type MQ interface {
	NewProducer() MQ                          // new producer
	NewConsumer(topic, channel string) MQ     // new consumer
	Stop()                                    // consumer Stop
	Pub(topic string, msg any) error          // pub any message
	MultiPub(topic string, msgs ...any) error // MultiPub ...any message
	Sub(handler HandlerFunc) error            // sub func to handle your work
}

type Message struct {
	*nsq.Message
}

type HandlerFunc func(message *Message) error

// HandlerMessage implements the Handler interface
func (h HandlerFunc) HandlerMessage(message *Message) error {
	return h(message)
}

// NewProducer new producer
func NewProducer(params ...any) (mq MQ) {
	// default
	if len(params) == 0 {
		mq = new(Nsg)
		mq.NewProducer()
		return
	}
	// init
	mq = params[0].(MQ)
	mq.NewProducer()
	return
}

// NewConsumer new consumer
func NewConsumer(topic, channel string, params ...any) (mq MQ) {
	// default
	if len(params) == 0 {
		mq = new(Nsg)
		mq.NewConsumer(topic, channel)
		return
	}
	// init
	mq = params[0].(MQ)
	mq.NewConsumer(topic, channel)
	return
}
