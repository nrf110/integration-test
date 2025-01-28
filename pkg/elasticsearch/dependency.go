package elasticsearch

import (
	"context"
	"crypto/tls"
	"github.com/testcontainers/testcontainers-go"
	"net/http"

	es "github.com/elastic/go-elasticsearch/v8"
	container "github.com/testcontainers/testcontainers-go/modules/elasticsearch"
)

const defaultImage = "docker.elastic.co/elasticsearch/elasticsearch:8.9.0"

type Dependency struct {
	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *container.ElasticsearchContainer
	client        *es.TypedClient
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
	return func(dep *Dependency) {
		dep.containerOpts = append(dep.containerOpts, opts...)
	}
}

func (dep *Dependency) Start(ctx context.Context) error {
	c, err := container.Run(ctx, dep.image, dep.containerOpts...)
	if err != nil {
		return err
	}

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	dep.container = c

	dep.env = map[string]string{
		"ELASTICSEARCH_ENDPOINT": c.Settings.Address,
		"ELASTICSEARCH_USERNAME": c.Settings.Username,
		"ELASTICSEARCH_PASSWORD": c.Settings.Password,
		"ELASTICSEARCH_CA_CERT":  string(c.Settings.CACert),
	}

	client, err := es.NewTypedClient(es.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Addresses: []string{c.Settings.Address},
		Username:  c.Settings.Username,
		Password:  c.Settings.Password,
		CACert:    c.Settings.CACert,
	})
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
		return dep.container.Terminate(ctx)
	}
	return nil
}
