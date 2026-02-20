package bedatabase

import (
	"time"

	"gorm.io/gorm/logger"
)

type DbConfig struct {
	Host           string
	Port           int
	User           string
	Password       string
	Database       string
	SSLMode        string
	ConnectTimeout int // seconds
	PostgresConfig
	MysqlConfig
	MongoConfig
}

type PostgresConfig struct {
	Logger           logger.Interface
	ApplicationName  string
	StatementTimeout int           // seconds
	MaxIdleConns     int           // Maximum idle connections
	MaxOpenConns     int           // Maximum open connections
	ConnMaxLifetime  time.Duration // Maximum connection lifetime
}

type MysqlConfig struct {
	Logger          logger.Interface
	WriteTimeout    int           // seconds
	ReadTimeout     int           // seconds
	MaxIdleConns    int           // Maximum idle connections
	MaxOpenConns    int           // Maximum open connections
	ConnMaxLifetime time.Duration // Maximum connection lifetime
}

type MongoConfig struct {
	ValidatePing bool
	QueryParams  map[string]string
	Mode         MongoDbMode
}
