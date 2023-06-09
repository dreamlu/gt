package tag

import (
	"github.com/dreamlu/gt/crud/dep/cons"
	. "github.com/dreamlu/gt/src/cons"
	"github.com/dreamlu/gt/src/cons/tag"
	. "github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/type/amap"
	"reflect"
)

// IsGtTagIgnore can determine gt-tags whether you do not need to parse
func IsGtTagIgnore(tag reflect.StructTag) bool {
	return IsTagIgnore(tag, GT, false, cons.GtIgnore, cons.GtSubSQL)
}

// IsJsonTagIgnore can determine gt-tags whether you do not need to parse
func IsJsonTagIgnore(t reflect.StructTag) bool {
	return IsTagIgnore(t, tag.Json, true)
}

// IsGormTagIgnore can determine gt-tags whether you do not need to parse
func IsGormTagIgnore(t reflect.StructTag) bool {
	return IsTagIgnore(t, tag.Gorm, true)
}

// GetGtTags
// gt:"-"
// gt:"ignore"
// gt:"sub_sql"
// gt:"excel:NAME"
// gt:"field:table.column"
// gt:"field:table.column;excel:NAME"
// GetGtTags use to analyze and obtain GT tags in the structure model
func GetGtTags(model any) GF[GtTags] {
	return ParseGtTags(TrueTypeof(model), IsGtTagIgnore)
}

// GetJsonTags analyze json tag
func GetJsonTags(model any) amap.AMap {
	return ParseJsonTags(TrueTypeof(model), IsJsonTagIgnore).ToAMap()
}

// GetGormTags analyze gorm tag
func GetGormTags(model any) amap.AMap {
	return ParseGormTags(TrueTypeof(model), IsGormTagIgnore).ToAMap()
}
