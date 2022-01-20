package msg

import "github.com/nsqio/go-nsq"

// Msg message platform
type Msg interface {
	NewProducer() Msg                                 // new producer
	NewConsumer(topic, channel string) Msg            // new consumer
	Stop()                                            // consumer Stop
	Pub(topic string, msg interface{}) error          // pub interface{} message
	MultiPub(topic string, msgs ...interface{}) error // MultiPub ...interface{} message
	Sub(handler HandlerFunc) error                    // sub func to handle your work
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
func NewProducer(params ...interface{}) (msg Msg) {
	// default
	if len(params) == 0 {
		msg = new(Nsg)
		msg.NewProducer()
		return
	}
	// init
	msg = params[0].(Msg)
	msg.NewProducer()
	return
}

// NewConsumer new consumer
func NewConsumer(topic, channel string, params ...interface{}) (msg Msg) {
	// default
	if len(params) == 0 {
		msg = new(Nsg)
		msg.NewConsumer(topic, channel)
		return
	}
	// init
	msg = params[0].(Msg)
	msg.NewConsumer(topic, channel)
	return
}
