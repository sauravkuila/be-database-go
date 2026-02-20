package bedatabase

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func ConnectMySql(config DbConfig) (*gorm.DB, error) {

	if config.ReadTimeout <= 0 {
		config.StatementTimeout = DEFAULT_READ_TIMEOUT
	}
	if config.WriteTimeout <= 0 {
		config.StatementTimeout = DEFAULT_WRITE_TIMEOUT
	}
	if config.ConnectTimeout <= 0 {
		config.ConnectTimeout = DEFAULT_CONNECT_TIMEOUT
	}
	// Init DB connection string params from config
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%ds&readTimeout=%ds&writeTimeout=%ds",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.ConnectTimeout,
		config.ReadTimeout,
		config.WriteTimeout,
	)

	option := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true, /// Removed defuault transaction used by GORM to faster the query and execution
	}
	if config.MysqlConfig.Logger == nil {
		fmt.Println("using default logger")
		option.Logger = logger.Default
	} else {
		option.Logger = config.MysqlConfig.Logger
	}

	// Connect to database
	dbc, err := gorm.Open(mysql.Open(dsn), &option)
	if err != nil {
		fmt.Printf("Failed to connect to database. err %v. dsn: %s\n", err, dsn)
		return nil, err
	}
	sqlDb, err := dbc.DB()
	if err != nil {
		fmt.Printf("Failed to get sql db. err %v. dsn: %s\n", err, dsn)
		return nil, err
	}

	if config.MysqlConfig.MaxIdleConns <= 0 {
		config.MysqlConfig.MaxIdleConns = DEFAULT_IDLE_CONNECTIONS
	}
	if config.MysqlConfig.MaxOpenConns <= 0 {
		config.MysqlConfig.MaxOpenConns = DEFAULT_OPEN_CONNECTIONS
	}
	sqlDb.SetMaxIdleConns(config.MysqlConfig.MaxIdleConns)
	sqlDb.SetMaxOpenConns(config.MysqlConfig.MaxOpenConns)

	log.Println("mysql database connected")
	return dbc, nil
}
