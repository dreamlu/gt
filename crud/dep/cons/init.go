package cons

// default Mysql
var (
	Driver                      = Mysql
	Backticks             uint8 = '`'    // different sql mark
	BackticksS                  = "`%s`" // different sql mark
	LikeKey                     = "like binary"
	GtClientPage                = "clientPage"
	GtClientPageUnderLine       = "client_page"
	GtEveryPage                 = "everyPage"
	GtEveryPageUnderLine        = "every_page"
)

func Init(driver string) {
	switch driver {
	case Postgres:
		Driver = Postgres
		//Backticks = '\''
		Backticks = '"'
		BackticksS = `"%s"`
		LikeKey = "like"
	default:
	}
}
