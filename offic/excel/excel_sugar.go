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

func Import[T comparable](r io.Reader, opts ...excelize.Options) (datas []T, err error) {
	e := NewExcel[T]()
	err, datas = e.Import(r)
	return
}
