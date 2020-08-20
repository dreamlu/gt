package gt

import (
	"context"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/str"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strconv"
)

// ===== search data =========

// scan data to mongo data
func (m *Mongo) CursorScan(cur *mongo.Cursor, data interface{}) {
	typ := reflect.TypeOf(data)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	switch typ.Kind() {
	case reflect.Slice:
		m.err = cur.All(context.TODO(), data)
	default:
		cur.Next(context.TODO())
		m.err = cur.Decode(data)
	}
}

func (m *Mongo) GetByDataSearch(params cmap.CMap) (cur *mongo.Cursor, err error) {
	var (
		//clientPage, everyPage int64
		filter = bson.M{}
		opt    = options.Find()
	)

	for k, v := range params {
		switch k {
		case str.GtClientPage, str.GtClientPageUnderLine:
			m.pager.ClientPage, _ = strconv.ParseInt(v[0], 10, 64)
			params.Del(k)
			continue
		case str.GtEveryPage, str.GtEveryPageUnderLine:
			m.pager.EveryPage, _ = strconv.ParseInt(v[0], 10, 64)
			params.Del(k)
			continue
		}
	}
	if m.pager.ClientPage > 0 {
		opt.SetLimit(m.pager.EveryPage)
		opt.SetSkip((m.pager.ClientPage - 1) * m.pager.EveryPage)
	}

	c := m.m.Collection(m.param.Table)

	// filter
	params.Struct(&filter)

	// pager
	m.pager.TotalNum, err = c.CountDocuments(context.TODO(), filter)
	return c.Find(context.TODO(), filter, options.Find(), opt)
}

// ====== delete =====
func (m *Mongo) DeleteData() {

}
