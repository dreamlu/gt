package msg

import (
	"github.com/dreamlu/gt/serv/log"
	"testing"
	"time"
)

type Notify struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func TestNsg(t *testing.T) {

	m := NewProducer() // or m := NewProducer(new(Nsg))
	//m.Pub("b", 123)
	m.Pub("b", Notify{
		Name:    "名称",
		Content: "内容",
	})
	m.MultiPub("b2", "哈", "呵")

	c := NewConsumer("b", "b-channel")
	err := c.Sub(B)
	if err != nil {
		t.Log(err)
		return
	}

	c = NewConsumer("b", "c-channel")
	err = c.Sub(B)
	if err != nil {
		t.Log(err)
		return
	}
	c.Stop()

	time.Sleep(15 * time.Second)
}

func B(message *Message) error {

	log.Info(string(message.Body))
	return nil
}
