// package der

package der

import (
	"bytes"
	"fmt"
	reflect2 "github.com/dreamlu/go-tool/tool/reflect"
	"github.com/dreamlu/go-tool/tool/result"
	sq "github.com/dreamlu/go-tool/tool/sql"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

//======================return tag=============================
//=============================================================

// select * replace
// select more tables
// tables : table name / table alias name
// 主表放在tables中第一个，紧接着为主表关联的外键表名(无顺序)
func GetMoreTableColumnSQL(model interface{}, tables ...string) (sql string) {
	var buf bytes.Buffer

	typ := reflect.TypeOf(model)
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("json")
		// foreign tables column
		for _, v := range tables {
			// tables
			switch {
			case !strings.Contains(tag, v+"_id") && strings.Contains(tag, v+"_"):
				//sql += "`" + v + "`.`" + string([]byte(tag)[len(v)+1:]) + "` as " + tag + ","
				buf.WriteString("`")
				buf.WriteString(v)
				buf.WriteString("`.`")
				buf.Write([]byte(tag)[len(v)+1:])
				buf.WriteString("` as ")
				buf.WriteString(tag)
				buf.WriteString(",")
				goto into
			}
		}
		//sql += "`" + tables[0] + "`.`" + tag + "`,"
		buf.WriteString("`")
		buf.WriteString(tables[0])
		buf.WriteString("`.`")
		buf.WriteString(tag)
		buf.WriteString("`,")
	into:
	}
	sql = string(buf.Bytes()[:buf.Len()-1]) //去点,
	return sql
}

// select *替换
// 两张表
func GetDoubleTableColumnSQL(model interface{}, table1, table2 string) (sql string) {
	var buf bytes.Buffer

	typ := reflect.TypeOf(model)
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("json")
		// table2的数据处理,去除table2_id
		if strings.Contains(tag, table2+"_") && !strings.Contains(tag, table2+"_id") {
			// sql += table2 + ".`" + string([]byte(tag)[len(table2)+1:]) + "` as " + tag + ","
			buf.WriteString(table2)
			buf.WriteString(".`")
			buf.Write([]byte(tag)[len(table2)+1:])
			buf.WriteString("` as ")
			buf.WriteString(tag)
			buf.WriteString(",")
			continue
		}
		buf.WriteString(table1)
		buf.WriteString(".`")
		buf.WriteString(tag)
		buf.WriteString("`,")
		//sql += table1 + ".`" + tag + "`,"
	}
	sql = string(buf.Bytes()[:buf.Len()-1]) //去掉点,
	return sql
}

// 根据model中表模型的json标签获取表字段
// 将select* 中'*'变为对应的字段名
// 增加别名,表连接问题
func GetColAliasSQL(model interface{}, alias string) (sql string) {
	var buf bytes.Buffer

	typ := reflect.TypeOf(model)
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("json")
		//sql += alias + ".`" + tag + "`,"
		buf.WriteString(alias)
		buf.WriteString(".`")
		buf.WriteString(tag)
		buf.WriteString("`,")
	}
	sql = string(buf.Bytes()[:buf.Len()-1]) //去掉点,
	return sql
}

// 根据model中表模型的json标签获取表字段
// 将select* 变为对应的字段名
func GetColSQL(model interface{}) (sql string) {
	var buf bytes.Buffer

	typ := reflect.TypeOf(model)
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("json")
		//sql += "`" + tag + "`,"
		buf.WriteString("`")
		buf.WriteString(tag)
		buf.WriteString("`,")
	}
	sql = string(buf.Bytes()[:buf.Len()-1]) //去掉点,
	return sql
}

// get col ?
func GetColParamSQL(model interface{}) (sql string) {
	var buf bytes.Buffer

	typ := reflect.TypeOf(model)
	for i := 0; i < typ.NumField(); i++ {
		buf.WriteString("?,")
	}
	sql = string(buf.Bytes()[:buf.Len()-1]) //去掉点,
	return sql
}

// get data value
// like GetColSQL
func GetParams(data interface{}) (params []interface{}) {

	typ := reflect.ValueOf(data)
	for i := 0; i < typ.NumField(); i++ {
		value := typ.Field(i).Interface() //.Tag.Get("json")
		params = append(params, value)
	}
	return
}

