package mongo

import (
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/result"
	"go.mongodb.org/mongo-driver/mongo"
)

// Crud
type Crud interface {
	// init crud
	Init(param *Params)
	// DB
	DB() *mongo.Database
	// new/replace param
	// return param
	Params(param ...Param) Crud
	// crud method

	// get url params
	// like form data
	GetBySearch(params cmap.CMap) Crud // search single table
	Get(params cmap.CMap) Crud         // get data no search
	GetByID(id interface{}) Crud       // by id

	// delete by id/ids
	Delete(id interface{}) Crud // delete

	// crud and search id
	// json data
	Update() Crud     // update
	Create() Crud     // create, include res insert id
	CreateMore() Crud // create more, data must array type, single table

	// select
	Error() error        // crud error
	RowsAffected() int64 // inflect rows
	Pager() result.Pager // search pager
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

	// count
	SubSQL string // SubQuery SQL
	// where
	WhereSQL string // SubWhere SQL
}

type Param func(*Params)

// new crud
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
