// package gt

/*
	gt is a fast go tool, help you dev project
*/

package gt

import (
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util"
	"strings"
)

const Version = "2.0.0+"

func init() {
	println("[gt version]:", Version)
}

// crud is db driver extend
type Crud interface {
	// init crud
	Init(param *Params)
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
	Get(params cmap.CMap) Crud             // get data no search
	GetMore(params cmap.CMap) Crud         // get data more table no search
	GetByID(id interface{}) Crud           // by id
	GetMoreBySearch(params cmap.CMap) Crud // more search

	// delete by id/ids
	Delete(id interface{}) Crud // delete

	// crud and search id
	// form data
	// [create/update] future all will use json replace form request
	// form will not update
	UpdateForm(params cmap.CMap) error // update
	CreateForm(params cmap.CMap) error // create

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
	SavePoint(name string) Crud                     // save a point
	RollbackTo(name string) Crud                    // rollback to point
}

// crud params
type Params struct {
	// attributes
	InnerTable []string    // inner join tables
	LeftTable  []string    // left join tables
	Table      string      // table name
	Model      interface{} // table model, like User{}
	KeyModel   interface{} // key like model
	Data       interface{} // table model data, like var user User{}, it is 'user', it store real data

	// sub query
	SubSQL string // SubQuery SQL
	// where
	WhereSQL string        // Where SQL
	wArgs    []interface{} // Where args

	// distinct
	distinct string
}

type Param func(*Params)

// new crud
func NewCrud(params ...Param) (crud Crud) {

	DB()
	crud = new(Mysql)
	crud.Init(newParam(params...))
	return
}

func newParam(params ...Param) *Params {
	param := &Params{}

	for _, p := range params {
		p(param)
	}
	return param
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

// Deprecated
// use WhereSQL replace
func SubWhereSQL(WhereSQL ...string) Param {

	return func(params *Params) {
		WhereSQL = util.RemoveStrings(WhereSQL, "")
		if len(WhereSQL) == 0 {
			return
		}
		params.WhereSQL = strings.Join(WhereSQL[:], " and ")
	}
}

// where sql and args, can not coexists with SubWhereSQL
func WhereSQL(WhereSQL string, args ...interface{}) Param {

	return func(params *Params) {
		if WhereSQL == "" {
			return
		}
		params.wArgs = args
		params.WhereSQL = WhereSQL
	}
}

func (p Param) WhereSQL(WhereSQL string, args ...interface{}) Param {

	return func(params *Params) {
		p(params)
		if WhereSQL == "" {
			return
		}
		if params.WhereSQL != "" {
			params.WhereSQL += " and "
		}
		params.wArgs = append(params.wArgs, args...)
		params.WhereSQL += WhereSQL
	}
}

// inner/left support distinct
func Distinct(Distinct string) Param {

	return func(params *Params) {
		params.distinct = Distinct
	}
}
