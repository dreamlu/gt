package excel

import (
	"github.com/xuri/excelize/v2"
	"io"
)

// Import excel data
// support type: string, int, int64, uint, float64
// if you want to handle imported data, please implement Handle interface
// eg: func (User) ExcelHandle(users []*User) {}
func Import[T comparable](r io.Reader, opts ...excelize.Options) (e []*T, err error) {
	excel, err := NewExcel[T]().Read(r, opts...)
	if err != nil {
		return nil, err
	}
	return excel.Import()
}

// Export excel data
func Export[T comparable](data any) (e *Excel[T], err error) {
	return NewExcel[T]().Export(data)
}

// ExportZip export excel data to zip
func ExportZip[T comparable](dst io.Writer, excels []*Excel[T]) error {
	return exportZip(dst, excels)
}
