// package der

package der

import "github.com/dreamlu/go-tool/tool/result"

// implement DBCrud
// form data
type DBCrud struct {
	// DBTool  tool
	DBTool *DBTool
	// crud param
	Param *CrudParam
}

// init DBTool tool
func (c *DBCrud) InitDBTool(dbTool *DBTool, param *CrudParam) {

	c.DBTool = dbTool
	c.Param = param
	return
}

// search
// pager info
// clientPage : default 1
// everyPage : default 10
func (c *DBCrud) GetBySearch(params map[string][]string) (pager result.Pager, err error) {

	return c.DBTool.GetDataBySearch(c.Param.Model, c.Param.ModelData, c.Param.Table, params)
}

// by id
func (c *DBCrud) GetByID(id string) error {

	//DB.AutoMigrate(&c.Model)
	return c.DBTool.GetDataByID(c.Param.ModelData, id)
}

// the same as search
// more tables
func (c *DBCrud) GetMoreBySearch(params map[string][]string) (pager result.Pager, err error) {

	return c.DBTool.GetMoreDataBySearch(c.Param.Model, c.Param.ModelData, params, c.Param.InnerTables, c.Param.LeftTables)
}

// common sql
// through sql get data
func (c *DBCrud) GetDataBySQL(sql string, args ...interface{}) error {

	return c.DBTool.GetDataBySQL(c.Param.ModelData, sql, args[:]...)
}

// common sql
// through sql get data
// args not include limit ?, ?
// args is sql and sqlNt common params
func (c *DBCrud) GetDataBySearchSQL(sql, sqlNt string, args ...interface{}) (pager result.Pager, err error) {

	return c.DBTool.GetDataBySQLSearch(c.Param.ModelData, sql, sqlNt, c.Param.ClientPage, c.Param.EveryPage, args)
}

// delete by sql
func (c *DBCrud) DeleteBySQL(sql string, args ...interface{}) error {

	return c.DBTool.DeleteDataBySQL(sql, args[:]...)
}

// update by sql
func (c *DBCrud) UpdateBySQL(sql string, args ...interface{}) error {

	return c.DBTool.UpdateDataBySQL(sql, args[:]...)
}

// create by sql
func (c *DBCrud) CreateBySQL(sql string, args ...interface{}) error {

	return c.DBTool.CreateDataBySQL(sql, args[:]...)
}

// delete
func (c *DBCrud) Delete(id string) error {

	return c.DBTool.DeleteDataByName(c.Param.Table, "id", id)
}

// === form data ===

// update
func (c *DBCrud) UpdateForm(params map[string][]string) error {

	return c.DBTool.UpdateData(c.Param.Table, params)
}

// create
func (c *DBCrud) CreateForm(params map[string][]string) error {

	return c.DBTool.CreateData(c.Param.Table, params)
}

// create res insert id
func (c *DBCrud) CreateResID(params map[string][]string) (ID, error) {

	return c.DBTool.CreateDataResID(c.Param.Table, params)
}

// == json data ==

// create
func (c *DBCrud) CreateMoreData(data interface{}) error {

	return c.DBTool.CreateMoreDataJ(c.Param.Table, c.Param.Model, data)
}

// update
func (c *DBCrud) Update(data interface{}) error {

	return c.DBTool.UpdateDataJ(data)
}

// create
func (c *DBCrud) Create(data interface{}) error {

	return c.DBTool.CreateDataJ(data)
}
