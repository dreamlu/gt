// package gt

/*
	gt is a fast go tool, help you dev project
*/

package gt

import (
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/util"
	"github.com/dreamlu/gt/tool/util/str"
	"strings"
)

const Version = "1.6.x"

// crud is db driver extend
type Crud interface {
	// init crud
	initCrud(dbTool *DBTool, param *Params)
	// DB
	DB() *DBTool
	// new/replace param
	// return param
	Params(param ...Param) Crud
	// crud method

	// get url params
	// like form data
	GetBySearch(params map[string][]string) Crud     // search
	GetByData(params map[string][]string) Crud       // get data no search
	GetByID(id interface{}) Crud                     // by id
	GetMoreBySearch(params map[string][]string) Crud // more search

	// delete by id
	Delete(id interface{}) Crud // delete

	// crud and search id
	// form data
	// [create/update] future all will use json replace form request
	// form will not update
	UpdateForm(params map[string][]string) error            // update
	CreateForm(params map[string][]string) error            // create
	CreateResID(params map[string][]string) (str.ID, error) // create res insert id

	// crud and search id
	// json data
	Update() Crud         // update
	Create() Crud         // create, include res insert id
	CreateMoreData() Crud // create more, data must array type, single table

	// select
	Select(query string, args ...interface{}) Crud // select sql
	From(query string) Crud                        // from sql, if use search, From must only once
	Group(query string) Crud                       // the last group by
	Search() Crud                                  // search pager
	Single() Crud                                  // no search
	Exec() Crud                                    // exec insert/update/delete select sql
	Error() error                                  // crud error
	RowsAffected() int64                           // inflect rows
	Pager() result.Pager                           // search pager
	Begin() Crud                                   // start a transaction
	Commit() Crud                                  // commit a transaction
	Rollback() Crud                                // rollback a transaction
}

// crud params
type Params struct {
	// attributes
	InnerTable []string    // inner join tables
	LeftTable  []string    // left join tables
	Table      string      // table name
	Model      interface{} // table model, like User{}
	Data       interface{} // table model data, like var user User{}, it is 'user', it store real data

	// pager info
	ClientPage int64 // page number
	EveryPage  int64 // Number of pages per page

	// count
	SubSQL string // SubQuery SQL
	// where
	SubWhereSQL string // SubWhere SQL
}

type Param func(*Params)

// new crud
func NewCrud(params ...Param) Crud {

	DBTooler()
	crud := new(DBCrud)
	crud.initCrud(dbTool, newParam(params...))
	return crud
}

func newParam(params ...Param) *Params {
	param := &Params{}

	for _, p := range params {
		p(param)
	}
	return param
}

func InnerTable(InnerTables []string) Param {

	return func(params *Params) {
		params.InnerTable = InnerTables
	}
}

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

func Data(Data interface{}) Param {

	return func(params *Params) {
		params.Data = Data
	}
}

func ClientPage(ClientPage int64) Param {

	return func(params *Params) {
		params.ClientPage = ClientPage
	}
}

func EveryPage(EveryPage int64) Param {

	return func(params *Params) {
		params.EveryPage = EveryPage
	}
}

func SubSQL(SubSQL ...string) Param {

	return func(params *Params) {
		SubSQL = util.RemoveStrings(SubSQL, "")
		params.SubSQL = "," + strings.Join(SubSQL[:], ",")
	}
}

func SubWhereSQL(SubWhereSQL ...string) Param {

	return func(params *Params) {
		SubWhereSQL = util.RemoveStrings(SubWhereSQL, "")
		params.SubWhereSQL = strings.Join(SubWhereSQL[:], " and ")
	}
}
