package tag

import (
	"github.com/dreamlu/gt/tool/util/cons"
	"reflect"
	"strings"
)

// GtTags All GT tags corresponding to a field
type GtTags struct {
	FieldName string   `json:"field_name"`
	GtTags    []*GtTag `json:"tags"`
}

// GtTag A GT tag
type GtTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// IsGtTagIgnore can determine gt-tags whether you do not need to parse
func IsGtTagIgnore(tag reflect.StructTag) bool {
	return IsTagIgnore(tag, cons.GT, cons.GtIgnore, cons.GtSubSQL)
}

// IsTagIgnore can determine tags whether you do not need to parse
func IsTagIgnore(tag reflect.StructTag, tagName string, extTags ...string) bool {
	tagValue, ok := tag.Lookup(tagName)
	return !ok || strings.EqualFold(tagValue, "-") || equalFolds(tagValue, extTags...)
}

// gt:"-"
// gt:"ignore"
// gt:"sub_sql"
// gt:"excel:NAME"
// gt:"field:table.column"
// gt:"field:table.column;excel:NAME"
// GetGtTags use to analyze and obtain GT tags in the structure model
func GetGtTags(model interface{}) map[string]GtTags {
	return ParseGtTags(reflect.TypeOf(model), IsGtTagIgnore)
}

// ParseGtTags use to get all gt tags
func ParseGtTags(ref reflect.Type, fs ...func(reflect.StructTag) bool) map[string]GtTags {
	var (
		tags map[string]string
		res  = make(map[string]GtTags)
	)
	tags = ObtainTags(ref, "gt", fs...)
	for k, v := range tags {
		gtTag := parseGtTag(v)
		gtTag.FieldName = k
		res[k] = gtTag
	}
	return res
}

// parseGtTag use to parse tag value of gt
func parseGtTag(tagValue string) GtTags {
	var tags GtTags
	tagValues := strings.Split(tagValue, ";")
	for _, value := range tagValues {
		var tag GtTag
		if strings.Contains(value, ":") {
			kv := strings.Split(value, ":")
			tag.Name = kv[0]
			tag.Value = kv[1]
			tags.GtTags = append(tags.GtTags, &tag)
		}
	}
	return tags
}

// gt:"-"
// gt:"ignore"
// gt:"sub_sql"
// GetJsonTags use to analyze and obtain JSON tags in the structure model, but it will ignore the ignored value of gt
func GetJsonTags(model interface{}) map[string]string {
	return ParseJsonTags(reflect.TypeOf(model), func(tag reflect.StructTag) bool {
		return IsTagIgnore(tag, "json")
	})
}

func ParseJsonTags(ref reflect.Type, fs ...func(reflect.StructTag) bool) map[string]string {
	var (
		tags map[string]string
		res  = make(map[string]string)
	)
	tags = ObtainTags(ref, "json", fs...)
	for k, v := range tags {
		res[k] = parseJsonTag(v)
	}
	return res
}

func parseJsonTag(tagValue string) string {
	if !strings.Contains(tagValue, ",") {
		return tagValue
	}
	return strings.Split(tagValue, ",")[0]
}

// ObtainTags use to get the specified tag in the structure
// fs use to filter specified tags, true means filtering
func ObtainTags(ref reflect.Type, tagName string, fs ...func(reflect.StructTag) bool) map[string]string {
	if ref.Kind() != reflect.Struct {
		return make(map[string]string)
	}
	var (
		field reflect.StructField
		tag   reflect.StructTag
		tags  = make(map[string]string)
	)
	for i := 0; i < ref.NumField(); i++ {
		field = ref.Field(i)
		if field.Anonymous {
			tags = mergeMap(tags, ObtainTags(field.Type, tagName, fs...))
			continue
		}
		tag = field.Tag
		if len(fs) == 0 {
			tags[field.Name] = tag.Get(tagName)
		} else {
			var b = true
			for _, f := range fs {
				b = b && !f(tag)
				if !b {
					break
				}
			}
			if b {
				tags[field.Name] = tag.Get(tagName)
			}
		}
	}
	return tags
}

// mergeMap use to merge more map slice
func mergeMap(ma ...map[string]string) map[string]string {
	m := make(map[string]string)
	for _, map1 := range ma {
		for k, v := range map1 {
			m[k] = v
		}
	}
	return m
}

// equalFolds Determine whether the strings are equal
func equalFolds(s string, str ...string) bool {
	for _, v := range str {
		if strings.EqualFold(s, v) {
			return true
		}
	}
	return false
}
