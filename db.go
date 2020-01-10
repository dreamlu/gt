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
func (db *DBTool) NewDB() *gorm.DB {

	DB := &gorm.DB{}
	dbS := &dba{}
	Configger().GetStruct("app.db", dbS)
	var (
		sql = fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local", dbS.User, dbS.Password, dbS.Host, dbS.Name)
	)

	db.once.Do(func() {
		DB = db.open(sql)
	})
	// Globally disable table names
	// use name replace names
	DB.SingularTable(true)

	DB.LogMode(dbS.Log)
	// connection pool
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	DB.DB().SetMaxIdleConns(dbS.MaxIdleConn)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	DB.DB().SetMaxOpenConns(dbS.MaxOpenConn)

	return DB
}

func (db *DBTool) open(sql string) *gorm.DB {
	// database, initialize once
	DB, err := gorm.Open("mysql", sql)
	//defer db.DB.Close()
	if err != nil {
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
	dbTool.DB = dbTool.NewDB()
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
func DBTooler() *DBTool {

	onceDB.Do(func() {
		dbTool = NewDBTool()
	})

	return dbTool
}

func (db *DBTool) clone() *DBTool {

	return &DBTool{
		DB: db.DB,
	}
}
