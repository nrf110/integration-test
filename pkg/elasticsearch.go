package integrationtest

import (
	"context"
	"crypto/tls"
	"github.com/testcontainers/testcontainers-go"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	container "github.com/testcontainers/testcontainers-go/modules/elasticsearch"
)

const defaultElasticsearchImage = "docker.elastic.co/elasticsearch/elasticsearch:8.9.0"

type ElasticsearchDependency struct {
	Dependency

	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *container.ElasticsearchContainer
	client        *elasticsearch.TypedClient
	env           map[string]string
}

func NewElasticsearchDependency(opts ...ElasticearchDependencyOpt) *ElasticsearchDependency {
	dep := &ElasticsearchDependency{
		image: defaultElasticsearchImage,
	}
	for _, opt := range opts {
		opt(dep)
	}
	return dep
}

type ElasticearchDependencyOpt func(dependency *ElasticsearchDependency)

func WithElasticsearchImage(image string) ElasticearchDependencyOpt {
	return func(dep *ElasticsearchDependency) {
		dep.image = image
	}
}

func WithElasticsearchContainerOpts(opts ...testcontainers.ContainerCustomizer) ElasticearchDependencyOpt {
	return func(d *ElasticsearchDependency) {
		d.containerOpts = opts
	}
}

func (osd *ElasticsearchDependency) Start(ctx context.Context) error {
	c, err := container.Run(ctx, osd.image)
	if err != nil {
		return err
	}

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	osd.container = c

	osd.env = map[string]string{
		"ELASTICSEARCH_ENDPOINT": c.Settings.Address,
		"ELASTICSEARCH_USERNAME": c.Settings.Username,
		"ELASTICSEARCH_PASSWORD": c.Settings.Password,
		"ELASTICSEARCH_CA_CERT":  string(c.Settings.CACert),
	}

	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
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

	osd.client = client

	return nil
}

func (osd *ElasticsearchDependency) Client() any {
	return osd.client
}

func (osd *ElasticsearchDependency) Env() map[string]string {
	return osd.env
}

func (osd *ElasticsearchDependency) Stop(ctx context.Context) error {
	if osd.container != nil {
		return osd.container.Terminate(ctx)
	}
	return nil
}
