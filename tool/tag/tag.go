package tag

import (
	"github.com/dreamlu/gt/tool/util/cons"
	"github.com/dreamlu/gt/tool/util/hump"
	"reflect"
	"strings"
)

// gt tag
// is gt ignore or sub_sql
func IsGtIgnore(tag reflect.StructTag) bool {

	gtTag := tag.Get(cons.GT)
	if strings.Contains(gtTag, cons.GtSubSQL) ||
		strings.Contains(gtTag, cons.GtIgnore) ||
		strings.Contains(gtTag, cons.Gt_) {
		return true
	}
	return false
}

// gt:"field:table.column"
func GtTag(sTag reflect.StructTag) (tag, tagTable string, b bool) {

	if IsGtIgnore(sTag) {
		b = true
		return
	}
	gtTag := sTag.Get(cons.GT)
	if gtTag == "" {
		return
	}
	gtFields := strings.Split(gtTag, ";")
	for _, v := range gtFields {
		if strings.Contains(v, cons.GtField) {
			tagTmp := strings.Split(v, ":")
			tag = tagTmp[1]
			if a := strings.Split(tag, "."); len(a) > 1 { // include table
				tag = a[1]
				tagTable = a[0]
			}
			return
		}
	}
	return
}

// get json field tag
// if no, use HumpToLine
func GetFieldTag(field reflect.StructField) string {

	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		tag = hump.HumpToLine(field.Name)
	}
	// json tag opt `json:"name,opt1,opt2,opts..."`
	tag = strings.Split(tag, ",")[0]
	return tag
}

// get struct fields tags via recursion
func getTags(ref reflect.Type) (tags []string) {
	if ref.Kind() != reflect.Struct {
		return
	}
	var (
		tag, tagTable string
		b             bool
	)
	for i := 0; i < ref.NumField(); i++ {
		if ref.Field(i).Anonymous {
			tags = append(tags, getTags(ref.Field(i).Type)...)
			continue
		}

		if tag, tagTable, b = GtTag(ref.Field(i).Tag); b {
			continue
		}
		if tag == "" {
			tag = GetFieldTag(ref.Field(i))
		}
		if tagTable != "" {
			tag = tagTable + "_" + tag
		}
		tags = append(tags, tag)
	}
	return tags
}

// get struct model fields tag []string
func GetTags(model interface{}) []string {
	return getTags(reflect.TypeOf(model))
}
