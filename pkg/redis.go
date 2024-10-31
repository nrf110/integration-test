package integrationtest

import (
	"context"
	"crypto/tls"
	"github.com/testcontainers/testcontainers-go"

	"github.com/redis/go-redis/v9"
	container "github.com/testcontainers/testcontainers-go/modules/redis"
)

const defaultRedisImage = "redis:7"

type RedisDependency struct {
	Dependency

	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *container.RedisContainer
	client        *redis.Client
	env           map[string]string
}

func NewRedisDependency(opts ...RedisDependencyOpt) *RedisDependency {
	dep := &RedisDependency{
		image: defaultRedisImage,
	}
	for _, opt := range opts {
		opt(dep)
	}
	return dep
}

type RedisDependencyOpt func(d *RedisDependency)

func WithRedisImage(image string) RedisDependencyOpt {
	return func(d *RedisDependency) {
		d.image = image
	}
}

func WithRedisContainerOpts(opts ...testcontainers.ContainerCustomizer) RedisDependencyOpt {
	return func(d *RedisDependency) {
		d.containerOpts = opts
	}
}

func (rd *RedisDependency) Start(ctx context.Context) error {
	c, err := container.Run(ctx, rd.image)
	if err != nil {
		return err
	}

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	rd.container = c

	addr, err := c.Endpoint(ctx, "")
	if err != nil {
		return err
	}

	rd.env = map[string]string{
		"REDIS_ADDRESS": addr,
	}

	rd.client = redis.NewClient(&redis.Options{
		Addr: addr,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	return nil
}

func (rd *RedisDependency) Client() any {
	return rd.client
}

func (rd *RedisDependency) Env() map[string]string {
	return rd.env
}

func (rd *RedisDependency) Stop(ctx context.Context) error {
	if rd.container != nil {
		return rd.container.Terminate(ctx)
	}
	return nil
}
