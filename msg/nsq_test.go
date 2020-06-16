package msg

import (
	"github.com/dreamlu/gt"
	"testing"
	"time"
)

func TestNsg(t *testing.T) {

	go func() {

		for {
			m := NewProducer() // or m := NewProducer(new(Nsg))
			m.Pub("b", 123)
			m.MultiPub("b2", "哈", "呵")
			time.Sleep(2 * time.Second)
		}
	}()

	c := NewConsumer("b", "b-channel")
	err := c.Sub(B)
	if err != nil {
		t.Log(err)
		return
	}
	time.Sleep(15 * time.Second)

}

func B(message *Message) error {

	gt.Logger().Info(string(message.Body))
	return nil
}
