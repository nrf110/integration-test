package integrationtest

import (
	"context"
	"github.com/nrf110/integration-test/pkg/elasticsearch"
	"github.com/nrf110/integration-test/pkg/gcs"
	"github.com/nrf110/integration-test/pkg/permify"
	"github.com/nrf110/integration-test/pkg/postgres"
	"github.com/nrf110/integration-test/pkg/pubsub"
	"github.com/nrf110/integration-test/pkg/redis"
	"github.com/pressly/goose/v3"
	"maps"
)

type TestSystem struct {
	env            map[string]string
	deps           []Dependency
	gooseProviders []GooseProviderFunc
	Elasticsearch  *elasticsearch.Dependency
	Redis          *redis.Dependency
	Postgres       *postgres.Dependency
	PubSub         *pubsub.Dependency
	GCS            *gcs.Dependency
	Permify        *permify.Dependency
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

func WithPubSub(opts ...pubsub.DependencyOpt) Option {
	return func(s *TestSystem) error {
		s.PubSub = pubsub.NewDependency(opts...)
		return WithDependency(s.PubSub)(s)
	}
}

func WithGCS(opts ...gcs.DependencyOpt) Option {
	return func(s *TestSystem) error {
		s.GCS = gcs.NewDependency(opts...)
		return WithDependency(s.GCS)(s)
	}
}

func WithPermify(opts ...permify.DependencyOpt) Option {
	return func(s *TestSystem) error {
		s.Permify = permify.NewDependency(opts...)
		return WithDependency(s.Permify)(s)
	}
}

type GooseProviderFunc func(s *TestSystem) (*goose.Provider, error)

func WithGooseProviders(providers ...GooseProviderFunc) Option {
	return func(s *TestSystem) error {
		s.gooseProviders = providers
		return nil
	}
}

func (sys *TestSystem) Start(ctx context.Context) error {
	for _, dep := range sys.deps {
		if err := dep.Start(ctx); err != nil {
			return err
		}
	}

	for _, providerFunc := range sys.gooseProviders {
		provider, err := providerFunc(sys)
		if err != nil {
			return err
		}

		if _, err = provider.Up(ctx); err != nil {
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
