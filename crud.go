// package gt

/*
	gt is a fast go tool, help you dev project
*/

package gt

import (
	"github.com/dreamlu/gt/lib"
	"github.com/dreamlu/gt/lib/result"
	"github.com/dreamlu/gt/src/type/cmap"
	"github.com/dreamlu/gt/third/log"
	"gorm.io/gorm"
	"strings"
)

func init() {
	log.Info("[welcome to gt (´∀｀)]")
}

// Crud interface
type Crud interface {
	// Init init crud
	Init(*Params)
	// DB db
	DB() *DB
	// Params new/replace param
	// return param
	Params(...Param) Crud
	// crud method

	Count() Crud          // count
	Find(cmap.CMap) Crud  // find data
	FindM(cmap.CMap) Crud // find data more table no search

	// Delete delete by id/ids/slice
	Delete(any) Crud // delete

	// Update crud and search id
	// json data
	Update() Crud // update
	Create() Crud // create (more), include res insert id

	Select(any, ...any) Crud // select sql
	From(string) Crud        // from sql, if use search, From must only once
	Group(string) Crud       // the last group by
	FindS(cmap.CMap) Crud    // Select origin sql find, params only support Pager and Mock
	Scan() Crud              // no search
	Exec() Crud              // exec insert/update/delete sql
	Error() error            // crud error
	RowsAffected() int64     // inflect rows
	Pager() result.Pager     // search pager
	Begin() Crud             // start a transaction
	Commit() Crud            // commit a transaction
	Rollback() Crud          // rollback a transaction
	SavePoint(string) Crud   // save a point
	RollbackTo(string) Crud  // rollback to point
}

// Params crud params
type Params struct {
	// attributes
	InnerTable []string // inner join tables
	LeftTable  []string // left join tables
	Table      string   // table name
	Model      any      // table model, like User{}
	KeyModel   any      // key like model
	Data       any      // table model data, like var user User{}, it is 'user', it store real data

	// sub query
	SubSQL string // SubQuery SQL
	// where
	WhereSQL string // Where SQL
	wArgs    []any  // Where args

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

	cusdb(db, log)
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

func Model(Model any) Param {

	return func(params *Params) {
		params.Model = Model
	}
}

func KeyModel(KeyModel any) Param {

	return func(params *Params) {
		params.KeyModel = KeyModel
	}
}

func Data(Data any) Param {

	return func(params *Params) {
		params.Data = Data
	}
}

func SubSQL(SubSQL ...string) Param {

	return func(params *Params) {
		SubSQL = lib.Remove(SubSQL, "")
		if len(SubSQL) == 0 {
			return
		}
		params.SubSQL = "," + strings.Join(SubSQL[:], ",")
	}
}

// WhereSQL where sql and args
func WhereSQL(WhereSQL string, args ...any) Param {

	return func(params *Params) {
		if WhereSQL == "" {
			return
		}
		params.wArgs = args
		params.WhereSQL = WhereSQL
	}
}

func (p Param) WhereSQL(WhereSQL string, args ...any) Param {

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
