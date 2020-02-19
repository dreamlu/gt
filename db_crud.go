// package gt

package gt

import "github.com/dreamlu/gt/tool/result"

// implement DBCrud
// form data
type DBCrud struct {
	// DBTool  tool
	dbTool *DBTool
	// crud param
	param *Params

	// select
	selectSQL string        // select/or if
	from      string        // from sql
	args      []interface{} // select args
	argsNt    []interface{} // select nt args, related from
	group     string        // the last group
}

// init DBTool tool
func (c *DBCrud) initCrud(dbTool *DBTool, param *Params) {

	c.dbTool = dbTool
	c.param = param
	return
}

func (c *DBCrud) DB() *DBTool {
	return c.dbTool
}

func (c *DBCrud) Params(params ...Param) Crud {

	for _, p := range params {
		p(c.param)
	}
	return c
}

// search
// pager info
func (c *DBCrud) GetBySearch(params map[string][]string) (pager result.Pager, err error) {
	return c.dbTool.GetDataBySearch(&GT{
		Table:       c.param.Table,
		Model:       c.param.Model,
		Data:        c.param.Data,
		Params:      params,
		SubSQL:      c.param.SubSQL,
		SubWhereSQL: c.param.SubWhereSQL,
	})
}

func (c *DBCrud) GetByData(params map[string][]string) error {
	return c.dbTool.GetData(&GT{
		Table:       c.param.Table,
		Model:       c.param.Model,
		Data:        c.param.Data,
		Params:      params,
		SubSQL:      c.param.SubSQL,
		SubWhereSQL: c.param.SubWhereSQL,
	})
}

// by id
func (c *DBCrud) GetByID(id string) error {

	return c.dbTool.GetDataByID(c.param.Data, id)
}

// the same as search
// more tables
func (c *DBCrud) GetMoreBySearch(params map[string][]string) (pager result.Pager, err error) {

	return c.dbTool.GetMoreDataBySearch(&GT{
		InnerTable:    c.param.InnerTable,
		LeftTable:     c.param.LeftTable,
		Model:         c.param.Model,
		Data:          c.param.Data,
		Params:        params,
		SubSQL:        c.param.SubSQL,
		SubWhereSQL:   c.param.SubWhereSQL,
		OtherTableSQL: c.param.OtherTableSQL,
	})
}

// delete
func (c *DBCrud) Delete(id string) (err error) {

	return c.dbTool.Delete(c.param.Table, id)
}

// === form data ===

// update
func (c *DBCrud) UpdateForm(params map[string][]string) error {

	return c.dbTool.UpdateFormData(c.param.Table, params)
}

// create
func (c *DBCrud) CreateForm(params map[string][]string) error {

	return c.dbTool.CreateFormData(c.param.Table, params)
}

// create res insert id
func (c *DBCrud) CreateResID(params map[string][]string) (ID, error) {

	return c.dbTool.CreateDataResID(c.param.Table, params)
}

// == json data ==

// create
func (c *DBCrud) CreateMoreData() error {

	return c.dbTool.CreateMoreData(c.param.Table, c.param.Model, c.param.Data)
}

// update
func (c *DBCrud) Update() (err error) {

	return c.dbTool.UpdateData(c.param.Data)
}

// create
func (c *DBCrud) Create() error {

	return c.dbTool.CreateData(c.param.Data)
}

// create
func (c *DBCrud) Select(query string, args ...interface{}) Crud {

	c.selectSQL += query + " "
	c.args = append(c.args, args...)
	if c.from != "" {
		c.argsNt = append(c.argsNt, args...)
	}
	return c
}

func (c *DBCrud) From(query string) Crud {

	c.from = query
	c.selectSQL += query + " "
	return c
}

func (c *DBCrud) Group(query string) Crud {

	c.group = query
	return c
}

func (c *DBCrud) Search() (pager result.Pager, err error) {

	if c.argsNt == nil {
		c.argsNt = c.args
	}
	return c.dbTool.GetDataBySelectSQLSearch(&GT{
		Data:       c.param.Data,
		ClientPage: c.param.ClientPage,
		EveryPage:  c.param.EveryPage,
		Select:     c.selectSQL,
		Args:       c.args,
		ArgsNt:     c.argsNt,
		From:       c.from,
		Group:      c.group,
	})
}

func (c *DBCrud) Single() (err error) {

	c.Select(c.group)
	return c.dbTool.GetDataBySQL(c.param.Data, c.selectSQL, c.args...)
}

func (c *DBCrud) Exec() (err error) {

	return c.dbTool.ExecSQL(c.selectSQL, c.args...)
}
