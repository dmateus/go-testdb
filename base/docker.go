package base

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"sync"
)

var (
	pool *dockertest.Pool
	once sync.Once
)

type Database interface {
	Connect(port string) error
	GetDockerConfigs() *DockerConfigs
	Ping() error
	GetPort() string
	resource
}

func GetPool() *dockertest.Pool {
	once.Do(func() {
		var err error
		pool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("could not construct pool: %s", err)
		}

		// uses pool to try to connect to Docker
		err = pool.Client.Ping()
		if err != nil {
			log.Fatalf("Could not connect to Docker: %s", err)
		}
	})
	return pool
}

func LaunchDocker(db Database) error {
	pool := GetPool()

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: db.GetDockerConfigs().Image,
		Tag:        db.GetDockerConfigs().Tag,
		Env:        db.GetDockerConfigs().EnvVars,
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return fmt.Errorf("could not start resource: %s", err.Error())
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		err := db.Connect(resource.GetPort(fmt.Sprintf("%s/tcp", db.GetPort())))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return fmt.Errorf("could not connect to database: %s", err.Error())
	}

	_ = resource.Expire(60 * 5) // in seconds

	db.setResource(resource)

	return nil
}

type resource interface {
	setResource(resource *dockertest.Resource)
	getResource() *dockertest.Resource
}

func StopDocker(r resource) error {
	if err := pool.Purge(r.getResource()); err != nil {
		return fmt.Errorf("could not purge resource: %s", err.Error())
	}
	return nil
}
