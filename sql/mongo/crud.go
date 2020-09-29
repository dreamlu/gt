package mongo

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
	SubWhereSQL string // SubWhere SQL
}

type Param func(*Params)

// new crud
func NewCrud(params ...Param) (crud *Mongo) {

	MongoDB()
	crud = &Mongo{}
	crud.initCrud(newParam(params...))
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
