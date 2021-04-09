package mongo

import (
	"context"
	"github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/hump"
	"github.com/dreamlu/gt/tool/util/result"
	sq "github.com/dreamlu/gt/tool/util/sql"
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

func (m *Mongo) Init(param *Params) {

	m.param = param
	m.m = mongoDB
	return
}

func (m *Mongo) Params(params ...Param) *Mongo {

	for _, p := range params {
		p(m.param)
	}
	return m
}

// search
// pager info
func (m *Mongo) GetBySearch(params cmap.CMap) *Mongo {
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

func (m *Mongo) GetByData(params cmap.CMap) *Mongo {
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

// must be mongodb primitive.ObjectID
// by id
func (m *Mongo) GetByID(id interface{}) *Mongo {
	clone := m.clone()

	res := clone.m.Collection(clone.param.Table).FindOne(context.TODO(), bson.M{"_id": id.(primitive.ObjectID)})
	m.err = res.Err()
	if m.err == nil {
		m.err = res.Decode(clone.param.Data)
	}
	return clone
}

// delete
func (m *Mongo) Delete(id interface{}) *Mongo {
	clone := m.clone()

	res, err := clone.m.Collection(clone.param.Table).DeleteMany(context.TODO(), bson.M{"_id": id.(primitive.ObjectID)})
	m.err = err
	if err == nil {
		m.rowANum = res.DeletedCount
	}
	return clone
}

// update
// must id string
func (m *Mongo) Update() *Mongo {
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
func (m *Mongo) Create() *Mongo {
	clone := m.clone()
	_, err := clone.m.Collection(clone.param.Table).InsertOne(context.TODO(), clone.param.Data)
	m.err = err
	if err == nil {
		//m.rowANum = res.InsertedID
	}
	return clone
}

// create more
func (m *Mongo) CreateMore() *Mongo {
	clone := m.clone()
	_, m.err = clone.m.Collection(clone.param.Table).InsertMany(context.TODO(), reflect.ToSlice(clone.param.Data))
	return clone
}

// create
func (m *Mongo) Select(q interface{}, args ...interface{}) *Mongo {

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

func (m *Mongo) Single() *Mongo {
	//m.m.Client().
	//m.err = m.m.RunCommand(context.TODO(), m.selectSQL).Decode(&m.param.Data)
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

// TODO print filter.. params
func (m *Mongo) clone() (mongo *Mongo) {

	// default table
	if m.param.Table == "" &&
		m.param.Model != nil {
		m.param.Table = hump.HumpToLine(reflect.StructName(m.param.Model))
	}

	mongo = &Mongo{
		param:     m.param,
		err:       m.err,
		m:         m.m,
		selectSQL: m.selectSQL,
	}
	return
}
