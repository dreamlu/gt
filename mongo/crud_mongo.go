package mongo

import (
	"context"
	"github.com/dreamlu/gt/crud/dep/result"
	"github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/type/cmap"
	"github.com/dreamlu/gt/src/type/tmap"
	"github.com/dreamlu/gt/src/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo implement Crud
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

func (m *Mongo) DB() *mongo.Database {

	return m.m
}

func (m *Mongo) Params(params ...Param) Crud {

	for _, p := range params {
		p(m.param)
	}
	return m
}

// FindSearch
// pager info
func (m *Mongo) FindSearch(params tmap.TMap[string, any]) Crud {
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

func (m *Mongo) Find(params cmap.CMap) Crud {
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

// FindID must be mongodb primitive.ObjectID
// by id
func (m *Mongo) FindID(id any) Crud {
	clone := m.clone()

	res := clone.m.Collection(clone.param.Table).FindOne(context.TODO(), bson.M{"_id": id.(primitive.ObjectID)})
	m.err = res.Err()
	if m.err == nil {
		m.err = res.Decode(clone.param.Data)
	}
	return clone
}

func (m *Mongo) Delete(id any) Crud {
	clone := m.clone()

	res, err := clone.m.Collection(clone.param.Table).DeleteMany(context.TODO(), bson.M{"_id": id.(primitive.ObjectID)})
	m.err = err
	if err == nil {
		m.rowANum = res.DeletedCount
	}
	return clone
}

// Update
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

func (m *Mongo) Create() Crud {
	clone := m.clone()
	_, err := clone.m.Collection(clone.param.Table).InsertOne(context.TODO(), clone.param.Data)
	m.err = err
	if err == nil {
		//m.rowANum = res.InsertedID
	}
	return clone
}

func (m *Mongo) CreateMore() Crud {
	clone := m.clone()
	_, m.err = clone.m.Collection(clone.param.Table).InsertMany(context.TODO(), reflect.ToSlice(clone.param.Data))
	return clone
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
		m.param.Table = util.HumpToLine(reflect.Name(m.param.Model))
	}

	mongo = &Mongo{
		param:     m.param,
		err:       m.err,
		m:         m.m,
		selectSQL: m.selectSQL,
	}
	return
}
