package mongo

import (
	"github.com/dreamlu/gt/lib/result"
	"github.com/dreamlu/gt/src/type/bmap"
	"github.com/dreamlu/gt/src/type/cmap"
	"go.mongodb.org/mongo-driver/mongo"
)

// Crud mongo
type Crud interface {
	// Init crud
	Init(*Params)
	// DB db
	DB() *mongo.Database
	// Params new/replace param
	// return param
	Params(...Param) Crud
	// crud method

	// FindSearch get url params
	// like form data
	FindSearch(bmap.BMap) Crud // search single table
	Find(cmap.CMap) Crud       // get data no search
	FindID(any) Crud           // by id

	// Delete delete by id/ids
	Delete(any) Crud // delete

	// Update crud and search id
	// json data
	Update() Crud     // update
	Create() Crud     // create, include res insert id
	CreateMore() Crud // create more, data must array type, single table

	Error() error        // crud error
	RowsAffected() int64 // inflect rows
	Pager() result.Pager // search pager
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
}

type Param func(*Params)

// NewCrud new crud
func NewCrud(params ...Param) (crud Crud) {

	MongoDB()
	crud = &Mongo{}
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
