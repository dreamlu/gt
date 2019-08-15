// package der

package der

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"strconv"
	"time"
)

// DB tool
type DBTool struct {
	// db driver
	DB *gorm.DB
	// crud interface
	Crud Crud
	// crud param
	Param CrudParam
}

// db params
type dba struct {
	user     string
	password string
	host     string
	name     string
}

func (db *DBTool) NewDB() *gorm.DB {

	dbS := &dba{
		user:     GetDevModeConfig("db.user"),
		password: GetDevModeConfig("db.password"),
		host:     GetDevModeConfig("db.host"),
		name:     GetDevModeConfig("db.name"),
	}
	var (
		err error
		sql = fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local", dbS.user, dbS.password, dbS.host, dbS.name)
	)

	//database, initialize once
	DB, err := gorm.Open("mysql", sql)
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
	// Globally disable table names
	// use name replace names
	DB.SingularTable(true)
	// sql print console log
	// or print sql err to file
	//LogMode("debug") // or sqlErr

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

	return DB
}

func (db *DBTool) NewDBTool() *DBTool {

	dbTool := &DBTool{
		DB:   db.NewDB(),
		Crud: &DBCrud{},
	}

	dbTool.Crud.InitDBTool(dbTool)
	return dbTool
}
