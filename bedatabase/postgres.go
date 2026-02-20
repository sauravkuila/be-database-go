package bedatabase

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// expected db connection info in config
// databases:
//   postgres:
//     host: localhost
//     port: 5432
//     user: postgres
//     password: postgres
//     db: test
//     sslmode: disable
//     connect_timeout: 5

func ConnectPostgres(config DbConfig) (*gorm.DB, error) {

	if config.StatementTimeout <= 0 {
		config.StatementTimeout = DEFAULT_STATEMENT_TIMEOUT
	}
	if config.ConnectTimeout <= 0 {
		config.ConnectTimeout = DEFAULT_CONNECT_TIMEOUT
	}
	// Init DB connection string params from config
	dsn := fmt.Sprintf("host=%s port=%d user=%s password='%s' dbname=%s sslmode=%s connect_timeout=%d statement_timeout=%d application_name=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		config.SSLMode,
		config.ConnectTimeout,
		config.StatementTimeout*1000, // convert seconds to milliseconds
		config.ApplicationName,
	)

	option := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true, /// Removed defuault transaction used by GORM to faster the query and execution
		// Logger:                 loggerObj,
	}
	if config.PostgresConfig.Logger == nil {
		fmt.Println("using default logger")
		option.Logger = logger.Default
	} else {
		option.Logger = config.PostgresConfig.Logger
	}

	// Connect to database
	dbc, err := gorm.Open(postgres.Open(dsn), &option)
	if err != nil {
		fmt.Printf("Failed to connect to database. err %v. dsn: %s\n", err, dsn)
		return nil, err
	}

	sqlDb, err := dbc.DB()
	if err != nil {
		fmt.Printf("Failed to get sql db. err %v. dsn: %s\n", err, dsn)
		return nil, err
	}

	if config.PostgresConfig.MaxIdleConns <= 0 {
		config.PostgresConfig.MaxIdleConns = DEFAULT_IDLE_CONNECTIONS
	}
	if config.PostgresConfig.MaxOpenConns <= 0 {
		config.PostgresConfig.MaxOpenConns = DEFAULT_OPEN_CONNECTIONS
	}
	sqlDb.SetMaxIdleConns(config.PostgresConfig.MaxIdleConns)
	sqlDb.SetMaxOpenConns(config.PostgresConfig.MaxOpenConns)

	log.Println("postgres database connected")
	return dbc, nil
}
