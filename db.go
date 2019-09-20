// package der

package der

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
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
	user     string
	password string
	host     string
	name     string
}

// new db driver
func (db *DBTool) NewDB() *gorm.DB {

	DB := &gorm.DB{}
	conf := Configger()

	dbS := &dba{
		user:     conf.GetString("app.db.user"),
		password: conf.GetString("app.db.password"),
		host:     conf.GetString("app.db.host"),
		name:     conf.GetString("app.db.name"),
	}
	var (
		err error
		sql = fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local", dbS.user, dbS.password, dbS.host, dbS.name)
	)

	db.once.Do(func() {
		//database, initialize once
		DB, err = gorm.Open("mysql", sql)
		//defer db.DB.Close()
		if err != nil {
			log.Println("[mysql连接错误]:", err)
			log.Println("[mysql开始尝试重连中]: try it every 5s...")
			// try to reconnect
			for {
				// go is so fast
				// try it every 5s
				time.Sleep(5 * time.Second)
				DB, err = gorm.Open("mysql", sql)
				//defer DB.Close()
				if err != nil {
					log.Println("[mysql连接错误]:", err)
					continue
				}
				log.Println("[mysql重连成功]")
				break
			}
		}
	})
	// Globally disable table names
	// use name replace names
	DB.SingularTable(true)
	// sql print console log
	// or print sql err to file
	//LogMode("debug") // or sqlErr

	// connection pool
	var maxIdle, maxOpen int
	var logMode bool
	if maxIdle = conf.GetInt("app.db.maxIdleConn"); maxIdle == 0 {
		maxIdle = 20
	}
	if maxOpen = conf.GetInt("app.db.maxOpenConn"); maxOpen == 0 {
		maxOpen = 100
	}
	logMode = conf.GetBool("app.db.log")
	DB.LogMode(logMode)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	DB.DB().SetMaxIdleConns(maxIdle)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	DB.DB().SetMaxOpenConns(maxOpen)

	return DB
}

// init DBTool
func NewDBTool() *DBTool {

	dbTool := &DBTool{}

	// init db
	dbTool.DB = dbTool.NewDB()
	return dbTool
}
