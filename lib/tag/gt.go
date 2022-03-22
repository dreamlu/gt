package tag

import (
	"github.com/dreamlu/gt/lib/cons"
	mr "github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/type/amap"
	"reflect"
	"strings"
)

func ParseGt(model any) (gt map[string]amap.AMap) {
	gt = make(map[string]amap.AMap)
	parseGt(mr.TrueTypeof(model), gt)
	return
}

func parseGt(typ reflect.Type, gt map[string]amap.AMap) {

	if !mr.IsStruct(typ) {
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			parseGt(typ.Field(i).Type, gt)
			continue
		}
		name, _, _, _ := ParseTag(field)
		gt[name] = ParseGtField(typ.Field(i))
	}

	return
}

func ParseGtField(field reflect.StructField) (gf amap.AMap) {
	gf = amap.NewAMap()
	tv := field.Tag.Get(cons.GT)
	fs := strings.Split(tv, ";")
	for _, f := range fs {
		gtFields := strings.Split(f, ":")
		if len(gtFields) > 1 {
			gf[gtFields[0]] = gtFields[1]
		} else {
			gf[gtFields[0]] = cons.GtExist
		}
	}
	return
}

func ParseGtFieldV(field reflect.StructField) string {
	f := ParseGtField(field)
	return f.Get(cons.GtField)
}

func ParseGtValidV(field reflect.StructField) string {
	f := ParseGtField(field)
	return f.Get(cons.GtValid)
}

func ParseGtTransV(field reflect.StructField) string {
	f := ParseGtField(field)
	return f.Get(cons.GtTrans)
}

func ParseGtLikeV(field reflect.StructField) string {
	f := ParseGtField(field)
	return f.Get(cons.GtLike)
}

// parseGtFieldRule gt:"field:table.column"
func parseGtFieldRule(tag string) (table, column string) {
	if a := strings.Split(tag, "."); len(a) > 1 { // include table
		table = a[0]
		column = a[1]
		return
	}
	// only tag
	return "", tag
}
