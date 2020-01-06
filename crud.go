// package gt

/*
	go-tool is a fast go tool, help you dev project
*/

package gt

import (
	"github.com/dreamlu/go-tool/tool/result"
	"strings"
)

const Version = "1.5.x"

// crud is db driver extend
type Crud interface {
	// init crud
	InitCrud(dbTool *DBTool, param *Params)
	// DB
	DB() *DBTool
	// new/replace param
	// return param
	Params(param ...Param) Crud
	// crud method

	// get url params
	// like form data
	GetBySearch(params map[string][]string) (pager result.Pager, err error)     // search
	GetByData(params map[string][]string) error                                 // get data no search
	GetByID(id string) error                                                    // by id
	GetMoreBySearch(params map[string][]string) (pager result.Pager, err error) // more search

	// common sql data
	// through sql, get the data
	GetDataBySQL(sql string, args ...interface{}) error // single data
	// page limit ?,?
	// args not include limit ?,?
	GetDataBySearchSQL(sql, sqlNt string, args ...interface{}) (pager result.Pager, err error) // more data
	DeleteBySQL(sql string, args ...interface{}) error
	UpdateBySQL(sql string, args ...interface{}) error
	CreateBySQL(sql string, args ...interface{}) error

	// delete
	Delete(id string) error // delete

	// crud and search id
	// form data
	UpdateForm(params map[string][]string) error        // update
	CreateForm(params map[string][]string) error        // create
	CreateResID(params map[string][]string) (ID, error) // create res insert id

	// crud and search id
	// json data
	Update(data interface{}) error         // update
	Create(data interface{}) error         // create, include res insert id
	CreateMoreData(data interface{}) error // create more

	// select
	Select(query string, args ...interface{}) Crud // select sql
	//Where(query string, args ...interface{}) Crud  // where condition sql
	Search() (pager result.Pager, err error) // search pager
	Single() error                           // no search
	//Where(query interface{}, args ...interface{}) Crud
}

// crud params
type Params struct {
	// attributes
	InnerTable []string    // inner join tables
	LeftTable  []string    // left join tables
	Table      string      // table name
	Model      interface{} // table model, like User{}
	ModelData  interface{} // table model data, like var user User{}, it is 'user'

	// pager info
	ClientPage int64 // page number
	EveryPage  int64 // Number of pages per page

	// count
	SubSQL string // SubQuery SQL
}

type Param func(*Params)

// new crud
func NewCrud(p ...Param) Crud {

	DBTooler()
	crud := new(DBCrud)
	crud.InitCrud(dbTool, newParam(p...))
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

func ModelData(ModelData interface{}) Param {

	return func(params *Params) {
		params.ModelData = ModelData
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
		params.SubSQL = "," + strings.Join(SubSQL[:], ",")
	}
}
