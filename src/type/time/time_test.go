package time

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestTime(t *testing.T) {

	ti := time.Now()
	t.Log(ti)
	te := CTime(time.Now()).String()
	t.Log(te)
	var tt CTime
	t.Log(tt.MarshalJSON())
	t.Log(tt.String())
	_ = tt.UnmarshalJSON([]byte(te))
	t.Log(tt)
	err := tt.UnmarshalJSON([]byte(`"2022-07-28"`))
	t.Log(err, tt)
	err = tt.UnmarshalJSON([]byte(`"2023-04-03T08:59:32.254Z"`))
	t.Log(err, tt)

	var h Hello
	err = json.Unmarshal([]byte(`{"msg":"hello","time":"2022-07-28 10"}`), &h)
	t.Log(err, h)
}

type Hello struct {
	Msg  string `json:"msg"`
	Time CTime  `json:"time"`
}

// test time Marshal
func TestWebTime(t *testing.T) {
	http.HandleFunc("/", sayhelloName)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w,  CTime(time.Now())) // 这个写入到 w 的是输出到客户端的
	b, err := json.Marshal(Hello{
		Msg:  "hello",
		Time: CTimeNow(),
	})
	log.Print(err)
	w.Write(b)
}

func TestCSTime(t *testing.T) {
	ti := time.Now()
	t.Log(ti)
	te := CSTime(time.Now())
	t.Log(te)
	t.Log(te.String())
	teb, _ := te.MarshalJSON()
	t.Log(string(teb))
	var tt CSTime
	t.Log(tt.MarshalJSON())
	t.Log(tt.String())
	_ = tt.UnmarshalJSON(teb)
	t.Log(tt)
}

func TestCYM(t *testing.T) {
	ti := time.Now()
	t.Log(ti)
	te := CYM(time.Now())
	t.Log(te)
	t.Log(te.String())
	teb, _ := te.MarshalJSON()
	t.Log(string(teb))
	var tt CYM
	t.Log(tt.MarshalJSON())
	t.Log(tt.String())
	_ = tt.UnmarshalJSON(teb)
	t.Log(tt)
}
