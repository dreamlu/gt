// package gt

package gt

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"time"
)

// DB tool
type DBTool struct {
	// once
	once sync.Once
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
	db.DB = db.open(sql)
	err := db.DB.Exec("create database if not exists `" + dbS.Name + "`").Error
	if err == nil {
		err = db.DB.Close()
		if err != nil {
			Logger().Info("[mysql根连接]:", err)
		}
	} else {
		Logger().Info("[mysql自动连接根数据库失败,尝试直接连接]")
	}

	sql = fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local", dbS.User, dbS.Password, dbS.Host, dbS.Name)
	db.once.Do(func() {
		db.DB = db.open(sql)
	})
	// Globally disable table names
	// use name replace names
	db.DB.SingularTable(true)

	db.DB.LogMode(dbS.Log)
	db.DB.SetLogger(defaultLog)
	// connection pool
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.DB.DB().SetMaxIdleConns(dbS.MaxIdleConn)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.DB.DB().SetMaxOpenConns(dbS.MaxOpenConn)

	return
}

func (db *DBTool) open(sql string) *gorm.DB {
	// database, initialize once
	DB, err := gorm.Open("mysql", sql)
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
			DB, err = gorm.Open("mysql", sql)
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

// init DBTool
func NewDBTool() *DBTool {

	dbTool := &DBTool{}

	// init db
	dbTool.NewDB()
	return dbTool
}

var (
	onceDB sync.Once
	// dbTool
	// dbTool was global
	dbTool *DBTool
	// config
	//config = NewConfig()
)

// single db
func DBTooler() {

	onceDB.Do(func() {
		dbTool = NewDBTool()
	})
}

func (db *DBTool) clone() *DBTool {

	return &DBTool{
		DB:  db.DB,
		log: db.log,
		res: db.res,
	}
}
