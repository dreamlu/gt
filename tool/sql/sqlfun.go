package sql

import (
	"bytes"
	"github.com/dreamlu/gt/tool/util/str"
	"reflect"
	"strings"
)

// get tag
func GtTag(structTag reflect.StructTag, curTag string) (oTag, tag string, b bool) {
	gtTag := structTag.Get("gt")
	tag = curTag
	if gtTag == "" {
		return
	}
	gtFields := strings.Split(gtTag, ";")
	for _, v := range gtFields {
		// gt:"sub_sql"
		if v == str.GtSubSQL {
			b = true
			return
		}
		// gt:"field:xx"
		oTag = curTag
		if strings.Contains(v, str.GtField) {
			tagTmp := strings.Split(v, ":")
			tag = tagTmp[1]
			return
		}
	}
	return
}

// 层级递增解析tag
func GetReflectTags(ref reflect.Type) (tags []string) {
	var (
		tag string
		b   bool
	)
	for i := 0; i < ref.NumField(); i++ {
		tag = ref.Field(i).Tag.Get("json")
		if tag == "" {
			tags = append(tags, GetReflectTags(ref.Field(i).Type)...)
			continue
		}
		if _, tag, b = GtTag(ref.Field(i).Tag, tag); b == true {
			continue
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
		//buf  bytes.Buffer
	)

	for _, key := range keys {
		if key == "" {
			continue
		}
		sqlKey += "("
		for _, tag := range tags {
			switch {
			// 排除_id结尾字段
			case !strings.HasSuffix(tag, "_id") &&
				!strings.HasPrefix(tag, "id"):
				sqlKey += "`" + alias + "`.`" + tag + "` like binary ? or "
				argsKey = append(argsKey, "%"+key+"%")
			}

		}
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