//=======================================sql语句处理==========================================
//===========================================================================================

// More Table
// params: innerTables is inner join tables
// params: leftTables is left join tables
// return: select sql
// table1 as main table, include other tables_id(foreign key)
func GetMoreSearchSQL(model interface{}, params map[string][]string, innerTables []string, leftTables []string) (sqlNt, sql string, clientPage, everyPage int64, args []interface{}) {

	var (
		order               = "desc"      // order by
		key                 = ""          // key like binary search
		tables              = innerTables // all tables
		bufNt, bufW, bufNtW bytes.Buffer  // sql bytes connect
	)

	tables = append(tables, leftTables...)
	// sql and sqlCount
	bufNt.WriteString("select ")
	bufNt.WriteString("count(`")
	bufNt.WriteString(tables[0])
	bufNt.WriteString("`.id) as total_num ")
	// bufNt.WriteString(GetMoreTableColumnSQL(model, tables[:]...))
	bufNt.WriteString("from `")
	bufNt.WriteString(tables[0])
	bufNt.WriteString("`")
	// inner join
	for i := 1; i < len(innerTables); i++ {
		bufNt.WriteString(" inner join `")
		bufNt.WriteString(innerTables[i])
		bufNt.WriteString("` on `")
		bufNt.WriteString(tables[0])
		bufNt.WriteString("`.")
		bufNt.WriteString(innerTables[i])
		bufNt.WriteString("_id=`")
		bufNt.WriteString(innerTables[i])
		bufNt.WriteString("`.id ")
		//sql += " inner join ·" + innerTables[i] + "`"
	}
	// left join
	for i := 0; i < len(leftTables); i++ {
		bufNt.WriteString(" left join `")
		bufNt.WriteString(leftTables[i])
		bufNt.WriteString("` on `")
		bufNt.WriteString(tables[0])
		bufNt.WriteString("`.")
		bufNt.WriteString(leftTables[i])
		bufNt.WriteString("_id=`")
		bufNt.WriteString(leftTables[i])
		bufNt.WriteString("`.id ")
		//sql += " inner join ·" + innerTables[i] + "`"
	}
	// bufNt.WriteString(" where 1=1 and ")

	// select* 变为对应的字段名
	sqlNt = bufNt.String()
	sql = strings.Replace(bufNt.String(), "count(`"+tables[0]+"`.id) as total_num", GetMoreTableColumnSQL(model, tables[:]...), 1)
	for k, v := range params {
		switch k {
		case "clientPage":
			clientPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case "everyPage":
			everyPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case "order":
			order = v[0]
			continue
		case "key":
			key = v[0]
			var tablens = append(tables, tables[:]...)
			for k, v := range tablens {
				tablens[k] += ":" + v
			}
			// more tables key search
			sqlKey, sqlNtKey, argsKey := sq.GetMoreKeySQL(key, model, tablens[:]...)
			bufW.WriteString(sqlKey)
			bufNtW.WriteString(sqlNtKey)
			args = append(args, argsKey[:]...)
			continue
		case "":
			continue
		}

		// other tables, except tables[0]
		for _, table := range tables[1:] {
			switch {
			case !strings.Contains(table, table+"_id") && strings.Contains(table, table+"_"):
				v[0] = strings.Replace(v[0], "'", "\\'", -1)
				bufW.WriteString("`" + table + "`.`" + string([]byte(k)[len(v)+1:]) + "`" + " = ? and ")
				bufNtW.WriteString("`" + table + "`.`" + string([]byte(k)[len(v)+1:]) + "`" + " = ? and ")
				args = append(args, v[0])
				goto into
			}
		}
		v[0] = strings.Replace(v[0], "'", "\\'", -1)
		bufW.WriteString("`" + tables[0] + "`." + k + " = ? and ")
		bufNtW.WriteString("`" + tables[0] + "`." + k + " = ? and ")
		args = append(args, v[0])
	into:
	}

	if len(params) != 0 {
		sql += fmt.Sprintf("where %s order by id %s ", bufW.Bytes()[:bufW.Len()-4], order)
		sqlNt += fmt.Sprintf("where %s", bufNtW.Bytes()[:bufNtW.Len()-4])
	}

	return
}

