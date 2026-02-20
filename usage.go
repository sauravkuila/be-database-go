package main

import (
	"fmt"

	"github.com/sauravkuila/be-database-go/bedatabase"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	////////////// PSQL CN //////////////
	psqlConfig := bedatabase.DbConfig{
		Host:           "localhost",
		Port:           5432,
		User:           "postgres",
		Password:       "postgres",
		Database:       "test",
		SSLMode:        "disable",
		ConnectTimeout: 5,
	}
	psqlConfig.ApplicationName = "usage_example"
	psqlConfig.StatementTimeout = 10
	psqlConfig.PostgresConfig.MaxIdleConns = 20
	psqlConfig.PostgresConfig.MaxOpenConns = 50
	//postgres
	psqlDbObj, err := bedatabase.ConnectPostgres(psqlConfig)
	if err != nil {
		panic(err)
	}
	psqlRows, err := psqlDbObj.Raw("SELECT 1").Rows()
	if err != nil {
		panic(err)
	}
	for psqlRows.Next() {
		var i int
		psqlRows.Scan(&i)
		println(i)
	}
	defer psqlRows.Close()

	//safely close connection
	psqlConn, _ := psqlDbObj.DB()
	psqlConn.Close()

	////////////// MONGO CN //////////////
	mongoConfig := bedatabase.DbConfig{
		Host:           "127.0.0.1",
		Port:           27017,
		User:           "mongo",
		Password:       "mongo",
		Database:       "admin",
		ConnectTimeout: 5,
		MongoConfig: bedatabase.MongoConfig{
			ValidatePing: true,
			QueryParams: map[string]string{
				"directconnection": "true",
				"authSource":       "admin",
				"directConnection": "true",
				"tls":              "true",
			},
			Mode: bedatabase.ModeTunnel,
		},
	}
	//mongo
	mongoDbObj, err := bedatabase.ConnectMongo(mongoConfig)
	if err != nil {
		panic(err)
	}
	dbNames, err := mongoDbObj.ListDatabaseNames(nil, bson.D{})
	if err != nil {
		panic(err)
	}
	fmt.Println("mongo db names: ", dbNames)
	//safely close connection
	mongoDbObj.Disconnect(nil)

	////////////// MYSQL CN //////////////
	mysqlConfig := bedatabase.DbConfig{
		Host:           "localhost",
		Port:           3306,
		User:           "root",
		Password:       "admin",
		Database:       "local",
		ConnectTimeout: 5,
	}
	mysqlConfig.ReadTimeout = 5
	mysqlConfig.WriteTimeout = 5
	mysqlConfig.MysqlConfig.MaxIdleConns = 20
	mysqlConfig.MysqlConfig.MaxOpenConns = 50
	//mysql
	mysqlDbObj, err := bedatabase.ConnectMySql(mysqlConfig)
	if err != nil {
		panic(err)
	}
	mysqlRows, err := mysqlDbObj.Raw("SELECT 1").Rows()
	if err != nil {
		panic(err)
	}
	for mysqlRows.Next() {
		var i int
		mysqlRows.Scan(&i)
		println(i)
	}
	defer mysqlRows.Close()

	//safely close connection
	mysqlConn, _ := mysqlDbObj.DB()
	mysqlConn.Close()
}
