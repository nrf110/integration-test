package gcs

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/testcontainers/testcontainers-go"
)

const defaultImage = "fsouza/fake-gcs-server:1.52.1"

type Dependency struct {
	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *Container
	client        *storage.Client
	env           map[string]string
}

func NewDependency(opts ...DependencyOpt) *Dependency {
	dep := &Dependency{
		image: defaultImage,
	}

	for _, opt := range opts {
		opt(dep)
	}

	return dep
}

type DependencyOpt func(dependency *Dependency)

func WithImage(image string) DependencyOpt {
	return func(dep *Dependency) {
		dep.image = image
	}
}

func WithContainerOpts(opts ...testcontainers.ContainerCustomizer) DependencyOpt {
	return func(d *Dependency) {
		d.containerOpts = append(d.containerOpts, opts...)
	}
}

func (dep *Dependency) Start(ctx context.Context) error {
	c, err := Run(ctx, defaultImage, dep.containerOpts...)
	if err != nil {
		return err
	}

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	dep.container = c

	hostPort, err := c.HostAndPort(ctx)
	if err != nil {
		return err
	}

	dep.env = map[string]string{
		"STORAGE_EMULATOR_HOST": hostPort,
	}

	os.Setenv("STORAGE_EMULATOR_HOST", hostPort)

	client, err := storage.NewClient(ctx, storage.WithJSONReads())

	if err != nil {
		return err
	}
	dep.client = client

	return nil
}

func (dep *Dependency) Client() any {
	return dep.client
}

func (dep *Dependency) Env() map[string]string {
	return dep.env
}

func (dep *Dependency) Stop(ctx context.Context) error {
	if dep.container != nil {
		err := dep.container.Terminate(ctx)
		if err != nil {
			log.Fatalf("failed to stop fake-gcs-server container: %v", err)
		}
	}
	return nil
}
