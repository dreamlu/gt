package crud

import (
	"bytes"
	"fmt"
	"github.com/dreamlu/gt/conf"
	depCons "github.com/dreamlu/gt/crud/dep/cons"
	"github.com/dreamlu/gt/crud/dep/result"
	"github.com/dreamlu/gt/src/cons"
	mr "github.com/dreamlu/gt/src/reflect"
)

func (db *DB) Find(gt *GT) (pager result.Pager) {

	// isMock
	if gt.isMock {
		return
	}
	gt.GetSQL()
	if gt.isCount {
		db.countSQL(gt)
		pager = db.count(gt)
		if pager.TotalNum == 0 {
			return
		}
	}
	db.getLimit(gt)
	return
}

// FindM no search
// params: innerTables is inner join tables
// params: leftTables is left join tables
// return search info
// table1 as main table, include other tables_id(foreign key)
func (db *DB) FindM(gt *GT) (pager result.Pager) {
	// isMock
	if gt.isMock {
		return
	}
	gt.GetMoreSQL()
	if gt.isCount {
		pager = db.count(gt)
		if pager.TotalNum == 0 {
			return
		}
	}
	db.getLimit(gt)
	return
}

// FindS select sql search
func (db *DB) FindS(gt *GT) (pager result.Pager) {
	// isMock
	if gt.isMock {
		return
	}
	gt.GetSelectSQL()
	if gt.isCount {
		pager = db.count(gt)
		if pager.TotalNum == 0 {
			return
		}
	}
	db.getLimit(gt)
	return
}

func (db *DB) countSQL(gt *GT) *DB {

	// default
	gt.order = fmt.Sprintf(depCons.OrderDesc, gt.tableT)

	gt.sqlNt = fmt.Sprintf(depCons.SelectCountFrom, gt.tableT)
	gt.whereCount()

	return db
}

func (db *DB) count(gt *GT) (pager result.Pager) {

	// if clientPage or everyPage < 0
	// return all data
	if gt.clientPage == 0 {
		gt.clientPage = depCons.ClientPage
	}
	if gt.everyPage == 0 {
		gt.everyPage = depCons.EveryPage
	}
	db.res = db.DB.Raw(gt.sqlNt, gt.Args...).Scan(&pager)
	if db.res.Error != nil || pager.TotalNum == 0 {
		return
	}
	pager.ClientPage = gt.clientPage
	pager.EveryPage = gt.everyPage
	return
}

// get data
func (db *DB) getLimit(gt *GT) {
	if gt.clientPage > 0 && gt.everyPage > 0 {
		gt.sql += fmt.Sprintf("limit %d offset %d", gt.everyPage, (gt.clientPage-1)*gt.everyPage)
	}
	db.get(gt)
}

// get data
func (db *DB) get(gt *GT) {
	db.res = db.DB.Raw(gt.sql, gt.Args...).Scan(gt.Data)
}

func (db *DB) exec(sql string, args ...any) {
	db.res = db.Exec(sql, args...)
}

func (db *DB) Delete(gt *GT, conds ...any) {
	gt.parse().common()
	if gt.sqlSoft != "" {
		db.exec(fmt.Sprintf("update %s set %s = now() where id in (?)", gt.tableT, gt.parses.GetS(gt.Table)), conds...)
		return
	}
	db.res = db.DB.Delete(gt.Data, conds...)
}

func (db *DB) Update(gt *GT) {
	if gt.Select != "" {
		db.res = db.Table(gt.Table).Where(gt.Select, gt.Args).Updates(gt.Data)
	} else {
		db.res = db.Table(gt.Table).Model(gt.Data).Updates(gt.Data)
	}
}

// Create single/array
func (db *DB) Create(table string, data any) {
	db.res = db.Table(table).Create(data)
}

// CreateMore data must array type
// more data create
// single table
// also can use Create array
func (db *DB) CreateMore(table string, model any, data any) {
	var (
		buf       bytes.Buffer
		params    []any
		p         = parse(model)
		arrayData = mr.ToSlice(data) // slice data
		colPSQL   = GetColParamSQL(p)
	)

	for _, v := range arrayData {
		// buf
		buf.WriteByte('(')
		buf.WriteString(colPSQL)
		buf.WriteString("),")
		// params
		p.Vs = nil
		parseV(p, v)
		params = append(params, p.Vs...)
	}
	values := string(buf.Bytes()[:buf.Len()-1])

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", ParseTable(table), GetColSQL(model), values)
	db.res = db.DB.Exec(sql, params...)
}

// InitColumns init db table columns map
func (db *DB) InitColumns(param *Params) {

	var (
		name   = conf.Get[string](cons.ConfDBName)
		tables = []string{param.Table}
	)

	tables = append(tables, param.InnerTable...)
	tables = append(tables, param.LeftTable...)

	for _, v := range tables {
		if v == "" {
			continue
		}
		if _, ok := TableCols[v]; ok {
			continue
		}
		var columns []string
		//param.Data = &columns
		tb := TableOnly(v)
		db.get(&GT{
			Params: &Params{Data: &columns},
			sql:    "SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_NAME = ? and TABLE_SCHEMA = ?",
			Args:   getDriverArgs(depCons.Driver, tb, name),
		})
		TableCols[tb] = columns
	}
}

func getDriverArgs(driver, tb, name string) []any {
	switch driver {
	case depCons.Postgres:
		return []any{tb, "public"}
	default:
		return []any{tb, name}
	}
}
