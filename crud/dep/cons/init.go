package cons

// default Mysql
var (
	Driver           = Mysql
	Backticks  uint8 = '`'    // different sql mark
	BackticksS       = "`%s`" // different sql mark
)

func Init(driver string) {
	switch driver {
	case Postgres:
		Driver = Postgres
		Backticks = '\''
		BackticksS = `"%s"`
	default:
	}
}
