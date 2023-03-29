package mq

// MQ message platform
type MQ interface {
	Pub(topic string, msg any) MQ                      // Pub any message`
	MultiPub(topic string, msgs ...any) MQ             // MultiPub ...any message
	Sub(topic, channel string, handler HandlerFunc) MQ // consumer Sub func to handle your work
	Stop(topic, channel string) MQ                     // consumer Stop
	Error() error
}

type Message struct {
	Body      []byte
	MessageID string
}

type HandlerFunc func(message *Message) error

// HandlerMessage implements the Handler interface
func (h HandlerFunc) HandlerMessage(message *Message) error {
	return h(message)
}

// NewMQ Sugar
func NewMQ(driver string, params ...any) (mq MQ) {
	switch driver {
	case "nsq":
		mq = NewNSQ(params[0].(string), params[1].(string))
		return
	}
	return
}
