package tag

import (
	mr "github.com/dreamlu/gt/src/reflect"
	"reflect"
	"strings"
)

// GetTypeTags use to get the specified tag in the structure
// fs use to filter specified tags, true means filtering
func GetTypeTags(typ reflect.Type, tagNames []string, fs ...func(reflect.StructTag) bool) map[string]GF[string] {
	var (
		tags = make(map[string]GF[string])
	)
	for _, tagName := range tagNames {
		tags[tagName] = GetTypeTag(typ, tagName, fs...)
	}
	return tags
}

func GetTypeTag(typ reflect.Type, tagName string, fs ...func(reflect.StructTag) bool) GF[string] {
	//if !mr.IsStruct(typ) {
	//	return nil
	//}
	//var (
	//	field reflect.StructField
	//	tags  = make(GF)
	//)
	//for i := 0; i < typ.NumField(); i++ {
	//	field = typ.Field(i)
	//	if field.Anonymous {
	//		tags = mergeMap(tags, GetTypeTag(field.Type, tagName, fs...))
	//		continue
	//	}
	//	tags = mergeMap(tags, GetFieldTag(field, tagName, fs...))
	//}
	//return tags
	return GetTypFor[string](typ,
		func(typ reflect.Type) GF[string] { return GetTypeTag(typ, tagName, fs...) },
		func(field reflect.StructField) GF[string] { return GetFieldTag(field, tagName, fs...) },
	)
}

func GetTypFor[T any](typ reflect.Type, tf TFunc[T], ff FFunc[T]) GF[T] {
	if !mr.IsStruct(typ) {
		return nil
	}
	var (
		field reflect.StructField
		tags  = make(GF[T])
	)
	for i := 0; i < typ.NumField(); i++ {
		field = typ.Field(i)
		if field.Anonymous {
			tags = mergeMap(tags, tf(field.Type)) // field.Type problem
			continue
		}
		tags = mergeMap(tags, ff(field))
	}
	return tags
}

func GetFieldTags(field reflect.StructField, tagNames []string, fs ...func(reflect.StructTag) bool) GF[string] {
	var tags = make(GF[string])
	for _, tagName := range tagNames {
		tags = mergeMap(tags, GetFieldTag(field, tagName, fs...))
	}
	return tags
}

func GetFieldTag(field reflect.StructField, tagName string, fs ...func(reflect.StructTag) bool) GF[string] {

	var (
		tag    reflect.StructTag
		oneTag = make(GF[string])
	)
	tag = field.Tag
	var b = true
	for _, f := range fs {
		b = b && !f(tag)
		if !b {
			break
		}
	}
	if b {
		oneTag[GtField{Field: field.Name, Type: field.Type.String()}] = tag.Get(tagName) // field.Type.Name() can not get not basic type, eg:*int
	}
	return oneTag
}

// mergeMap use to merge more map slice
func mergeMap[T any](ma ...GF[T]) GF[T] {
	m := make(GF[T])
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
