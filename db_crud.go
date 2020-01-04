// package gt

package gt

import "github.com/dreamlu/go-tool/tool/result"

// implement DBCrud
// form data
type DBCrud struct {
	// DBTool  tool
	dbTool *DBTool
	// crud param
	param *Params
}

// init DBTool tool
func (c *DBCrud) InitCrud(dbTool *DBTool, param *Params) {

	c.dbTool = dbTool
	c.param = param
	return
}

func (c *DBCrud) DB() *DBTool {
	return c.dbTool
}

func (c *DBCrud) Params(params ...Param) *Params {

	for _, p := range params {
		p(c.param)
	}
	return c.param
}

// search
// pager info
// clientPage : default 1
// everyPage : default 10
func (c *DBCrud) GetBySearch(params map[string][]string) (pager result.Pager, err error) {
	//c.param.Model, c.param.ModelData, c.param.Table, params
	return c.dbTool.GetDataBySearch(&GT{
		Table:     c.param.Table,
		Model:     c.param.Model,
		ModelData: c.param.ModelData,
		Params:    params,
	})
}

// by id
func (c *DBCrud) GetByID(id string) error {

	//DB.AutoMigrate(&c.Model)
	return c.dbTool.GetDataByID(c.param.ModelData, id)
}

// the same as search
// more tables
func (c *DBCrud) GetMoreBySearch(params map[string][]string) (pager result.Pager, err error) {

	return c.dbTool.GetMoreDataBySearch(&GT{
		InnerTable: c.param.InnerTable,
		LeftTable:  c.param.LeftTable,
		Model:      c.param.Model,
		ModelData:  c.param.ModelData,
		Params:     params,
	})
}

// common sql
// through sql get data
func (c *DBCrud) GetDataBySQL(sql string, args ...interface{}) error {

	return c.dbTool.GetDataBySQL(c.param.ModelData, sql, args[:]...)
}

// common sql
// through sql get data
// args not include limit ?, ?
// args is sql and sqlNt common params
func (c *DBCrud) GetDataBySearchSQL(sql, sqlNt string, args ...interface{}) (pager result.Pager, err error) {

	return c.dbTool.GetDataBySQLSearch(c.param.ModelData, sql, sqlNt, c.param.ClientPage, c.param.EveryPage, args)
}

// delete by sql
func (c *DBCrud) DeleteBySQL(sql string, args ...interface{}) error {

	return c.dbTool.DeleteDataBySQL(sql, args[:]...)
}

// update by sql
func (c *DBCrud) UpdateBySQL(sql string, args ...interface{}) error {

	return c.dbTool.UpdateDataBySQL(sql, args[:]...)
}

// create by sql
func (c *DBCrud) CreateBySQL(sql string, args ...interface{}) error {

	return c.dbTool.CreateDataBySQL(sql, args[:]...)
}

// delete
func (c *DBCrud) Delete(id string) error {

	return c.dbTool.DeleteDataByName(c.param.Table, "id", id)
}

// === form data ===

// update
func (c *DBCrud) UpdateForm(params map[string][]string) error {

	return c.dbTool.UpdateData(c.param.Table, params)
}

// create
func (c *DBCrud) CreateForm(params map[string][]string) error {

	return c.dbTool.CreateData(c.param.Table, params)
}

// create res insert id
func (c *DBCrud) CreateResID(params map[string][]string) (ID, error) {

	return c.dbTool.CreateDataResID(c.param.Table, params)
}

// == json data ==

// create
func (c *DBCrud) CreateMoreData(data interface{}) error {

	return c.dbTool.CreateMoreDataJ(c.param.Table, c.param.Model, data)
}

// update
func (c *DBCrud) Update(data interface{}) error {

	return c.dbTool.UpdateDataJ(data)
}

// create
func (c *DBCrud) Create(data interface{}) error {

	return c.dbTool.CreateDataJ(data)
}
