package time

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {

	te := CTime(time.Now()).String()
	t.Log(te)
	var tt CTime
	_ = tt.UnmarshalJSON([]byte(te))
	t.Log(tt.String())
}
