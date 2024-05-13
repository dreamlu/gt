package excel

import (
	"github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/type/tmap"
	"github.com/xuri/excelize/v2"
	"io"
)

func (f *Excel[T]) Read(r io.Reader, opts ...excelize.Options) (*Excel[T], error) {
	var err error
	f.File, err = excelize.OpenReader(r, opts...)
	if err != nil {
		return nil, err
	}
	f.rows, err = f.GetRows(f.sheet)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f, nil
}

func (f *Excel[T]) Import() (datas []*T, err error) {
	var (
		titles = f.rows[0]
		title  = tmap.NewTMap[string, int]()
		max    = len(titles)
	)
	for k, colCell := range titles {
		title.Set(colCell, k)
	}

	for i := 1; i < len(f.rows); i++ {

		row := f.rows[i]
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
