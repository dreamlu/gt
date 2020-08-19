package gt

import (
	"context"
	"github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/hump"
	"github.com/dreamlu/gt/tool/util/str"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Mongo struct {
	// crud param
	param *Params
	// mongo
	m   *mongo.Database
	err error
}

func (c *Mongo) initCrud(param *Params) {

	c.param = param
	c.m = mongoDB
	return
}

func (c *Mongo) DB() *DBTool {
	return nil
}

func (c *Mongo) AutoMigrate(values ...interface{}) Crud {
	return c
}

func (c *Mongo) Params(params ...Param) Crud {

	for _, p := range params {
		p(c.param)
	}
	return c
}

// search
// pager info
func (c *Mongo) GetBySearch(params cmap.CMap) Crud {

	return c
}

func (c *Mongo) GetByData(params cmap.CMap) Crud {
	return c
}

// by id
func (c *Mongo) GetByID(id interface{}) Crud {
	return c
}

// the same as search
// more tables
func (c *Mongo) GetMoreBySearch(params cmap.CMap) Crud {
	return c
}

// delete
func (c *Mongo) Delete(id interface{}) Crud {
	return c
}

// === form data ===

// update
func (c *Mongo) UpdateForm(params cmap.CMap) error {
	return nil
}

// create
func (c *Mongo) CreateForm(params cmap.CMap) error {
	return nil
}

// create res insert id
func (c *Mongo) CreateResID(params cmap.CMap) (str.ID, error) {
	return str.ID{}, nil
}

// == json data ==

// create more
func (c *Mongo) CreateMore() Crud {
	return c
}

// update
func (c *Mongo) Update() Crud {
	return c
}

// create
func (c *Mongo) Create() Crud {
	clone := c.clone()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	clone.m.Collection(clone.param.Table).InsertOne(ctx, clone.param.Data)
	return clone
}

// create
func (c *Mongo) Select(q interface{}, args ...interface{}) Crud {

	return c
}

func (c *Mongo) From(query string) Crud {

	return c
}

func (c *Mongo) Group(query string) Crud {

	return c
}

func (c *Mongo) Search(params cmap.CMap) Crud {
	return c
}

func (c *Mongo) Single() Crud {
	return c
}

func (c *Mongo) Exec() Crud {
	return c
}

func (c *Mongo) Error() error {

	return nil
}

func (c *Mongo) RowsAffected() int64 {

	return 0
}

func (c *Mongo) Pager() result.Pager {

	return result.Pager{}
}

func (c *Mongo) Begin() Crud {
	return c
}

func (c *Mongo) Commit() Crud {
	return c
}

func (c *Mongo) Rollback() Crud {
	return c
}

func (c *Mongo) clone() (mongo *Mongo) {

	// default table
	if c.param.Table == "" &&
		c.param.Model != nil {
		c.param.Table = hump.HumpToLine(reflect.StructToString(c.param.Model))
	}

	mongo = &Mongo{
		param: c.param,
		err:   c.err,
		m:     c.m,
	}
	return
}
