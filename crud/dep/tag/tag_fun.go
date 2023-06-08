package tag

import (
	"github.com/dreamlu/gt/crud/dep/cons"
	mr "github.com/dreamlu/gt/src/reflect"
	"reflect"
	"strings"
)

// gtTags All GT tags corresponding to a field
type gtTags struct {
	Field  GtField
	GtTags []*GtTag
}

type GtField struct {
	Field string
	Type  string
}

// GtTag A GT tag
type GtTag struct {
	Name  string
	Value string
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

// GetGtTags
// gt:"-"
// gt:"ignore"
// gt:"sub_sql"
// gt:"excel:NAME"
// gt:"field:table.column"
// gt:"field:table.column;excel:NAME"
// GetGtTags use to analyze and obtain GT tags in the structure model
func GetGtTags(model any) map[string]gtTags {
	return ParseGtTags(reflect.TypeOf(model), IsGtTagIgnore)
}

// ParseGtTags use to get all gt tags
func ParseGtTags(ref reflect.Type, fs ...func(reflect.StructTag) bool) map[string]gtTags {
	var (
		tags map[GtField]string
		res  = make(map[string]gtTags)
	)
	tags = ObtainTags(ref, cons.GT, fs...)
	for k, v := range tags {
		tag := parseGtTag(v)
		tag.Field = k
		res[k.Field] = tag
	}
	return res
}

// parseGtTag use to parse tag value of gt
func parseGtTag(tagValue string) gtTags {
	var tags gtTags
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

// GetJsonTags use to analyze and obtain JSON tags in the structure model, but it will ignore the ignored value of json
func GetJsonTags(model any) map[string]string {
	return ParseJsonTags(reflect.TypeOf(model), func(tag reflect.StructTag) bool {
		return IsTagIgnore(tag, "json", true)
	})
}

func ParseJsonTags(ref reflect.Type, fs ...func(reflect.StructTag) bool) map[string]string {
	var (
		tags map[GtField]string
		res  = make(map[string]string)
	)
	tags = ObtainTags(ref, "json", fs...)
	for k, v := range tags {
		res[k.Field] = parseJsonTag(v)
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
func ObtainTags(ref reflect.Type, tagName string, fs ...func(reflect.StructTag) bool) map[GtField]string {
	return ObtainMoreTags(ref, []string{tagName}, fs...)[tagName]
}

// ObtainMoreTags use to get the specified tag in the structure
// fs use to filter specified tags, true means filtering
func ObtainMoreTags(typ reflect.Type, tagNames []string, fs ...func(reflect.StructTag) bool) map[string]map[GtField]string {
	if !mr.IsStruct(typ) {
		return nil
	}
	var (
		field reflect.StructField
		tag   reflect.StructTag
		res   = make(map[string]map[GtField]string)
	)
	for _, tagName := range tagNames {
		var tags = make(map[GtField]string)
		for i := 0; i < typ.NumField(); i++ {
			field = typ.Field(i)
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
				tags[GtField{Field: field.Name, Type: field.Type.Name()}] = tag.Get(tagName)
			}
		}
		res[tagName] = tags
	}
	return res
}

// mergeMap use to merge more map slice
func mergeMap(ma ...map[GtField]string) map[GtField]string {
	m := make(map[GtField]string)
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

func parseFieldTag(tagValue string) (table, column string) {
	ts := strings.Split(tagValue, ":")
	if len(ts) == 1 {
		return
	}
	return parseGtFieldRule(ts[1])
}

// get struct fields tags via recursion
// include gt tag rule
func getTags(typ reflect.Type) (tags []string) {
	typ = mr.TrueType(typ)
	if !mr.IsStruct(typ) {
		return
	}
	var (
		tag string
	)
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Anonymous {
			tags = append(tags, getTags(typ.Field(i).Type)...)
			continue
		}

		if tag, _ = ParseGtFieldTag(typ.Field(i)); tag == "" {
			continue
		}
		if tag == "" {
			tag = ParseJsonFieldTag(typ.Field(i))
		}
		tags = append(tags, tag)
	}
	return tags
}
