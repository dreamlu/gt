package sql

import (
	"bytes"
	"fmt"
	reflect2 "github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/cons"
	"github.com/dreamlu/gt/tool/util/hump"
	"github.com/dreamlu/gt/tool/util/tag"
	"reflect"
	"strings"
)

type Key struct {
	sql     string
	argsNum int
}

var (
	// TableCols table columns map
	TableCols = cmap.NewCMap()
	sqlBuffer = make(map[string]Key)
)

// copy and
func keyAnd(keys []string, buf *bytes.Buffer, num int) (argsKey []interface{}) {
	var (
		sqlKey = buf.String()
		kn     = len(keys)
	)
	for i := 0; i < kn; i++ {
		if i > 0 {
			buf.WriteString(sqlKey)
		}
		for k := 0; k < num; k++ {
			argsKey = append(argsKey, "%"+keys[i]+"%")
		}
	}
	return
}

// GetKeySQL key search sql
func GetKeySQL(key string, model interface{}, alias string) (sqlKey string, argsKey []interface{}) {

	var (
		keys = strings.Fields(key)
		typ  = reflect.TypeOf(model)
		ks   = typ.PkgPath() + typ.Name()
	)
	sqlKey = sqlBuffer[ks].sql
	if sqlKey != "" {
		var buf = bytes.NewBuffer([]byte(sqlKey))
		argsKey = keyAnd(keys, buf, sqlBuffer[ks].argsNum)
		sqlKey = buf.String()
		return
	}

	var (
		tags = tag.GetPartTags(model)
		buf  = bytes.NewBuffer(nil)
		v    = "%" + keys[0] + "%"
	)

	buf.WriteString("(")
	for _, t := range tags {
		buf.WriteString(alias)
		buf.WriteString(".`")
		buf.WriteString(t)
		buf.WriteString("` like binary ? or ")
		argsKey = append(argsKey, v)
	}
	buf = bytes.NewBuffer(buf.Bytes()[:buf.Len()-4])
	buf.WriteString(") and ")

	// add sqlBuffer
	sqlBuffer[ks] = Key{
		sql:     buf.String(),
		argsNum: len(argsKey),
	}

	// copy and
	argsKey = keyAnd(keys, buf, len(argsKey))
	sqlKey = buf.String()
	return
}

// GetMoreKeySQL more tables
// key search sql
// tables [table1:table1_alias]
// searModel : Model
func GetMoreKeySQL(key string, model interface{}, tables ...string) (sqlKey string, argsKey []interface{}) {

	var (
		keys = strings.Split(key, " ") // 空格隔开
		typ  = reflect.TypeOf(model)
		ks   = typ.PkgPath() + "/more/" + typ.Name()
	)
	sqlKey = sqlBuffer[ks].sql
	if sqlKey != "" {
		var buf = bytes.NewBuffer([]byte(sqlKey))
		argsKey = keyAnd(keys, buf, sqlBuffer[ks].argsNum)
		sqlKey = buf.String()
		return
	}

	var (
		//tags = tag.GetPartTags(model)
		buf = bytes.NewBuffer(nil)
		v   = "%" + keys[0] + "%"
	)
	buf.WriteString("(")
	getTagMore(reflect.TypeOf(model), v, &argsKey, buf, tables...)
	buf = bytes.NewBuffer(buf.Bytes()[:buf.Len()-4])
	buf.WriteString(") and ")

	// add sqlBuffer
	sqlBuffer[ks] = Key{
		sql:     buf.String(),
		argsNum: len(argsKey),
	}

	// copy and
	argsKey = keyAnd(keys, buf, len(argsKey))
	sqlKey = buf.String()
	return
}

// more tables
// get sql tag alias recursion
func getTagMore(ref reflect.Type, v string, argsKey *[]interface{}, buf *bytes.Buffer, tables ...string) {

	var (
		tg, tagTable string
		b            bool
	)

	if ref.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < ref.NumField(); i++ {
		if ref.Field(i).Anonymous {
			getTagMore(ref.Field(i).Type, v, argsKey, buf, tables[:]...)
			continue
		}
		if tg, tagTable, _, b = tag.ParseTag(ref.Field(i)); b {
			continue
		}
		// gt tg rule
		if tagTable != "" {
			writeTagString(buf, tagTable, tg)
			*argsKey = append(*argsKey, v)
			continue
		}

		// sql tg rule
		tb := UniqueTagTable(tg, tables...)
		if tb != "" {
			writeTagString(buf, tb, tg)
			*argsKey = append(*argsKey, v)
			continue
		}

		// default tg rule
		if b = otherTableTagSQL(tg, argsKey, buf, tables...); !b {
			writeTagString(buf, tables[0], tg)
			*argsKey = append(*argsKey, v)
		}
	}
	return
}

// if there is tag gt and json, select json tag first
// more tables, other tables sql tag
func otherTableTagSQL(tag string, argsKey *[]interface{}, buf *bytes.Buffer, tables ...string) bool {
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

			writeTagString(buf, v, string([]byte(tag)[len(v)+1:]))
			*argsKey = append(*argsKey, v)
			return true
		}
	}
	return false
}

// write tag sql
func writeTagString(buf *bytes.Buffer, tb, tag string) {
	buf.WriteString("`")
	buf.WriteString(tb)
	buf.WriteString("`.`")
	buf.WriteString(tag)
	buf.WriteString("` like binary ? or ")
}

// StructWhereSQL struct to where sql
// return key1 = value1 and key2 = value2...
func StructWhereSQL(st interface{}) (sql string, args []interface{}) {
	var (
		buf bytes.Buffer
		m   = reflect2.ToMap(st)
	)

	for k, v := range m {
		if reflect2.IsZero(v) {
			continue
		}
		buf.WriteString(hump.HumpToLine(k))
		buf.WriteString(cons.ParamAnd)
		args = append(args, v)
	}
	if len(m) > 0 {
		sql = string(buf.Bytes()[:len(buf.Bytes())-5])
	}
	return
}

// Table table parse
func Table(table string) string {

	if table == "" {
		return table
	}
	if table[0] == '`' {
		return table
	}
	if strings.Contains(table, ".") {
		ts := strings.Split(table, ".")
		table = fmt.Sprintf("`%s`.`%s`", ts[0], ts[1])
		return table
	}

	return fmt.Sprintf("`%s`", table)
}

// TableOnly only table name
func TableOnly(table string) string {

	if table == "" {
		return table
	}
	if strings.Contains(table, ".") {
		table = strings.Split(table, ".")[1]
	}
	if strings.Contains(table, ":") {
		table = strings.Split(table, ":")[0]
	}
	return table
}

// UniqueTagTable return unique tag table
func UniqueTagTable(tag string, tables ...string) string {
	tbs := TagTables(tag, tables...)
	if len(tbs) == 1 {
		return tbs[0]
	}
	return ""
}

// TagTables return tag tables
func TagTables(tag string, tables ...string) (tbs []string) {
	for _, v := range tables {
		tags := TableCols[v]
		for _, t := range tags {
			if t == tag {
				tbs = append(tbs, v)
			}
		}
	}
	return
}
