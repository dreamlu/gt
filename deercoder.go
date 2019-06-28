// @author  dreamlu
package der

import (
	"github.com/dreamlu/go-tool/util/result"
)

const Version = "1.0.x"

// db database
type DataBase interface {
	// nothing
}

// crud
type Crud interface {
	// crud method

	// get url params
	// like form data
	GetBySearch(args map[string][]string) result.GetInfoPager     // search
	GetByID(id string) result.GetInfo                             // by id
	GetMoreBySearch(args map[string][]string) result.GetInfoPager // more search

	// common sql data
	// through sql, get the data
	GetDataBySQL(sql string, args ...interface{}) result.GetInfo // single data
	// page limit ?,?
	// args not include limit ?,?
	GetDataBySearchSQL(sql, sqlnolimit string, args ...interface{}) result.GetInfoPager // more data
	DeleteBySQL(sql string, args ...interface{}) result.MapData
	UpdateBySQL(sql string, args ...interface{}) result.MapData
	CreateBySQL(sql string, args ...interface{}) result.MapData
}

// common crud
// detail impl, ==>DbCrud, implement DBCrud
// form data
type DBCruder interface {
	// crud and search id
	Create(args map[string][]string) result.MapData      // create
	CreateResID(args map[string][]string) result.GetInfo // create res insert id
	Update(args map[string][]string) result.MapData      // update
	Delete(id string) result.MapData                     // delete

	// common sql data
	Crud
}

// common crud
// json data
type DBCrudJer interface {
	// crud and search id
	Create(data interface{}) result.MapData      // create
	CreateResID(data interface{}) result.GetInfo // create res insert id
	Update(data interface{}) result.MapData      // update
	Delete(id string) result.MapData             // delete

	// common sql data
	Crud
}
