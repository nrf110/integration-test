package integrationtest

import (
	"context"
	"maps"
)

type TestSystem struct {
	env        map[string]string
	deps       []Dependency
	OpenSearch *OpenSearchDependency
	Redis      *RedisDependency
	Postgres   *PostgresDependency
	Kafka      *KafkaDependency
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

func WithOpenSearch(opts ...OpenSearchDependencyOpt) Option {
	return func(s *TestSystem) error {
		s.OpenSearch = NewOpenSearchDependency(opts...)
		return WithDependency(s.OpenSearch)(s)
	}
}

func WithPostgres(config *PostgresConfig, opts ...PostgresDependencyOpt) Option {
	return func(s *TestSystem) error {
		s.Postgres = NewPostgresDependency(config, opts...)
		return WithDependency(s.Postgres)(s)
	}
}

func WithKafka(opts ...KafkaDependencyOpt) Option {
	return func(s *TestSystem) error {
		if kafka, err := NewKafkaDependency(opts...); err != nil {
			return err
		} else {
			s.Kafka = kafka
			return WithDependency(s.Kafka)(s)
		}
	}
}

func WithRedis(opts ...RedisDependencyOpt) Option {
	return func(s *TestSystem) error {
		s.Redis = NewRedisDependency()
		return WithDependency(s.Redis)(s)
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
