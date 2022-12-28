package testmongo

import (
	"context"
	"fmt"
	"github.com/dmateus/go-testdb/base"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

type Mongo struct {
	base.Base
	db *mongoDriver.Database
}

func NewMongo() *Mongo {
	m := Mongo{}
	m.DockerConfigs = &base.DockerConfigs{
		Image: "mongo",
		Tag:   "5.0",
		EnvVars: []string{
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_ROOT_PASSWORD=password",
		}}
	return &m
}

// WithTag Sets the image tag. Default: 5.0
func (m *Mongo) WithTag(tag string) *Mongo {
	m.DockerConfigs.Tag = tag
	return m
}

func (m *Mongo) WithTest(t *testing.T) *Mongo {
	t.Cleanup(func() {
		m.Stop()
	})
	return m
}

func (m *Mongo) Ping() error {
	return m.db.Client().Ping(context.Background(), nil)
}

func (m *Mongo) GetPort() string {
	return "27017"
}

func (m *Mongo) Connect(port string) error {
	client, err := mongoDriver.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://root:password@localhost:%s", port)))
	if err != nil {
		return err
	}
	m.db = client.Database("defaultdb")
	return nil
}

func (m *Mongo) MustStart() *Mongo {
	err := base.LaunchDocker(m)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func (m *Mongo) GetDB() *mongoDriver.Database {
	return m.db
}

func (m *Mongo) ResetDB() error {
	return m.db.Drop(context.Background())
}
