## Golang Database Test Kits

This library is a testing tool that runs various databases through docker.

### Installation
Only install the required database
```shell
go get github.com/dmateus/go-testdb/mongo
go get github.com/dmateus/go-testdb/cockroachdb
```

### Mongo Usage
```go
import "github.com/dmateus/go-testdb/mongo"

func TestSomething(t *testing.T) {
    c := NewMongo().
        WithTag("5.0").
        MustStart()
	defer c.Stop()
	
    // Run tests that use the database here
    db := c.Database("my-database")
}
```

### API
```go
WithTag(tag string)                 -- Select the version of the database you want to run.
WithMigrations(migrationsFS fs.FS)  -- Runs the migration files in the given folder. Available in SQL databases.
WithTest(t *testing.T)              -- Instead of handling the database termination, you can rely on `WithTest` to close it in the end.
```
