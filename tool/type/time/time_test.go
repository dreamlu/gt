package time

import (
	"encoding/json"
	"fmt"
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
	t.Log(tt)
	fmt.Println(tt)
	_ = tt.UnmarshalJSON([]byte(te))
	t.Log(tt)
}

// test time Marshal
func TestWebTime(t *testing.T) {
	http.HandleFunc("/", sayhelloName)       // 设置访问的路由
	err := http.ListenAndServe(":9090", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w,  CTime(time.Now())) // 这个写入到 w 的是输出到客户端的
	b, err := json.Marshal(CTime{})
	log.Print(err)
	w.Write(b)
}
