// package gt

package crud

import (
	"bytes"
	"fmt"
	cons2 "github.com/dreamlu/gt/crud/dep/cons"
	"github.com/dreamlu/gt/crud/dep/tag"
	mr "github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/type/cmap"
	"github.com/dreamlu/gt/src/util"
	"github.com/dreamlu/gt/third/mock"
	"reflect"
	"strconv"
	"strings"
)

// GT SQL struct
type GT struct {
	*Params
	// CMap
	CMaps cmap.CMap // params

	// select sql
	tableT     string // `table`
	Select     string // select sql
	From       string // only once
	Group      string // the last group
	Args       []any
	sql        string
	sqlNt      string
	clientPage int64
	everyPage  int64
	order      string // order by
	sqlW       string // where sql
	sqlSoft    string // soft_del sql

	// mock data
	isMock bool

	// gt parses
	parses tag.Parse

	// count
	isCount bool
}

func (gt *GT) parse() *GT {
	var (
		typ = mr.TrueTypeof(gt.Model)
		key = mr.Path(typ, cons2.GT)
		v   = buffer.Get(key)
	)
	if v != "" {
		gt.parses = tag.Parse{}
		gt.parses.Marshal(v)
		return gt
	}
	gt.parses = tag.ParseGt(gt.Model)
	v = gt.parses.String()
	buffer.Set(key, v)
	return gt
}

func (gt *GT) common() {
	for table, soffDelField := range gt.parses.Sd {
		gt.sqlSoft += fmt.Sprintf(cons2.SoftDel, table, soffDelField)
	}
	if gt.sqlSoft != "" {
		gt.sqlSoft = gt.sqlSoft[:len(gt.sqlSoft)-5]
	}

	gt.tableT = ParseTable(gt.Table)
}

//=======================================sql script==========================================
//===========================================================================================

// GetMoreSQL more table
// params: innerTables is inner join tables, must even number
// params: leftTables is left join tables
// return: select sql
// table1 as main table, include other tables_id(foreign key)
func (gt *GT) GetMoreSQL() {
	gt.parse().common()

	var (
		tables = gt.moreSql()
		sql    = GetMoreColSQL(gt.Model, tables...)
		bufW   bytes.Buffer // gt.sql bytes connect
		count  = cons2.Count
	)
	if gt.distinct != "" {
		count = fmt.Sprintf(cons2.CountDistinct, gt.distinct)
		sql = cons2.Distinct + sql
	}
	gt.sql = strings.Replace(gt.sqlNt, count, sql+gt.SubSQL, 1)
	// default
	gt.order = fmt.Sprintf(cons2.OrderDesc, ParseTable(tables[0]))

	gt.whereParams()
	for k, v := range gt.CMaps {
		if k == cons2.GtKey {
			if gt.KeyModel == nil {
				gt.KeyModel = gt.Model
			}
			// more tables key search
			sqlKey, argsKey := GetMoreKeySQL(v[0], gt.KeyModel, tables...)
			bufW.WriteString(sqlKey)
			gt.Args = append(gt.Args, argsKey...)
			continue
		}

		tbs := TagTables(k, tables...)
		if len(tbs) > 0 {
			// unique or first table where condition
			gt.whereTbKv(&bufW, tbs[0], k, v[0])
			continue
		}

		if !gt.otherTableWhere(&bufW, tables[1:], k, v[0]) {
			gt.whereTbKv(&bufW, tables[0], k, v[0])
		}
	}

	gt.whereAllSQL(&bufW)
	return
}

// GetSQL get single sql
func (gt *GT) GetSQL() {
	gt.parse().common()

	var (
		bufW bytes.Buffer // where sql, sqlNt bytes sql
	)
	// select* replace
	gt.sql = fmt.Sprintf(cons2.SelectFrom, GetColSQL(gt.Model)+gt.SubSQL, gt.tableT)

	gt.whereParams()
	for k, v := range gt.CMaps {
		if k == cons2.GtKey {
			if gt.KeyModel == nil {
				gt.KeyModel = gt.Model
			}
			sqlKey, argsKey := GetKeySQL(v[0], gt.KeyModel, gt.tableT)
			bufW.WriteString(sqlKey)
			gt.Args = append(gt.Args, argsKey...)
			continue
		}
		gt.whereKv(&bufW, k, v[0])
	}

	gt.whereSQLNt(&bufW)
	return
}

