// package der

/*
	go-tool is a fast go tool, help you dev project
*/

package der

import (
	"github.com/dreamlu/go-tool/tool/result"
)

const Version = "1.1.x"

// crud
type Crud interface {
	// init db tool
	InitDBTool(dbTool *DBTool)
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

	// crud and search id
	// form data
	UpdateForm(args map[string][]string) error        // update
	CreateForm(args map[string][]string) error        // create
	CreateResID(args map[string][]string) (ID, error) // create res insert id

	// crud and search id
	// json data
	Update(data interface{}) error          // update
	Create(data interface{}) error          // create, include res insert id
	CreateMoreDataJ(data interface{}) error // create more
}

// crud params
type CrudParam struct {
	// attributes
	InnerTables []string    // inner join tables
	LeftTables  []string    // left join tables
	Table       string      // table name
	Model       interface{} // table model, like User{}
	ModelData   interface{} // table model data, like var user User{}, it is 'user'

	// pager info
	ClientPage int64 // page number
	EveryPage  int64 // Number of pages per page
}