// 分页参数不传, 查询所有
// 默认根据id倒序
// 单张表
func GetSearchSQL(model interface{}, table string, params map[string][]string) (sqlNt, sql string, clientPage, everyPage int64, args []interface{}) {

	var (
		order        = "desc"     // order by
		key          = ""         // key like binary search
		bufW, bufNtW bytes.Buffer // where sql, sqlNt bytes sql
	)

	// select* replace
	sql = fmt.Sprintf("select %s from `%s` ", GetColSQL(model), table)
	sqlNt = fmt.Sprintf("select count(id) as total_num from `%s` ", table)
	for k, v := range params {
		switch k {
		case "clientPage":
			clientPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case "everyPage":
			everyPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case "order":
			order = v[0]
			continue
		case "key":
			key = v[0]
			sqlKey, sqlNtKey, argsKey := sq.GetKeySQL(key, model, table)
			bufW.WriteString(sqlKey)
			bufNtW.WriteString(sqlNtKey)
			args = append(args, argsKey[:]...)
			continue
		case "":
			continue
		}

		//v[0] = strings.Replace(v[0], "'", "\\'", -1) //转义
		//sql += k + " = ? and "
		//sqlNt += k + " = ? and "
		bufW.WriteString(k)
		bufW.WriteString(" = ? and ")
		bufNtW.WriteString(k)
		bufNtW.WriteString(" = ? and ")
		args = append(args, v[0]) // args
	}

	if len(params) != 0 {
		sql += fmt.Sprintf("where %s order by id %s ", bufW.Bytes()[:bufW.Len()-4], order)
		sqlNt += fmt.Sprintf("where %s", bufNtW.Bytes()[:bufNtW.Len()-4])
	}

	return
}

// 传入数据库表名
// 更新语句拼接
func GetUpdateSQL(table string, params map[string][]string) (sql string, args []interface{}) {

	// sql connect
	var (
		id  string       // id
		buf bytes.Buffer // sql bytes connect
	)
	buf.WriteString("update `")
	buf.WriteString(table)
	buf.WriteString("` set ")
	for k, v := range params {
		if k == "id" {
			id = v[0]
			continue
		}
		buf.WriteString("`")
		buf.WriteString(k)
		buf.WriteString("` = ?,")
		args = append(args, v[0])
	}
	args = append(args, id)
	sql = string(buf.Bytes()[:buf.Len()-1]) + " where id = ?"
	return sql, args
}

// 传入数据库表名
// 插入语句拼接
func GetInsertSQL(table string, params map[string][]string) (sql string, args []interface{}) {

	// sql connect
	var (
		sqlv string
		buf  bytes.Buffer // sql bytes connect
	)
	buf.WriteString("insert `")
	buf.WriteString(table)
	buf.WriteString("`(")
	//sql = "insert `" + table + "`("

	for k, v := range params {
		buf.WriteString("`")
		buf.WriteString(k)
		buf.WriteString("`,")

		args = append(args, v[0])
		sqlv += "?,"
	}
	//sql = buf.Bytes()[:buf.Len()-1]
	sql = buf.String()
	sql = string([]byte(sql)[:len(sql)-1]) + ") value(" + sqlv
	sql = string([]byte(sql)[:len(sql)-1]) + ")" // remove ','

	return sql, args
}

// ===================================================================================
// ==========================common crud=========made=by=lucheng======================
// ===================================================================================

// get
// relation get
////////////////

// 获得数据,根据sql语句,无分页
func GetDataBySQL(data interface{}, sql string, args ...interface{}) (err error) {

	dba := DB.Raw(sql, args[:]...).Scan(data)
	// 有数据是返回相应信息
	if dba.Error != nil {
		err = sq.GetSQLError(dba.Error.Error())
	} else {
		// get data
		err = nil
	}
	return
}

