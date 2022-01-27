// package cons

package cons

// devMode const
// key words
const (
	Dev  = "dev"
	Prod = "prod"
	// ConfPath default config path
	ConfPath = "conf/app.yaml"

	GtSubSQL              = "sub_sql"
	GtClientPage          = "clientPage"
	GtClientPageUnderLine = "client_page"
	GtEveryPage           = "everyPage"
	GtEveryPageUnderLine  = "every_page"
	GtOrder               = "order"
	GtKey                 = "key"
	GtMock                = "mock"
	// GT tag
	GT           = "gt"
	GtField      = "field"
	GtValid      = "valid"
	GtTrans      = "trans"
	GtIgnore     = "ignore"
	Gt_          = "-"
	GtComma      = ","
	GtGorm       = "gorm"
	GtGormColumn = "column"
)

// default page
const (
	ClientPage = 1
	EveryPage  = 10
)

// part sql
const (
	SQL_                = "SQL_"
	WhereS              = "where %s "
	AndS                = "and %s "
	OrderDesc           = "%s.id desc"
	OrderS              = "order by %s "
	ParamAnd            = " = ? and "
	ParamInAnd          = " in (?) and "
	SelectFrom          = "select %s from %s "
	Distinct            = "distinct "
	Count               = "count(*) as total_num"
	CountDistinct       = "count(distinct %s) as total_num"
	SelectCount         = "select " + Count + " "
	SelectCountDistinct = "select " + CountDistinct + " "
	SelectCountFrom     = SelectCount + "from %s "
)

var (
	Backticks uint8 = '`' // different sql mark
)
