// package gt

package gt

import (
	"bytes"
	"fmt"
	"github.com/dreamlu/gt/lib/cons"
	"github.com/dreamlu/gt/lib/result"
	mr "github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/third/conf"
	"github.com/dreamlu/gt/third/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"sync"
	"time"
)

// DB tool
type DB struct {
	// db driver
	*gorm.DB
	res *gorm.DB
	// db log mode
	log bool
}

// db params
type dba struct {
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	Name        string `yaml:"name"`
	MaxIdleConn int    `yaml:"maxIdleConn"`
	MaxOpenConn int    `yaml:"maxOpenConn"`
	// db log mode
	Log bool `yaml:"log"`
}

// NewDB new db driver
func (db *DB) NewDB() {

	dbS := &dba{}
	conf.UnmarshalField(cons.ConfDB, dbS)
	db.log = dbS.Log
	var (
		sql = fmt.Sprintf("%s:%s@%s/?charset=utf8mb4&parseTime=True&loc=Local", dbS.User, dbS.Password, dbS.Host)
	)

	// auto create database
	db.DB = db.open(sql, dbS)
	err := db.DB.Exec("create database if not exists `" + dbS.Name + "`").Error
	if err != nil {
		log.Info("[mysql connect database error, try connect direct]")
	}

	sql = fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local", dbS.User, dbS.Password, dbS.Host, dbS.Name)
	db.DB = db.open(sql, dbS)
	// Globally disable table names
	// use name replace names
	db.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
	}
	//db.db.SingularTable(true)

	//if l := dbS.Log; l {
	//	db.db.Logger.LogMode(logger2.Error)
	//}
	//db.db.SetLogger(defaultLog)
	// connection pool
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sdb, _ := db.DB.DB()
	sdb.SetMaxIdleConns(dbS.MaxIdleConn)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sdb.SetMaxOpenConns(dbS.MaxOpenConn)

	return
}

