// package gt

package gt

import (
	"bytes"
	"fmt"
	"github.com/dreamlu/gt/tool/mock"
	sq "github.com/dreamlu/gt/tool/sql"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util"
	"github.com/dreamlu/gt/tool/util/cons"
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
	Select     string // select sql
	From       string // only once
	Group      string // the last group
	Args       []interface{}
	sql        string
	sqlNt      string
	clientPage int64
	everyPage  int64
	order      string // order by

	// mock data
	isMock bool
}

//=======================================sql script==========================================
//===========================================================================================

// more table
// params: innerTables is inner join tables, must even number
// params: leftTables is left join tables
// return: select sql
// table1 as main table, include other tables_id(foreign key)
func (gt *GT) GetMoreSQL() {

	var (
		tables []string
		bufW   bytes.Buffer // gt.sql bytes connect
		count  = cons.Count
	)
	tables = gt.moreSql()
	var sql = GetMoreTableColumnSQL(gt.Model, tables[:]...)
	if gt.distinct != "" {
		count = fmt.Sprintf(cons.CountDistinct, gt.distinct)
		sql = cons.Distinct + sql
	}
	gt.sql = strings.Replace(gt.sqlNt, count, sql+gt.SubSQL, 1)
	// default
	gt.order = fmt.Sprintf(cons.OrderDesc, sq.Table(tables[0]))

	gt.whereParams()
	for k, v := range gt.CMaps {
		if k == cons.GtKey {
			if gt.KeyModel == nil {
				gt.KeyModel = gt.Model
			}
			// more tables key search
			sqlKey, argsKey := sq.GetMoreKeySQL(v[0], gt.KeyModel, tables...)
			bufW.WriteString(sqlKey)
			gt.Args = append(gt.Args, argsKey...)
			continue
		}

		if b := otherTableWhere(&bufW, tables[1:], k); !b {
			v[0] = strings.Replace(v[0], "'", "\\'", -1)
			bufW.WriteString("`")
			bufW.WriteString(tables[0])
			bufW.WriteString("`.`")
			bufW.WriteString(k)
			bufW.WriteString("` = ? and ")
		}
		gt.Args = append(gt.Args, v[0])
	}

	gt.whereSQL(&bufW)
	return
}

// search sql
// default order by id desc
func (gt *GT) GetSearchSQL() {

	var (
		bufW  bytes.Buffer // where sql, sqlNt bytes sql
		table = sq.Table(gt.Table)
	)
	// default
	gt.order = fmt.Sprintf(cons.OrderDesc, table)

	// select* replace
	gt.sql = fmt.Sprintf(cons.SelectFrom, GetColSQL(gt.Model)+gt.SubSQL, table)
	gt.sqlNt = fmt.Sprintf(cons.SelectCountFrom, table)

	gt.whereParams()
	for k, v := range gt.CMaps {
		if k == cons.GtKey {
			if gt.KeyModel == nil {
				gt.KeyModel = gt.Model
			}
			sqlKey, argsKey := sq.GetKeySQL(v[0], gt.KeyModel, table)
			bufW.WriteString(sqlKey)
			gt.Args = append(gt.Args, argsKey...)
			continue
		}
		bufW.WriteString(k)
		bufW.WriteString(cons.ParamAnd)
		gt.Args = append(gt.Args, v[0]) // args
	}

	gt.whereSQL(&bufW)
	return
}

// get single sql
func (gt *GT) GetSQL() {

	var (
		bufW  bytes.Buffer // where sql, sqlNt bytes sql
		table = sq.Table(gt.Table)
	)

	// select* replace
	gt.sql = fmt.Sprintf(cons.SelectFrom, GetColSQL(gt.Model)+gt.SubSQL, table)

	gt.whereParams()
	for k, v := range gt.CMaps {
		if k == cons.GtKey {
			if gt.KeyModel == nil {
				gt.KeyModel = gt.Model
			}
			sqlKey, argsKey := sq.GetKeySQL(v[0], gt.KeyModel, table)
			bufW.WriteString(sqlKey)
			gt.Args = append(gt.Args, argsKey...)
			continue
		}
		bufW.WriteString(k)
		bufW.WriteString(cons.ParamAnd)
		gt.Args = append(gt.Args, v[0]) // args
	}

	gt.whereSQLNt(&bufW)
	return
}

// select sql
func (gt *GT) GetSelectSearchSQL() {

	gt.whereParams()
	gt.sql = gt.Select
	if gt.From == "" {
		gt.From = "from"
	}
	gt.sqlNt = cons.SelectCount + gt.From + strings.Join(strings.Split(gt.sql, gt.From)[1:], "")
	if gt.Group != "" {
		gt.sql += gt.Group
	}
	return
}

// other tables where
func otherTableWhere(bufW *bytes.Buffer, tables []string, k string) (b bool) {

	tb := sq.UniqueTagTable(k, tables...)
	if tb != "" {
		writeBufWhere(bufW, tb, k)
		return
	}

	// other tables, except tables[0]
	for _, v := range tables {
		if !strings.Contains(k, v+"_id") && strings.Contains(k, v+"_") {
			writeBufWhere(bufW, v, string([]byte(k)[len(v)+1:]))
			b = true
			return
		}
	}
	return
}

