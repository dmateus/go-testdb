## Postgres

This library is a testing tool that runs postgres through docker.

### Installation
```shell
go get github.com/dmateus/go-testdb/testcrdb
```

### Usage
```go
import (
    "embed"
    "testing"
    "github.com/dmateus/go-testdb/testcrdb"
)

//go:embed migrations
var migrationsFolder embed.FS

func TestSomething(t *testing.T) {
    db := testpostgres.NewPostgres().
        WithMigrations(migrationsFolder).
        WithTest(t).
        MustStart().
        GetDB()

    // Run tests that use the database here
}
```
