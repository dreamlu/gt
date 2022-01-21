package gt

import (
	"bytes"
	mr "github.com/dreamlu/gt/tool/reflect"
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

// GetMoreTableColumnSQL select * replace
// select more tables
// tables : table name / table alias name
// first table must main table, like from a inner join b, first table is a
func GetMoreTableColumnSQL(model interface{}, tables ...string) (sql string) {
	var (
		typ = mr.TrueTypeof(model)
		key = mr.Path(typ, "more")
	)
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
func getTagMore(typ reflect.Type, buf *bytes.Buffer, tables ...string) {

	var (
		oTag, tag, tagTable string
		b                   bool
	)
	typ = mr.TrueType(typ)
	if !mr.IsStruct(typ) {
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Anonymous {
			getTagMore(typ.Field(i).Type, buf, tables[:]...)
			continue
		}
		if tag, tagTable, oTag, b = ParseTag(typ.Field(i)); b {
			continue
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

// GetColSQLAlias select * replace
// add alias
func GetColSQLAlias(model interface{}, alias string) (sql string) {
	var (
		typ = mr.TrueTypeof(model)
		key = mr.Path(typ, alias)
	)
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
func getTagAlias(typ reflect.Type, buf *bytes.Buffer, alias string) {

	if !mr.IsStruct(typ) {
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Anonymous {
			getTagAlias(typ.Field(i).Type, buf, alias)
			continue
		}
		if tag, b := getRTag(typ, i); !b {
			buf.WriteString(alias)
			buf.WriteString(".`")
			buf.WriteString(tag)
			buf.WriteString("`,")
		} // continue
	}
}

// GetColSQL select * replace
func GetColSQL(model interface{}) (sql string) {
	var (
		typ = mr.TrueTypeof(model)
		key = mr.Path(typ)
	)
	sql = sqlBuffer.Get(key)
	if sql != "" {
		return
	}
	var buf bytes.Buffer
	getTag(typ, &buf)
	sql = string(buf.Bytes()[:buf.Len()-1]) // remove ,
	sqlBuffer.Set(key, sql)
	return
}

// single table
// get sql tag recursion
func getTag(typ reflect.Type, buf *bytes.Buffer) {

	if !mr.IsStruct(typ) {
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Anonymous {
			getTag(typ.Field(i).Type, buf)
			continue
		}
		if tag, b := getRTag(typ, i); !b {
			buf.WriteString("`")
			buf.WriteString(tag)
			buf.WriteString("`,")
		} // continue
	}
}

func getRTag(ref reflect.Type, i int) (tag string, b bool) {
	if tag, _, _, b = ParseTag(ref.Field(i)); b {
		return "", true
	}
	return tag, b
}

// GetColParamSQL get col ?
func GetColParamSQL(model interface{}) (sql string) {
	var buf bytes.Buffer
	getColParamSQLByType(reflect.TypeOf(model), &buf)
	return string(buf.Bytes()[:buf.Len()-1])
}

// get col ? type
func getColParamSQLByType(typ reflect.Type, buf *bytes.Buffer) {
	typ = mr.TrueType(typ)
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Anonymous {
			getColParamSQLByType(typ.Field(i).Type, buf)
			continue
		}
		buf.WriteString("?,")
	}
}

// GetParams get single struct data value
func GetParams(data interface{}) (params []interface{}) {
	params = append(params, getParams(reflect.ValueOf(data))...)
	return
}

// get params ? params
func getParams(typ reflect.Value) (params []interface{}) {
	// todo mr.TrueType() use go1.18 replace
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	// todo mr.IsStruct() use go1.18 replace
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
