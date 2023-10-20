package excel

import (
	"bytes"
	"github.com/dreamlu/gt/src/type/time"
	"os"
	"strconv"
	"testing"
)

type User struct {
	ID     int        `json:"id" gt:"excel:id"`
	Name   string     `json:"name" gt:"excel:名称"`
	Gender int        `json:"gender"`
	Date   time.CDate `json:"date" gt:"excel:日期"`
}

func TestExport(t *testing.T) {
	var arr []*User
	for i := 0; i < 10; i++ {
		arr = append(arr, &User{
			ID:   i,
			Name: "测试" + strconv.Itoa(i),
			Date: time.CDateNow(),
		})
	}
	e, err := Export[User](arr)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(e.SaveAs("test.xlsx"))
}

func (User) Namet(users []*User) error {
	for _, user := range users {
		user.Gender = 1
	}
	return nil
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

func TestExportZip(t *testing.T) {
	var arr []*User
	for i := 0; i < 10; i++ {
		arr = append(arr, &User{
			ID:   i,
			Name: "测试" + strconv.Itoa(i),
			Date: time.CDateNow(),
		})
	}
	e1, _ := Export[User](arr)
	e1.FileName = "e1.xlsx"
	e2, _ := Export[User](arr)
	e2.FileName = "e2.xlsx"

	// 1.bytes file stream
	var bf = bytes.NewBuffer(nil)
	t.Log(ExportZip[User](bf, []*Excel[User]{e1, e2}))

	f, _ := os.Create("test.zip")
	t.Log(ExportZip[User](f, []*Excel[User]{e1, e2}))
}
