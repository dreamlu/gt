package excel

import (
	"strconv"
	"testing"
)

func TestExportExcel(t *testing.T) {
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
	e, err := ExportExcel(User{}, arr)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(e.SaveAs("1.xlsx"))
}
