// package cons

package cons

// devMode const
// key words
const (
	// devMode
	Dev  = "dev"
	Prod = "prod"
	// default config file dir
	ConfDir = "conf/"
	// db sql const
	GtSubSQL              = "sub_sql"
	GtClientPage          = "clientPage"
	GtClientPageUnderLine = "client_page"
	GtEveryPage           = "everyPage"
	GtEveryPageUnderLine  = "every_page"
	GtOrder               = "order"
	GtKey                 = "key"
	GtMock                = "mock"
	// gt tag
	GT       = "gt"
	GtField  = "field"
	GtValid  = "valid"
	GtTrans  = "trans"
	GtIgnore = "ignore"
	Gt_      = "-"
)

// default page
const (
	ClientPage = 1
	EveryPage  = 10
)

// part sql
const (
	WhereS          = "where %s "
	AndS            = "and %s "
	OrderDesc       = "%s.id desc"
	OrderS          = "order by %s "
	ParamAnd        = " = ? and "
	SelectFrom      = "select %s from %s "
	Count           = "count(*) as total_num"
	SelectCount     = "select " + Count + " "
	SelectCountFrom = SelectCount + "from %s "
)
