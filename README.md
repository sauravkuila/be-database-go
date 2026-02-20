# README #

Database connection management provides functionality to connect to a database.

* **v1.0.0**: 
- Basic PostgreSQL connection wrapper. Using `ConnectPostgres()` to establish a connection to the database.
- Added support for Mongo connection. Updated yaml config pull to custom DB model.
- Updated config to allow custom logger and ping check for databases.
- Added support for MySQL connection.
- Added support for statement timeout in psql, write/read/connect timeout in mysql. updated default connection pool for idle and open connections with configurable options.
- Added support for query params in mongo with connection modes.


* **Usage**:
```go
    psqlConfig := bedatabase.DbConfig{
		Host:           "localhost",
		Port:           5432,
		User:           "postgres",
		Password:       "postgres",
		Database:       "dev",
		SSLMode:        "disable",
		ConnectTimeout: 5,
		PostgresConfig: bedatabase.PostgresConfig{
			ApplicationName:  "test",
			StatementTimeout: 5,
			MaxIdleConns:     20,
			MaxOpenConns:     50,
			ConnMaxLifetime:  0,
		},
	}
	//postgres
	bedatabase.ConnectPostgres(psqlConfig)

    mongoConfig := bedatabase.DbConfig{
		Host:           "localhost",
		Port:           27017,
		User:           "admin",
		Password:       "admin",
		Database:       "dev",
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
	bedatabase.ConnectMongo(mongoConfig)

	
	mysqlConfig := bedatabase.DbConfig{
		Host:           "localhost",
		Port:           3306,
		User:           "root",
		Password:       "admin",
		Database:       "local",
		ConnectTimeout: 5,
		MysqlConfig: bedatabase.MysqlConfig{
			WriteTimeout: 5,
			ReadTimeout:  5,
			MaxIdleConns: 20,
			MaxOpenConns: 50,
		},
	}
	//mysql
	bedatabase.ConnectMySql(mysqlConfig)
```

### How do I contribute? ###

* branch out from `main`
    - `git checkout main`
    - `git pull`
    - `git checkout -b <branchname>`
* branch for hotfixes or older versions
    - if you need to patch older versions, create a branch from the tagged commit
    - `git checkout -b <newer_version> <older_version>`
    - eg. `git checkout -b v1.0.x v1.0.0`
* add changes to repo
* raise a PR to `main` branch
* once merged to `main`, make a release tag based on contribution. ask admin to tag the change if main changes are not accesible
    - `git tag <version>`
    - `git push origin <version>`
* check tag is created
    - `git ls-remote --tags origin`

### Who do I talk to? ###

* @sauravkuila