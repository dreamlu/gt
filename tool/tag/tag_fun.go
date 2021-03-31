package tag

import (
	"github.com/dreamlu/gt/tool/util/cons"
	"github.com/dreamlu/gt/tool/util/hump"
	"reflect"
	"strings"
)

// gt tag
// is gt ignore or sub_sql
// Deprecated: please use IsGtTagIgnore()
func IsGtIgnore(tag reflect.StructTag) bool {

	return IsGtTagIgnore(tag)
}

// gt:"field:table.column"
func ParseGtTag(sTag reflect.StructTag) (tag, tagTable string, b bool) {

	if IsGtTagIgnore(sTag) {
		b = true
		return
	}
	tagValue := sTag.Get(cons.GT)
	if tagValue == "" {
		return
	}
	gtFields := strings.Split(tagValue, ";")
	for _, v := range gtFields {
		if strings.Contains(v, cons.GtField) {
			tagTable, tag = parseFieldTag(v)
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

// get struct model fields tag []string
func GetTags(model interface{}) (arr []string) {
	tags := ObtainMoreTags(reflect.TypeOf(model), []string{"json", cons.GT}, IsGtTagIgnore)
	jsonTags := tags["json"]
	gts := tags[cons.GT]
	for k, jt := range jsonTags {
		var (
			tag = jt
			gt  = gts[k]
		)
		if strings.Contains(gt, cons.GtField) {
			t, c := parseFieldTag(gt)
			tag = t + "_" + c
		} else if tag == "" || tag == cons.Gt_ {
			tag = hump.HumpToLine(k)
		}
		arr = append(arr, tag)
	}
	//arr = getTags(reflect.TypeOf(model))
	return
}

func parseFieldTag(tagValue string) (table, column string) {
	tagTmp := strings.Split(tagValue, ":")
	tag := tagTmp[1]
	if a := strings.Split(tag, "."); len(a) > 1 { // include table
		table = a[0]
		column = a[1]
	}
	return
}

// get struct fields tags via recursion
//
// Deprecated
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

		if tag, tagTable, b = ParseGtTag(ref.Field(i).Tag); b {
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
