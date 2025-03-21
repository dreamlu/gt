package excel

import (
	"github.com/dreamlu/gt/src/cons/excel"
	"github.com/dreamlu/gt/src/reflect"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func (f *Excel[T]) Export(data any) (e *Excel[T], err error) {

	ch, preCh, pre := f.sheetCellCharInit()
	f.File = excelize.NewFile()

	for _, header := range f.Headers {
		err = f.SetCellValue(excel.Sheet, pre+string(ch)+"1", header)
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
			v = reflect.TrueField(value, f.HeaderMapper[col])
			err = f.SetCellValue(excel.Sheet, pre+string(ch)+num, v)
			if err != nil {
				return
			}

			f.sheetCellCharChange(&ch, &preCh, &pre)
		}
	}
	e = f
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
