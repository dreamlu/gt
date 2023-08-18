// package gt

package crud

import (
	"fmt"
	"github.com/dreamlu/gt/conf"
	crudCons "github.com/dreamlu/gt/crud/dep/cons"
	"github.com/dreamlu/gt/log"
	"github.com/dreamlu/gt/src/cons"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
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
	Driver      string `yaml:"driver"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
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
	db.DB = db.open(dbS)
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

func dialector(d *dba) gorm.Dialector {
	switch d.Driver {
	case "postgres":
		return postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", d.Host, d.User, d.Password, d.Name, d.Port))
	default: // mysql
		return mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", d.User, d.Password, d.Host, d.Port, d.Name))
	}
}

func (db *DB) open(dbS *dba) *gorm.DB {
	// database, initialize once
	cf := &gorm.Config{
		Logger:                                   logInfo(dbS),
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	crudCons.Init(dbS.Driver)
	dial := dialector(dbS)
	DB, err := gorm.Open(dial, cf)
	//defer db.db.Close()
	if err != nil {
		log.Error("[db connect error]:", err)
		log.Error("[db try connect again]: try it every 5s...")
		// try to reconnect
		for {
			// go is so fast
			// try it every 5s
			time.Sleep(5 * time.Second)
			DB, err = gorm.Open(dial, cf)
			//defer db.Close()
			if err != nil {
				log.Error("[db connect error]:", err)
				continue
			}
			log.Info("[db connect success]")
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
	return newDBLog(
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
