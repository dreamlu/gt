package excel

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/util/tag"
	"strconv"
)

type Excel struct {
	*excelize.File
	Data    interface{}
	Model   interface{}
	Headers []string
	Mapper  map[string]string
	sheet   string
	index   int
	x       rune
	y       rune
	Value   Mapper
}

type Mapper func(interface{}, string) (interface{}, error)

func NewExcel(model interface{}) *Excel {
	f := excelize.NewFile()
	m, h := getMapper(model)
	return &Excel{
		File:    f,
		Model:   model,
		Mapper:  m,
		Headers: h,
		sheet:   "Sheet1",
		x:       'A',
		y:       '1',
		Value:   reflect.GetDataByFieldName,
	}
}

// Point use to set the starting point for exporting data
func (f *Excel) Point(str string) {
	if len(str) > 0 {
		f.x = rune(str[0])
		if len(str) > 1 {
			f.y = rune(str[1])
		}
	}
}

func (f *Excel) Export(data interface{}) (err error) {

	ch := f.x
	row := string(f.y)
	for _, header := range f.Headers {
		err = f.SetCellValue(f.sheet, string(ch)+row, header)
		if err != nil {
			return
		}
		ch++
	}

	arr := reflect.ToSlice(data)
	// 设置宽度样式
	//_ = f.SetColWidth(St, "B", "I", 18)

	// Set active sheet of the workbook.
	if f.index != 0 {
		f.SetActiveSheet(f.index)
	}

	for i, value := range arr {
		num := strconv.Itoa(i + 2)
		ch = f.x
		for _, col := range f.Headers {
			var v interface{}
			v, err = f.Value(value, f.Mapper[col])
			if err != nil {
				return
			}
			err = f.SetCellValue(f.sheet, string(ch)+num, v)
			if err != nil {
				return
			}
			ch++
		}
	}
	return
}

func getMapper(model interface{}) (map[string]string, []string) {
	var (
		mapper  = make(map[string]string)
		headers []string
	)

	tags := tag.GetGtTags(model)
	for k, v := range tags {
		for _, t := range v.GtTags {
			if t.Name == "excel" {
				mapper[k] = t.Value
				mapper[t.Value] = k
				headers = append(headers, mapper[k])
			}
		}
	}
	return mapper, headers
}
