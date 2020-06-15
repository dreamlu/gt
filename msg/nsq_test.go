package msg

import (
	"testing"
	"time"
)

func TestNsg(t *testing.T) {

	go func() {

		for {
			m := NewProducer() // or m := NewProducer(new(Nsg))
			m.Pub("b", 123)
			time.Sleep(2 * time.Second)
		}
	}()

	c := NewConsumer("b", "b-channel")
	err := c.Sub(HandlerFunc(func(message *Message) error {

		t.Log(string(message.Body))
		return nil
	}))
	if err != nil {
		t.Log(err)
		return
	}
	time.Sleep(15 * time.Second)

}
