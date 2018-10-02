package docker

import (
	"github.com/drone/drone-runtime/engine"

	"docker.io/go-docker"
)

type dockerFactory struct {
	client docker.APIClient
}

// NewEnv returns a new Factory that creates Docker-runtime Engines
// from the environment.
func NewEnv() (engine.Factory, error) {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return New(cli), nil
}

// New returns a new Factory that creates Docker-runtime Engines.
func New(client docker.APIClient) engine.Factory {
	return &dockerFactory{
		client: client,
	}
}

func (f *dockerFactory) Create(spec *engine.Spec) engine.Engine {
	return &dockerEngine{
		spec:   spec,
		client: f.client,
	}
}
