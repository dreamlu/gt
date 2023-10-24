package reflect

import (
	"testing"
)

// order
type Order struct {
	ID         int   `json:"id"`
	UserID     int64 `json:"user_id"`     // user id
	ServiceID  int64 `json:"service_id"`  // service table id
	CreateTime int64 `json:"create_time"` // createtime
}

func TestReflect(t *testing.T) {
	or := New(Order{})
	t.Log(or)
	ors := NewArray(Order{})
	t.Log(ors)
}

func TestGetDataID(t *testing.T) {
	or := Order{} //new(Order)
	//var a = 23
	or.ID = 23
	id := Field(or, "ID")
	t.Log(id)
	id = TrueField(or, "ID")
	t.Log(id)
}

func TestStructToString(t *testing.T) {
	type TestDA struct {
	}
	t.Log(Name(TestDA{}))
}

func TestPath(t *testing.T) {
	typ := TrueTypeof(Order{})
	t.Log(Path(typ, "A", "B"))
}

func TestTrueValueOf(t *testing.T) {
	t.Log(TrueValueOf(&Order{}))
	t.Log(TrueValueOf(Order{}))
	t.Log(TrueValueOf([]Order{}))
}

func TestSet(t *testing.T) {
	or := Order{} //new(Order)
	Set(&or, "ID", int(3))
	t.Log(or)

	var i any
	ot := Order{} //new(Order)
	i = &ot
	Set(i, "ID", int(4))
	t.Log(ot)

	oi := Order{} //new(Order)
	i = &oi
	SetByIndex(i, 0, int(4))
	t.Log(oi)
}

func TestCall(t *testing.T) {
	or := Order{}
	Call(or, "NoMethod")
}
