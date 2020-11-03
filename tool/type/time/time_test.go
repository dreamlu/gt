package time

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {

	ti := time.Now()
	t.Log(ti)
	te := CTime(time.Now()).String()
	t.Log(te)
	var tt CTime
	t.Log(tt)
	fmt.Println(tt)
	_ = tt.UnmarshalJSON([]byte(te))
	t.Log(tt)
}
