![test workflow](https://github.com/dmateus/go-testdb/actions/workflows/test.yml/badge.svg)
## Golang Database Test Kits

This library is a testing tool that runs various databases through docker.

### Databases
- [PostgreSQL](testpostgres)
- [MySQL](testmysql)
- [CockroachDB](testcrdb)
- [Mongo](testmongo)

### Why use this library
- Launching test databases with code simplifies the process of running tests.
- Make use a useful utils like: run migrations, disconnect database automatically and reset database.
- If the test crashes, the database container will still be stopped and removed.

### Installation
```shell
go get github.com/dmateus/go-testdb/testcrdb
go get github.com/dmateus/go-testdb/testmysql
go get github.com/dmateus/go-testdb/testpostgres
go get github.com/dmateus/go-testdb/testmongo
```

### CockroachDB Usage
```go
import (
    "embed"
    "github.com/dmateus/go-testdb/testcrdb"
    "testing"
)

//go:embed migrations
var migrationsFolder embed.FS

func TestSomething(t *testing.T) {
    db := testcrdb.NewCockroachDB().
        WithMigrations(migrationsFolder).
        WithTest(t).
        MustStart().
        GetDB()
	
    // Run tests that use the database here
}
```

For more examples, check the tests in each of the database packages.

### API
```
Stop()                              -- Stops and removes the docker container.
WithTag(tag string)                 -- Select the docker tag you want to run.
WithMigrations(migrationsFS fs.FS)  -- Runs the migration files in the given folder. Available in SQL databases.
WithTest(t *testing.T)              -- Instead of handling the database termination with `Stop()`, you can rely on `WithTest` to close it in the end. It also works with testify's Suites.
MustStart()                         -- Starts the docker container.
GetDB()                             -- Returns a database from the launched container.
ResetDB()                           -- Cleans the database.
```