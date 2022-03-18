package gt

import (
	"encoding/json"
	"fmt"
	"github.com/dreamlu/gt/serv/log"
	mr "github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/tool/hump"
	. "github.com/dreamlu/gt/tool/tag"
	"reflect"
	"strings"
)

type Parses struct {
	Table  string // main table
	Key    string
	Tags   []string
	Vs     []any
	OTags  map[string]string
	TagTb  map[string]string // part: tag->tb
	TagTag map[string]string // part: tag->tag
}

func (r *Parses) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

func (r *Parses) Marshal(v string) {
	err := json.Unmarshal([]byte(v), r)
	if err != nil {
		log.Info("Parses Marshal error ", err)
		return
	}
	return
}

func parse(model any, tables ...string) (r *Parses) {
	var (
		typ = mr.TrueTypeof(model)
		key = mr.Path(typ, tables...)
		v   = buffer.Get(key)
	)
	r = &Parses{}
	if v != "" {
		r.Marshal(v)
		return
	}
	r.Table = hump.HumpToLine(typ.Name())
	if tables != nil {
		r.Table = tables[0]
	}
	r.Key = key
	r.OTags = make(map[string]string)
	r.TagTb = make(map[string]string)
	r.TagTag = make(map[string]string)
	parseTag(r, typ, tables...)
	v = r.String()
	buffer.Set(key, v)
	return
}

func parseTag(r *Parses, typ reflect.Type, tables ...string) {
	var (
		oTag, tag, tagTable, t string
		b                      bool
	)
	if !mr.IsStruct(typ) {
		log.Info("Parses model Not Struct")
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Anonymous {
			parseTag(r, typ.Field(i).Type, tables...)
			continue
		}
		if tag, tagTable, oTag, b = ParseTag(typ.Field(i)); b {
			continue
		}

		t = tag
		if tagTable != "" {
			tag = fmt.Sprintf("%s.%s", tagTable, tag)
			r.TagTb[tag] = tagTable
		}
		r.Tags = append(r.Tags, tag)
		r.TagTag[tag] = t

		if tag != oTag {
			r.OTags[tag] = oTag
		}

		// UniqueTagTable
		tagTable = UniqueTagTable(tag, tables...)
		if tagTable != "" {
			r.TagTb[tag] = tagTable
			continue
		}
		parseOtherTag(r, tag, t, tables...)
	}
}

func parseOtherTag(r *Parses, tag, t string, tables ...string) {
	if tables == nil {
		return
	}
	tables = tables[:len(tables)-1]
	// foreign tables column
	for _, tb := range tables {
		if strings.Contains(tag, tb+"_id") {
			break
		}
		// tables
		if strings.HasPrefix(tag, tb+"_") &&
			!strings.Contains(tag, "_id") &&
			!strings.Contains(tag, ".") {
			r.TagTag[tag] = t[len(tb)+1:]
			r.TagTb[tag] = tb
		}
	}
}

// parseV value may be different
func parseV(r *Parses, v any) {
	var (
		typ = mr.TrueValueOf(v)
	)
	if !mr.IsStruct(typ) {
		log.Info("Parses model Not Struct")
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		v = typ.Field(i).Interface()
		if typ.Type().Field(i).Anonymous {
			parseV(r, v)
			continue
		}
		r.Vs = append(r.Vs, v)
	}
}