// 获得数据,根据name条件
func GetDataByName(data interface{}, name, value string) (err error) {

	dba := DB.Where(name+" = ?", value).Find(data) //只要一行数据时使用 LIMIT 1,增加查询性能

	//有数据是返回相应信息
	if dba.Error != nil {
		err = sq.GetSQLError(dba.Error.Error())
	} else {
		// get data
		err = nil
	}
	return
}

// inner join
// 查询数据约定,表名_字段名(若有重复)
// 获得数据,根据id,两张表连接尝试
func GetDoubleTableDataByID(model, data interface{}, id, table1, table2 string) error {
	sql := fmt.Sprintf("select %s from `%s` inner join `%s` "+
		"on `%s`.%s_id=`%s`.id where `%s`.id=? limit 1", GetDoubleTableColumnSQL(model, table1, table2), table1, table2, table1, table2, table2, table1)

	return GetDataBySQL(data, sql, id)
}

// left join
// 查询数据约定,表名_字段名(若有重复)
// 获得数据,根据id,两张表连接
func GetLeftDoubleTableDataByID(model, data interface{}, id, table1, table2 string) error {

	sql := fmt.Sprintf("select %s from `%s` left join `%s` on `%s`.%s_id=`%s`.id where `%s`.id=? limit 1", GetDoubleTableColumnSQL(model, table1, table2), table1, table2, table1, table2, table2, table1)

	return GetDataBySQL(data, sql, id)
}

// 获得数据,根据id
func GetDataByID(data interface{}, id string) (err error) {

	dba := DB.First(data, id) //只要一行数据时使用 LIMIT 1,增加查询性能

	//有数据是返回相应信息
	if dba.Error != nil {
		err = sq.GetSQLError(dba.Error.Error())
	} else {
		// get data
		err = nil
	}
	return
}

// More Table
// params: innerTables is inner join tables
// params: leftTables is left join tables
// return: search info
// table1 as main table, include other tables_id(foreign key)
func GetMoreDataBySearch(model, data interface{}, params map[string][]string, innerTables []string, leftTables []string, args ...interface{}) (pager result.Pager, err error) {
	// more table search
	sqlNt, sql, clientPage, everyPage, args := GetMoreSearchSQL(model, params, innerTables, leftTables)

	return GetDataBySQLSearch(data, sql, sqlNt, clientPage, everyPage, args[:]...)
}

// 获得数据,根据sql语句,分页
// args : sql参数'？'
// sql, sqlNt args 相同, 共用args
func GetDataBySQLSearch(data interface{}, sql, sqlNt string, clientPage, everyPage int64, args ...interface{}) (pager result.Pager, err error) {

	// if no clientPage or everyPage
	// return all data
	if clientPage != 0 && everyPage != 0 {
		sql += fmt.Sprintf("limit %d, %d", (clientPage-1)*everyPage, everyPage)
	}

	// sqlNt += limit
	dba := DB.Raw(sqlNt, args[:]...).Scan(&pager)
	if dba.Error != nil {
		err = sq.GetSQLError(dba.Error.Error())
	} else {
		// DB.Debug().Find(&dest)
		dba = DB.Raw(sql, args[:]...).Scan(data)

		if dba.Error != nil {
			err = sq.GetSQLError(dba.Error.Error())
			return pager, err
		}
		// pager data
		pager.ClientPage = clientPage
		pager.EveryPage = everyPage
		err = nil
	}
	return
}

// 获得数据,分页/查询
func GetDataBySearch(model, data interface{}, table string, params map[string][]string) (pager result.Pager, err error) {

	sqlNt, sql, clientPage, everyPage, args := GetSearchSQL(model, table, params)

	return GetDataBySQLSearch(data, sql, sqlNt, clientPage, everyPage, args[:]...)
}

// delete
///////////////////

// delete by sql
func DeleteDataBySQL(sql string, args ...interface{}) (err error) {

	dba := DB.Exec(sql, args[:]...)
	switch {
	case dba.Error != nil:
		err = sq.GetSQLError(dba.Error.Error())
	default:
		err = nil
	}
	return
}

// 删除通用,任意参数
func DeleteDataByName(table string, key, value string) error {
	sql := fmt.Sprintf("delete from `%s` where %s=?", table, key)

	return DeleteDataBySQL(sql, value)
}

