package excel

import (
	"bytes"
	"os"
	"strconv"
	"testing"
)

func TestExport(t *testing.T) {
	type User struct {
		ID     int    `json:"id" gt:"excel:id"`
		Name   string `json:"name" gt:"excel:名称"`
		Gender int    `json:"gender"`
	}
	var arr []*User
	for i := 0; i < 10; i++ {
		arr = append(arr, &User{
			ID:   i,
			Name: "测试" + strconv.Itoa(i),
		})
	}
	e, err := Export[User](arr)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(e.SaveAs("test.xlsx"))
}

type User struct {
	ID     int64  `json:"id" gt:"excel:id"`
	Name   string `json:"name" gt:"excel:名称"`
	Gender int    `json:"gender"`
}

func (User) ExcelHandle(users []*User) error {
	for _, user := range users {
		user.Gender = 1
	}
	return nil
}

func TestImport(t *testing.T) {
	bs, _ := os.ReadFile("test.xlsx")
	r := bytes.NewReader(bs)
	datas, err := Import[User](r)
	if err != nil {
		t.Log(err)
		return
	}
	for _, user := range datas {
		t.Log(user)
	}
}
