package integrationtest

import (
	"context"
	"github.com/nrf110/integration-test/pkg/elasticsearch"
	"github.com/nrf110/integration-test/pkg/postgres"
	"github.com/nrf110/integration-test/pkg/pubsub"
	"github.com/nrf110/integration-test/pkg/redis"
	"github.com/testcontainers/testcontainers-go"
	"maps"
)

type TestSystem struct {
	env           map[string]string
	deps          []Dependency
	Elasticsearch *elasticsearch.Dependency
	Redis         *redis.Dependency
	Postgres      *postgres.Dependency
	PubSub        *pubsub.Dependency
}

type Option func(s *TestSystem) error

func NewTestSystem(options ...Option) (*TestSystem, error) {
	ts := &TestSystem{}
	for _, option := range options {
		if err := option(ts); err != nil {
			return nil, err
		}
	}
	return ts, nil
}

func WithDependency(dep Dependency) Option {
	return func(s *TestSystem) error {
		s.deps = append(s.deps, dep)
		maps.Copy(s.env, dep.Env())
		return nil
	}
}

func WithElasticsearch(opts ...elasticsearch.DependencyOpt) Option {
	return func(s *TestSystem) error {
		s.Elasticsearch = elasticsearch.NewDependency(opts...)
		return WithDependency(s.Elasticsearch)(s)
	}
}

func WithPostgres(config *postgres.Config, opts ...postgres.DependencyOpt) Option {
	return func(s *TestSystem) error {
		s.Postgres = postgres.NewDependency(config, opts...)
		return WithDependency(s.Postgres)(s)
	}
}

func WithRedis(opts ...redis.DependencyOpt) Option {
	return func(s *TestSystem) error {
		s.Redis = redis.NewDependency()
		return WithDependency(s.Redis)(s)
	}
}

func WithPubSub(opts ...testcontainers.ContainerCustomizer) Option {
	return func(s *TestSystem) error {
		s.PubSub = pubsub.NewDependency(opts...)
		return WithDependency(s.PubSub)(s)
	}
}

func (sys *TestSystem) Start(ctx context.Context) error {
	for _, dep := range sys.deps {
		if err := dep.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (sys *TestSystem) Stop(ctx context.Context) error {
	for _, dep := range sys.deps {
		if err := dep.Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (sys *TestSystem) Run(ctx context.Context) (err error) {
	if err = sys.Start(ctx); err != nil {
		return err
	}

	if err = sys.Stop(ctx); err != nil {
		return err
	}
	return nil
}
