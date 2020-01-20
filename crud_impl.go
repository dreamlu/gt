// package gt

package gt

import (
	"github.com/dreamlu/gt/tool/result"
	sq "github.com/dreamlu/gt/tool/sql"
	"github.com/dreamlu/gt/tool/util/str"
)

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
	// pager
	pager result.Pager
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
func (c *DBCrud) GetBySearch(params map[string][]string) Crud {
	clone := c.clone()
	clone.dbTool.GetDataBySearch(&GT{
		Table:  c.param.Table,
		Model:  c.param.Model,
		Data:   c.param.Data,
		Params: params,
		SubSQL: c.param.SubSQL,
	})
	return clone
}

func (c *DBCrud) GetByData(params map[string][]string) Crud {
	clone := c.clone()
	clone.dbTool.GetData(&GT{
		Table:  c.param.Table,
		Model:  c.param.Model,
		Data:   c.param.Data,
		Params: params,
		SubSQL: c.param.SubSQL,
	})
	return clone
}

// by id
func (c *DBCrud) GetByID(id string) Crud {

	clone := c.clone()
	clone.dbTool.GetDataByID(c.param.Data, id)
	return clone
}

// the same as search
// more tables
func (c *DBCrud) GetMoreBySearch(params map[string][]string) Crud {

	clone := c.clone()
	clone.dbTool.GetMoreDataBySearch(&GT{
		InnerTable: c.param.InnerTable,
		LeftTable:  c.param.LeftTable,
		Model:      c.param.Model,
		Data:       c.param.Data,
		Params:     params,
		SubSQL:     c.param.SubSQL,
	})
	return clone
}

// delete
func (c *DBCrud) Delete(id string) Crud {

	clone := c.clone()
	clone.dbTool.Delete(c.param.Table, id)
	return clone
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
func (c *DBCrud) CreateResID(params map[string][]string) (str.ID, error) {

	return c.dbTool.CreateDataResID(c.param.Table, params)
}

// == json data ==

// create
func (c *DBCrud) CreateMoreData() error {

	return c.dbTool.CreateMoreData(c.param.Table, c.param.Model, c.param.Data)
}

// update
func (c *DBCrud) Update() Crud {
	clone := c.clone()
	clone.dbTool.UpdateData(c.param.Data)
	return clone
}

// create
func (c *DBCrud) Create() Crud {
	clone := c.clone()
	clone.dbTool.CreateData(c.param.Data)
	return clone
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

func (c *DBCrud) Search() Crud {

	if c.argsNt == nil {
		c.argsNt = c.args
	}
	clone := c.clone()
	clone.dbTool.GetDataBySelectSQLSearch(&GT{
		Data:       c.param.Data,
		ClientPage: c.param.ClientPage,
		EveryPage:  c.param.EveryPage,
		Select:     c.selectSQL,
		Args:       c.args,
		ArgsNt:     c.argsNt,
		From:       c.from,
		Group:      c.group,
	})
	return clone
}

func (c *DBCrud) Single() Crud {

	c.Select(c.group)

	clone := c.clone()
	clone.dbTool.GetDataBySQL(c.param.Data, c.selectSQL, c.args...)
	return clone
}

func (c *DBCrud) Exec() Crud {

	clone := c.clone()
	clone.dbTool.ExecSQL(c.selectSQL, c.args...)
	return clone
}

func (c *DBCrud) Error() error {

	if c.dbTool.Error != nil {
		c.dbTool.Error = sq.GetSQLError(c.dbTool.Error.Error())
	}
	return c.dbTool.Error
}

func (c *DBCrud) RowsAffected() int64 {

	return c.dbTool.RowsAffected
}

func (c *DBCrud) Pager() result.Pager {

	return c.pager
}

func (c *DBCrud) clone() *DBCrud {

	dbCrud := &DBCrud{
		dbTool:    c.dbTool.clone(),
		param:     c.param,
		selectSQL: c.selectSQL,
		from:      c.from,
		args:      c.args,
		argsNt:    c.argsNt,
		group:     c.group,
	}
	return dbCrud
}