// more sql
func (gt *GT) moreSql() (tables []string) {
	typ := reflect.TypeOf(gt.Model)
	keyNt := typ.PkgPath() + "/sqlNt/" + typ.Name()
	keyTs := typ.PkgPath() + "/sqlNtTables/" + typ.Name()
	gt.sqlNt = sqlBuffer.Get(keyNt)
	if tables = strings.Split(sqlBuffer.Get(keyTs), ","); tables[0] == "" {
		tables = []string{}
	}
	if gt.sqlNt != "" {
		return
	}

	innerTables, leftTables, innerField, leftField, DBS := gt.moreTables()
	tables = append(tables, innerTables...)
	tables = append(tables, leftTables...)
	tables = util.RemoveDuplicateString(tables)

	var (
		bufNt bytes.Buffer // sql bytes connect
		count = cons.SelectCount
	)
	if gt.distinct != "" {
		count = fmt.Sprintf(cons.SelectCountDistinct, gt.distinct)
	}
	// sql and sqlCount
	bufNt.WriteString(count)
	bufNt.WriteString("from ")
	if tb := DBS[tables[0]]; tb != "" {
		bufNt.WriteString("`")
		bufNt.WriteString(tb)
		bufNt.WriteString("`.")
	}
	bufNt.WriteString("`")
	bufNt.WriteString(tables[0])
	bufNt.WriteString("` ")
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
	sqlBuffer.Set(keyNt, gt.sqlNt)
	sqlBuffer.Set(keyTs, strings.Join(tables, ","))
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
		case cons.GtClientPage, cons.GtClientPageUnderLine:
			gt.clientPage, _ = strconv.ParseInt(v[0], 10, 64)
			gt.CMaps.Del(k)
			continue
		case cons.GtEveryPage, cons.GtEveryPageUnderLine:
			gt.everyPage, _ = strconv.ParseInt(v[0], 10, 64)
			gt.CMaps.Del(k)
			continue
		case cons.GtOrder:
			gt.order = v[0]
			gt.CMaps.Del(k)
			continue
		case cons.GtMock:
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
func (gt *GT) whereSQL(bufW *bytes.Buffer) {

	gt.whereSQLNt(bufW)
	if bufW.Len() != 0 {
		gt.sqlNt += fmt.Sprintf(cons.WhereS, bufW.Bytes()[:bufW.Len()-5])
		if gt.WhereSQL != "" {
			gt.sqlNt += fmt.Sprintf(cons.AndS, gt.WhereSQL)
		}
	} else if gt.WhereSQL != "" {
		gt.sqlNt += fmt.Sprintf(cons.WhereS, gt.WhereSQL)
	}
	return
}

// sql where sql
func (gt *GT) whereSQLNt(bufW *bytes.Buffer) {
	if bufW.Len() != 0 {
		gt.sql += fmt.Sprintf(cons.WhereS, bufW.Bytes()[:bufW.Len()-5])
		if gt.WhereSQL != "" {
			gt.Args = append(gt.Args, gt.wArgs...)
			gt.sql += fmt.Sprintf(cons.AndS, gt.WhereSQL)
		}
	} else if gt.WhereSQL != "" {
		gt.Args = append(gt.Args, gt.wArgs...)
		gt.sql += fmt.Sprintf(cons.WhereS, gt.WhereSQL)
	}
	if gt.order != "" {
		gt.sql += fmt.Sprintf(cons.OrderS, gt.order)
	}
	return
}

// form-data create/update
// future will remove
// use json replace

// get update sql
func GetUpdateSQL(table string, params cmap.CMap) (sql string, args []interface{}) {

	// sql connect
	var (
		id  string       // id
		buf bytes.Buffer // sql bytes connect
	)
	buf.WriteString("update `")
	buf.WriteString(table)
	buf.WriteString("` set ")
	for k, v := range params {
		if k == "id" {
			id = v[0]
			continue
		}
		buf.WriteString("`")
		buf.WriteString(k)
		buf.WriteString("` = ?,")
		args = append(args, v[0])
	}
	args = append(args, id)
	sql = string(buf.Bytes()[:buf.Len()-1]) + " where id = ?"
	return sql, args
}

// get insert sql
func GetInsertSQL(table string, params cmap.CMap) (sql string, args []interface{}) {

	// sql connect
	var (
		sqlv string
		buf  bytes.Buffer // sql bytes connect
	)
	buf.WriteString("insert `")
	buf.WriteString(table)
	buf.WriteString("`(")
	//sql = "insert `" + table + "`("

	for k, v := range params {
		buf.WriteString("`")
		buf.WriteString(k)
		buf.WriteString("`,")

		args = append(args, v[0])
		sqlv += "?,"
	}
	//sql = buf.Bytes()[:buf.Len()-1]
	sql = buf.String()
	sql = string([]byte(sql)[:len(sql)-1]) + ") value(" + sqlv
	sql = string([]byte(sql)[:len(sql)-1]) + ")" // remove ','

	return sql, args
}