func (db *DB) open(sql string, dbS *dba) *gorm.DB {
	// database, initialize once
	cf := &gorm.Config{
		Logger:                                   logInfo(dbS),
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	DB, err := gorm.Open(mysql.Open(sql), cf)
	//defer db.db.Close()
	if err != nil {
		log.Error("[mysql connect error]:", err)
		log.Error("[mysql try connect again]: try it every 5s...")
		// try to reconnect
		for {
			// go is so fast
			// try it every 5s
			time.Sleep(5 * time.Second)
			DB, err = gorm.Open(mysql.Open(sql), cf)
			//defer db.Close()
			if err != nil {
				log.Error("[mysql connect error]:", err)
				continue
			}
			log.Info("[mysql connect success]")
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
	return newMysqlLog(
		Config{
			SlowThreshold: 200 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:      lv,                     // Log level
			Colorful:      true,                   // 彩色打印
		},
	)
}

var (
	onceDB sync.Once
	// dbTool is global
	dbTool *DB
)

// cusdb
func cusdb(db *gorm.DB, log bool) *DB {
	onceDB.Do(func() {
		dbTool = &DB{
			DB:  db,
			res: nil,
			log: log,
		}
	})
	return dbTool
}

// db single db
func db() *DB {

	onceDB.Do(func() {
		dbTool = &DB{}
		// init db
		dbTool.NewDB()
	})
	return dbTool
}

func (db *DB) clone() *DB {

	return &DB{
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
func (db *DB) getBySQL(data any, sql string, args ...any) {

	db.res = db.DB.Raw(sql, args...).Scan(data)
}

func (db *DB) GetByID(gt *GT, id any) {

	db.getBySQL(gt.Data, fmt.Sprintf(cons.SelectFrom+"where id = ?", GetColSQL(gt.Model), ParseTable(gt.Table)), id)
}

// GetMoreBySearch more table
// params: innerTables is inner join tables
// params: leftTables is left join tables
// return search info
// table1 as main table, include other tables_id(foreign key)
func (db *DB) GetMoreBySearch(gt *GT) (pager result.Pager) {
	// more table search
	gt.GetMoreSQL()
	// isMock
	if gt.isMock {
		return
	}
	return db.GetBySQLSearch(gt.Data, gt.sql, gt.sqlNt, gt.clientPage, gt.everyPage, gt.Args)
}

// GetBySearch single table
// return search info
func (db *DB) GetBySearch(gt *GT) (pager result.Pager) {

	gt.GetSearchSQL()
	// isMock
	if gt.isMock {
		return
	}
	return db.GetBySQLSearch(gt.Data, gt.sql, gt.sqlNt, gt.clientPage, gt.everyPage, gt.Args)
}

// Get no search
func (db *DB) Get(gt *GT) {

	gt.GetSQL()
	// isMock
	if gt.isMock {
		return
	}
	db.getBySQL(gt.Data, gt.sql, gt.Args...)
}

// GetMoreData no search
func (db *DB) GetMoreData(gt *GT) {

	gt.GetMoreSQL()
	// isMock
	if gt.isMock {
		return
	}
	db.getBySQL(gt.Data, gt.sql, gt.Args...)
}

// GetDataBySelectSQLSearch select sql search
func (db *DB) GetDataBySelectSQLSearch(gt *GT) (pager result.Pager) {

	gt.GetSelectSearchSQL()
	// isMock
	if gt.isMock {
		return
	}
	return db.GetBySQLSearch(gt.Data, gt.sql, gt.sqlNt, gt.clientPage, gt.everyPage, gt.Args)
}

// GetBySQLSearch get sql search data
// clientPage: default 1
// everyPage: default 10
// if clientPage or everyPage < 0, return all
func (db *DB) GetBySQLSearch(data any, sql, sqlNt string, clientPage, everyPage int64, args []any) (pager result.Pager) {

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
	db.res = db.DB.Raw(sqlNt, args...).Scan(&pager)
	if db.res.Error == nil && pager.TotalNum > 0 {
		db.res = db.DB.Raw(sql, args...).Scan(data)
		// pager data
		pager.ClientPage = clientPage
		pager.EveryPage = everyPage
		return
	}
	return
}

// exec common
////////////////////

func (db *DB) ExecSQL(sql string, args ...any) {
	db.res = db.Exec(sql, args...)
}

// delete
///////////////////

func (db *DB) Delete(table string, id any) {
	switch id.(type) {
	case string:
		if strings.Contains(id.(string), ",") {
			id = strings.Split(id.(string), ",")
		}
	}
	db.ExecSQL(fmt.Sprintf("delete from %s where id in (?)", ParseTable(table)), id)
}

// update
///////////////////

func (db *DB) Update(gt *GT) {
	if gt.Select != "" {
		db.res = db.Table(gt.Table).Where(gt.Select, gt.Args).Updates(gt.Data)
	} else {
		db.res = db.Table(gt.Table).Model(gt.Data).Updates(gt.Data)
	}
}

// create
////////////////////

// Create single/array
func (db *DB) Create(table string, data any) {
	db.res = db.Table(table).Create(data)
}

// CreateMore data must array type
// more data create
// single table
// also can use Create array
func (db *DB) CreateMore(table string, model any, data any) {
	var (
		buf       bytes.Buffer
		params    []any
		p         = parse(model)
		arrayData = mr.ToSlice(data) // slice data
		colPSQL   = GetColParamSQL(p)
	)

	for _, v := range arrayData {
		// buf
		buf.WriteByte('(')
		buf.WriteString(colPSQL)
		buf.WriteString("),")
		// params
		p.Vs = nil
		parseV(p, v)
		params = append(params, p.Vs...)
	}
	values := string(buf.Bytes()[:buf.Len()-1])

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", ParseTable(table), GetColSQL(model), values)
	db.res = db.DB.Exec(sql, params...)
}

// InitColumns init db table columns map
func (db *DB) InitColumns(param *Params) {

	var (
		name   = conf.Get[string](cons.ConfDBName)
		tables = []string{param.Table}
	)

	tables = append(tables, param.InnerTable...)
	tables = append(tables, param.LeftTable...)

	for _, v := range tables {
		if v == "" {
			continue
		}
		if _, ok := TableCols[v]; ok {
			continue
		}
		var columns []string
		tb := TableOnly(v)
		db.getBySQL(&columns, "SELECT COLUMN_NAME FROM `information_schema`.`COLUMNS` WHERE TABLE_NAME = ? and TABLE_SCHEMA = ?", tb, name)
		TableCols[tb] = columns
	}
}
