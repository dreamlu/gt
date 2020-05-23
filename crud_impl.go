// package gt

package gt

import (
	"fmt"
	"github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/result"
	sq "github.com/dreamlu/gt/tool/sql"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/hump"
	"github.com/dreamlu/gt/tool/util/str"
	"runtime"
	"strings"
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
	isTrans byte // open(0), close(1)
}

// init DBTool tool
func (c *DBCrud) initCrud(dbTool *DBTool, param *Params) {

	c.dbTool = dbTool
	c.param = param
	return
}

func (c *DBCrud) DB() *DBTool {
	c.common()
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
	c.common()
	clone := c.clone()
	clone.pager = clone.dbTool.GetDataBySearch(&GT{
		Params: clone.param,
		CMaps:  params,
	})

	return clone
}

func (c *DBCrud) GetByData(params cmap.CMap) Crud {
	c.common()
	clone := c.clone()
	clone.dbTool.GetData(&GT{
		Params: clone.param,
		CMaps:  params,
	})
	return clone
}

// by id
func (c *DBCrud) GetByID(id interface{}) Crud {
	c.common()

	clone := c.clone()
	clone.dbTool.GetDataByID(clone.param.Data, id)
	return clone
}

// the same as search
// more tables
func (c *DBCrud) GetMoreBySearch(params cmap.CMap) Crud {
	c.common()

	clone := c.clone()
	clone.pager = clone.dbTool.GetMoreDataBySearch(&GT{
		CMaps:  params,
		Params: clone.param,
	})
	return clone
}

// delete
func (c *DBCrud) Delete(id interface{}) Crud {
	c.common()

	clone := c.clone()
	clone.dbTool.Delete(clone.param.Table, id)
	return clone
}

// === form data ===

// update
func (c *DBCrud) UpdateForm(params cmap.CMap) error {
	c.common()

	return c.dbTool.UpdateFormData(c.param.Table, params)
}

// create
func (c *DBCrud) CreateForm(params cmap.CMap) error {
	c.common()

	return c.dbTool.CreateFormData(c.param.Table, params)
}

// create res insert id
func (c *DBCrud) CreateResID(params cmap.CMap) (str.ID, error) {
	c.common()

	return c.dbTool.CreateDataResID(c.param.Table, params)
}

// == json data ==

// create
func (c *DBCrud) CreateMoreData() Crud {

	clone := c.clone()
	clone.dbTool.CreateMoreData(clone.param.Table, clone.param.Model, clone.param.Data)
	return clone
}

// create more
func (c *DBCrud) CreateMore() Crud {
	c.common()

	clone := c.clone()
	clone.dbTool.CreateMoreData(clone.param.Table, clone.param.Model, clone.param.Data)
	return clone
}

// update
func (c *DBCrud) Update() Crud {
	c.common()
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
	c.common()
	clone := c.clone()
	clone.dbTool.CreateData(clone.param.Table, clone.param.Data)
	return clone
}

// create
func (c *DBCrud) Select(query string, args ...interface{}) Crud {

	clone := c
	if c.selectSQL == "" {
		clone = c.clone()
	}

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

func (c *DBCrud) Search(params cmap.CMap) Crud {
	c.common()

	if c.argsNt == nil {
		c.argsNt = c.args
	}
	//clone := c
	c.pager = c.dbTool.GetDataBySelectSQLSearch(&GT{
		Params: c.param,
		Select: c.selectSQL,
		Args:   c.args,
		ArgsNt: c.argsNt,
		From:   c.from,
		Group:  c.group,
		CMaps:  params,
	})
	return c
}

func (c *DBCrud) Single() Crud {
	c.common()

	c.Select(c.group)

	//clone := c.clone()
	c.dbTool.GetDataBySQL(c.param.Data, c.selectSQL, c.args...)
	return c
}

func (c *DBCrud) Exec() Crud {
	c.common()

	//clone := c.clone()
	c.dbTool.ExecSQL(c.selectSQL, c.args...)
	return c
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

func (c *DBCrud) common() {
	if c.dbTool.log {
		c.line()
	}
}

func (c *DBCrud) line() {
	_, fullFile, line, ok := runtime.Caller(3) // 3 skip
	file := fullFile
	if file != "" {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
	}
	if ok {
		buf := new(strings.Builder)

		fmt.Fprintf(buf, "\n\033[35m[gt]\033[0m: ")
		fmt.Fprintf(buf, "%s:%d", fullFile, line)
		fmt.Print(buf.String())
	}
}
