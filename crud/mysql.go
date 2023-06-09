// package gt

package crud

import (
	"bytes"
	"fmt"
	"github.com/dreamlu/gt/conf"
	depCons "github.com/dreamlu/gt/crud/dep/cons"
	"github.com/dreamlu/gt/crud/dep/result"
	"github.com/dreamlu/gt/log"
	"github.com/dreamlu/gt/src/cons"
	mr "github.com/dreamlu/gt/src/reflect"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
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
	//	db.db.GetLog.LogMode(gormLog.Error)
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
func logInfo(dbS *dba) gormLog.Interface {
	lv := gormLog.Error
	if l := dbS.Log; l {
		lv = gormLog.Info
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
// 设计模式--单例模式[懒汉式]
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

func (db *DB) Find(gt *GT) (pager result.Pager) {

	// isMock
	if gt.isMock {
		return
	}
	gt.GetSQL()
	if gt.isCount {
		db.countSQL(gt)
		pager = db.count(gt)
		if pager.TotalNum == 0 {
			return
		}
	}
	db.get(gt)
	return
}

// FindM no search
// params: innerTables is inner join tables
// params: leftTables is left join tables
// return search info
// table1 as main table, include other tables_id(foreign key)
func (db *DB) FindM(gt *GT) (pager result.Pager) {
	// isMock
	if gt.isMock {
		return
	}
	gt.GetMoreSQL()
	if gt.isCount {
		pager = db.count(gt)
		if pager.TotalNum == 0 {
			return
		}
	}
	db.get(gt)
	return
}

// FindS select sql search
func (db *DB) FindS(gt *GT) (pager result.Pager) {
	// isMock
	if gt.isMock {
		return
	}
	gt.GetSelectSQL()
	if gt.isCount {
		pager = db.count(gt)
		if pager.TotalNum == 0 {
			return
		}
	}
	db.get(gt)
	return
}

func (db *DB) countSQL(gt *GT) *DB {

	// default
	gt.order = fmt.Sprintf(depCons.OrderDesc, gt.tableT)

	gt.sqlNt = fmt.Sprintf(depCons.SelectCountFrom, gt.tableT)
	gt.whereCount()

	return db
}

func (db *DB) count(gt *GT) (pager result.Pager) {

	// if clientPage or everyPage < 0
	// return all data
	if gt.clientPage == 0 {
		gt.clientPage = depCons.ClientPage
	}
	if gt.everyPage == 0 {
		gt.everyPage = depCons.EveryPage
	}
	db.res = db.DB.Raw(gt.sqlNt, gt.Args...).Scan(&pager)
	if db.res.Error != nil || pager.TotalNum == 0 {
		return
	}
	pager.ClientPage = gt.clientPage
	pager.EveryPage = gt.everyPage
	// sqlNt += limit
	if gt.clientPage > 0 && gt.everyPage > 0 {
		gt.sql += fmt.Sprintf("limit %d, %d", (gt.clientPage-1)*gt.everyPage, gt.everyPage)
	}
	return
}

// get data
func (db *DB) get(gt *GT) {

	db.res = db.DB.Raw(gt.sql, gt.Args...).Scan(gt.Data)
}

func (db *DB) exec(sql string, args ...any) {
	db.res = db.Exec(sql, args...)
}

func (db *DB) Delete(gt *GT, conds ...any) {
	gt.parse().common()
	if gt.sqlSoft != "" {
		db.exec(fmt.Sprintf("update %s set %s = now() where id in (?)", gt.tableT, gt.parses.GetS(gt.Table)), conds...)
		return
	}
	db.res = db.DB.Delete(gt.Data, conds)
}

func (db *DB) Update(gt *GT) {
	if gt.Select != "" {
		db.res = db.Table(gt.Table).Where(gt.Select, gt.Args).Updates(gt.Data)
	} else {
		db.res = db.Table(gt.Table).Model(gt.Data).Updates(gt.Data)
	}
}

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
		param.Data = &columns
		tb := TableOnly(v)
		db.get(&GT{
			Params: &Params{Data: &columns},
			sql:    "SELECT COLUMN_NAME FROM `information_schema`.`COLUMNS` WHERE TABLE_NAME = ? and TABLE_SCHEMA = ?",
			Args:   []any{tb, name},
		})
		TableCols[tb] = columns
	}
}
