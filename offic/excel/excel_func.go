package excel

import (
	"github.com/dreamlu/gt/crud/dep/tag"
	"github.com/dreamlu/gt/src/cons/excel"
	"github.com/dreamlu/gt/src/type/amap"
	"strconv"
)

func getMapper(model any) ([]string, amap.AMap, map[tag.GtField]string) {
	var (
		headerMapper = amap.AMap{}
		excelMapper  = make(map[tag.GtField]string)
		headers      []string
	)

	tags := tag.GetGtTags(model)
	for k, v := range tags {
		for _, t := range v.GtTags {
			if t.Name == excel.Excel {
				headerMapper[t.Value] = k
				excelMapper[v.Field] = t.Value
				headers = append(headers, t.Value)
			}
		}
	}
	return headers, headerMapper, excelMapper
}

func string2any(typ, cell string) any {
	var value any
	switch typ {
	case "int":
		value, _ = strconv.Atoi(cell)
	case "int64":
		value, _ = strconv.ParseInt(cell, 10, 64)
	case "uint":
		value, _ = strconv.ParseUint(cell, 10, 64)
	case "float64":
		value, _ = strconv.ParseFloat(cell, 64)
	default:
		value = cell
	}
	return value
}
