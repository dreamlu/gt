// package gt

package gt

import (
	"bytes"
	"fmt"
	"github.com/dreamlu/gt/tool/conf"
	log3 "github.com/dreamlu/gt/tool/log"
	reflect2 "github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/result"
	sq "github.com/dreamlu/gt/tool/sql"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/cons"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
	"sync"
	"time"
)

// DB tool
type DBTool struct {
	// db driver
	*gorm.DB
	res *gorm.DB
	// db log mode
	log bool
}

// db params
type dba struct {
	User        string
	Password    string
	Host        string
	Name        string
	MaxIdleConn int
	MaxOpenConn int
	// db log mode
	Log bool
}

// new db driver
func (db *DBTool) NewDB() {

	dbS := &dba{}
	conf.GetStruct("app.db", dbS)
	db.log = dbS.Log
	var (
		sql = fmt.Sprintf("%s:%s@%s/?charset=utf8mb4&parseTime=True&loc=Local", dbS.User, dbS.Password, dbS.Host)
	)

	// auto create database
	db.DB = db.open(sql, dbS)
	err := db.DB.Exec("create database if not exists `" + dbS.Name + "`").Error
	if err != nil {
		log3.Info("[mysql自动连接根数据库失败,尝试直接连接]")
	}

	sql = fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local", dbS.User, dbS.Password, dbS.Host, dbS.Name)
	db.DB = db.open(sql, dbS)
	// Globally disable table names
	// use name replace names
	db.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
	}
	//db.DB.SingularTable(true)

	//if l := dbS.Log; l {
	//	db.DB.Logger.LogMode(logger2.Error)
	//}
	//db.DB.SetLogger(defaultLog)
	// connection pool
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sdb, _ := db.DB.DB()
	sdb.SetMaxIdleConns(dbS.MaxIdleConn)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sdb.SetMaxOpenConns(dbS.MaxOpenConn)

	return
}

