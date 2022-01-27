// package gt

/*
	gt is a fast go tool, help you dev project
*/

package gt

import (
	"github.com/dreamlu/gt/tool/log"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util"
	"github.com/dreamlu/gt/tool/util/result"
	"gorm.io/gorm"
	"strings"
)

func init() {
	log.Info("[welcome to gt (´∀｀)]")
}

// Crud interface
type Crud interface {
	// Init init crud
	Init(param *Params)
	// DB db
	DB() *DB
	// Params new/replace param
	// return param
	Params(param ...Param) Crud
	// crud method

	// GetBySearch get url params
	// like form data
	GetBySearch(params cmap.CMap) Crud     // search single table
	Get(params cmap.CMap) Crud             // get data no search
	GetMore(params cmap.CMap) Crud         // get data more table no search
	GetByID(id interface{}) Crud           // by id
	GetMoreBySearch(params cmap.CMap) Crud // more search, more tables inner/left join

	// Delete delete by id/ids/slice
	Delete(id interface{}) Crud // delete

	// Update crud and search id
	// json data
	Update() Crud     // update
	Create() Crud     // create, include res insert id
	CreateMore() Crud // create more, data must array type, single table

	Select(q interface{}, args ...interface{}) Crud // select sql
	From(query string) Crud                         // from sql, if use search, From must only once
	Group(query string) Crud                        // the last group by
	Search(params cmap.CMap) Crud                   // Select Search pager, params only support Pager and Mock
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

// Params crud params
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

	valid bool
}

type Param func(*Params)

// NewCrud new crud
func NewCrud(params ...Param) (crud Crud) {

	db()
	crud = new(Mysql)
	crud.Init(newParam(params...))
	return
}

// NewCusCrud new your custom db crud
func NewCusCrud(db *gorm.DB, log bool, params ...Param) (crud Crud) {

	newDB(db, log)
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

func Left(LeftTable ...string) Param {

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

// WhereSQL where sql and args
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

// Distinct inner/left support distinct
func Distinct(Distinct string) Param {

	return func(params *Params) {
		params.distinct = Distinct
	}
}

func Valid(valid bool) Param {

	return func(params *Params) {
		params.valid = valid
	}
}
