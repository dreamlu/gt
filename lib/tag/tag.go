package tag

import (
	"github.com/dreamlu/gt/lib/cons"
	"github.com/dreamlu/gt/lib/hump"
	mr "github.com/dreamlu/gt/src/reflect"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
)

// IsGtTagIgnore can determine gt-tags whether you do not need to parse
func IsGtTagIgnore(tag reflect.StructTag) bool {
	return IsTagIgnore(tag, cons.GT, false, cons.GtIgnore, cons.GtSubSQL)
}

// ParseTag parse tag
// GtTagIgnore and ParseTagOnly
func ParseTag(field reflect.StructField) (tag, tagTable, jsonTag string, b bool) {

	// ignore
	if IsGtTagIgnore(field.Tag) {
		b = true
		return
	}
	tag, tagTable, jsonTag = ParseTagOnly(field)
	return
}

// ParseTagOnly parse tag
// gt:"field:table.field"
// gorm:"column:field"
// json:"field"
// gt > gorm > json > struct field
func ParseTagOnly(field reflect.StructField) (tag, tagTable, jsonTag string) {

	// gt
	tag, tagTable = ParseGtFieldTag(field)
	// gorm
	if tag == "" {
		tag, tagTable = ParseGormFieldTag(field)
	}
	// json
	jsonTag = ParseJsonFieldTag(field)
	// tag still empty
	if tag == "" {
		tag = jsonTag
	}
	return
}

// ParseGtFieldTag gt:"field:table.column"
func ParseGtFieldTag(field reflect.StructField) (tag, tagTable string) {

	if v := ParseGtFieldV(field); v != "" {
		tagTable, tag = parseGtFieldRule(v)
	}
	return
}

// ParseGormFieldTag gorm:"column:field"
func ParseGormFieldTag(sTag reflect.StructField) (tag, tagTable string) {
	tagValue := sTag.Tag.Get(cons.GtGorm)
	if tagValue == "" {
		return
	}
	gtFields := strings.Split(tagValue, ";")
	for _, v := range gtFields {
		if strings.Contains(v, cons.GtGormColumn) {
			tagTable, tag = parseFieldTag(v)
		}
	}
	return
}

// ParseJsonFieldTag get json field tag
// if no, use HumpToLine
func ParseJsonFieldTag(field reflect.StructField) string {

	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		tag = hump.HumpToLine(field.Name)
	}
	// json tag opt `json:"name,opt1,opt2,opts..."`
	tag = strings.Split(tag, ",")[0]
	return tag
}

// GetTags get struct model fields tag []string
// include gt tag rule
func GetTags(model any) (arr []string) {
	return getTags(reflect.TypeOf(model))
}

// GetPartTags remove some like id,_id
func GetPartTags(model any) (arr []string) {
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

func ModelTable(model any) string {
	if t, ok := model.(schema.Tabler); ok {
		return t.TableName()
	}
	return hump.HumpToLine(mr.Name(model))
}
