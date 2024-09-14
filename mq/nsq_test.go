package mq

import (
	"github.com/dreamlu/gt/conf"
	"github.com/dreamlu/gt/src/cons"
	"log"
	"testing"
	"time"
)

type Notify struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func TestNsg(t *testing.T) {

	m := NewNSQ(conf.Get[string](cons.ConfNsqProducerAddr), conf.Get[string](cons.ConfNsqConsumerAddr)) // or m := NewProducer(new(NSQ))
	//m.Pub("b", 123)
	m.Pub("b", Notify{
		Name:    "名称",
		Content: "内容",
	})
	m.MultiPub("b2", "哈", "呵")

	err := m.Sub("b", "b-channel", B).Error()
	if err != nil {
		t.Error(err)
		return
	}
	err = m.Sub("b", "c-channel", B).Error()
	if err != nil {
		t.Error(err)
		return
	}
	//m.Stop()
	time.Sleep(15 * time.Second)
}

func B(message *Message) error {

	log.Println(string(message.Body))
	return nil
}
