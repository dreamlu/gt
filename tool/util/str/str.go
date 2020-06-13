// package gt

package str

// max upload file size
var MaxUploadMemory int64

// struct value
type Value struct {
	Value string `json:"value"`
}

type Num struct {
	Num int `json:"num"`
}

// ID struct
type ID struct {
	ID int64 `json:"id"`
}

// string
type String interface {
	String() (string, error)
}

// devMode const
// key words
const (
	// devMode
	Dev  = "dev"
	Prod = "prod"
	// default config file dir
	ConfDir = "conf/"
	// db sql const
	GtSubSQL = "sub_sql"
	// Deprecated
	GtClientPage          = "clientPage"
	GtClientPageUnderLine = "client_page" // 支持下划线
	// Deprecated
	GtEveryPage          = "everyPage"
	GtEveryPageUnderLine = "every_page" // 支持下划线
	GtOrder              = "order"
	GtKey                = "key"
	GtMock               = "mock"
	// gt tag
	GtField = "field"
)
