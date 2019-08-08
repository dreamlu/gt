// package der

/*
	go-tool is a fast go tool, help you dev project
*/

package der

import (
	"github.com/dreamlu/go-tool/tool/result"
)

const Version = "1.1.x"

// crud
type Crud interface {
	InitDBTool() *DBTool
	// crud method

	// get url params
	// like form data
	GetBySearch(args map[string][]string) (pager result.Pager, err error)     // search
	GetByID(id string) error                                                  // by id
	GetMoreBySearch(args map[string][]string) (pager result.Pager, err error) // more search

	// common sql data
	// through sql, get the data
	GetDataBySQL(sql string, args ...interface{}) error // single data
	// page limit ?,?
	// args not include limit ?,?
	GetDataBySearchSQL(sql, sqlNt string, args ...interface{}) (pager result.Pager, err error) // more data
	DeleteBySQL(sql string, args ...interface{}) error
	UpdateBySQL(sql string, args ...interface{}) error
	CreateBySQL(sql string, args ...interface{}) error

	// delete
	Delete(id string) error // delete
}

// common crud
// detail impl, ==>DbCrud, implement DBCrud
// form data
type DBCruder interface {
	// common sql data
	Crud

	// crud and search id
	Update(args map[string][]string) error            // update
	Create(args map[string][]string) error            // create
	CreateResID(args map[string][]string) (ID, error) // create res insert id
}

// common crud
// json data
type DBCrudJer interface {
	// common sql data
	Crud

	// crud and search id
	Update(data interface{}) error          // update
	Create(data interface{}) error          // create, include res insert id
	CreateMoreDataJ(data interface{}) error // create more
}

// crud params
type CrudParam struct {
	// attributes
	InnerTables []string    // inner join tables
	LeftTables  []string    // left join tables
	Table       string      // table name
	Model       interface{} // table model, like User{}
	ModelData   interface{} // table model data, like var user User{}, it is 'user'

	// pager info
	ClientPage int64 // page number
	EveryPage  int64 // Number of pages per page
}

// crud common
type DBCrudCom struct {
	*DBTool
	Param  CrudParam
}

// init db tool
func (c *DBCrudCom) InitDBTool() *DBTool {

	c.DBTool = c.DBTool.NewDB()
	return c.DBTool
}

// delete
func (c *DBCrudCom) Delete(id string) error {

	return c.DBTool.DeleteDataByName(c.Param.Table, "id", id)
}

// search
// pager info
// clientPage : default 1
// everyPage : default 10
func (c *DBCrudCom) GetBySearch(params map[string][]string) (pager result.Pager, err error) {

	return c.DBTool.GetDataBySearch(c.Param.Model, c.Param.ModelData, c.Param.Table, params)
}

// by id
func (c *DBCrudCom) GetByID(id string) error {

	//DB.AutoMigrate(&c.Model)
	return c.DBTool.GetDataByID(c.Param.ModelData, id)
}

// the same as search
// more tables
func (c *DBCrudCom) GetMoreBySearch(params map[string][]string) (pager result.Pager, err error) {

	return c.DBTool.GetMoreDataBySearch(c.Param.Model, c.Param.ModelData, params, c.Param.InnerTables, c.Param.LeftTables)
}

// common sql
// through sql get data
func (c *DBCrudCom) GetDataBySQL(sql string, args ...interface{}) error {

	return c.DBTool.GetDataBySQL(c.Param.ModelData, sql, args[:]...)
}

// common sql
// through sql get data
// args not include limit ?, ?
// args is sql and sqlNt common params
func (c *DBCrudCom) GetDataBySearchSQL(sql, sqlNt string, args ...interface{}) (pager result.Pager, err error) {

	return c.DBTool.GetDataBySQLSearch(c.Param.ModelData, sql, sqlNt, c.Param.ClientPage, c.Param.EveryPage, args)
}

// delete by sql
func (c *DBCrudCom) DeleteBySQL(sql string, args ...interface{}) error {

	return c.DBTool.DeleteDataBySQL(sql, args[:]...)
}

// update by sql
func (c *DBCrudCom) UpdateBySQL(sql string, args ...interface{}) error {

	return c.DBTool.UpdateDataBySQL(sql, args[:]...)
}

// create by sql
func (c *DBCrudCom) CreateBySQL(sql string, args ...interface{}) error {

	return c.DBTool.CreateDataBySQL(sql, args[:]...)
}
