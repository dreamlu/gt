package crud

import (
	"bytes"
	"fmt"
	cons2 "github.com/dreamlu/gt/crud/dep/cons"
	"github.com/dreamlu/gt/src/type/amap"
	"strings"
)

// sql tag reflect
// resolve go struct field from model

// sql memory
var buffer = amap.NewAMap()

// GetMoreColSQL select * replace
// select more tables : table name / table alias name
// first table must main table, like from a inner join b, first table is a
func GetMoreColSQL(model any, tables ...string) (sql string) {
	var (
		p   = parse(model, tables...)
		key = fmt.Sprintf("%s%s_more", cons2.SQL_, p.Key)
	)
	sql = buffer.Get(key)
	if sql != "" {
		return
	}
	var buf bytes.Buffer
	getTagMoreSQL(p, &buf)
	sql = string(buf.Bytes()[:buf.Len()-1])
	buffer.Set(key, sql)
	return
}

// get more table tag sql
func getTagMoreSQL(p *Parses, buf *bytes.Buffer) {

	for _, tag := range p.Tags {
		tb := p.TagTb[tag]
		if tb == "" {
			tb = p.Table
		}
		buf.WriteByte(cons2.Backticks)
		buf.WriteString(tb)
		buf.WriteByte(cons2.Backticks)
		buf.WriteByte('.')
		buf.WriteByte(cons2.Backticks)
		buf.WriteString(p.TagTag[tag])
		buf.WriteByte(cons2.Backticks)
		if p.OTags[tag] != "" {
			buf.WriteString(" as ")
			buf.WriteString(p.OTags[tag])
		}
		buf.WriteByte(',')
	}
}

// GetColSQL select * replace
func GetColSQL(model any) (sql string) {
	var (
		r   = parse(model)
		key = fmt.Sprintf("%s%s", cons2.SQL_, r.Key)
	)
	sql = buffer.Get(key)
	if sql != "" {
		return
	}
	var buf bytes.Buffer
	getTagSQL(r, &buf, "")
	sql = string(buf.Bytes()[:buf.Len()-1]) // remove ,
	buffer.Set(key, sql)
	return
}

// GetColSQLAlias select * replace
// add alias
func GetColSQLAlias(model any, alias string) (sql string) {
	var (
		r   = parse(model)
		key = fmt.Sprintf("%s%s_%s", cons2.SQL_, r.Key, alias)
	)
	sql = buffer.Get(key)
	if sql != "" {
		return
	}
	var buf bytes.Buffer
	getTagSQL(r, &buf, alias)
	sql = string(buf.Bytes()[:buf.Len()-1])
	buffer.Set(key, sql)
	return
}

// get tag sql
func getTagSQL(p *Parses, buf *bytes.Buffer, alias string) {
	for _, tag := range p.Tags {
		if alias != "" {
			buf.WriteString(alias)
			buf.WriteByte('.')
		}
		buf.WriteByte(cons2.Backticks)
		buf.WriteString(tag)
		buf.WriteByte(cons2.Backticks)
		buf.WriteByte(',')
	}
}

// GetColParamSQL get col ?
func GetColParamSQL(p *Parses) (sql string) {
	var (
		buf bytes.Buffer
	)
	for range p.Tags {
		buf.WriteString("?,")
	}
	return buf.String()[:buf.Len()-1]
}

func innerLeftSQL(bufNt *bytes.Buffer, DBS map[string]string, tables, fields []string, i int) {
	if tb := DBS[tables[i]]; tb != "" {
		bufNt.WriteByte(cons2.Backticks)
		bufNt.WriteString(tb)
		bufNt.WriteByte(cons2.Backticks)
		bufNt.WriteByte('.')
	}
	bufNt.WriteByte(cons2.Backticks)
	bufNt.WriteString(tables[i])
	bufNt.WriteByte(cons2.Backticks)
	bufNt.WriteString(" on ")
	fieldSQL(bufNt, tables[i-1], tables[i], fields[i-1], fields[i])
}

// fieldSQL field analyze
func fieldSQL(bufNt *bytes.Buffer, leftTable, rightTable, left, right string) {

	var (
		ils  = strings.Split(left, ",")  // left condition column
		irs  = strings.Split(right, ",") // right condition column
		ilts []string                    // left table field condition
		irts []string                    // right table field condition
	)

	for k := 0; k < len(ils); k++ {
		is := strings.Split(ils[k], "=")
		if len(is) > 1 {
			ilts = append(ilts, ils[k])
			ils = append(ils[:k], ils[k+1:]...)
			k--
		}
	}

	for k := 0; k < len(irs); k++ {
		is := strings.Split(irs[k], "=")
		if len(is) > 1 {
			irts = append(irts, irs[k])
			irs = append(irs[:k], irs[k+1:]...)
			k--
		}
	}

	for k := 0; k < len(ils); k++ {
		writeField(bufNt, leftTable, ils[k], " = ")
		writeField(bufNt, rightTable, irs[k], cons2.And)
	}

	for _, v := range ilts {
		is := strings.Split(v, "=")
		if len(is) > 1 {
			writeField(bufNt, leftTable, is[0], " = ")
			bufNt.WriteString(is[1])
			bufNt.WriteString(cons2.And)
		}
	}
	for _, v := range irts {
		is := strings.Split(v, "=")
		if len(is) > 1 {
			writeField(bufNt, rightTable, is[0], " = ")
			bufNt.WriteString(is[1])
			bufNt.WriteString(cons2.And)
		}
	}
	nBuf := bytes.NewBuffer(bufNt.Bytes()[:bufNt.Len()-4])
	defer nBuf.Reset()
	bufNt.Reset()
	bufNt.Write(nBuf.Bytes())
}

func writeField(bufNt *bytes.Buffer, tb, v, c string) {
	bufNt.WriteByte(cons2.Backticks)
	bufNt.WriteString(tb)
	bufNt.WriteByte(cons2.Backticks)
	bufNt.WriteByte('.')
	bufNt.WriteByte(cons2.Backticks)
	bufNt.WriteString(v)
	bufNt.WriteByte(cons2.Backticks)
	bufNt.WriteString(c)
}
