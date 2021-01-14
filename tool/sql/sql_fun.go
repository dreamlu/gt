package sql

import (
	"bytes"
	"fmt"
	reflect2 "github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/tag"
	"reflect"
	"strings"
)

type Key struct {
	sql     string
	argsNum int
}

var sqlBuffer = make(map[string]Key)

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
		keys = strings.Split(key, " ") // 空格隔开
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
		// 排除_id结尾字段
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
// searModel : 搜索字段模型
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
		// 排除_id结尾字段
		if !strings.HasSuffix(t, "_id") &&
			!strings.HasPrefix(t, "id") {

			if b := otherTableKeySql(t, buf, tables...); b == true {
				argsKey = append(argsKey, "%"+v+"%")
				continue
			}

			// 主表
			ts := strings.Split(tables[0], ":")
			alias := ts[1]
			buf.WriteString("`")
			buf.WriteString(alias)
			buf.WriteString("`.`")
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

func otherTableKeySql(tag string, buf *bytes.Buffer, tables ...string) (b bool) {
	// 多表判断
	for _, v := range tables {
		ts := strings.Split(v, ":")
		table := ts[0]
		alias := ts[1]
		if strings.Contains(tag, table+"_") && !strings.Contains(tag, table+"_id") {
			buf.WriteString("`")
			buf.WriteString(alias)
			buf.WriteString("`.`")
			buf.WriteString(string([]byte(tag)[len(table)+1:]))
			buf.WriteString("` like binary ? or ")
			b = true
			return
		}
	}
	return
}

// struct to where sql
// return key1 = value1 and key2 = value2...
func StructWhereSQL(st interface{}) (sql string, args []interface{}) {
	var (
		buf bytes.Buffer
		m   = reflect2.ToMap(st)
	)

	for k, v := range m {
		//if v == "" {
		//	continue
		//}
		buf.WriteString(k)
		buf.WriteString(" = ? and ")
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
	if []byte(table)[0] == '`' {
		return table
	}
	if strings.Contains(table, ".") {
		ts := strings.Split(table, ".")
		table = fmt.Sprintf("`%s`.`%s`", ts[0], ts[1])
		return table
	}

	return fmt.Sprintf("`%s`", table)
}
