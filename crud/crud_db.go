// package gt

package crud

import (
	"errors"
	"fmt"
	"github.com/dreamlu/gt/crud/dep/result"
	"github.com/dreamlu/gt/crud/dep/tag"
	"github.com/dreamlu/gt/src/type/cmap"
	"github.com/dreamlu/gt/src/valid"
	"runtime"
	"strings"
)

// CrudGo implement Crud
type CrudGo struct {
	// DB  tool
	dbTool *DB
	// error
	err error
	// crud param
	param *Params

	// select
	selectSQL string // select/or if
	from      string // from sql
	args      []any  // select args
	argsNt    []any  // select nt args, related from
	group     string // the last group
	// pager
	pager result.Pager

	// transaction
	isTrans bool

	// count
	isCount bool
}

func (c *CrudGo) Init(param *Params) {

	c.dbTool = dbTool
	c.param = param
	c.dbTool.InitColumns(c.param)
	return
}

func (c *CrudGo) DB() *DB {
	c.common()
	return c.dbTool
}

func (c *CrudGo) Params(params ...Param) Crud {

	for _, p := range params {
		p(c.param)
	}
	c.dbTool.InitColumns(c.param)
	return c
}

func (c *CrudGo) Count() Crud {
	c.common()
	clone := c.clone()
	clone.isCount = true
	return clone
}

func (c *CrudGo) Find(p cmap.CMap) Crud {
	c.common()
	clone := c.clone()
	clone.pager = clone.dbTool.Find(&GT{
		Params:  clone.param,
		CMaps:   p,
		isCount: c.isCount,
	})
	return clone
}

func (c *CrudGo) FindM(params cmap.CMap) Crud {
	c.common()
	clone := c.clone()
	clone.pager = clone.dbTool.FindM(&GT{
		Params:  clone.param,
		CMaps:   params,
		isCount: c.isCount,
	})
	return clone
}

func (c *CrudGo) Delete(conds ...any) Crud {
	c.common()

	clone := c.clone()
	clone.dbTool.Delete(&GT{
		Params: clone.param,
	}, conds)
	return clone
}

func (c *CrudGo) Update() Crud {
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

func (c *CrudGo) Create() Crud {
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

func (c *CrudGo) Select(q any, args ...any) Crud {

	clone := c
	if c.selectSQL == "" {
		clone = c.clone()
	}

	var query string
	switch q.(type) {
	case string:
		query = q.(string)
	case any:
		query, args = StructWhereSQL(q)
	}

	clone.selectSQL += query + " "
	clone.args = append(clone.args, args...)
	if clone.from != "" {
		clone.argsNt = append(clone.argsNt, args...)
	}
	return clone
}

func (c *CrudGo) From(query string) Crud {

	c.from = query
	c.selectSQL += query + " "
	return c
}

func (c *CrudGo) Group(query string) Crud {

	c.group = query
	return c
}

func (c *CrudGo) FindS(params cmap.CMap) Crud {
	c.common()

	if c.argsNt == nil {
		c.argsNt = c.args
	}
	c.pager = c.dbTool.FindS(&GT{
		Params:  c.param,
		Select:  c.selectSQL,
		Args:    c.args,
		From:    c.from,
		Group:   c.group,
		isCount: c.isCount,
		CMaps:   params,
	})
	return c
}

func (c *CrudGo) Scan() Crud {
	c.common()
	c.Select(c.group)
	c.dbTool.get(&GT{
		sql:    c.selectSQL,
		Params: c.param,
		Args:   c.args,
	})
	return c
}

func (c *CrudGo) Exec() Crud {
	c.common()
	c.dbTool.exec(c.selectSQL, c.args...)
	return c
}

func (c *CrudGo) Error() error {

	if c.err != nil {
		return c.err
	}
	if c.dbTool.res != nil {
		c.err = c.dbTool.res.Error
	}
	return c.err
}

func (c *CrudGo) RowsAffected() int64 {

	if c.dbTool.res == nil {
		return 0
	}
	return c.dbTool.res.RowsAffected
}

func (c *CrudGo) Pager() result.Pager {

	return c.pager
}

func (c *CrudGo) Begin() Crud {
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

func (c *CrudGo) Commit() Crud {
	if c.dbTool.res == nil || c.dbTool.res.Error != nil {
		c.dbTool.Rollback()
	}
	c.dbTool.Commit()
	c.isTrans = false
	return c
}

func (c *CrudGo) Rollback() Crud {
	c.dbTool.Rollback()
	return c
}

func (c *CrudGo) SavePoint(name string) Crud {
	c.dbTool.SavePoint(name)
	return c
}

func (c *CrudGo) RollbackTo(name string) Crud {
	c.dbTool.RollbackTo(name)
	return c
}

func (c *CrudGo) clone() (dbCrud *CrudGo) {

	// default table
	if c.param.Table == "" &&
		c.param.Model != nil {
		c.param.Table = tag.ModelTable(c.param.Model)
	}

	dbCrud = &CrudGo{
		dbTool:    c.dbTool,
		param:     c.param,
		selectSQL: c.selectSQL,
		from:      c.from,
		args:      c.args,
		argsNt:    c.argsNt,
		group:     c.group,
		err:       c.err,
		isCount:   c.isCount,
	}

	// isTrans
	if c.isTrans {
		return
	}
	dbCrud.dbTool = c.dbTool.clone()
	return
}

func (c *CrudGo) common() {
	if c.dbTool.log {
		c.line()
	}
}

func (c *CrudGo) valid(data any) error {

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

func (c *CrudGo) line() {
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
