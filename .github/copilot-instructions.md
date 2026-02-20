# Copilot Instructions: be-database-go

## Project Overview
A Go library providing unified database connection management for PostgreSQL, MySQL, and MongoDB with production-ready defaults. Part of the `github.com/sauravkuila/be-database-go` module.

## Architecture

### Core Components
- **bedatabase/model.go**: Central `DbConfig` struct with embedded database-specific configs (`PostgresConfig`, `MysqlConfig`, `MongoConfig`)
- **bedatabase/postgres.go**, **mysql.go**: Return `*gorm.DB` instances with GORM ORM integration
- **bedatabase/mongo.go**: Returns native `*mongo.Client` with mode-based connection strategies
- **bedatabase/constant.go**: Default values for timeouts and connection pools

### Key Design Patterns

#### 1. Embedded Configuration Pattern
```go
type DbConfig struct {
    Host, Port, User, Password, Database string
    ConnectTimeout int
    PostgresConfig  // embedded
    MysqlConfig     // embedded
    MongoConfig     // embedded
}
```
Use embedded structs to access database-specific settings like `config.MaxIdleConns` (Postgres/MySQL) or `config.Mode` (Mongo).

#### 2. MongoDB Connection Modes (mongo.go)
Three distinct modes handle different deployment scenarios:
- **ModeTunnel**: SSH/SSM tunnels to localhost with TLS skip verify, direct connection
- **ModeAtlas**: `mongodb+srv://` for Atlas clusters with majority write concern
- **ModePrivate**: VPC/peering setups with custom query params

Mode selection drives URI format, client options, and default QueryParams. Each mode strips specific params from `cfg.QueryParams` before applying them programmatically (e.g., `directConnection`, `retryWrites`, `tls`).

#### 3. Default Value Injection
All `Connect*` functions check for zero values and apply defaults from constants:
```go
if config.ConnectTimeout <= 0 {
    config.ConnectTimeout = DEFAULT_CONNECT_TIMEOUT
}
```
Always use `<= 0` checks before setting defaults to avoid overriding user-provided values.

#### 4. GORM Configuration Standards
Both Postgres and MySQL use:
- `SingularTable: true` (no plural table names)
- `SkipDefaultTransaction: true` (performance optimization)
- Custom logger support via `config.Logger` or fallback to `logger.Default`

## Critical Workflows

### Adding a New Database Type
1. Create `bedatabase/<dbtype>.go` with `Connect<DbType>(DbConfig) (*client, error)` signature
2. Add `<DbType>Config` struct to `model.go` and embed in `DbConfig`
3. Implement timeout defaults and connection pooling logic
4. Add example to `usage.go` demonstrating full connection lifecycle (connect → query → close)

### Testing Database Connections
No test files exist. When adding tests:
- Mock database connections or use testcontainers
- Test default value injection separately from connection logic
- Verify QueryParams handling in Mongo modes (especially param deletion behavior)

### Version Tagging
Follow strict git tag workflow from README:
1. Branch from `main`, merge PR
2. Tag releases as `v<major>.<minor>.<patch>` matching version history in README
3. Push tags: `git push origin <version>`
4. Verify: `git ls-remote --tags origin`

## Common Conventions

### Error Handling
- Wrap errors with context: `fmt.Errorf("mongo ping error: %w", err)`
- Log connection success: `log.Println("postgres database connected")`
- Print detailed failure info: `fmt.Printf("Failed to connect. err %v. dsn: %s\n", err, dsn)`

### Connection Cleanup
Always demonstrate proper cleanup in examples:
```go
sqlDb, _ := gormDb.DB()
defer sqlDb.Close()
// Or for Mongo:
defer mongoClient.Disconnect(nil)
```

### Query Parameter Handling (Mongo)
`prepareUriWithParams()` appends map as query string. Mode-specific logic removes params that become client options (e.g., `directConnection` → `SetDirect()`). Check mode implementation before adding new params.

## File References
- Connection entry points: [postgres.go](bedatabase/postgres.go), [mysql.go](bedatabase/mysql.go), [mongo.go](bedatabase/mongo.go)
- Configuration model: [model.go](bedatabase/model.go)
- Defaults: [constant.go](bedatabase/constant.go)
- Usage examples: [usage.go](usage.go)