// update
///////////////////

// 修改数据,通用
func UpdateDataBySQL(sql string, args ...interface{}) (err error) {
	var num int64 //返回影响的行数

	dba := DB.Exec(sql, args[:]...)
	num = dba.RowsAffected
	switch {
	case dba.Error != nil:
		err = sq.GetSQLError(dba.Error.Error())
	case num == 0 && dba.Error == nil:
		err = errors.New(result.MsgExistOrNo)
	default:
		err = nil
	}
	return
}

// 修改数据,通用
func UpdateData(table string, params map[string][]string) error {

	sql, args := GetUpdateSQL(table, params)

	return UpdateDataBySQL(sql, args[:]...)
}

// 结合struct修改
func UpdateStructData(data interface{}) (err error) {
	var num int64 //返回影响的行数

	dba := DB.Save(data)
	num = dba.RowsAffected
	switch {
	case dba.Error != nil:
		err = sq.GetSQLError(dba.Error.Error())
	case num == 0 && dba.Error == nil:
		err = errors.New(result.MsgExistOrNo)
	default:
		err = nil
	}
	return
}

// create
////////////////////

// Create data by sql
func CreateDataBySQL(sql string, args ...interface{}) (err error) {

	dba := DB.Exec(sql, args[:]...)
	switch {
	case dba.Error != nil:
		err = sq.GetSQLError(dba.Error.Error())
	default:
		err = nil
	}
	return err
}

// 创建数据,通用
func CreateData(table string, params map[string][]string) (err error) {

	sql, args := GetInsertSQL(table, params)

	return CreateDataBySQL(sql, args[:]...)
}

// map[string][]string 形式批量创建问题: 顺序对应, 使用json形式批量创建

// 创建数据,通用
// 返回id,事务,慎用
// 业务少可用
func CreateDataResID(table string, params map[string][]string) (id ID, err error) {

	//开启事务
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	sql, args := GetInsertSQL(table, params)
	dba := tx.Exec(sql, args[:]...)

	tx.Raw("select max(id) as id from `%s`", table).Scan(&id)

	switch {
	case dba.Error != nil:
		err = sq.GetSQLError(dba.Error.Error())
	default:
		err = nil
	}

	if tx.Error != nil {
		tx.Rollback()
	}

	tx.Commit()
	return
}

// select检查是否存在
// == nil 即存在
func ValidateSQL(sql string) (err error) {
	var num int64 //返回影响的行数

	var ve Value
	dba := DB.Raw(sql).Scan(&ve)
	num = dba.RowsAffected
	switch {
	case dba.Error != nil:
		err = sq.GetSQLError(dba.Error.Error())
	case num == 0 && dba.Error == nil:
		err = errors.New(result.MsgExistOrNo)
	default:
		err = nil
	}
	return
}

//==============================================================
// json处理(struct data)
//==============================================================

// create
func CreateDataJ(data interface{}) (err error) {

	dba := DB.Create(data)

	switch {
	case dba.Error != nil:
		err = sq.GetSQLError(dba.Error.Error())
	default:
		err = nil
	}
	return err
}

// data must array type
// more data create
// single table
func CreateMoreDataJ(table string, model interface{}, data interface{}) error {
	var (
		//buf    bytes.Buffer
		buf    bytes.Buffer
		params []interface{}
	)

	// array data
	arrayData := reflect2.ToSlice(data)

	for _, v := range arrayData {
		// buf
		buf.WriteString("(")
		buf.WriteString(GetColParamSQL(model))
		buf.WriteString("),")
		// params
		params = append(params, GetParams(v)[:]...)
	}
	values := string(buf.Bytes()[:buf.Len()-1])

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", table, GetColSQL(model), values)
	err := DB.Exec(sql, params[:]...).Error
	return err
}

// update
func UpdateDataJ(data interface{}) (err error) {

	dba := DB.Model(data).Update(data)

	switch {
	case dba.Error != nil:
		err = sq.GetSQLError(dba.Error.Error())
	default:
		err = nil
	}
	return
}
