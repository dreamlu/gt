package mq

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"strings"
)

// NSQ nsq
type NSQ struct {
	producer     *nsq.Producer
	consumerAddr string
	consumers    map[string]*nsq.Consumer
	bs           []byte
	err          error
}

func NewNSQ(producerAddr, consumerAddr string) MQ {

	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(producerAddr, config)
	if err != nil {
		panic(err)
		return nil
	}
	return &NSQ{
		producer:     producer,
		consumerAddr: consumerAddr,
		consumers:    make(map[string]*nsq.Consumer),
	}
}

func (n *NSQ) Stop(topic, channel string) MQ {
	n.consumers[topic+channel].Stop()
	return n
}

func (n *NSQ) Pub(topic string, msg any) MQ {

	b, err := json.Marshal(msg)
	if err != nil {
		n.err = err
		return n
	}
	n.err = n.producer.Publish(topic, b)
	return n
}

func (n *NSQ) MultiPub(topic string, msgs ...any) MQ {

	var bs [][]byte
	for _, v := range msgs {
		b, err := json.Marshal(v)
		if err != nil {
			n.err = err
			return n
		}
		bs = append(bs, b)
	}

	n.err = n.producer.MultiPublish(topic, bs)
	return n
}

func (n *NSQ) Sub(topic, channel string, handler HandlerFunc) MQ {

	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		n.err = err
		return n
	}
	n.consumers[topic+channel] = consumer

	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		msg := &Message{}
		msg.Body = message.Body
		msg.MessageID = string(message.ID[:])
		err := handler(msg)
		return err
	}))
	n.err = consumer.ConnectToNSQDs(strings.Split(n.consumerAddr, ","))
	return n
}

func (n *NSQ) Error() error {
	return n.err
}
