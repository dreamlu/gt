package tag

import (
	mr "github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/tag"
	"github.com/dreamlu/gt/src/util"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
)

// ParseTag parse tag
// GtTagIgnore and ParseTagOnly
func ParseTag(field reflect.StructField) (_tag, tagTable, jsonTag string, b bool) {

	// ignore
	if tag.IsGtTagIgnore(field.Tag) {
		b = true
		return
	}
	_tag, tagTable, jsonTag = ParseTagOnly(field)
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
		tag = ParseGormFieldTag(field)
	}
	// json
	jsonTag = ParseJsonFieldTag(field)

	if tag == "" {
		tag = jsonTag
	}
	if tag == "" {
		tag = ParseDefaultFieldTag(field)
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
func ParseGormFieldTag(field reflect.StructField) string {
	return tag.ParseGormFieldTag(field, tag.IsGormTagIgnore).Top()
}

// ParseJsonFieldTag get json field tag
func ParseJsonFieldTag(field reflect.StructField) string {
	return tag.ParseJsonFieldTag(field, tag.IsJsonTagIgnore).Top()
}

func ParseDefaultFieldTag(field reflect.StructField) string {
	return util.HumpToLine(field.Name)
}

// GetTags get struct model fields tag []string
// include gt tag rule
func GetTags(model any) (arr []string) {
	return getTags(reflect.TypeOf(model))
}

// GetKeyTags get sql key tags
func GetKeyTags(model any) (arr []string) {
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
	return util.HumpToLine(mr.Name(model))
}

// ============ split ============

// get struct fields tags via recursion
// include gt tag rule
func getTags(typ reflect.Type) (tags []string) {
	typ = mr.TrueType(typ)
	if !mr.IsStruct(typ) {
		return
	}
	var (
		tag string
		b   bool
	)
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Anonymous {
			tags = append(tags, getTags(typ.Field(i).Type)...)
			continue
		}
		tag, _, _, b = ParseTag(typ.Field(i))
		if b {
			continue
		}
		tags = append(tags, tag)
	}
	return tags
}
