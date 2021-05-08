package tag

import (
	"github.com/dreamlu/gt/tool/util/cons"
	"github.com/dreamlu/gt/tool/util/hump"
	"reflect"
	"strings"
)

// IsGtTagIgnore can determine gt-tags whether you do not need to parse
func IsGtTagIgnore(tag reflect.StructTag) bool {
	return IsTagIgnore(tag, cons.GT, false, cons.GtIgnore, cons.GtSubSQL)
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
	// problem: map no sort -> struct field sort
	//tags := ObtainMoreTags(reflect.TypeOf(model), []string{"json", cons.GT}, IsGtTagIgnore)
	//jsonTags := tags["json"]
	//gts := tags[cons.GT]
	//for k, jt := range jsonTags {
	//	var (
	//		tag = jt
	//		gt  = gts[k]
	//	)
	//	if strings.Contains(gt, cons.GtField) {
	//		t, c := parseFieldTag(gt)
	//		tag = t + "_" + c
	//	} else if tag == "" || tag == cons.Gt_ {
	//		tag = hump.HumpToLine(k)
	//	}
	//	arr = append(arr, tag)
	//}
	return getTags(reflect.TypeOf(model))
}

// GetPartTags remove some like id,_id
func GetPartTags(model interface{}) (arr []string) {
	arr = GetTags(model)
	for i := 0; i < len(arr); i++ {
		v := arr[i]
		if strings.HasSuffix(v, "_id") ||
			strings.HasPrefix(v, "id") {
			arr = append(arr[:i], arr[i+1:]...)
			i--
		}
	}
	return
}
