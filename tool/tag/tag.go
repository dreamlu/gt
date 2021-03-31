package tag

import (
	"github.com/dreamlu/gt/tool/util/cons"
	"reflect"
	"strings"
)

// gtTags All GT tags corresponding to a field
type gtTags struct {
	FieldName string
	GtTags    []*gtTag
}

// gtTag A GT tag
type gtTag struct {
	Name  string
	Value string
}

// IsGtTagIgnore can determine gt-tags whether you do not need to parse
func IsGtTagIgnore(tag reflect.StructTag) bool {
	return IsTagIgnore(tag, cons.GT, false, cons.GtIgnore, cons.GtSubSQL)
}

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

// gt:"-"
// gt:"ignore"
// gt:"sub_sql"
// gt:"excel:NAME"
// gt:"field:table.column"
// gt:"field:table.column;excel:NAME"
// GetGtTags use to analyze and obtain GT tags in the structure model
func GetGtTags(model interface{}) map[string]gtTags {
	return ParseGtTags(reflect.TypeOf(model), IsGtTagIgnore)
}

// ParseGtTags use to get all gt tags
func ParseGtTags(ref reflect.Type, fs ...func(reflect.StructTag) bool) map[string]gtTags {
	var (
		tags map[string]string
		res  = make(map[string]gtTags)
	)
	tags = ObtainTags(ref, cons.GT, fs...)
	for k, v := range tags {
		tag := parseGtTag(v)
		tag.FieldName = k
		res[k] = tag
	}
	return res
}

// parseGtTag use to parse tag value of gt
func parseGtTag(tagValue string) gtTags {
	var tags gtTags
	tagValues := strings.Split(tagValue, ";")
	for _, value := range tagValues {
		var tag gtTag
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
		return IsTagIgnore(tag, "json", true)
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
	return ObtainMoreTags(ref, []string{tagName}, fs...)[tagName]
}

// ObtainTags use to get the specified tag in the structure
// fs use to filter specified tags, true means filtering
func ObtainMoreTags(ref reflect.Type, tagNames []string, fs ...func(reflect.StructTag) bool) map[string]map[string]string {
	if ref.Kind() != reflect.Struct {
		return make(map[string]map[string]string)
	}
	var (
		field reflect.StructField
		tag   reflect.StructTag
		res   = make(map[string]map[string]string)
	)
	for _, tagName := range tagNames {
		var tags = make(map[string]string)
		for i := 0; i < ref.NumField(); i++ {
			field = ref.Field(i)
			if field.Anonymous {
				tags = mergeMap(tags, ObtainTags(field.Type, tagName, fs...))
				continue
			}
			tag = field.Tag
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
		res[tagName] = tags
	}
	return res
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
