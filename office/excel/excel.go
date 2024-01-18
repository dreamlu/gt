package excel

import (
	"github.com/dreamlu/gt/src/cons/excel"
	"github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/tag"
	"github.com/dreamlu/gt/src/type/amap"
	"github.com/dreamlu/gt/src/type/tmap"
	"github.com/xuri/excelize/v2"
	"io"
	"strconv"
)

type Excel[T comparable] struct {
	*excelize.File
	FileName     string
	Data         any
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
	f := excelize.NewFile()
	var model T
	h, m, e := getMapper(model)
	return &Excel[T]{
		File:         f,
		HeaderMapper: m,
		ExcelMapper:  e,
		Headers:      h,
		sheet:        excel.Sheet,
		dict:         tmap.NewTMap[string, dict](),
	}
}

func (f *Excel[T]) Export(data any) (err error) {

	ch, preCh, pre := f.sheetCellCharInit()
	f.File = excelize.NewFile()

	for _, header := range f.Headers {
		err = f.SetCellValue(f.sheet, pre+string(ch)+"1", header)
		if err != nil {
			return
		}
		f.sheetCellCharChange(&ch, &preCh, &pre)
	}

	arr := reflect.ToSlice(data)
	//_ = f.SetColWidth(St, "B", "I", 18)

	for i, value := range arr {
		num := strconv.Itoa(i + 2)
		ch, preCh, pre = f.sheetCellCharInit()
		for _, col := range f.Headers {
			var v any
			v = reflect.Field(value, f.HeaderMapper[col])
			err = f.SetCellValue(f.sheet, pre+string(ch)+num, v)
			if err != nil {
				return
			}

			f.sheetCellCharChange(&ch, &preCh, &pre)
		}
	}
	return
}

func (f *Excel[T]) sheetCellCharInit() (ch, preCh int32, pre string) {
	return 'A', 'A', ""
}

func (f *Excel[T]) sheetCellCharChange(ch, preCh *int32, pre *string) {
	*ch++
	if *ch > 'Z' {
		*ch = 'A'
		if *pre != "" {
			*preCh++
			*pre = string(*preCh)
			return
		}
		*pre = string(*preCh)
	}
	return
}

// ========== import ===========

func (f *Excel[T]) AddDict(key string, value dict) *Excel[T] {
	f.dict.Set(key, value)
	return f
}

func (f *Excel[T]) Import(r io.Reader, opts ...excelize.Options) (datas []*T, err error) {

	f.File, err = excelize.OpenReader(r, opts...)
	if err != nil {
		return
	}
	defer f.Close()
	rows, err := f.GetRows(f.sheet)
	if err != nil {
		return nil, err
	}

	var (
		title = tmap.NewTMap[string, int]()
		max   = len(rows[0])
	)
	for k, colCell := range rows[0] {
		title.Set(colCell, k)
	}

	for i := 1; i < len(rows); i++ {

		row := rows[i]
		for len(row) < max {
			row = append(row, "")
		}
		var data T
		for k, v := range f.ExcelMapper {
			if fc := f.dict.Get(v); fc != nil {
				var value any
				value, err = fc(v, row[title.Get(v)])
				if err != nil {
					return
				}
				reflect.Set(&data, k.Field, value)
				continue
			}
			value := string2any(k.Type, row[title.Get(v)])
			reflect.Set(&data, k.Field, value)
		}
		datas = append(datas, &data)
	}

	// after import
	var data T
	if reflect.IsImplements(data, new(Handle[T])) {
		err = reflect.Call(data, "ExcelHandle", datas)
		if err != nil {
			return
		}
	}

	return
}
