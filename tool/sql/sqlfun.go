package sql

import (
	"bytes"
	"fmt"
	reflect2 "github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/util/cons"
	"reflect"
	"strings"
)

// get tag
func GtTag(structTag reflect.StructTag, curTag string) (oTag, tag, tagTable string, b bool) {
	gtTag := structTag.Get("gt")
	tag = curTag
	if gtTag == "" {
		return
	}
	gtFields := strings.Split(gtTag, ";")
	for _, v := range gtFields {
		// gt:"sub_sql"
		if v == cons.GtSubSQL ||
			v == cons.GtIgnore {
			b = true
			return
		}
		// gt:"field:xx"
		oTag = curTag
		if strings.Contains(v, cons.GtField) {
			tagTmp := strings.Split(v, ":")
			tag = tagTmp[1]
			if a := strings.Split(tag, "."); len(a) > 1 { // include table
				tag = a[1]
				tagTable = a[0]
			}
			return
		}
	}
	return
}

// 层级递增解析tag
func GetReflectTags(ref reflect.Type) (tags []string) {
	if ref.Kind() != reflect.Struct {
		return
	}
	var (
		tag, tagTable string
		b             bool
	)
	for i := 0; i < ref.NumField(); i++ {
		tag = ref.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			tags = append(tags, GetReflectTags(ref.Field(i).Type)...)
			continue
		}
		if _, tag, tagTable, b = GtTag(ref.Field(i).Tag, tag); b == true {
			continue
		}
		if tagTable != "" {
			tag = tagTable + "_" + tag
		}
		tags = append(tags, tag)
	}
	return tags
}

// 根据model中表模型的json标签获取表字段
// 将select* 变为对应的字段名
func GetTags(model interface{}) []string {
	return GetReflectTags(reflect.TypeOf(model))
}

// key search sql
func GetKeySQL(key string, model interface{}, alias string) (sqlKey string, argsKey []interface{}) {

	var (
		tags = GetTags(model)
		keys = strings.Split(key, " ") // 空格隔开
		buf  bytes.Buffer
	)

	for _, key := range keys {
		if key == "" {
			continue
		}
		buf.WriteString("(")
		for _, tag := range tags {
			switch {
			// 排除_id结尾字段
			case !strings.HasSuffix(tag, "_id") &&
				!strings.HasPrefix(tag, "id"):
				buf.WriteString(alias)
				buf.WriteString(".`")
				buf.WriteString(tag)
				buf.WriteString("` like binary ? or ")
				argsKey = append(argsKey, "%"+key+"%")
			}

		}
		sqlKey = buf.String()
		sqlKey = string([]byte(sqlKey)[:len(sqlKey)-4]) + ") and "
	}
	return
}

// 多张表, 第一个表为主表
// key search sql
// tables [table1:table1_alias]
// searModel : 搜索字段模型
func GetMoreKeySQL(key string, searModel interface{}, tables ...string) (sqlKey string, argsKey []interface{}) {

	var (
		tags = GetTags(searModel)
		keys = strings.Split(key, " ") // 空格隔开
		buf  bytes.Buffer
	)
	for _, key := range keys {
		if key == "" {
			continue
		}
		buf.WriteString("(")
		for _, tag := range tags {
			// 排除_id结尾字段
			if !strings.HasSuffix(tag, "_id") &&
				!strings.HasPrefix(tag, "id") {

				if b := otherTableKeySql(tag, &buf, tables...); b == true {
					argsKey = append(argsKey, "%"+key+"%")
					continue
				}

				// 主表
				ts := strings.Split(tables[0], ":")
				alias := ts[1]
				buf.WriteString("`")
				buf.WriteString(alias)
				buf.WriteString("`.`")
				buf.WriteString(tag)
				buf.WriteString("` like binary ? or ")
				argsKey = append(argsKey, "%"+key+"%")
			}
		}
		sqlKey = buf.String()
		sqlKey = string([]byte(sqlKey)[:len(sqlKey)-4]) + ") and "
	}
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

//// cmap to where sql
//// return key1 = value1 and key2 = value2...
//func CMapWhereSQL(cm cmap.CMap) (sql string, args []interface{}) {
//	var buf bytes.Buffer
//	for k, v := range cm {
//		buf.WriteString(k)
//		buf.WriteString(" = ? and ")
//		args = append(args, v)
//	}
//	if len(cm) > 0 {
//		sql = string(buf.Bytes()[:len(buf.Bytes())-5])
//	}
//	return
//}
