package excel

import (
	"github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/type/tmap"
	"github.com/xuri/excelize/v2"
	"io"
)

func (f *Excel[T]) Read(r io.Reader, opts ...excelize.Options) (*Excel[T], error) {
	return f.read(r, false, opts...)
}

func (f *Excel[T]) ReadAll(r io.Reader, opts ...excelize.Options) (*Excel[T], error) {
	return f.read(r, true, opts...)
}

func (f *Excel[T]) read(r io.Reader, allSheets bool, opts ...excelize.Options) (*Excel[T], error) {
	var err error
	f.File, err = excelize.OpenReader(r, opts...)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetList := f.GetSheetList()
	if allSheets {
		f.sheets = sheetList
	} else {
		f.sheets = append(f.sheets, sheetList[0])
	}
	for i, sheet := range f.sheets {
		var rows [][]string
		rows, err = f.GetRows(sheet)
		f.rows[sheet] = rows
		if i == 0 {
			f.Titles = append(f.Titles, rows[0]...)
		}
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

func (f *Excel[T]) open(r io.Reader, opts ...excelize.Options) (err error) {
	f.File, err = excelize.OpenReader(r, opts...)
	return
}

func (f *Excel[T]) Import() (datas []*T, err error) {
	for _, sheet := range f.sheets {
		var ds []*T
		ds, err = f.importSheet(sheet)
		if err != nil {
			return
		}
		datas = append(datas, ds...)
	}
	return
}

func (f *Excel[T]) importSheet(sheet string) (datas []*T, err error) {
	var (
		title = tmap.NewTMap[string, int]()
		max   = len(f.Titles)
	)
	for k, colCell := range f.Titles {
		title.Set(colCell, k)
	}

	for i := 1; i < len(f.rows[sheet]); i++ {

		row := f.rows[sheet][i]
		for len(row) < max {
			row = append(row, "")
		}
		var data T
		for k, v := range f.ExcelMapper {
			var cell = row[title.Get(v)]
			if cell == "" {
				continue // zero value
			}
			if fc := f.dict.Get(v); fc != nil {
				var value any
				value, err = fc(v, cell)
				if err != nil {
					return
				}
				reflect.Set(&data, k.Field, value)
				continue
			}
			if !title.IsExist(v) {
				continue
			}
			reflect.Set(&data, k.Field, string2any(k.Type, cell))
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

func (f *Excel[T]) AddDict(key string, value dict) *Excel[T] {
	f.dict.Set(key, value)
	return f
}
