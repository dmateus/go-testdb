package mongo

import (
	"context"
	"fmt"
	"github.com/dmateus/go-testdb/base"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

type mongo struct {
	base.Base
	client *mongoDriver.Client
}

func NewMongo() *mongo {
	m := mongo{}
	m.DockerConfigs = base.NewDockerConfigs("mongo", "5.0", []string{
		"MONGO_INITDB_ROOT_USERNAME=root",
		"MONGO_INITDB_ROOT_PASSWORD=password",
	})
	return &m
}

// WithTag Sets the image tag. Default: 5.0
func (m *mongo) WithTag(tag string) *mongo {
	m.DockerConfigs.Tag = tag
	return m
}

func (m *mongo) WithTest(t *testing.T) *mongo {
	t.Cleanup(func() {
		m.Stop()
	})
	return m
}

func (m *mongo) Ping() error {
	return m.client.Ping(context.Background(), nil)
}

func (m *mongo) GetPort() string {
	return "27017"
}

func (m *mongo) Connect(port string) error {
	client, err := mongoDriver.Connect(
		context.Background(),
		options.Client().ApplyURI(
			fmt.Sprintf("mongodb://root:password@localhost:%s", port),
		),
	)
	if err != nil {
		return err
	}
	m.client = client
	return nil
}

func (m *mongo) MustStart() *mongoDriver.Client {
	err := base.LaunchDocker(m)
	if err != nil {
		log.Fatal(err)
	}
	return m.client
}
