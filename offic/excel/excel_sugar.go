package excel

import (
	"github.com/xuri/excelize/v2"
	"io"
)

func Export[T comparable](data any) (e *Excel[T], err error) {
	e = NewExcel[T]()
	err = e.Export(data)
	return
}

// Import excel data
// support type: string, int, int64, uint, float64
// if you want to handle imported data, please implement Handle interface
// eg: func (User) ExcelHandle(users []*User) {}
func Import[T comparable](r io.Reader, opts ...excelize.Options) (datas []*T, err error) {
	e := NewExcel[T]()
	err, datas = e.Import(r)
	return
}
