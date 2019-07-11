package sql

import (
	"reflect"
	"strings"
)

// 根据model中表模型的json标签获取表字段
// 将select* 变为对应的字段名
func GetTags(model interface{}) (tags []string) {
	typ := reflect.TypeOf(model)
	//var buffer bytes.Buffer
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("json")
		tags = append(tags, tag)
	}
	return tags
}

// key search sql
func GetKeySQL(key string, model interface{}, alias string) (sqlKey, sqlNtKey string, argsKey []interface{}) {

	var (
		tags = GetTags(model)
		keys = strings.Split(key, " ") // 空格隔开
		//buf,bufNt bytes.Buffer
	)

	for _, key := range keys {
		if key == "" {
			continue
		}
		sqlKey += "("
		sqlNtKey += "("
		for _, tag := range tags {
			switch {
			// 排除id结尾字段
			// 排除date,time结尾字段
			case !strings.HasSuffix(tag, "id"): //&& !strings.HasSuffix(tag, "date") && !strings.HasSuffix(tag, "time"):
				sqlKey += "`" + alias + "`.`" + tag + "` like binary ? or "
				sqlNtKey += "`" + alias + "`.`" + tag + "` like binary ? or "
				argsKey = append(argsKey, "%"+key+"%")
			}

		}
		sqlKey = string([]byte(sqlKey)[:len(sqlKey)-4]) + ") and "
		sqlNtKey = string([]byte(sqlNtKey)[:len(sqlNtKey)-4]) + ") and "
	}
	return
}

// 多张表, 第一个表为主表
// key search sql
// tables [table1:table1_alias]
// searModel : 搜索字段模型
func GetMoreKeySQL(key string, searModel interface{}, tables ...string) (sqlKey, sqlNtKey string, argsKey []interface{}) {

	// 搜索字段
	tags := GetTags(searModel)
	// 多表
	keys := strings.Split(key, " ") // 空格隔开
	for _, key := range keys {
		if key == "" {
			continue
		}
		sqlKey += "("
		sqlNtKey += "("
		for _, tag := range tags {
			switch {
			// 排除id结尾字段
			// 排除date,time结尾字段
			case !strings.HasSuffix(tag, "id") && !strings.HasSuffix(tag, "date") && !strings.HasSuffix(tag, "time"):

				// 多表判断
				for _, v := range tables {
					ts := strings.Split(v, ":")
					table := ts[0]
					alias := ts[1]
					if strings.Contains(tag, table+"_") && !strings.Contains(tag, table+"_id") {
						sqlKey += "`" + alias + "`.`" + string([]byte(tag)[len(table)+1:]) + "` like binary ? or "
						sqlNtKey += "`" + alias + "`.`" + string([]byte(tag)[len(table)+1:]) + "` like binary ? or "
						argsKey = append(argsKey, "%"+key+"%")
						goto into
					}
				}

				// 主表
				ts := strings.Split(tables[0], ":")
				alias := ts[1]
				sqlKey += "`" + alias + "`.`" + tag + "` like binary ? or "
				sqlNtKey += "`" + alias + "`.`" + tag + "` like binary ? or "
				argsKey = append(argsKey, "%"+key+"%")
			}
		into:
		}
		sqlKey = string([]byte(sqlKey)[:len(sqlKey)-4]) + ") and "
		sqlNtKey = string([]byte(sqlNtKey)[:len(sqlNtKey)-4]) + ") and "
	}
	return
}
