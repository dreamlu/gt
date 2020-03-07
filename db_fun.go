// package gt

package gt

import (
	"bytes"
	"fmt"
	reflect2 "github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/result"
	sq "github.com/dreamlu/gt/tool/sql"
	"github.com/dreamlu/gt/tool/type/te"
	"github.com/dreamlu/gt/tool/util"
	"github.com/dreamlu/gt/tool/util/str"
	"reflect"
	"strconv"
	"strings"
)

//======================return tag=============================
//=============================================================

// select * replace
// select more tables
// tables : table name / table alias name
// 主表放在tables中第一个, 紧接着为主表关联的外键表名(无顺序)
func GetMoreTableColumnSQL(model interface{}, tables ...string) (sql string) {
	var buf bytes.Buffer

	//typ := reflect.TypeOf(model)
	GetReflectTagMore(reflect.TypeOf(model), &buf, tables[:]...)
	sql = string(buf.Bytes()[:buf.Len()-1]) //去点,
	return sql
}

// 层级递增解析tag, more tables
func GetReflectTagMore(reflectType reflect.Type, buf *bytes.Buffer, tables ...string) {

	if reflectType.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < reflectType.NumField(); i++ {
		tag := reflectType.Field(i).Tag.Get("json")
		if tag == "" {
			GetReflectTagMore(reflectType.Field(i).Type, buf, tables[:]...)
			continue
		}
		// sub sql
		gtTag := reflectType.Field(i).Tag.Get("gt")
		gtFields := strings.Split(gtTag, ";")
		for _, v := range gtFields {
			if v == str.GtSubSQL {
				goto into
			}
			if tag == "-" && strings.Contains(v, "field") {
				tagTmp := strings.Split(v, ":")
				tag = tagTmp[1]
			}
		}

		// foreign tables column
		for _, v := range tables {
			if strings.Contains(tag, v+"_id") {
				break
			}
			// tables
			switch {
			case strings.Contains(tag, v+"_"):
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
}

// 根据model中表模型的json标签获取表字段
// 将select* 中'*'变为对应的字段名
// 增加别名,表连接问题
func GetColSQLAlias(model interface{}, alias string) (sql string) {
	var buf bytes.Buffer

	//typ := reflect.TypeOf(model)
	GetReflectTagAlias(reflect.TypeOf(model), &buf, alias)
	sql = string(buf.Bytes()[:buf.Len()-1]) //去掉点,
	return sql
}

// 层级递增解析tag, 别名
func GetReflectTagAlias(reflectType reflect.Type, buf *bytes.Buffer, alias string) {

	if reflectType.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < reflectType.NumField(); i++ {
		tag := reflectType.Field(i).Tag.Get("json")
		if tag == "" {
			GetReflectTagAlias(reflectType.Field(i).Type, buf, alias)
			continue
		}
		// sub sql
		gtTag := reflectType.Field(i).Tag.Get("gt")
		if strings.Contains(gtTag, str.GtSubSQL) {
			continue
		}
		buf.WriteString(alias)
		buf.WriteString("`")
		buf.WriteString(tag)
		buf.WriteString("`,")
	}
}

// 根据model中表模型的json标签获取表字段
// 将select* 变为对应的字段名
func GetColSQL(model interface{}) (sql string) {
	var buf bytes.Buffer

	//typ := reflect.TypeOf(model)
	GetReflectTag(reflect.TypeOf(model), &buf)
	sql = string(buf.Bytes()[:buf.Len()-1]) // remove ,
	return sql
}

// 层级递增解析tag
func GetReflectTag(reflectType reflect.Type, buf *bytes.Buffer) {

	if reflectType.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < reflectType.NumField(); i++ {
		tag := reflectType.Field(i).Tag.Get("json")
		if tag == "" {
			GetReflectTag(reflectType.Field(i).Type, buf)
			continue
		}
		// sub sql
		gtTag := reflectType.Field(i).Tag.Get("gt")
		if strings.Contains(gtTag, str.GtSubSQL) {
			continue
		}
		buf.WriteString("`")
		buf.WriteString(tag)
		buf.WriteString("`,")
	}
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

// GT SQL struct
type GT struct {
	// attributes
	InnerTable []string    // inner join tables
	LeftTable  []string    // left join tables
	Table      string      // table name
	Model      interface{} // table model, like User{}
	Data       interface{} // table model data, like var user User{}, it is 'user'

	// pager info
	ClientPage int64 // page number
	EveryPage  int64 // Number of pages per page

	// count
	SubSQL string // SubQuery SQL
	// where
	SubWhereSQL string // SubWhere SQL
	// maybe future will use gt.params replace params
	Params map[string][]string // params

	// select sql
	Select string // select sql
	From   string // only once
	Group  string // the last group
	Args   []interface{}
	ArgsNt []interface{}
}

//=======================================sql语句处理==========================================
//===========================================================================================

// More Table
// params: innerTables is inner join tables, must even number
// params: leftTables is left join tables
// return: select sql
// table1 as main table, include other tables_id(foreign key)
func GetMoreSearchSQL(gt *GT) (sqlNt, sql string, clientPage, everyPage int64, args []interface{}) {

	var (
		order               = gt.InnerTable[0] + ".id desc" // order by
		key                 = ""                            // key like binary search
		tables              = gt.InnerTable                 // all tables
		bufNt, bufW, bufNtW bytes.Buffer                    // sql bytes connect
	)

	tables = append(tables, gt.LeftTable...)
	tables = util.RemoveDuplicateString(tables)
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
	for i := 1; i < len(gt.InnerTable); i += 2 {
		bufNt.WriteString(" inner join `")
		bufNt.WriteString(gt.InnerTable[i])
		bufNt.WriteString("` on `")
		bufNt.WriteString(gt.InnerTable[i-1])
		bufNt.WriteString("`.")
		bufNt.WriteString(gt.InnerTable[i])
		bufNt.WriteString("_id=`")
		bufNt.WriteString(gt.InnerTable[i])
		bufNt.WriteString("`.id ")
		//sql += " inner join ·" + innerTables[i] + "`"
	}
	// left join
	for i := 1; i < len(gt.LeftTable); i += 2 {
		bufNt.WriteString(" left join `")
		bufNt.WriteString(gt.LeftTable[i])
		bufNt.WriteString("` on `")
		bufNt.WriteString(gt.InnerTable[i-1])
		bufNt.WriteString("`.")
		bufNt.WriteString(gt.LeftTable[i])
		bufNt.WriteString("_id=`")
		bufNt.WriteString(gt.LeftTable[i])
		bufNt.WriteString("`.id ")
		//sql += " inner join ·" + innerTables[i] + "`"
	}
	// bufNt.WriteString(" where 1=1 and ")

	// select* 变为对应的字段名
	sqlNt = bufNt.String()
	sql = strings.Replace(bufNt.String(), "count(`"+tables[0]+"`.id) as total_num", GetMoreTableColumnSQL(gt.Model, tables[:]...)+gt.SubSQL, 1)
	for k, v := range gt.Params {
		switch k {
		case str.GtClientPage:
			clientPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case str.GtEveryPage:
			everyPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case str.GtOrder:
			order = v[0]
			continue
		case str.GtKey:
			key = v[0]
			var tablens = append(tables, tables[:]...)
			for k, v := range tablens {
				tablens[k] += ":" + v
			}
			// more tables key search
			sqlKey, sqlNtKey, argsKey := sq.GetMoreKeySQL(key, gt.Model, tablens[:]...)
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

	if bufW.Len() != 0 {
		sql += fmt.Sprintf("where %s%s ", bufW.Bytes()[:bufW.Len()-4], gt.SubWhereSQL)
		sqlNt += fmt.Sprintf("where %s%s", bufNtW.Bytes()[:bufNtW.Len()-4], gt.SubWhereSQL)
	}
	sql += fmt.Sprintf(" order by %s ", order)

	return
}

// 分页参数不传, 查询所有
// 默认根据id倒序
// 单张表
func GetSearchSQL(gt *GT) (sqlNt, sql string, clientPage, everyPage int64, args []interface{}) {

	var (
		order        = "id desc"  // default order by
		key          = ""         // key like binary search
		bufW, bufNtW bytes.Buffer // where sql, sqlNt bytes sql
	)

	// select* replace
	sql = fmt.Sprintf("select %s%s from `%s`", GetColSQL(gt.Model), gt.SubSQL, gt.Table)
	sqlNt = fmt.Sprintf("select count(id) as total_num from `%s`", gt.Table)
	for k, v := range gt.Params {
		switch k {
		case str.GtClientPage:
			clientPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case str.GtEveryPage:
			everyPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case str.GtOrder:
			order = v[0]
			continue
		case str.GtKey:
			key = v[0]
			sqlKey, sqlNtKey, argsKey := sq.GetKeySQL(key, gt.Model, gt.Table)
			bufW.WriteString(sqlKey)
			bufNtW.WriteString(sqlNtKey)
			args = append(args, argsKey[:]...)
			continue
		case "":
			continue
		}
		bufW.WriteString(k)
		bufW.WriteString(" = ? and ")
		bufNtW.WriteString(k)
		bufNtW.WriteString(" = ? and ")
		args = append(args, v[0]) // args
	}

	if bufW.Len() != 0 {
		sql += fmt.Sprintf(" where %s%s ", bufW.Bytes()[:bufW.Len()-4], gt.SubWhereSQL)
		sqlNt += fmt.Sprintf(" where %s%s", bufW.Bytes()[:bufW.Len()-4], gt.SubWhereSQL)
	}
	sql += fmt.Sprintf(" order by %s ", order)
	return
}

// get data sql
func GetDataSQL(gt *GT) (sql string, args []interface{}) {

	var (
		order        = "id desc"  // default order by
		key          = ""         // key like binary search
		bufW, bufNtW bytes.Buffer // where sql, sqlNt bytes sql
	)

	// select* replace
	sql = fmt.Sprintf("select %s%s from `%s`", GetColSQL(gt.Model), gt.SubSQL, gt.Table)
	for k, v := range gt.Params {
		switch k {
		case str.GtOrder:
			order = v[0]
			continue
		case str.GtKey:
			key = v[0]
			sqlKey, sqlNtKey, argsKey := sq.GetKeySQL(key, gt.Model, gt.Table)
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

	if bufW.Len() != 0 {
		sql += fmt.Sprintf(" where %s%s ", bufW.Bytes()[:bufW.Len()-4], gt.SubWhereSQL)
	}
	sql += fmt.Sprintf(" order by %s ", order)
	return
}

// select sql
func GetSelectSearchSQL(gt *GT) (sqlNt, sql string) {

	sql = gt.Select
	if gt.From == "" {
		gt.From = "from"
	}
	sqlNt = "select count(*) as total_num " + gt.From + strings.Join(strings.Split(sql, gt.From)[1:], "")
	if gt.Group != "" {
		sql += gt.Group
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
func (db *DBTool) GetDataBySQL(data interface{}, sql string, args ...interface{}) {

	db.DB = db.DB.Raw(sql, args[:]...).Scan(data)
}

// 获得数据,根据name条件
func (db *DBTool) GetDataByName(data interface{}, name, value string) (err error) {

	dba := db.DB.Where(name+" = ?", value).Find(data) //只要一行数据时使用 LIMIT 1,增加查询性能

	//有数据是返回相应信息
	if dba.Error != nil {
		err = sq.GetSQLError(dba.Error.Error())
	} else {
		// get data
		err = nil
	}
	return
}

//// inner join
//// 查询数据约定,表名_字段名(若有重复)
//// 获得数据,根据id,两张表连接尝试
//func (db *DBTool) GetDoubleTableDataByID(model, data interface{}, id, table1, table2 string) error {
//	sql := fmt.Sprintf("select %s from `%s` inner join `%s` "+
//		"on `%s`.%s_id=`%s`.id where `%s`.id=? limit 1", GetDoubleTableColumnSQL(model, table1, table2), table1, table2, table1, table2, table2, table1)
//
//	return db.GetDataBySQL(data, sql, id)
//}
//
//// left join
//// 查询数据约定,表名_字段名(若有重复)
//// 获得数据,根据id,两张表连接
//func (db *DBTool) GetLeftDoubleTableDataByID(model, data interface{}, id, table1, table2 string) error {
//
//	sql := fmt.Sprintf("select %s from `%s` left join `%s` on `%s`.%s_id=`%s`.id where `%s`.id=? limit 1", GetDoubleTableColumnSQL(model, table1, table2), table1, table2, table1, table2, table2, table1)
//
//	return db.GetDataBySQL(data, sql, id)
//}

// 获得数据,根据id
func (db *DBTool) GetDataByID(data interface{}, id string) {

	db.DB = db.DB.First(data, id) // limit 1
}

// More Table
// params: innerTables is inner join tables
// params: leftTables is left join tables
// return: search info
// table1 as main table, include other tables_id(foreign key)
func (db *DBTool) GetMoreDataBySearch(gt *GT) (pager result.Pager) {
	// more table search
	sqlNt, sql, clientPage, everyPage, args := GetMoreSearchSQL(gt)

	return db.GetDataBySQLSearch(gt.Data, sql, sqlNt, clientPage, everyPage, args, args)
}

// 获得数据,分页/查询
func (db *DBTool) GetDataBySearch(gt *GT) (pager result.Pager) {

	sqlNt, sql, clientPage, everyPage, args := GetSearchSQL(gt)

	return db.GetDataBySQLSearch(gt.Data, sql, sqlNt, clientPage, everyPage, args, args)
}

// 获得数据, no search
func (db *DBTool) GetData(gt *GT) {

	sql, args := GetDataSQL(gt)
	db.GetDataBySQL(gt.Data, sql, args[:]...)
}

// select sql search
func (db *DBTool) GetDataBySelectSQLSearch(gt *GT) (pager result.Pager) {

	sqlNt, sql := GetSelectSearchSQL(gt)

	return db.GetDataBySQLSearch(gt.Data, sql, sqlNt, gt.ClientPage, gt.EveryPage, gt.Args, gt.ArgsNt)
}

// 获得数据,根据sql语句,分页
// args : sql参数'？'
// sql, sqlNt args 相同, 共用args
func (db *DBTool) GetDataBySQLSearch(data interface{}, sql, sqlNt string, clientPage, everyPage int64, args []interface{}, argsNt []interface{}) (pager result.Pager) {

	// if no clientPage or everyPage
	// return all data
	if clientPage != 0 && everyPage != 0 {
		sql += fmt.Sprintf("limit %d, %d", (clientPage-1)*everyPage, everyPage)
	}
	// sqlNt += limit
	dba := db.DB.Raw(sqlNt, argsNt[:]...).Scan(&pager)
	if db.Error == nil {
		db.DB = db.DB.Raw(sql, args[:]...).Scan(data)
		// pager data
		pager.ClientPage = clientPage
		pager.EveryPage = everyPage
		return
	}
	db.DB = dba
	return
}

// exec common
////////////////////

// exec sql
func (db *DBTool) ExecSQL(sql string, args ...interface{}) {

	db.DB = db.Exec(sql, args...)
	//return db
}

// delete
///////////////////

// delete
func (db *DBTool) Delete(table string, id string) {
	// sql := fmt.Sprintf("delete from `%s` where id=?", table)
	db.ExecSQL(fmt.Sprintf("delete from `%s` where id=?", table), id)
}

// update
///////////////////

// via form data update
func (db *DBTool) UpdateFormData(table string, params map[string][]string) (err error) {

	sql, args := GetUpdateSQL(table, params)
	db.ExecSQL(sql, args...)
	return db.Error
}

// 结合struct修改
func (db *DBTool) UpdateStructData(data interface{}) (err error) {
	var num int64

	dba := db.DB.Save(data)
	num = dba.RowsAffected
	switch {
	case dba.Error != nil:
		err = sq.GetSQLError(dba.Error.Error())
	case num == 0 && dba.Error == nil:
		err = &te.TextError{Msg: result.MsgExistOrNo}
	default:
		err = nil
	}
	return
}

// create
////////////////////

// via form data create
func (db *DBTool) CreateFormData(table string, params map[string][]string) error {

	sql, args := GetInsertSQL(table, params)
	db.ExecSQL(sql, args...)
	return db.Error
}

// map[string][]string 形式批量创建问题: 顺序对应, 使用json形式批量创建

// 创建数据,通用
// 返回id,事务,慎用
// 业务少可用
func (db *DBTool) CreateDataResID(table string, params map[string][]string) (id str.ID, err error) {

	//开启事务
	tx := db.DB.Begin()
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
func (db *DBTool) ValidateSQL(sql string) (err error) {
	var num int64 //返回影响的行数

	var ve str.Value
	dba := db.DB.Raw(sql).Scan(&ve)
	num = dba.RowsAffected
	switch {
	case dba.Error != nil:
		err = sq.GetSQLError(dba.Error.Error())
	case num == 0 && dba.Error == nil:
		err = &te.TextError{Msg: result.MsgExistOrNo}
	default:
		err = nil
	}
	return
}

//==============================================================
// json处理(struct data)
//==============================================================

// create
func (db *DBTool) CreateData(data interface{}) *DBTool {

	db.DB = db.Create(data)
	return db
}

// data must array type
// more data create
// single table
func (db *DBTool) CreateMoreData(table string, model interface{}, data interface{}) error {
	var (
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
	err := db.DB.Exec(sql, params[:]...).Error
	return err
}

// update
func (db *DBTool) UpdateData(data interface{}) *DBTool {

	db.DB = db.DB.Model(data).Update(data)
	return db
}
