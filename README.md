## Golang Database Test Kits

This library is a testing tool that runs various databases through docker.

### Installation
```shell
go get github.com/StreamElements/cockroachdb-test-kit
```

### Mongo Usage
```go
import "github.com/dmateus/go-testdb/mongo"

func TestSomething(t *testing.T) {
    client, err := NewMongo().Start()
    require.NoError(t, err)
    defer m.Stop()
	
	// Run tests that use the database here
	db := client.Database("my-database")
}
```

### API
```go
WithTag(tag string)                 -- Select the version of the database you want to run.
WithMigrations(migrationsFS fs.FS)  -- Runs the migration files in the given folder. Available in SQL databases.
```