// GetSelectSQL select sql
func (gt *GT) GetSelectSQL() {

	gt.whereParams()
	gt.sql = gt.Select
	if gt.From == "" {
		gt.From = "from"
	}
	gt.sqlNt = cons2.SelectCount + gt.From + strings.Join(strings.Split(gt.sql, gt.From)[1:], "")
	if gt.Group != "" {
		gt.sql += gt.Group
	}
	return
}

// other tables where
func (gt *GT) otherTableWhere(bufW *bytes.Buffer, tables []string, k, v string) (b bool) {
	// other tables, except tables[0]
	for _, tb := range tables {
		if !strings.Contains(k, tb+"_id") && strings.Contains(k, tb+"_") {
			gt.whereTbKv(bufW, tb, string([]byte(k)[len(tb)+1:]), v)
			b = true
			return
		}
	}
	return
}

// more sql
func (gt *GT) moreSql() (tables []string) {
	var (
		typ   = reflect.TypeOf(gt.Model)
		keyNt = mr.Path(typ, "sqlNt")
		keyTs = mr.Path(typ, "sqlNtTables")
	)
	gt.sqlNt = buffer.Get(keyNt)
	if tables = strings.Split(buffer.Get(keyTs), ","); tables[0] == "" {
		tables = []string{}
	}
	if gt.sqlNt != "" {
		return
	}

	innerTables, leftTables, innerField, leftField, DBS := gt.moreTables()
	tables = append(tables, innerTables...)
	tables = append(tables, leftTables...)
	tables = util.RemoveDuplicate(tables)

	var (
		bufNt bytes.Buffer // sql bytes connect
		count = cons2.SelectCount
	)
	if gt.distinct != "" {
		count = fmt.Sprintf(cons2.SelectCountDistinct, gt.distinct)
	}
	// sql and sqlCount
	bufNt.WriteString(count)
	bufNt.WriteString("from ")
	if tb := DBS[tables[0]]; tb != "" {
		bufNt.WriteByte(cons2.Backticks)
		bufNt.WriteString(tb)
		bufNt.WriteByte(cons2.Backticks)
		bufNt.WriteByte('.')
	}
	bufNt.WriteByte(cons2.Backticks)
	bufNt.WriteString(tables[0])
	bufNt.WriteByte(cons2.Backticks)
	bufNt.WriteByte(' ')
	// inner join
	for i := 1; i < len(innerTables); i += 2 {
		bufNt.WriteString("inner join ")
		innerLeftSQL(&bufNt, DBS, innerTables, innerField, i)
	}
	// left join
	for i := 1; i < len(leftTables); i += 2 {
		bufNt.WriteString("left join ")
		innerLeftSQL(&bufNt, DBS, leftTables, leftField, i)
	}
	gt.sqlNt = bufNt.String()
	buffer.Set(keyNt, gt.sqlNt)
	buffer.Set(keyTs, strings.Join(tables, ","))
	return
}

// more sql tables
// can read by ram
func (gt *GT) moreTables() (innerTables, leftTables, innerField, leftField []string, DBS map[string]string) {

	for k, v := range gt.InnerTable {
		st := strings.Split(v, ":")

		if strings.Contains(st[0], ".") {
			sts := strings.Split(st[0], ".")
			if DBS == nil {
				DBS = make(map[string]string)
			}
			DBS[sts[1]] = sts[0]
			st[0] = sts[1]
		}
		innerTables = append(innerTables, st[0])
		if len(st) == 1 { // default
			field := "id"
			if k%2 == 0 {
				// default other table_id
				otb := strings.Split(gt.InnerTable[k+1], ":")[0]
				if strings.Contains(otb, ".") {
					otb = strings.Split(otb, ".")[1]
				}
				field = otb + "_id"
			}
			innerField = append(innerField, field)
		} else {
			innerField = append(innerField, st[1])
		}
	}
	// left
	for k, v := range gt.LeftTable {
		st := strings.Split(v, ":")

		if strings.Contains(st[0], ".") {
			sts := strings.Split(st[0], ".")
			if DBS == nil {
				DBS = make(map[string]string)
			}
			DBS[sts[1]] = sts[0]
			st[0] = sts[1]
		}
		leftTables = append(leftTables, st[0])
		if len(st) == 1 {
			field := "id"
			if k%2 == 0 {
				// default other table_id
				otb := strings.Split(gt.LeftTable[k+1], ":")[0]
				if strings.Contains(otb, ".") {
					otb = strings.Split(otb, ".")[1]
				}
				field = otb + "_id"
			}
			leftField = append(leftField, field)
		} else {
			leftField = append(leftField, st[1])
		}
	}
	return
}

