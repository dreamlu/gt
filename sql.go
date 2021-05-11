package gt

import (
	"bytes"
	"fmt"
	"github.com/dreamlu/gt/tool/type/cmap"
	sq "github.com/dreamlu/gt/tool/util/sql"
	. "github.com/dreamlu/gt/tool/util/tag"
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
	sqlBuffer.Set(key, sql)
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
		if tag, tagTable, b = ParseGtTag(ref.Field(i).Tag); b {
			continue
		}
		oTag = GetFieldTag(ref.Field(i))
		if tag == "" {
			tag = oTag
		}
		// gt tag rule
		if tagTable != "" {
			writeBufTag(buf, tagTable, tag, oTag)
			continue
		}

		// sql tag rule
		tb := sq.UniqueTagTable(tag, tables...)
		if tb != "" {
			writeBufTag(buf, tb, tag, "")
			continue
		}

		// default tag rule
		if b = otherTableTagSQL(oTag, tag, buf, tables...); !b {
			writeBufTag(buf, tables[0], tag, "")
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
			// next two condition, eg: db_test.go==>TestGetReflectTagMore()
			!strings.Contains(tag, "_id") &&
			!strings.EqualFold(v, tables[0]) {

			writeBufTag(buf, v, string([]byte(tag)[len(v)+1:]), oTag)
			return true
		}
	}
	return false
}

// write tag sql
func writeBufTag(buf *bytes.Buffer, tb, tag, alias string) {
	buf.WriteString("`")
	buf.WriteString(tb)
	buf.WriteString("`.`")
	buf.WriteString(tag)
	buf.WriteString("`")
	if alias != "" && alias != "-" {
		buf.WriteString(" as ")
		buf.WriteString(alias)
	}
	buf.WriteString(",")
}

// write where tag sql
func writeBufWhere(buf *bytes.Buffer, tb, tag string) {
	buf.WriteString("`")
	buf.WriteString(tb)
	buf.WriteString("`.`")
	buf.WriteString(tag)
	buf.WriteString("` = ? and ")
}

// select * replace
// add alias
func GetColSQLAlias(model interface{}, alias string) (sql string) {
	typ := reflect.TypeOf(model)
	key := fmt.Sprintf("%s%s-%s", typ.PkgPath(), typ.Name(), alias)
	sql = sqlBuffer.Get(key)
	if sql != "" {
		return
	}
	var buf bytes.Buffer
	getTagAlias(typ, &buf, alias)
	sql = string(buf.Bytes()[:buf.Len()-1]) //去掉点,
	sqlBuffer.Set(key, sql)
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
		if IsGtTagIgnore(ref.Field(i).Tag) {
			continue
		}

		tag := GetFieldTag(ref.Field(i))
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
		return
	}
	var buf bytes.Buffer
	getTag(reflect.TypeOf(model), &buf)
	sql = string(buf.Bytes()[:buf.Len()-1]) // remove ,
	sqlBuffer.Set(key, sql)
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
		if IsGtTagIgnore(ref.Field(i).Tag) {
			continue
		}
		tag := GetFieldTag(ref.Field(i))
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
			getColParamSQLByType(typ.Field(i).Type, buf)
			continue
		}
		buf.WriteString("?,")
	}
}

// get single struct data value
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

func innerLeftSQL(bufNt *bytes.Buffer, DBS map[string]string, tables, fields []string, i int) {
	if tb := DBS[tables[i]]; tb != "" {
		bufNt.WriteString("`" + tb + "`.")
	}
	bufNt.WriteString("`")
	bufNt.WriteString(tables[i])
	bufNt.WriteString("` on ")
	fieldSQL(bufNt, tables[i-1], tables[i], fields[i-1], fields[i])
}

// fieldSQL field analyze
func fieldSQL(bufNt *bytes.Buffer, leftTable, rightTable, left, right string) {

	var (
		ils  = strings.Split(left, ",")  // left condition column
		irs  = strings.Split(right, ",") // right condition column
		ilts []string                    // left table field condition
		irts []string                    // right table field condition
	)

	for k := 0; k < len(ils); k++ {
		is := strings.Split(ils[k], "=")
		if len(is) > 1 {
			ilts = append(ilts, ils[k])
			ils = append(ils[:k], ils[k+1:]...)
			k--
		}
	}

	for k := 0; k < len(irs); k++ {
		is := strings.Split(irs[k], "=")
		if len(is) > 1 {
			irts = append(irts, irs[k])
			irs = append(irs[:k], irs[k+1:]...)
			k--
		}
	}

	for k := 0; k < len(ils); k++ {
		bufNt.WriteByte('`')
		bufNt.WriteString(leftTable)
		bufNt.WriteString("`.`")
		bufNt.WriteString(ils[k])
		bufNt.WriteString("`=`")
		bufNt.WriteString(rightTable)
		bufNt.WriteString("`.`")
		bufNt.WriteString(irs[k])
		bufNt.WriteString("` and ")
	}

	for _, v := range ilts {
		is := strings.Split(v, "=")
		if len(is) > 1 {
			bufNt.WriteByte('`')
			bufNt.WriteString(leftTable)
			bufNt.WriteString("`.`")
			bufNt.WriteString(is[0])
			bufNt.WriteByte('`')
			bufNt.WriteByte('=')
			bufNt.WriteString(is[1])
			bufNt.WriteString(" and ")
		}
	}
	for _, v := range irts {
		is := strings.Split(v, "=")
		if len(is) > 1 {
			bufNt.WriteByte('`')
			bufNt.WriteString(rightTable)
			bufNt.WriteString("`.`")
			bufNt.WriteString(is[0])
			bufNt.WriteByte('`')
			bufNt.WriteByte('=')
			bufNt.WriteString(is[1])
			bufNt.WriteString(" and ")
		}
	}
	nBuf := bytes.NewBuffer(bufNt.Bytes()[:bufNt.Len()-4])
	defer nBuf.Reset()
	bufNt.Reset()
	bufNt.Write(nBuf.Bytes())
}
