// package gt

/*
	gt is a fast go tool, help you dev project
*/

package gt

import (
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util"
	"github.com/dreamlu/gt/tool/util/str"
	"strings"
)

const Version = "1.20.0+"

func init() {
	println("[gt version]:", Version)
}

// crud is db driver extend
type Crud interface {
	// init crud
	initCrud(param *Params)
	// DB
	// Deprecated, use gt.DB() replace
	DB() *DBTool
	// new/replace param
	// return param
	Params(param ...Param) Crud
	// crud method

	// get url params
	// like form data
	GetBySearch(params cmap.CMap) Crud     // search
	GetByData(params cmap.CMap) Crud       // get data no search
	GetMoreByData(params cmap.CMap) Crud   // get data more table no search
	GetByID(id interface{}) Crud           // by id
	GetMoreBySearch(params cmap.CMap) Crud // more search

	// delete by id
	Delete(id interface{}) Crud // delete

	// crud and search id
	// form data
	// [create/update] future all will use json replace form request
	// form will not update
	UpdateForm(params cmap.CMap) error            // update
	CreateForm(params cmap.CMap) error            // create
	CreateResID(params cmap.CMap) (str.ID, error) // create res insert id

	// crud and search id
	// json data
	Update() Crud     // update
	Create() Crud     // create, include res insert id
	CreateMore() Crud // create more, data must array type, single table

	// select
	Select(q interface{}, args ...interface{}) Crud // select sql
	From(query string) Crud                         // from sql, if use search, From must only once
	Group(query string) Crud                        // the last group by
	Search(params cmap.CMap) Crud                   // search pager
	Single() Crud                                   // no search
	Exec() Crud                                     // exec insert/update/delete sql
	Error() error                                   // crud error
	RowsAffected() int64                            // inflect rows
	Pager() result.Pager                            // search pager
	Begin() Crud                                    // start a transaction
	Commit() Crud                                   // commit a transaction
	Rollback() Crud                                 // rollback a transaction
}

// crud params
type Params struct {
	// db type
	// default ""(mysql)
	D string

	// attributes
	InnerTable []string    // inner join tables
	LeftTable  []string    // left join tables
	Table      string      // table name
	Model      interface{} // table model, like User{}
	KeyModel   interface{} // key like model
	Data       interface{} // table model data, like var user User{}, it is 'user', it store real data

	// count
	SubSQL string // SubQuery SQL
	// where
	SubWhereSQL string // SubWhere SQL
}

type Param func(*Params)

// new crud
func NewCrud(params ...Param) (crud Crud) {

	p := newParam(params...)
	switch p.D {
	case "mongo":
		MongoDB()
		crud = new(Mongo)
		crud.initCrud(p)
	default:
		DB()
		crud = new(Mysql)
		crud.initCrud(p)
	}
	return
}

func newParam(params ...Param) *Params {
	param := &Params{}

	for _, p := range params {
		p(param)
	}
	return param
}

func D(d string) Param {

	return func(params *Params) {
		params.D = d
	}
}

func Inner(InnerTables ...string) Param {

	return func(params *Params) {
		params.InnerTable = InnerTables
	}
}

// Deprecated
func InnerTable(InnerTables []string) Param {

	return func(params *Params) {
		params.InnerTable = InnerTables
	}
}

func Left(LeftTable ...string) Param {

	return func(params *Params) {
		params.LeftTable = LeftTable
	}
}

// Deprecated
func LeftTable(LeftTable []string) Param {

	return func(params *Params) {
		params.LeftTable = LeftTable
	}
}

func Table(Table string) Param {

	return func(params *Params) {
		params.Table = Table
	}
}

func Model(Model interface{}) Param {

	return func(params *Params) {
		params.Model = Model
	}
}

func KeyModel(KeyModel interface{}) Param {

	return func(params *Params) {
		params.KeyModel = KeyModel
	}
}

func Data(Data interface{}) Param {

	return func(params *Params) {
		params.Data = Data
	}
}

func SubSQL(SubSQL ...string) Param {

	return func(params *Params) {
		SubSQL = util.RemoveStrings(SubSQL, "")
		if len(SubSQL) == 0 {
			return
		}
		params.SubSQL = "," + strings.Join(SubSQL[:], ",")
	}
}

func SubWhereSQL(SubWhereSQL ...string) Param {

	return func(params *Params) {
		SubWhereSQL = util.RemoveStrings(SubWhereSQL, "")
		if len(SubWhereSQL) == 0 {
			return
		}
		params.SubWhereSQL = strings.Join(SubWhereSQL[:], " and ")
	}
}