func (db *DBTool) open(sql string, dbS *dba) *gorm.DB {
	// database, initialize once
	cf := &gorm.Config{
		Logger:                                   logInfo(dbS),
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	DB, err := gorm.Open(mysql.Open(sql), cf)
	//defer db.DB.Close()
	if err != nil {
		log3.Error("[mysql连接错误]:", err)
		log3.Error("[mysql开始尝试重连中]: try it every 5s...")
		// try to reconnect
		for {
			// go is so fast
			// try it every 5s
			time.Sleep(5 * time.Second)
			DB, err = gorm.Open(mysql.Open(sql), cf)
			//defer DB.Close()
			if err != nil {
				log3.Error("[mysql连接错误]:", err)
				continue
			}
			log3.Info("[mysql重连成功]")
			break
		}
	}
	return DB
}

// log info
func logInfo(dbS *dba) logger2.Interface {
	lv := logger2.Error
	if l := dbS.Log; l {
		lv = logger2.Info
	}
	return New(
		Config{
			SlowThreshold: 200 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:      lv,                     // Log level
			Colorful:      true,                   // 彩色打印
		},
	)
}

// init DBTool
func newDBTool() *DBTool {

	dbTool := &DBTool{}

	// init db
	dbTool.NewDB()
	return dbTool
}

var (
	onceDB sync.Once
	// dbTool is global
	dbTool *DBTool
)

// single db
func DB() *DBTool {

	onceDB.Do(func() {
		dbTool = newDBTool()
	})
	return dbTool
}

func (db *DBTool) clone() *DBTool {

	return &DBTool{
		DB:  db.DB,
		log: db.log,
		res: db.res,
	}
}

// ===================================================================================
// ==========================common crud==============================================
// ===================================================================================

// get
////////////////

// get single data
func (db *DBTool) getBySQL(data interface{}, sql string, args ...interface{}) {

	typ := reflect.TypeOf(data)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	db.res = db.DB.Raw(sql, args[:]...).Scan(data)
}

// get by id
func (db *DBTool) GetByID(gt *GT, id interface{}) {

	db.getBySQL(gt.Data, fmt.Sprintf(cons.SelectFrom+"where id = ?", GetColSQL(gt.Model), sq.Table(gt.Table)), id)
}

// more table
// params: innerTables is inner join tables
// params: leftTables is left join tables
// return search info
// table1 as main table, include other tables_id(foreign key)
func (db *DBTool) GetMoreBySearch(gt *GT) (pager result.Pager) {
	// more table search
	gt.GetMoreSQL()
	// isMock
	if gt.isMock {
		return
	}
	return db.GetBySQLSearch(gt.Data, gt.sql, gt.sqlNt, gt.clientPage, gt.everyPage, gt.Args)
}

// single table
// return search info
func (db *DBTool) GetBySearch(gt *GT) (pager result.Pager) {

	gt.GetSearchSQL()
	// isMock
	if gt.isMock {
		return
	}
	return db.GetBySQLSearch(gt.Data, gt.sql, gt.sqlNt, gt.clientPage, gt.everyPage, gt.Args)
}

// 获得数据, no search
func (db *DBTool) Get(gt *GT) {

	gt.GetSQL()
	// isMock
	if gt.isMock {
		return
	}
	db.getBySQL(gt.Data, gt.sql, gt.Args...)
}

// 获得数据, no search
func (db *DBTool) GetMoreData(gt *GT) {

	gt.GetMoreSQL()
	// isMock
	if gt.isMock {
		return
	}
	db.getBySQL(gt.Data, gt.sql, gt.Args...)
}

// select sql search
func (db *DBTool) GetDataBySelectSQLSearch(gt *GT) (pager result.Pager) {

	gt.GetSelectSearchSQL()
	return db.GetBySQLSearch(gt.Data, gt.sql, gt.sqlNt, gt.clientPage, gt.everyPage, gt.Args)
}

// get sql search data
// clientPage: default 1
// everyPage: default 10
// if clientPage or everyPage < 0, return all
func (db *DBTool) GetBySQLSearch(data interface{}, sql, sqlNt string, clientPage, everyPage int64, args []interface{}) (pager result.Pager) {

	// if clientPage or everyPage < 0
	// return all data
	if clientPage == 0 {
		clientPage = cons.ClientPage
	}
	if everyPage == 0 {
		everyPage = cons.EveryPage
	}
	if clientPage > 0 && everyPage > 0 {
		sql += fmt.Sprintf("limit %d, %d", (clientPage-1)*everyPage, everyPage)
	}
	// sqlNt += limit
	db.res = db.DB.Raw(sqlNt, args[:]...).Scan(&pager)
	if db.res.Error == nil && pager.TotalNum > 0 {
		db.res = db.DB.Raw(sql, args[:]...).Scan(data)
		// pager data
		pager.ClientPage = clientPage
		pager.EveryPage = everyPage
		return
	}
	return
}

// exec common
////////////////////

// exec sql
func (db *DBTool) ExecSQL(sql string, args ...interface{}) {
	db.res = db.Exec(sql, args...)
}

// delete
///////////////////

// delete
func (db *DBTool) Delete(table string, id interface{}) {
	switch id.(type) {
	case string:
		if strings.Contains(id.(string), ",") {
			id = strings.Split(id.(string), ",")
		}
	}
	db.ExecSQL(fmt.Sprintf("delete from %s where id in (?)", sq.Table(table)), id)
}

// update
///////////////////

// update
func (db *DBTool) Update(gt *GT) {

	if gt.Model == nil {
		gt.Model = gt.Data
	}

	if gt.Select != "" {
		db.res = db.Table(gt.Table).Where(gt.Select, gt.Args).Updates(gt.Data)
	} else {
		db.res = db.Table(gt.Table).Model(gt.Data).Updates(gt.Data)
	}
}

// create
////////////////////

// create single/array
func (db *DBTool) Create(table string, data interface{}) {
	db.res = db.Table(table).Create(data)
}

// data must array type
// more data create
// single table
// also can use Create array
func (db *DBTool) CreateMore(table string, model interface{}, data interface{}) {
	var (
		buf    bytes.Buffer
		params []interface{}
	)

	// array data
	arrayData := reflect2.ToSlice(data)
	colPSQL := GetColParamSQL(model)

	for _, v := range arrayData {
		// buf
		buf.WriteString("(")
		buf.WriteString(colPSQL)
		buf.WriteString("),")
		// params
		params = append(params, GetParams(v)[:]...)
	}
	values := string(buf.Bytes()[:buf.Len()-1])

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", sq.Table(table), GetColSQL(model), values)
	db.res = db.DB.Exec(sql, params[:]...)
}

// form-data create/update
// future will remove
// use json replace

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

	tx.Raw("select max(id) as id from ?", table).Scan(&id)

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

// init db table columns map
func (db *DBTool) InitColumns(param *Params) {

	var (
		name   = conf.GetString("app.db.name")
		tables = []string{param.Table}
	)

	tables = append(tables, param.InnerTable...)
	tables = append(tables, param.LeftTable...)

	for _, v := range param.InnerTable {
		if v == "" {
			continue
		}
		if _, ok := sq.TableCols[v]; ok {
			continue
		}
		var columns []string
		tb := sq.TableOnly(v)
		db.getBySQL(&columns, "SELECT COLUMN_NAME FROM `information_schema`.`COLUMNS` WHERE TABLE_NAME = ? and TABLE_SCHEMA = ?", tb, name)
		sq.TableCols[tb] = columns
	}
}
