// package gt

package gt

import (
	"fmt"
	log2 "github.com/dreamlu/gt/sql/mysql/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
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
	Configger().GetStruct("app.db", dbS)
	db.log = dbS.Log
	var (
		sql = fmt.Sprintf("%s:%s@%s/?charset=utf8mb4&parseTime=True&loc=Local", dbS.User, dbS.Password, dbS.Host)
	)

	// auto create database
	db.DB = db.open(sql, dbS)
	err := db.DB.Exec("create database if not exists `" + dbS.Name + "`").Error
	if err != nil {
		Logger().Info("[mysql自动连接根数据库失败,尝试直接连接]")
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
		Logger:                 logInfo(dbS),
		SkipDefaultTransaction: true,
	}
	DB, err := gorm.Open(mysql.Open(sql), cf)
	//defer db.DB.Close()
	if err != nil {
		//if strings.Contains(err.Error(), "Unknown database"){
		//	DB.Exec("create database 'coupon'")
		//}
		Logger().Error("[mysql连接错误]:", err)
		Logger().Error("[mysql开始尝试重连中]: try it every 5s...")
		// try to reconnect
		for {
			// go is so fast
			// try it every 5s
			time.Sleep(5 * time.Second)
			DB, err = gorm.Open(mysql.Open(sql), cf)
			//defer DB.Close()
			if err != nil {
				Logger().Error("[mysql连接错误]:", err)
				continue
			}
			Logger().Info("[mysql重连成功]")
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
	return log2.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		log2.Config{
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
