// package der

package der

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"strconv"
	"time"
)

var (
	DB *gorm.DB
)

func NewDB() {
	var err error
	//database, initialize once
	DB, err = gorm.Open("mysql", GetDevModeConfig("db.user")+":"+GetDevModeConfig("db.password")+"@"+GetDevModeConfig("db.host")+"/"+GetDevModeConfig("db.name")+"?charset=utf8&parseTime=True&loc=Local")
	//defer DB.Close()
	if err != nil {
		log.Println("[mysql连接错误]:", err)
		log.Println("[mysql开始尝试重连中]: try it every 5s...")
		// try to reconnect
		for {
			// go is so fast
			// try it every 5s
			time.Sleep(5 * time.Second)
			DB, err = gorm.Open("mysql", GetDevModeConfig("db.user")+":"+GetDevModeConfig("db.password")+"@"+GetDevModeConfig("db.host")+"/"+GetDevModeConfig("db.name")+"?charset=utf8&parseTime=True&loc=Local")
			//defer DB.Close()
			if err != nil {
				log.Println("[mysql连接错误]:", err)
				continue
			}
			log.Println("[mysql重连成功]")
			break
		}
	}
	// Globally disable table names
	// use name replace names
	DB.SingularTable(true)
	// sql print console log
	// or print sql err to file
	LogMode("debug") // or sqlErr

	// connection pool
	var maxIdle, maxOpen int
	maxIdleConn := GetDevModeConfig("db.maxIdleConn")
	if maxIdleConn == "" {
		maxIdle = 20
	}
	maxIdle, _ = strconv.Atoi(maxIdleConn)

	maxOpenConn := GetDevModeConfig("db.maxOpenConn")
	if maxOpenConn == "" {
		maxOpen = 100
	}
	maxOpen, _ = strconv.Atoi(maxIdleConn)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	DB.DB().SetMaxIdleConns(maxIdle)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	DB.DB().SetMaxOpenConns(maxOpen)
}
