## Golang Database Test Kits

This library is a testing tool that runs various databases through docker.

### Installation
Only install the required database
```shell
go get github.com/dmateus/go-testdb/testmongo
go get github.com/dmateus/go-testdb/testcrdb
go get github.com/dmateus/go-testdb/testmysql
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
        WithTag("v21.2.4").
        WithMigrations(migrationsFolder).
        WithTest(t).
        MustStart().
        GetDB()
	
    // Run tests that use the database here
}
```

For more examples, check the tests in each of the database packages.

### API
```go
WithTag(tag string)                 -- Select the version of the database you want to run.
WithMigrations(migrationsFS fs.FS)  -- Runs the migration files in the given folder. Available in SQL databases.
WithTest(t *testing.T)              -- Instead of handling the database termination with `Stop()`, you can rely on `WithTest` to close it in the end. It also works with testify's Suites. 
```