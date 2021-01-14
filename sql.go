package gt

import (
	"bytes"
	"fmt"
	. "github.com/dreamlu/gt/tool/tag"
	"github.com/dreamlu/gt/tool/type/cmap"
	"reflect"
	"strings"
)

// sql tag reflect
// resolve go struct field from model

// sql memory
var sqlBuffer = cmap.NewCMap()

// select * replace
// select more tables
// tables : table name / table alias name
// first table must main table, like from a inner join b, first table is a
func GetMoreTableColumnSQL(model interface{}, tables ...string) (sql string) {
	typ := reflect.TypeOf(model)
	key := typ.PkgPath() + "/more/" + typ.Name()
	sql = sqlBuffer.Get(key)
	if sql != "" {
		return
	}
	var buf bytes.Buffer
	getTagMore(typ, &buf, tables[:]...)
	sql = string(buf.Bytes()[:buf.Len()-1])
	sqlBuffer.Add(key, sql)
	return
}

// more tables
// get sql tag alias recursion
func getTagMore(ref reflect.Type, buf *bytes.Buffer, tables ...string) {

	var (
		oTag, tag, tagTable string
		b                   bool
	)
	if ref.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < ref.NumField(); i++ {
		if ref.Field(i).Anonymous {
			getTagMore(ref.Field(i).Type, buf, tables[:]...)
			continue
		}
		if tag, tagTable, b = GtTag(ref.Field(i).Tag); b {
			continue
		}
		oTag = GetSQLField(ref.Field(i))
		if tag == "" {
			tag = oTag
		}
		if tagTable != "" {
			buf.WriteString("`")
			buf.WriteString(tagTable)
			buf.WriteString("`.`")
			buf.WriteString(tag)
			//buf.WriteString("`,")
			buf.WriteString("` as ")
			if oTag != "" && oTag != "-" {
				buf.WriteString(oTag)
			} else {
				buf.WriteString(tag)
			}
			buf.WriteString(",")
			continue
		}

		if b = otherTableTagSQL(oTag, tag, buf, tables...); !b {
			buf.WriteString("`")
			buf.WriteString(tables[0])
			buf.WriteString("`.`")
			buf.WriteString(tag)
			buf.WriteString("`,")
		}
	}
}

// if there is tag gt and json, select json tag first
// more tables, other tables sql tag
func otherTableTagSQL(oTag, tag string, buf *bytes.Buffer, tables ...string) bool {
	// foreign tables column
	for _, v := range tables {
		if strings.Contains(tag, v+"_id") {
			break
		}
		// tables
		if strings.HasPrefix(tag, v+"_") &&
			// next two condition see db_test.go==>TestGetReflectTagMore()
			!strings.Contains(tag, "_id") &&
			!strings.EqualFold(v, tables[0]) {
			buf.WriteString("`")
			buf.WriteString(v)
			buf.WriteString("`.`")
			buf.Write([]byte(tag)[len(v)+1:])
			buf.WriteString("` as ")
			if oTag != "" && oTag != "-" {
				buf.WriteString(oTag)
			} else {
				buf.WriteString(tag)
			}
			buf.WriteString(",")
			return true
		}
	}
	return false
}

// select * replace
// add alias
func GetColSQLAlias(model interface{}, alias string) (sql string) {
	typ := reflect.TypeOf(model)
	key := fmt.Sprintf("%s%s-%s", typ.PkgPath(), typ.Name(), alias)
	sql = sqlBuffer.Get(key)
	if sql != "" {
		//Logger().Info("[USE sqlBuffer GET ColumnSQL]")
		return
	}
	var buf bytes.Buffer
	getTagAlias(typ, &buf, alias)
	sql = string(buf.Bytes()[:buf.Len()-1]) //去掉点,
	sqlBuffer.Add(key, sql)
	return
}

// single table
// get sql tag alias recursion
func getTagAlias(ref reflect.Type, buf *bytes.Buffer, alias string) {

	if ref.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < ref.NumField(); i++ {
		if ref.Field(i).Anonymous {
			getTagAlias(ref.Field(i).Type, buf, alias)
			continue
		}
		if IsGtIgnore(ref.Field(i).Tag) {
			continue
		}

		tag := GetSQLField(ref.Field(i))
		buf.WriteString(alias)
		buf.WriteString(".`")
		buf.WriteString(tag)
		buf.WriteString("`,")
	}
}

// select * replace
func GetColSQL(model interface{}) (sql string) {
	typ := reflect.TypeOf(model)
	key := typ.PkgPath() + typ.Name()
	sql = sqlBuffer.Get(key)
	if sql != "" {
		//Logger().Info("[USE sqlBuffer GET ColumnSQL]")
		return
	}
	var buf bytes.Buffer
	//typ := reflect.TypeOf(model)
	getTag(reflect.TypeOf(model), &buf)
	sql = string(buf.Bytes()[:buf.Len()-1]) // remove ,
	sqlBuffer.Add(key, sql)
	return
}

// single table
// get sql tag recursion
func getTag(ref reflect.Type, buf *bytes.Buffer) {

	if ref.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < ref.NumField(); i++ {
		if ref.Field(i).Anonymous {
			getTag(ref.Field(i).Type, buf)
			continue
		}
		if IsGtIgnore(ref.Field(i).Tag) {
			continue
		}
		tag := GetSQLField(ref.Field(i))
		buf.WriteString("`")
		buf.WriteString(tag)
		buf.WriteString("`,")
	}
}

// get col ?
func GetColParamSQL(model interface{}) (sql string) {
	var buf bytes.Buffer
	getColParamSQLByType(reflect.TypeOf(model), &buf)
	return string(buf.Bytes()[:buf.Len()-1])
}

// get col ? type
func getColParamSQLByType(typ reflect.Type, buf *bytes.Buffer) {
	if typ.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Anonymous {
			// why this is error ?
			// typ = typ.Field(i).Type
			// getColParamSQLByType(typ.Field(i).Type, buf)
			getColParamSQLByType(typ.Field(i).Type, buf)
			continue
		}
		buf.WriteString("?,")
	}
}

// get data value
// like GetColSQL
func GetParams(data interface{}) (params []interface{}) {
	params = append(params, getParams(reflect.ValueOf(data))...)
	return
}

// get params ? params
func getParams(typ reflect.Value) (params []interface{}) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		if typ.Type().Field(i).Anonymous {
			params = append(params, getParams(typ.Field(i))...)
			continue
		}
		value := typ.Field(i).Interface()
		params = append(params, value)
	}
	return
}
