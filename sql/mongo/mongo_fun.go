package mongo

import (
	"context"
	"github.com/dreamlu/gt/tool/type/bmap"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/cons"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

// ===== search data =========

// CursorScan scan data to mongo data
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

func (m *Mongo) GetByDataSearch(params bmap.BMap) (cur *mongo.Cursor, err error) {
	var (
		//clientPage, everyPage int64
		filter = bson.M{}
		opt    = options.Find()
		order  string
	)

	for k, _ := range params {
		switch k {
		case cons.GtClientPage, cons.GtClientPageUnderLine:
			m.pager.ClientPage = params.Pop(k).(int64)
			continue
		case cons.GtEveryPage, cons.GtEveryPageUnderLine:
			m.pager.EveryPage = params.Pop(k).(int64)
			continue
		case cons.GtOrder:
			order = params.Pop(k).(string)
			continue
		}
	}
	if m.pager.ClientPage > 0 {
		opt.SetLimit(m.pager.EveryPage)
		opt.SetSkip((m.pager.ClientPage - 1) * m.pager.EveryPage)
	}

	c := m.m.Collection(m.param.Table)

	// filter
	_ = dataToBSON(params, &filter)
	if order != "" {
		i := 1
		os := strings.Split(order, " ")
		o := os[0]
		if len(os) > 1 && os[1] == "desc" {
			i = -1
		}
		if o == "id" {
			o = "_id"
		}
		opt.SetSort(bson.D{{o, i}})
	}

	// pager
	m.pager.TotalNum, err = c.CountDocuments(context.TODO(), filter)
	return c.Find(context.TODO(), filter, opt)
}

// ====== other =====
// data to bson value
func dataToBSON(data, value interface{}) error {

	switch data.(type) {
	case cmap.CMap:
		data = data.(cmap.CMap).BSON()
	}

	b, err := bson.Marshal(data)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(b, value)
	if err != nil {
		return err
	}
	return nil
}
