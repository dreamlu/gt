package msg

import (
	"encoding/json"
	"fmt"
	"github.com/dreamlu/gt/tool/conf"
	"github.com/dreamlu/gt/tool/log"
	"github.com/nsqio/go-nsq"
	"strings"
)

// Nsg nsq
type Nsg struct {
	producer *nsq.Producer
	consumer *nsq.Consumer
	bs       []byte
}

func (n *Nsg) NewProducer() Msg {

	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(conf.Configger().GetString("app.nsq.producer_addr"), config)
	if err != nil {
		fmt.Printf("[gt]:create producer failed, err:%v\n", err)
		return nil
	}
	n.producer = producer
	return n
}

func (n *Nsg) NewConsumer(topic, channel string) Msg {

	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		fmt.Printf("[gt]:create producer failed, err:%v\n", err)
		return nil
	}
	n.consumer = consumer
	return n
}

func (n *Nsg) Stop() {
	n.consumer.Stop()
}

func (n *Nsg) Pub(topic string, msg interface{}) error {

	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("[gt]:MSG Pub err: ", err)
		return err
	}
	err = n.producer.Publish(topic, b)
	return err
}

func (n *Nsg) MultiPub(topic string, msgs ...interface{}) error {

	var bs [][]byte
	for _, v := range msgs {
		b, err := json.Marshal(v)
		if err != nil {
			fmt.Println("[gt]:MSG MultiPub err: ", err)
			return err
		}
		bs = append(bs, b)
	}

	err := n.producer.MultiPublish(topic, bs)
	return err
}

func (n *Nsg) Sub(handler HandlerFunc) error {

	n.consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		msg := &Message{}
		msg.Message = message
		err := handler(msg)
		return err
	}))
	// use ',' split address
	// ConnectToNSQD/ConnectToNSQLookupd
	err := n.consumer.ConnectToNSQDs(strings.Split(conf.Configger().GetString("app.nsq.consumer_addr"), ","))
	if err != nil {
		log.Error("[gt]:MSG Consumer ConnectToNSQD err: ", err)
		return err
	}

	return nil
}
