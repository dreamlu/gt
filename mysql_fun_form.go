package gt

import (
	"bytes"
	sq "github.com/dreamlu/gt/tool/sql"
	"github.com/dreamlu/gt/tool/type/cmap"
)

// form-data create/update

// get update sql
func GetUpdateSQL(table string, params cmap.CMap) (sql string, args []interface{}) {

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

// get insert sql
func GetInsertSQL(table string, params cmap.CMap) (sql string, args []interface{}) {

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

// via form data update
func (db *DBTool) UpdateFormData(table string, params cmap.CMap) (err error) {

	sql, args := GetUpdateSQL(table, params)
	db.ExecSQL(sql, args...)
	return db.res.Error
}

// via form data create
func (db *DBTool) CreateFormData(table string, params cmap.CMap) error {

	sql, args := GetInsertSQL(table, params)
	db.ExecSQL(sql, args...)
	return db.res.Error
}

// param.CMap 形式批量创建问题: 顺序对应, 使用json形式批量创建

// create data return id
func (db *DBTool) CreateDataResID(table string, params cmap.CMap) (id uint64, err error) {

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
