// @author  dreamlu
package der

import (
	"github.com/dreamlu/go-tool/tool/result"
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
	GetBySearch(args map[string][]string) (pager result.Pager, err error)     // search
	GetByID(id string) error                                                  // by id
	GetMoreBySearch(args map[string][]string) (pager result.Pager, err error) // more search

	// common sql data
	// through sql, get the data
	GetDataBySQL(sql string, args ...interface{}) error // single data
	// page limit ?,?
	// args not include limit ?,?
	GetDataBySearchSQL(sql, sqlNt string, args ...interface{}) (pager result.Pager, err error) // more data
	DeleteBySQL(sql string, args ...interface{}) error
	UpdateBySQL(sql string, args ...interface{}) error
	CreateBySQL(sql string, args ...interface{}) error

	// delete
	Delete(id string) error // delete
}

// common crud
// detail impl, ==>DbCrud, implement DBCrud
// form data
type DBCruder interface {
	// common sql data
	Crud

	// crud and search id
	Update(args map[string][]string) error            // update
	Create(args map[string][]string) error            // create
	CreateResID(args map[string][]string) (ID, error) // create res insert id
}

// common crud
// json data
type DBCrudJer interface {
	// common sql data
	Crud

	// crud and search id
	Update(data interface{}) error // update
	Create(data interface{}) error // create, include res insert id
}
