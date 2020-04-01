// package gt

package str

// max upload file size
var MaxUploadMemory int64

// struct value
type Value struct {
	Value string `json:"value"`
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
	Dev  = "dev"
	Prod = "prod"
	// default config file dir
	ConfDir = "conf/"
	// db sql const
	GtSubSQL     = "sub_sql"
	GtClientPage = "clientPage"
	GtEveryPage  = "everyPage"
	GtOrder      = "order"
	GtKey        = "key"
	// gt tag
	GtField = "field"
)
