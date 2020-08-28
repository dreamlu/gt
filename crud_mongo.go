package gt

import (
	"context"
	"github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/result"
	sq "github.com/dreamlu/gt/tool/sql"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/hump"
	"github.com/dreamlu/gt/tool/util/str"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// implement Crud
type Mongo struct {
	// crud param
	param *Params
	// mongo
	m       *mongo.Database
	rowANum int64
	// err
	err error

	// select
	selectSQL string // select/or if

	// pager
	pager result.Pager
}

func (m *Mongo) initCrud(param *Params) {

	m.param = param
	m.m = mongoDB
	return
}

func (m *Mongo) DB() *DBTool {
	return nil
}

func (m *Mongo) Params(params ...Param) Crud {

	for _, p := range params {
		p(m.param)
	}
	return m
}

// search
// pager info
func (m *Mongo) GetBySearch(params cmap.CMap) Crud {
	clone := m.clone()

	cur, err := m.GetByDataSearch(params)
	if err != nil {
		clone.err = err
		return clone
	}
	m.CursorScan(cur, clone.param.Data)
	m.err = cur.Close(context.TODO())
	return clone
}

func (m *Mongo) GetByData(params cmap.CMap) Crud {
	clone := m.clone()

	filter := bson.M{}
	m.err = params.Struct(&filter)
	cur, err := clone.m.Collection(clone.param.Table).Find(context.TODO(), filter)
	if err != nil {
		clone.err = err
		return clone
	}
	m.CursorScan(cur, clone.param.Data)
	m.err = cur.Close(context.TODO())
	return clone
}

func (m *Mongo) GetMoreByData(params cmap.CMap) Crud {
	return m
}

// must be mongodb primitive.ObjectID
// by id
func (m *Mongo) GetByID(id interface{}) Crud {
	clone := m.clone()

	res := clone.m.Collection(clone.param.Table).FindOne(context.TODO(), bson.M{"_id": id.(primitive.ObjectID)})
	m.err = res.Err()
	if m.err == nil {
		m.err = res.Decode(clone.param.Data)
	}
	return clone
}

// the same as search
// more tables
func (m *Mongo) GetMoreBySearch(params cmap.CMap) Crud {
	return m
}

// delete
func (m *Mongo) Delete(id interface{}) Crud {
	clone := m.clone()

	res, err := clone.m.Collection(clone.param.Table).DeleteMany(context.TODO(), bson.M{"_id": id.(primitive.ObjectID)})
	m.err = err
	if err == nil {
		m.rowANum = res.DeletedCount
	}
	return clone
}

// === form data ===

// update
func (m *Mongo) UpdateForm(params cmap.CMap) error {
	return nil
}

// create
func (m *Mongo) CreateForm(params cmap.CMap) error {
	return nil
}

// create res insert id
func (m *Mongo) CreateResID(params cmap.CMap) (str.ID, error) {
	return str.ID{}, nil
}

// == json data ==

// update
// must id string
func (m *Mongo) Update() Crud {
	clone := m.clone()

	data := bson.M{}
	_ = dataToBSON(m.param.Data, &data)
	_id, _ := primitive.ObjectIDFromHex(data["_id"].(string))
	delete(data, "_id")
	res, err := clone.m.Collection(clone.param.Table).
		UpdateOne(
			context.TODO(),
			bson.M{"_id": _id},
			bson.D{{"$set", data}},
		)
	m.err = err
	if err == nil {
		m.rowANum = res.UpsertedCount
	}
	return clone
}

// create
func (m *Mongo) Create() Crud {
	clone := m.clone()
	_, err := clone.m.Collection(clone.param.Table).InsertOne(context.TODO(), clone.param.Data)
	m.err = err
	if err == nil {
		//m.rowANum = res.InsertedID
	}
	return clone
}

// create more
func (m *Mongo) CreateMore() Crud {
	clone := m.clone()
	_, m.err = clone.m.Collection(clone.param.Table).InsertMany(context.TODO(), reflect.ToSlice(clone.param.Data))
	return clone
}

// create
func (m *Mongo) Select(q interface{}, args ...interface{}) Crud {

	clone := m
	if m.selectSQL == "" {
		clone = m.clone()
	}

	var query string
	switch q.(type) {
	case string:
		query = q.(string)
	//case cmap.CMap:
	//	query, args = sq.CMapWhereSQL(q.(cmap.CMap))
	case interface{}:
		query, args = sq.StructWhereSQL(q)
	}

	clone.selectSQL += query + " "
	//clone.args = append(clone.args, args...)
	return clone
}

func (m *Mongo) From(query string) Crud {

	return m
}

func (m *Mongo) Group(query string) Crud {

	return m
}

func (m *Mongo) Search(params cmap.CMap) Crud {
	return m
}

func (m *Mongo) Single() Crud {
	//m.m.Client().
	//m.err = m.m.RunCommand(context.TODO(), m.selectSQL).Decode(&m.param.Data)
	return m
}

func (m *Mongo) Exec() Crud {
	return m
}

func (m *Mongo) Error() error {

	return m.err
}

func (m *Mongo) RowsAffected() int64 {

	return m.rowANum
}

func (m *Mongo) Pager() result.Pager {

	return m.pager
}

func (m *Mongo) Begin() Crud {
	return m
}

func (m *Mongo) Commit() Crud {
	return m
}

func (m *Mongo) Rollback() Crud {
	return m
}

// TODO print filter.. params
func (m *Mongo) clone() (mongo *Mongo) {

	// default table
	if m.param.Table == "" &&
		m.param.Model != nil {
		m.param.Table = hump.HumpToLine(reflect.StructToString(m.param.Model))
	}

	mongo = &Mongo{
		param:     m.param,
		err:       m.err,
		m:         m.m,
		selectSQL: m.selectSQL,
	}
	return
}
