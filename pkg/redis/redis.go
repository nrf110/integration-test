package redis

import (
	"context"
	"github.com/testcontainers/testcontainers-go"

	"github.com/redis/go-redis/v9"
	container "github.com/testcontainers/testcontainers-go/modules/redis"
)

const defaultRedisImage = "redis:7"

type Dependency struct {
	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *container.RedisContainer
	client        *redis.Client
	env           map[string]string
}

func NewDependency(opts ...DependencyOpt) *Dependency {
	dep := &Dependency{
		image: defaultRedisImage,
	}
	for _, opt := range opts {
		opt(dep)
	}
	return dep
}

type DependencyOpt func(d *Dependency)

func WithImage(image string) DependencyOpt {
	return func(dep *Dependency) {
		dep.image = image
	}
}

func WithContainerOpts(opts ...testcontainers.ContainerCustomizer) DependencyOpt {
	return func(dep *Dependency) {
		dep.containerOpts = opts
	}
}

func (dep *Dependency) Start(ctx context.Context) error {
	c, err := container.Run(ctx, dep.image)
	if err != nil {
		return err
	}

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	dep.container = c

	url, err := c.ConnectionString(ctx)
	if err != nil {
		return err
	}

	options, err := redis.ParseURL(url)
	if err != nil {
		return err
	}
	dep.env = map[string]string{
		"REDIS_ADDRESS": options.Addr,
	}

	dep.client = redis.NewClient(options)

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
		return dep.container.Terminate(ctx)
	}
	return nil
}
