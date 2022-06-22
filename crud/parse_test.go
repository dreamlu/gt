package crud

import "testing"

func TestParse(t *testing.T) {
	r := parse(OrderD{}, "order", "service")
	r = parse(OrderD{}, "order", "service")
	t.Log(r)
}
