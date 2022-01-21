// package gt

package gt

import (
	"errors"
	"fmt"
	"github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/hump"
	"github.com/dreamlu/gt/tool/util/result"
	sq "github.com/dreamlu/gt/tool/util/sql"
	"github.com/dreamlu/gt/tool/valid"
	"runtime"
	"strings"
)

// Mysql implement Crud
type Mysql struct {
	// DB  tool
	dbTool *DB
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
	isTrans bool
}

func (c *Mysql) Init(param *Params) {

	c.dbTool = dbTool
	c.param = param
	c.dbTool.InitColumns(c.param)
	return
}

func (c *Mysql) DB() *DB {
	c.common()
	return c.dbTool
}

func (c *Mysql) Params(params ...Param) Crud {

	for _, p := range params {
		p(c.param)
	}
	c.dbTool.InitColumns(c.param)
	return c
}

// GetBySearch
// pager info
func (c *Mysql) GetBySearch(params cmap.CMap) Crud {
	c.common()
	clone := c.clone()
	clone.pager = clone.dbTool.GetBySearch(&GT{
		Params: clone.param,
		CMaps:  params,
	})

	return clone
}

func (c *Mysql) Get(params cmap.CMap) Crud {
	c.common()
	clone := c.clone()
	clone.dbTool.Get(&GT{
		Params: clone.param,
		CMaps:  params,
	})
	return clone
}

func (c *Mysql) GetMore(params cmap.CMap) Crud {
	c.common()
	clone := c.clone()
	clone.dbTool.GetMoreData(&GT{
		Params: clone.param,
		CMaps:  params,
	})
	return clone
}

func (c *Mysql) GetByID(id interface{}) Crud {
	c.common()

	clone := c.clone()
	clone.dbTool.GetByID(&GT{
		Params: clone.param,
	}, id)
	return clone
}

// GetMoreBySearch the same as search
// more tables
func (c *Mysql) GetMoreBySearch(params cmap.CMap) Crud {
	c.common()

	clone := c.clone()
	clone.pager = clone.dbTool.GetMoreBySearch(&GT{
		CMaps:  params,
		Params: clone.param,
	})
	return clone
}

func (c *Mysql) Delete(id interface{}) Crud {
	c.common()

	clone := c.clone()
	clone.dbTool.Delete(clone.param.Table, id)
	return clone
}

// === form data ===

func (c *Mysql) UpdateForm(params cmap.CMap) error {
	c.common()

	return c.dbTool.UpdateFormData(c.param.Table, params)
}

func (c *Mysql) CreateForm(params cmap.CMap) error {
	c.common()

	return c.dbTool.CreateFormData(c.param.Table, params)
}

// == json data ==

// CreateMore can use Create replace
func (c *Mysql) CreateMore() Crud {
	c.common()
	clone := c.clone()
	if c.param.valid {
		clone.err = c.valid(clone.param.Data)
		if clone.err != nil {
			return clone
		}
	}
	clone.dbTool.CreateMore(clone.param.Table, clone.param.Model, clone.param.Data)
	return clone
}

func (c *Mysql) Update() Crud {
	c.common()
	clone := c.clone()
	if c.param.valid {
		clone.err = c.valid(clone.param.Data)
		if clone.err != nil {
			return clone
		}
	}
	clone.dbTool.Update(&GT{
		Params: clone.param,
		Select: clone.selectSQL,
		Args:   clone.args,
	})
	return clone
}

func (c *Mysql) Create() Crud {
	c.common()
	clone := c.clone()
	if c.param.valid {
		clone.err = c.valid(clone.param.Data)
		if clone.err != nil {
			return clone
		}
	}
	clone.dbTool.Create(clone.param.Table, clone.param.Data)
	return clone
}

func (c *Mysql) Select(q interface{}, args ...interface{}) Crud {

	clone := c
	if c.selectSQL == "" {
		clone = c.clone()
	}

	var query string
	switch q.(type) {
	case string:
		query = q.(string)
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
	c.dbTool.getBySQL(c.param.Data, c.selectSQL, c.args...)
	return c
}

func (c *Mysql) Exec() Crud {
	c.common()
	c.dbTool.ExecSQL(c.selectSQL, c.args...)
	return c
}

func (c *Mysql) Error() error {

	if c.err != nil {
		return c.err
	}
	if c.dbTool.res != nil {
		c.err = c.dbTool.res.Error
	}
	return c.err
}

func (c *Mysql) RowsAffected() int64 {

	if c.dbTool.res == nil {
		return 0
	}
	return c.dbTool.res.RowsAffected
}

func (c *Mysql) Pager() result.Pager {

	return c.pager
}

func (c *Mysql) Begin() Crud {
	clone := c.clone()
	clone.isTrans = true
	clone.dbTool.DB = clone.dbTool.Begin()
	defer func() {
		if r := recover(); r != nil {
			clone.dbTool.Rollback()
		}
	}()
	return clone
}

func (c *Mysql) Commit() Crud {
	if c.dbTool.res == nil || c.dbTool.res.Error != nil {
		c.dbTool.Rollback()
	}
	c.dbTool.Commit()
	c.isTrans = false
	return c
}

func (c *Mysql) Rollback() Crud {
	c.dbTool.Rollback()
	return c
}

func (c *Mysql) SavePoint(name string) Crud {
	c.dbTool.SavePoint(name)
	return c
}

func (c *Mysql) RollbackTo(name string) Crud {
	c.dbTool.RollbackTo(name)
	return c
}

func (c *Mysql) clone() (dbCrud *Mysql) {

	// default table
	if c.param.Table == "" &&
		c.param.Model != nil {
		c.param.Table = hump.HumpToLine(reflect.Name(c.param.Model))
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
	if c.isTrans {
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

func (c *Mysql) valid(data interface{}) error {

	ves := valid.Valid(data)
	if len(ves) > 0 {
		var s string
		for k, v := range ves {
			s += fmt.Sprintf("%s:%s;", k, v)
		}
		return errors.New(strings.TrimSuffix(s, ";"))
	}
	return nil
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
		_, _ = fmt.Fprintf(buf, "\n\033[35m[gt]\033[0m: ")
		_, _ = fmt.Fprintf(buf, "%s:%d\n", fullFile, line)
		fmt.Print(buf.String())
	}
}
