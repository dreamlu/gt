package cons

var (
	Backticks           uint8 = '`' // different sql mark
	DefaultDevMode            = "app.devMode"
	ConfDB                    = "app.db"
	ConfDBName                = "app.db.name"
	ConfRedis                 = "app.redis"
	ConfNsqProducerAddr       = "app.nsq.producer_addr"
	ConfNsqConsumerAddr       = "app.nsq.consumer_addr"
	ConfMongo                 = "app.mongo"
	ConfFile                  = "app.filepath"
	ConfTaskNum               = "app.daemon.task_num"
	ConfLogLevel              = "app.log.level"
)
