package bedatabase

const (
	DEFAULT_CONNECT_TIMEOUT   = 3 // value in seconds
	DEFAULT_STATEMENT_TIMEOUT = 5 // value in seconds
	DEFAULT_READ_TIMEOUT      = 5 // value in seconds
	DEFAULT_WRITE_TIMEOUT     = 5 // value in seconds
	DEFAULT_OPEN_CONNECTIONS  = 10
	DEFAULT_IDLE_CONNECTIONS  = 10
)

type MongoDbMode string

const (
	ModeTunnel  MongoDbMode = "tunnel"
	ModeAtlas   MongoDbMode = "direct"
	ModePrivate MongoDbMode = "private"
)
