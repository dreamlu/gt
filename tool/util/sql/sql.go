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
	// table columns map
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

// key search sql
func GetKeySQL(key string, model interface{}, alias string) (sqlKey string, argsKey []interface{}) {

	var (
		keys = strings.Split(key, " ")
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
		tags = tag.GetTags(model)
		buf  = bytes.NewBuffer(nil)
		v    = keys[0]
	)

	buf.WriteString("(")
	for _, t := range tags {
		switch {
		case !strings.HasSuffix(t, "_id") &&
			!strings.HasPrefix(t, "id"):
			buf.WriteString(alias)
			buf.WriteString(".`")
			buf.WriteString(t)
			buf.WriteString("` like binary ? or ")
			argsKey = append(argsKey, "%"+v+"%")
		}

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

// more tables
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
		tags = tag.GetTags(model)
		buf  = bytes.NewBuffer(nil)
		v    = keys[0]
	)
	buf.WriteString("(")
	for _, t := range tags {
		if !strings.HasSuffix(t, "_id") &&
			!strings.HasPrefix(t, "id") {

			tb := UniqueTagTable(t)
			if tb != "" {
				writeTagString(buf, tb, t)
				argsKey = append(argsKey, "%"+v+"%")
				continue
			}

			if b := otherTableKeySql(t, buf, tables...); b == true {
				argsKey = append(argsKey, "%"+v+"%")
				continue
			}

			// 主表
			ts := strings.Split(tables[0], ":")
			alias := ts[1]
			writeTagString(buf, alias, t)
			argsKey = append(argsKey, "%"+v+"%")
		}
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

// other table key search sql
func otherTableKeySql(tag string, buf *bytes.Buffer, tables ...string) (b bool) {
	for _, v := range tables {
		ts := strings.Split(v, ":")
		table := ts[0]
		alias := ts[1]
		if strings.Contains(tag, table+"_") && !strings.Contains(tag, table+"_id") {
			writeTagString(buf, alias, string([]byte(tag)[len(table)+1:]))
			b = true
			return
		}
	}
	return
}

// write tag sql
func writeTagString(buf *bytes.Buffer, tb, tag string) {
	buf.WriteString("`")
	buf.WriteString(tb)
	buf.WriteString("`.`")
	buf.WriteString(tag)
	buf.WriteString("` like binary ? or ")
}

// struct to where sql
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

// table parse
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

// only table name
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

// return unique table tag
func UniqueTagTable(tag string) (table string) {
	var (
		i = 0
	)

	for k, tags := range TableCols {
		for _, t := range tags {
			if t == tag {
				i++
				table = k
			}
		}
	}
	if i == 1 {
		return
	}
	return ""
}
