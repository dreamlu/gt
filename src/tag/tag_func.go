package tag

import (
	"github.com/dreamlu/gt/crud/dep/cons"
	"github.com/dreamlu/gt/src/cons/tag"
	"reflect"
	"strings"
)

// IsTagIgnore can determine tags whether you do not need to parse
func IsTagIgnore(tag reflect.StructTag, tagName string, exist bool, extTags ...string) bool {
	// exist ok return
	// 1     1  -
	// 1     0  ture
	// 0     1  -
	// 0     0  -
	tagValue, ok := tag.Lookup(tagName)
	if exist && !ok {
		return true
	}
	return strings.EqualFold(tagValue, "-") || equalFolds(tagValue, extTags...)
}

// ParseGtTags use to get all gt tags
func ParseGtTags(typ reflect.Type, fs ...func(reflect.StructTag) bool) GF[GtTags] {
	return GetTypFor[GtTags](typ,
		func(typ reflect.Type) GF[GtTags] { return ParseGtTags(typ, fs...) },
		func(field reflect.StructField) GF[GtTags] { return ParseGtFieldTag(field, fs...) },
	)
}

func ParseGtFieldTag(field reflect.StructField, fs ...func(reflect.StructTag) bool) GF[GtTags] {
	return GetFieldTag(field, tag.Gt, fs...).ToGtTag().Parse(parseGtTag)
}

// parseGtTag use to parse tag value of gt
func parseGtTag(g GtTags) GtTags {
	var tags GtTags
	tagValues := strings.Split(g.origin, tag.Semicolon)
	for _, value := range tagValues {
		var t GtTag
		kv := strings.Split(value, tag.Colon)
		t.Name = kv[0]
		if len(kv) == 1 {
			t.Value = cons.GtExist
		} else {
			t.Value = kv[1]
		}
		tags.GtTags = append(tags.GtTags, &t)
	}
	return tags
}

// ========= json ==========

func ParseJsonTags(typ reflect.Type, fs ...func(reflect.StructTag) bool) GF[string] {
	return GetTypFor[string](typ,
		func(typ reflect.Type) GF[string] { return ParseJsonTags(typ, fs...) },
		func(field reflect.StructField) GF[string] { return ParseJsonFieldTag(field, fs...) },
	)
}

func ParseJsonFieldTag(field reflect.StructField, fs ...func(reflect.StructTag) bool) GF[string] {
	return GetFieldTag(field, tag.Json, fs...).Parse(parseJsonTag)
}

func parseJsonTag(tagValue string) string {
	if !strings.Contains(tagValue, tag.Comma) {
		return tagValue
	}
	return strings.Split(tagValue, tag.Comma)[0]
}

// ========= gorm ==========

func ParseGormTags(typ reflect.Type, fs ...func(reflect.StructTag) bool) GF[string] {
	return GetTypFor[string](typ,
		func(typ reflect.Type) GF[string] { return ParseGormTags(typ, fs...) },
		func(field reflect.StructField) GF[string] { return ParseGormFieldTag(field, fs...) },
	)
}

// ParseGormFieldTag only support gorm:column
func ParseGormFieldTag(field reflect.StructField, fs ...func(reflect.StructTag) bool) GF[string] {
	return GetFieldTag(field, tag.Gorm, fs...).Parse(parseGormTag)
}

func parseGormTag(tagValue string) string {
	gtFields := strings.Split(tagValue, tag.Semicolon)
	for _, v := range gtFields {
		if strings.Contains(v, cons.GtGormColumn) {
			ts := strings.Split(tagValue, tag.Colon)
			if len(ts) == 1 {
				break
			}
			return ts[1]
		}
	}
	return ""
}
