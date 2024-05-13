package excel

import (
	"github.com/dreamlu/gt/src/cons/excel"
	"github.com/dreamlu/gt/src/tag"
	"github.com/dreamlu/gt/src/type/amap"
	"github.com/dreamlu/gt/src/type/tmap"
	"github.com/xuri/excelize/v2"
)

type Excel[T comparable] struct {
	*excelize.File
	rows         [][]string
	FileName     string
	Headers      []string
	HeaderMapper amap.AMap
	ExcelMapper  map[tag.GtField]string
	sheet        string
	index        int
	dict         tmap.TMap[string, dict]
}

type dict func(string, string) (any, error)

type Handle[T comparable] interface {
	ExcelHandle([]*T) error
}

func NewExcel[T comparable]() *Excel[T] {
	var model T
	h, m, e := getMapper(model)
	return &Excel[T]{
		HeaderMapper: m,
		ExcelMapper:  e,
		Headers:      h,
		sheet:        excel.Sheet,
		dict:         tmap.NewTMap[string, dict](),
	}
}

func (f *Excel[T]) SetSheet(sheet string) *Excel[T] {
	f.sheet = sheet
	return f
}
