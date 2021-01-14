// package gt

package gt

import (
	"bytes"
	"fmt"
	"github.com/dreamlu/gt/tool/mock"
	reflect2 "github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/result"
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
func GetMoreSQL(gt *GT) {

	var tables []string
	gt.sqlNt, tables = moreSql(gt)
	// select* 变为对应的字段名
	gt.sql = strings.Replace(gt.sqlNt, "count(`"+tables[0]+"`.id) as total_num", GetMoreTableColumnSQL(gt.Model, tables[:]...)+gt.SubSQL, 1)
	var (
		order = "`" + tables[0] + "`.id desc" // order by
		key   = ""                            // key like binary search
		bufW  bytes.Buffer                    // gt.sql bytes connect
	)
	for k, v := range gt.CMaps {
		if v[0] == "" {
			continue
		}
		switch k {
		case cons.GtClientPage, cons.GtClientPageUnderLine:
			gt.clientPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case cons.GtEveryPage, cons.GtEveryPageUnderLine:
			gt.everyPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case cons.GtOrder:
			order = v[0]
			continue
		case cons.GtKey:
			key = v[0]
			if gt.KeyModel == nil {
				gt.KeyModel = gt.Model
			}
			var tablens = append(tables, tables[:]...)
			for k, v := range tablens {
				tablens[k] += ":" + v
			}
			// more tables key search
			sqlKey, argsKey := sq.GetMoreKeySQL(key, gt.KeyModel, tablens[:]...)
			bufW.WriteString(sqlKey)
			gt.Args = append(gt.Args, argsKey[:]...)
			continue
		case cons.GtMock:
			mock.Mock(gt.Data)
			gt.isMock = true
			return
		case "":
			continue
		}

		if b := otherTableWhere(&bufW, tables[1:], k); b != true {
			v[0] = strings.Replace(v[0], "'", "\\'", -1)
			bufW.WriteString("`")
			bufW.WriteString(tables[0])
			bufW.WriteString("`.`")
			bufW.WriteString(k)
			bufW.WriteString("` = ? and ")
		}
		gt.Args = append(gt.Args, v[0])
	}

	if bufW.Len() != 0 {
		gt.sql += fmt.Sprintf("where %s ", bufW.Bytes()[:bufW.Len()-4])
		gt.sqlNt += fmt.Sprintf("where %s", bufW.Bytes()[:bufW.Len()-4])
		if gt.WhereSQL != "" {
			gt.Args = append(gt.Args, gt.wArgs...)
			gt.sql += fmt.Sprintf("and %s ", gt.WhereSQL)
			gt.sqlNt += fmt.Sprintf("and %s", gt.WhereSQL)
		}
	} else if gt.WhereSQL != "" {
		gt.Args = append(gt.Args, gt.wArgs...)
		gt.sql += fmt.Sprintf("where %s ", gt.WhereSQL)
		gt.sqlNt += fmt.Sprintf("where %s", gt.WhereSQL)
	}
	gt.sql += fmt.Sprintf(" order by %s ", order)

	return
}

func otherTableWhere(bufW *bytes.Buffer, tables []string, k string) (b bool) {
	// other tables, except tables[0]
	for _, v := range tables {
		switch {
		case !strings.Contains(k, v+"_id") && strings.Contains(k, v+"_"):
			//bufW.WriteString("`" + table + "`.`" + string([]byte(k)[len(v)+1:]) + "` = ? and ")
			bufW.WriteString("`")
			bufW.WriteString(v)
			bufW.WriteString("`.`")
			bufW.WriteString(string([]byte(k)[len(v)+1:]))
			bufW.WriteString("` = ? and ")
			//args = append(args, v[0])
			b = true
			return
		}
	}
	return
}

// more sql
func moreSql(gt *GT) (sqlNt string, tables []string) {

	// read ram
	typ := reflect.TypeOf(gt.Model)
	keyNt := typ.PkgPath() + "/sqlNt/" + typ.Name()
	keyTs := typ.PkgPath() + "/sqlNtTables/" + typ.Name()
	sqlNt = sqlBuffer.Get(keyNt)
	if tables = strings.Split(sqlBuffer.Get(keyTs), ","); tables[0] == "" {
		tables = []string{}
	}
	if sqlNt != "" {
		//Logger().Info("[USE sqlBuffer GET ColumnSQL]")
		return
	}

	innerTables, leftTables, innerField, leftField, DBS := moreTables(gt)
	tables = append(tables, innerTables...)
	tables = append(tables, leftTables...)
	tables = util.RemoveDuplicateString(tables)

	var (
		//tables = innerTables // all tables
		bufNt bytes.Buffer // sql bytes connect
	)
	// sql and sqlCount
	bufNt.WriteString("select count(`")
	bufNt.WriteString(tables[0])
	bufNt.WriteString("`.id) as total_num from ")
	if tb := DBS[tables[0]]; tb != "" {
		bufNt.WriteString("`" + tb + "`.")
	}
	bufNt.WriteString("`")
	bufNt.WriteString(tables[0])
	bufNt.WriteString("` ")
	// inner join
	for i := 1; i < len(innerTables); i += 2 {
		bufNt.WriteString("inner join ")
		// innerDB not support ``
		if tb := DBS[innerTables[i]]; tb != "" {
			bufNt.WriteString("`" + tb + "`.")
		}
		bufNt.WriteString("`")
		bufNt.WriteString(innerTables[i])
		bufNt.WriteString("` on `")
		bufNt.WriteString(innerTables[i-1])
		bufNt.WriteString("`.`")
		bufNt.WriteString(innerField[i-1])
		bufNt.WriteString("`=`")
		bufNt.WriteString(innerTables[i])
		bufNt.WriteString("`.`")
		bufNt.WriteString(innerField[i])
		bufNt.WriteString("` ")
	}
	// left join
	for i := 1; i < len(leftTables); i += 2 {
		bufNt.WriteString("left join ")
		if tb := DBS[leftTables[i]]; tb != "" {
			bufNt.WriteString("`" + tb + "`.")
		}
		bufNt.WriteString("`")
		bufNt.WriteString(leftTables[i])
		bufNt.WriteString("` on `")
		bufNt.WriteString(leftTables[i-1])
		bufNt.WriteString("`.`")
		bufNt.WriteString(leftField[i-1])
		bufNt.WriteString("`=`")
		bufNt.WriteString(leftTables[i])
		bufNt.WriteString("`.`")
		bufNt.WriteString(leftField[i])
		bufNt.WriteString("` ")
	}
	sqlNt = bufNt.String()
	sqlBuffer.Add(keyNt, sqlNt)
	sqlBuffer.Add(keyTs, strings.Join(tables, ","))
	return
}

// more sql tables
// can read by ram
func moreTables(gt *GT) (innerTables, leftTables, innerField, leftField []string, DBS map[string]string) {

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

// search sql
// default order by id desc
func GetSearchSQL(gt *GT) (sqlNt, sql string, clientPage, everyPage int64, args []interface{}) {

	var (
		order        = "id desc"  // default order by
		key          = ""         // key like binary search
		bufW, bufNtW bytes.Buffer // where sql, sqlNt bytes sql
		table        = sq.Table(gt.Table)
	)

	// select* replace
	sql = fmt.Sprintf("select %s%s from %s ", GetColSQL(gt.Model), gt.SubSQL, table)
	sqlNt = fmt.Sprintf("select count(id) as total_num from %s ", table)
	for k, v := range gt.CMaps {
		if v[0] == "" {
			continue
		}
		switch k {
		case cons.GtClientPage, cons.GtClientPageUnderLine:
			clientPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case cons.GtEveryPage, cons.GtEveryPageUnderLine:
			everyPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case cons.GtOrder:
			order = v[0]
			continue
		case cons.GtKey:
			key = v[0]
			if gt.KeyModel == nil {
				gt.KeyModel = gt.Model
			}
			sqlKey, argsKey := sq.GetKeySQL(key, gt.KeyModel, table)
			bufW.WriteString(sqlKey)
			bufNtW.WriteString(sqlKey)
			args = append(args, argsKey[:]...)
			continue
		case cons.GtMock:
			mock.Mock(gt.Data)
			gt.isMock = true
			return
		case "":
			continue
		}
		bufW.WriteString(k)
		bufW.WriteString(" = ? and ")
		bufNtW.WriteString(k)
		bufNtW.WriteString(" = ? and ")
		args = append(args, v[0]) // args
	}

	if bufW.Len() != 0 {
		sql += fmt.Sprintf("where %s ", bufW.Bytes()[:bufW.Len()-4])
		sqlNt += fmt.Sprintf("where %s", bufNtW.Bytes()[:bufNtW.Len()-4])
		if gt.WhereSQL != "" {
			gt.Args = append(gt.Args, gt.wArgs...)
			sql += fmt.Sprintf("and %s ", gt.WhereSQL)
			sqlNt += fmt.Sprintf("and %s", gt.WhereSQL)
		}
	} else if gt.WhereSQL != "" {
		gt.Args = append(gt.Args, gt.wArgs...)
		sql += fmt.Sprintf(" where %s ", gt.WhereSQL)
		sqlNt += fmt.Sprintf(" where %s", gt.WhereSQL)
	}
	sql += fmt.Sprintf(" order by %s ", order)
	return
}

// get single sql
func GetSQL(gt *GT) (sql string, args []interface{}) {

	var (
		order = ""         // default no order by
		key   = ""         // key like binary search
		bufW  bytes.Buffer // where sql, sqlNt bytes sql
		table = sq.Table(gt.Table)
	)

	// select* replace
	sql = fmt.Sprintf("select %s%s from %s ", GetColSQL(gt.Model), gt.SubSQL, table)
	for k, v := range gt.CMaps {
		if v[0] == "" {
			continue
		}
		switch k {
		case cons.GtOrder:
			order = v[0]
			continue
		case cons.GtKey:
			key = v[0]
			if gt.KeyModel == nil {
				gt.KeyModel = gt.Model
			}
			sqlKey, argsKey := sq.GetKeySQL(key, gt.KeyModel, table)
			bufW.WriteString(sqlKey)
			args = append(args, argsKey[:]...)
			continue
		case cons.GtMock:
			mock.Mock(gt.Data)
			gt.isMock = true
			return
		case "":
			continue
		}
		bufW.WriteString(k)
		bufW.WriteString(" = ? and ")
		args = append(args, v[0]) // args
	}

	if bufW.Len() != 0 {
		sql += fmt.Sprintf(" where %s ", bufW.Bytes()[:bufW.Len()-4])
		if gt.WhereSQL != "" {
			gt.Args = append(gt.Args, gt.wArgs...)
			sql += fmt.Sprintf("and %s ", gt.WhereSQL)
		}
	} else if gt.WhereSQL != "" {
		gt.Args = append(gt.Args, gt.wArgs...)
		sql += fmt.Sprintf(" where %s ", gt.WhereSQL)
	}
	if order != "" {
		sql += fmt.Sprintf(" order by %s ", order)
	}
	return
}

// select sql
func GetSelectSearchSQL(gt *GT) (sqlNt, sql string, clientPage, everyPage int64) {

	for k, v := range gt.CMaps {
		switch k {
		case cons.GtClientPage, cons.GtClientPageUnderLine:
			clientPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		case cons.GtEveryPage, cons.GtEveryPageUnderLine:
			everyPage, _ = strconv.ParseInt(v[0], 10, 64)
			continue
		}
	}

	sql = gt.Select
	if gt.From == "" {
		gt.From = "from"
	}
	sqlNt = "select count(*) as total_num " + gt.From + strings.Join(strings.Split(sql, gt.From)[1:], "")
	if gt.Group != "" {
		sql += gt.Group
	}
	return
}

// ===================================================================================
// ==========================common crud=========== dreamlu ==========================
// ===================================================================================

// get
// relation get
////////////////

// get single data
func (db *DBTool) GetBySQL(data interface{}, sql string, args ...interface{}) {

	typ := reflect.TypeOf(data)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	db.res = db.DB.Raw(sql, args[:]...).Scan(data)
}

// get by id
func (db *DBTool) GetByID(gt *GT, id interface{}) {

	db.GetBySQL(gt.Data, fmt.Sprintf("select %s from %s where id = ?", GetColSQL(gt.Model), sq.Table(gt.Table)), id)
}

// more table
// params: innerTables is inner join tables
// params: leftTables is left join tables
// return search info
// table1 as main table, include other tables_id(foreign key)
func (db *DBTool) GetMoreBySearch(gt *GT) (pager result.Pager) {
	// more table search
	GetMoreSQL(gt)
	// isMock
	if gt.isMock {
		return
	}
	return db.GetBySQLSearch(gt.Data, gt.sql, gt.sqlNt, gt.clientPage, gt.everyPage, gt.Args)
}

// single table
// return search info
func (db *DBTool) GetBySearch(gt *GT) (pager result.Pager) {

	sqlNt, sql, clientPage, everyPage, args := GetSearchSQL(gt)
	// isMock
	if gt.isMock {
		return
	}
	return db.GetBySQLSearch(gt.Data, sql, sqlNt, clientPage, everyPage, args)
}

// 获得数据, no search
func (db *DBTool) Get(gt *GT) {

	sql, args := GetSQL(gt)
	// isMock
	if gt.isMock {
		return
	}
	db.GetBySQL(gt.Data, sql, args[:]...)
}

// 获得数据, no search
func (db *DBTool) GetMoreData(gt *GT) {

	GetMoreSQL(gt)
	// isMock
	if gt.isMock {
		return
	}
	db.GetBySQL(gt.Data, gt.sql, gt.Args...)
}

// select sql search
func (db *DBTool) GetDataBySelectSQLSearch(gt *GT) (pager result.Pager) {

	sqlNt, sql, clientPage, everyPage := GetSelectSearchSQL(gt)
	return db.GetBySQLSearch(gt.Data, sql, sqlNt, clientPage, everyPage, gt.Args)
}

// get sql search data
// clientPage: default 1
// everyPage: default 10
// if clientPage or everyPage < 0, return all
func (db *DBTool) GetBySQLSearch(data interface{}, sql, sqlNt string, clientPage, everyPage int64, args []interface{}) (pager result.Pager) {

	// if clientPage or everyPage < 0
	// return all data
	if clientPage == 0 {
		clientPage = cons.ClientPage
	}
	if everyPage == 0 {
		everyPage = cons.EveryPage
	}
	if clientPage > 0 && everyPage > 0 {
		sql += fmt.Sprintf("limit %d, %d", (clientPage-1)*everyPage, everyPage)
	}
	// sqlNt += limit
	db.res = db.DB.Raw(sqlNt, args[:]...).Scan(&pager)
	if db.res.Error == nil {
		db.res = db.DB.Raw(sql, args[:]...).Scan(data)
		// pager data
		pager.ClientPage = clientPage
		pager.EveryPage = everyPage
		return
	}
	return
}

// exec common
////////////////////

// exec sql
func (db *DBTool) ExecSQL(sql string, args ...interface{}) {
	db.res = db.Exec(sql, args...)
}

// delete
///////////////////

// delete
func (db *DBTool) Delete(table string, id interface{}) {
	switch id.(type) {
	case string:
		if strings.Contains(id.(string), ",") {
			id = strings.Split(id.(string), ",")
		}
	}
	db.ExecSQL(fmt.Sprintf("delete from %s where id in (?)", sq.Table(table)), id)
}

// update
///////////////////

// update
func (db *DBTool) Update(gt *GT) {

	if gt.Model == nil {
		gt.Model = gt.Data
	}

	if gt.Select != "" {
		db.res = db.Table(gt.Table).Model(gt.Model).Where(gt.Select, gt.Args).Updates(gt.Data)
	} else {
		db.res = db.Table(gt.Table).Model(gt.Data).Updates(gt.Data)
	}
}

// create
////////////////////

// create
func (db *DBTool) Create(table string, data interface{}) {
	db.res = db.Table(table).Create(data)
}

// data must array type
// more data create
// single table
func (db *DBTool) CreateMore(table string, model interface{}, data interface{}) {
	var (
		buf    bytes.Buffer
		params []interface{}
	)

	// array data
	arrayData := reflect2.ToSlice(data)
	colPSQL := GetColParamSQL(model)

	for _, v := range arrayData {
		// buf
		buf.WriteString("(")
		buf.WriteString(colPSQL)
		buf.WriteString("),")
		// params
		params = append(params, GetParams(v)[:]...)
	}
	values := string(buf.Bytes()[:buf.Len()-1])

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", sq.Table(table), GetColSQL(model), values)
	db.res = db.DB.Exec(sql, params[:]...)
}
