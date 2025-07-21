package log

const (
	SuccessLevel = "success"
	DebugLevel   = "debug" // default level
	InfoLevel    = "info"
	WarnLevel    = "warn"
	ErrorLevel   = "error"
	InConsole    = "console"
	InFile       = "file"
	InAll        = "all"
)

// default log config
var (
	confLogLevel    = DebugLevel
	confLogDirector = "log"
	confLogMaxAge   = 7
	logIn           = InConsole
)
