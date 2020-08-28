// package gt

package gt

import (
	"fmt"
	"github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/result"
	sq "github.com/dreamlu/gt/tool/sql"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/errors"
	"github.com/dreamlu/gt/tool/util/hump"
	"github.com/dreamlu/gt/tool/util/str"
	"github.com/dreamlu/gt/tool/valid"
	"runtime"
	"strings"
)

// implement Crud
type Mysql struct {
	// DBTool  tool
	dbTool *DBTool
	// error
	err error
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
func (c *Mysql) initCrud(param *Params) {

	c.dbTool = dbTool
	c.param = param
	return
}

func (c *Mysql) DB() *DBTool {
	c.common()
	return c.dbTool
}

func (c *Mysql) Params(params ...Param) Crud {

	for _, p := range params {
		p(c.param)
	}
	return c
}

// search
// pager info
func (c *Mysql) GetBySearch(params cmap.CMap) Crud {
	c.common()
	clone := c.clone()
	clone.pager = clone.dbTool.GetDataBySearch(&GT{
		Params: clone.param,
		CMaps:  params,
	})

	return clone
}

func (c *Mysql) GetByData(params cmap.CMap) Crud {
	c.common()
	clone := c.clone()
	clone.dbTool.GetData(&GT{
		Params: clone.param,
		CMaps:  params,
	})
	return clone
}

func (c *Mysql) GetMoreByData(params cmap.CMap) Crud {
	c.common()
	clone := c.clone()
	clone.dbTool.GetMoreData(&GT{
		Params: clone.param,
		CMaps:  params,
	})
	return clone
}

// by id
func (c *Mysql) GetByID(id interface{}) Crud {
	c.common()

	clone := c.clone()
	clone.dbTool.GetDataByID(&GT{
		Params: clone.param,
	}, id)
	return clone
}

// the same as search
// more tables
func (c *Mysql) GetMoreBySearch(params cmap.CMap) Crud {
	c.common()

	clone := c.clone()
	clone.pager = clone.dbTool.GetMoreDataBySearch(&GT{
		CMaps:  params,
		Params: clone.param,
	})
	return clone
}

// delete
func (c *Mysql) Delete(id interface{}) Crud {
	c.common()

	clone := c.clone()
	clone.dbTool.Delete(clone.param.Table, id)
	return clone
}

// === form data ===

// update
func (c *Mysql) UpdateForm(params cmap.CMap) error {
	c.common()

	return c.dbTool.UpdateFormData(c.param.Table, params)
}

// create
func (c *Mysql) CreateForm(params cmap.CMap) error {
	c.common()

	return c.dbTool.CreateFormData(c.param.Table, params)
}

// create res insert id
func (c *Mysql) CreateResID(params cmap.CMap) (str.ID, error) {
	c.common()

	return c.dbTool.CreateDataResID(c.param.Table, params)
}

// == json data ==

// create more
func (c *Mysql) CreateMore() Crud {
	c.common()
	clone := c.clone()
	clone.err = check(clone.param.Data)
	if clone.err != nil {
		return clone
	}
	clone.dbTool.CreateMoreData(clone.param.Table, clone.param.Model, clone.param.Data)
	return clone
}

// update
func (c *Mysql) Update() Crud {
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
func (c *Mysql) Create() Crud {
	c.common()
	clone := c.clone()
	clone.err = check(clone.param.Data)
	if clone.err != nil {
		return clone
	}
	clone.dbTool.CreateData(clone.param.Table, clone.param.Data)
	return clone
}

// create
func (c *Mysql) Select(q interface{}, args ...interface{}) Crud {

	clone := c
	if c.selectSQL == "" {
		clone = c.clone()
	}

	var query string
	switch q.(type) {
	case string:
		query = q.(string)
	//case cmap.CMap:
	//	query, args = sq.CMapWhereSQL(q.(cmap.CMap))
	case interface{}:
		query, args = sq.StructWhereSQL(q)
	}

	clone.selectSQL += query + " "
	clone.args = append(clone.args, args...)
	if clone.from != "" {
		clone.argsNt = append(clone.argsNt, args...)
	}
	return clone
}

func (c *Mysql) From(query string) Crud {

	c.from = query
	c.selectSQL += query + " "
	return c
}

func (c *Mysql) Group(query string) Crud {

	c.group = query
	return c
}

func (c *Mysql) Search(params cmap.CMap) Crud {
	c.common()

	if c.argsNt == nil {
		c.argsNt = c.args
	}
	//clone := c
	c.pager = c.dbTool.GetDataBySelectSQLSearch(&GT{
		Params: c.param,
		Select: c.selectSQL,
		Args:   c.args,
		From:   c.from,
		Group:  c.group,
		CMaps:  params,
	})
	return c
}

func (c *Mysql) Single() Crud {
	c.common()

	c.Select(c.group)

	//clone := c.clone()
	c.dbTool.GetDataBySQL(c.param.Data, c.selectSQL, c.args...)
	return c
}

func (c *Mysql) Exec() Crud {
	c.common()

	//clone := c.clone()
	c.dbTool.ExecSQL(c.selectSQL, c.args...)
	return c
}

func (c *Mysql) Error() error {

	if c.err != nil {
		return c.err
	}
	if c.dbTool.res != nil {
		c.err = c.dbTool.res.Error
		if c.err != nil {
			c.err = sq.GetSQLError(c.err.Error())
		}
	}
	return c.err
}

func (c *Mysql) RowsAffected() int64 {

	return c.dbTool.RowsAffected
}

func (c *Mysql) Pager() result.Pager {

	return c.pager
}

func (c *Mysql) Begin() Crud {
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

func (c *Mysql) Commit() Crud {
	if c.dbTool.res.Error != nil {
		c.dbTool.Rollback()
	}
	c.dbTool.Commit()
	c.isTrans = 0
	return c
}

func (c *Mysql) Rollback() Crud {
	c.dbTool.Rollback()
	return c
}

func (c *Mysql) clone() (dbCrud *Mysql) {

	// default table
	if c.param.Table == "" &&
		c.param.Model != nil {
		c.param.Table = hump.HumpToLine(reflect.StructToString(c.param.Model))
	}

	dbCrud = &Mysql{
		dbTool:    c.dbTool,
		param:     c.param,
		selectSQL: c.selectSQL,
		from:      c.from,
		args:      c.args,
		argsNt:    c.argsNt,
		group:     c.group,
		err:       c.err,
	}

	// isTrans
	if c.isTrans == 1 {
		return
	}
	dbCrud.dbTool = c.dbTool.clone()
	return
}

func (c *Mysql) common() {
	if c.dbTool.log {
		c.line()
	}
}

func (c *Mysql) line() {
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
		fmt.Fprintf(buf, "%s:%d\n", fullFile, line)
		fmt.Print(buf.String())
	}
}

// auto valid data
func check(data interface{}) error {

	ves := valid.Valid(data)
	if len(ves) > 0 {
		var s string
		for k, v := range ves {
			s += fmt.Sprintf("%s:%s;", k, v)
		}
		return errors.New(strings.TrimRight(s, ";"))
	}
	return nil
}
