package base

import (
	"github.com/ory/dockertest/v3"
	"log"
)

type Base struct {
	resource      *dockertest.Resource
	DockerConfigs *DockerConfigs
}

func (b *Base) setResource(resource *dockertest.Resource) {
	b.resource = resource
}

func (b *Base) getResource() *dockertest.Resource {
	return b.resource
}

func (b *Base) GetDockerConfigs() *DockerConfigs {
	return b.DockerConfigs
}

func (b *Base) Stop() {
	err := StopDocker(b)
	if err != nil {
		log.Fatal(err)
	}
}
