// package der

package der

import "github.com/dreamlu/go-tool/tool/result"

// implement DBCrud
// form data
type DBCrud struct {
	// db  tool
	db    *DBTool
	//Param CrudParam
}

// init db tool
func (c *DBCrud) InitDBTool(dbTool *DBTool) {

	c.db = dbTool
	return
}

// search
// pager info
// clientPage : default 1
// everyPage : default 10
func (c *DBCrud) GetBySearch(params map[string][]string) (pager result.Pager, err error) {

	return c.db.GetDataBySearch(c.db.Param.Model, c.db.Param.ModelData, c.db.Param.Table, params)
}

// by id
func (c *DBCrud) GetByID(id string) error {

	//DB.AutoMigrate(&c.Model)
	return c.db.GetDataByID(c.db.Param.ModelData, id)
}

// the same as search
// more tables
func (c *DBCrud) GetMoreBySearch(params map[string][]string) (pager result.Pager, err error) {

	return c.db.GetMoreDataBySearch(c.db.Param.Model, c.db.Param.ModelData, params, c.db.Param.InnerTables, c.db.Param.LeftTables)
}

// common sql
// through sql get data
func (c *DBCrud) GetDataBySQL(sql string, args ...interface{}) error {

	return c.db.GetDataBySQL(c.db.Param.ModelData, sql, args[:]...)
}

// common sql
// through sql get data
// args not include limit ?, ?
// args is sql and sqlNt common params
func (c *DBCrud) GetDataBySearchSQL(sql, sqlNt string, args ...interface{}) (pager result.Pager, err error) {

	return c.db.GetDataBySQLSearch(c.db.Param.ModelData, sql, sqlNt, c.db.Param.ClientPage, c.db.Param.EveryPage, args)
}

// delete by sql
func (c *DBCrud) DeleteBySQL(sql string, args ...interface{}) error {

	return c.db.DeleteDataBySQL(sql, args[:]...)
}

// update by sql
func (c *DBCrud) UpdateBySQL(sql string, args ...interface{}) error {

	return c.db.UpdateDataBySQL(sql, args[:]...)
}

// create by sql
func (c *DBCrud) CreateBySQL(sql string, args ...interface{}) error {

	return c.db.CreateDataBySQL(sql, args[:]...)
}

// delete
func (c *DBCrud) Delete(id string) error {

	return c.db.DeleteDataByName(c.db.Param.Table, "id", id)
}

// === form data ===

// update
func (c *DBCrud) UpdateForm(params map[string][]string) error {

	return c.db.UpdateData(c.db.Param.Table, params)
}

// create
func (c *DBCrud) CreateForm(params map[string][]string) error {

	return c.db.CreateData(c.db.Param.Table, params)
}

// create res insert id
func (c *DBCrud) CreateResID(params map[string][]string) (ID, error) {

	return c.db.CreateDataResID(c.db.Param.Table, params)
}

// == json data ==

// create
func (c *DBCrud) CreateMoreDataJ(data interface{}) error {

	return c.db.CreateMoreDataJ(c.db.Param.Table, c.db.Param.Model, data)
}

// update
func (c *DBCrud) Update(data interface{}) error {

	return c.db.UpdateDataJ(data)
}

// create
func (c *DBCrud) Create(data interface{}) error {

	return c.db.CreateDataJ(data)
}
