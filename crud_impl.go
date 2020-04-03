// package gt

package gt

import (
	"github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/result"
	sq "github.com/dreamlu/gt/tool/sql"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/hump"
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

	// transaction
	isTrans int8 // open(0), close(1)
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
func (c *DBCrud) GetBySearch(params cmap.CMap) Crud {
	clone := c.clone()
	clone.pager = clone.dbTool.GetDataBySearch(&GT{
		Params: clone.param,
		CMaps:  params,
	})

	return clone
}

func (c *DBCrud) GetByData(params cmap.CMap) Crud {
	clone := c.clone()
	clone.dbTool.GetData(&GT{
		Params: clone.param,
		CMaps:  params,
	})
	return clone
}

// by id
func (c *DBCrud) GetByID(id interface{}) Crud {

	clone := c.clone()
	clone.dbTool.GetDataByID(clone.param.Data, id)
	return clone
}

// the same as search
// more tables
func (c *DBCrud) GetMoreBySearch(params cmap.CMap) Crud {

	clone := c.clone()
	clone.pager = clone.dbTool.GetMoreDataBySearch(&GT{
		CMaps:  params,
		Params: clone.param,
	})
	return clone
}

// delete
func (c *DBCrud) Delete(id interface{}) Crud {

	clone := c.clone()
	clone.dbTool.Delete(clone.param.Table, id)
	return clone
}

// === form data ===

// update
func (c *DBCrud) UpdateForm(params cmap.CMap) error {

	return c.dbTool.UpdateFormData(c.param.Table, params)
}

// create
func (c *DBCrud) CreateForm(params cmap.CMap) error {

	return c.dbTool.CreateFormData(c.param.Table, params)
}

// create res insert id
func (c *DBCrud) CreateResID(params cmap.CMap) (str.ID, error) {

	return c.dbTool.CreateDataResID(c.param.Table, params)
}

// == json data ==

// create
func (c *DBCrud) CreateMoreData() Crud {

	clone := c.clone()
	clone.dbTool.CreateMoreData(clone.param.Table, clone.param.Model, clone.param.Data)
	return clone
}

// update
func (c *DBCrud) Update() Crud {
	clone := c.clone()
	clone.dbTool.UpdateData(&GT{
		Params: clone.param,
		Select: clone.selectSQL,
		Args:   clone.args,
	})
	return clone
}

// create
func (c *DBCrud) Create() Crud {
	clone := c.clone()
	clone.dbTool.CreateData(clone.param.Data)
	return clone
}

// create
func (c *DBCrud) Select(query string, args ...interface{}) Crud {

	clone := c.clone()
	clone.selectSQL += query + " "
	clone.args = append(clone.args, args...)
	if clone.from != "" {
		clone.argsNt = append(clone.argsNt, args...)
	}
	return clone
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
	clone.pager = clone.dbTool.GetDataBySelectSQLSearch(&GT{
		Params: clone.param,
		Select: clone.selectSQL,
		Args:   clone.args,
		ArgsNt: clone.argsNt,
		From:   clone.from,
		Group:  clone.group,
	})
	return clone
}

func (c *DBCrud) Single() Crud {

	c.Select(c.group)

	clone := c.clone()
	clone.dbTool.GetDataBySQL(clone.param.Data, c.selectSQL, c.args...)
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

func (c *DBCrud) Begin() Crud {
	clone := c.clone()
	clone.isTrans = 1
	clone.dbTool.DB = clone.dbTool.Begin()
	defer func() {
		if r := recover(); r != nil {
			clone.dbTool.Rollback()
		}
	}()
	return clone
}

func (c *DBCrud) Commit() Crud {
	if c.dbTool.Error != nil {
		c.dbTool.Rollback()
	}
	c.dbTool.Commit()
	c.isTrans = 0
	return c
}

func (c *DBCrud) Rollback() Crud {
	c.dbTool.Rollback()
	return c
}

func (c *DBCrud) clone() (dbCrud *DBCrud) {

	// default table
	if c.param.Table == "" &&
		c.param.Model != nil {
		c.param.Table = hump.HumpToLine(reflect.StructToString(c.param.Model))
	}

	dbCrud = &DBCrud{
		dbTool:    c.dbTool,
		param:     c.param,
		selectSQL: c.selectSQL,
		from:      c.from,
		args:      c.args,
		argsNt:    c.argsNt,
		group:     c.group,
	}

	// isTrans
	if c.isTrans == 1 {
		return
	}
	dbCrud.dbTool = c.dbTool.clone()
	return
}
