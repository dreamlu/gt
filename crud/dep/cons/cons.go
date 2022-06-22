// package cons

package cons

// devMode const
// key words
const (
	Dev  = "dev"
	Prod = "prod"

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
	GtLike       = "like"
	Gt_          = "-"
	GtComma      = ","
	GtGorm       = "gorm"
	GtGormColumn = "column"
	GtExist      = "exist"
	GtSoftDel    = "soft_del"
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
	And                 = " and "
	ParamAnd            = " = ?" + And
	ParamLike           = " like binary ? and "
	ParamInAnd          = " in (?) and "
	SelectFrom          = "select %s from %s "
	Distinct            = "distinct "
	Count               = "count(*) as total_num"
	CountDistinct       = "count(distinct %s) as total_num"
	SelectCount         = "select " + Count + " "
	SelectCountDistinct = "select " + CountDistinct + " "
	SelectCountFrom     = SelectCount + "from %s "
	SoftDel             = "%s.%s is not null and "
)