// gt some params
func (gt *GT) whereParams() {

	for k, v := range gt.CMaps {
		if v[0] == "" {
			gt.CMaps.Del(k)
			continue
		}
		switch k {
		case cons2.GtClientPage, cons2.GtClientPageUnderLine:
			gt.clientPage, _ = strconv.ParseInt(v[0], 10, 64)
			gt.CMaps.Del(k)
			continue
		case cons2.GtEveryPage, cons2.GtEveryPageUnderLine:
			gt.everyPage, _ = strconv.ParseInt(v[0], 10, 64)
			gt.CMaps.Del(k)
			continue
		case cons2.GtOrder:
			gt.order = v[0]
			gt.CMaps.Del(k)
			continue
		case cons2.GtMock:
			mock.Mock(gt.Data)
			gt.isMock = true
			gt.CMaps.Del(k)
			return
		case "":
			gt.CMaps.Del(k)
			continue
		}
	}
}

// sql and sqlNt where sql
func (gt *GT) whereAllSQL(bufW *bytes.Buffer) {

	gt.whereSQLNt(bufW)
	gt.whereCount()
	return
}

func (gt *GT) whereCount() {
	gt.sqlNt += gt.sqlW
	return
}

// sql where sql
func (gt *GT) whereSQLNt(bufW *bytes.Buffer) {
	gt.whereSQL(bufW)
	gt.sql += gt.sqlW
	if gt.order != "" {
		gt.sql += fmt.Sprintf(cons2.OrderS, gt.order)
	}
	return
}

func (gt *GT) whereSQL(bufW *bytes.Buffer) {
	var softConnect = cons2.WhereS
	if bufW.Len() != 0 {
		gt.sqlW += fmt.Sprintf(cons2.WhereS, bufW.Bytes()[:bufW.Len()-5])
		if gt.WhereSQL != "" {
			gt.Args = append(gt.Args, gt.wArgs...)
			gt.sqlW += fmt.Sprintf(cons2.AndS, gt.WhereSQL)
		}
		softConnect = cons2.AndS
	} else if gt.WhereSQL != "" {
		gt.Args = append(gt.Args, gt.wArgs...)
		gt.sqlW += fmt.Sprintf(cons2.WhereS, gt.WhereSQL)
		softConnect = cons2.AndS
	}

	if gt.sqlSoft != "" {
		gt.sqlW += fmt.Sprintf(softConnect, gt.sqlSoft)
	}
}

// where k =/in v
func (gt *GT) whereKv(bufW *bytes.Buffer, k, v string) {
	bufW.WriteString(k)
	gt.where(bufW, k, v)
}

// where k =/in v
func (gt *GT) whereTbKv(bufW *bytes.Buffer, tb, k, v string) {
	bufW.WriteByte(cons2.Backticks)
	bufW.WriteString(tb)
	bufW.WriteByte(cons2.Backticks)
	bufW.WriteByte('.')
	bufW.WriteByte(cons2.Backticks)
	bufW.WriteString(k)
	bufW.WriteByte(cons2.Backticks)
	gt.where(bufW, k, v)
}

func (gt *GT) where(bufW *bytes.Buffer, k, v string) {
	if strings.Contains(v, cons2.GtComma) {
		bufW.WriteString(cons2.ParamInAnd)
		gt.Args = append(gt.Args, strings.Split(v, cons2.GtComma)) // args
	} else {
		if p := gt.parses.Get(k); p != nil && p.Get(cons2.GtLike) == cons2.GtExist {
			bufW.WriteString(cons2.ParamLike)
			v = "%" + v + "%"
		} else {
			bufW.WriteString(cons2.ParamAnd)
		}
		gt.Args = append(gt.Args, v) // args
	}
}
