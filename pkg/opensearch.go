package integrationtest

import (
	"context"
	"crypto/tls"
	"github.com/testcontainers/testcontainers-go"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2"
	container "github.com/testcontainers/testcontainers-go/modules/opensearch"
)

const defaultOpenSearchImage = "opensearchproject/opensearch:2.11.1"

type OpenSearchDependency struct {
	Dependency

	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *container.OpenSearchContainer
	client        *opensearch.Client
	env           map[string]string
}

func NewOpenSearchDependency(opts ...OpenSearchDependencyOpt) *OpenSearchDependency {
	dep := &OpenSearchDependency{
		image: defaultOpenSearchImage,
	}
	for _, opt := range opts {
		opt(dep)
	}
	return dep
}

type OpenSearchDependencyOpt func(*OpenSearchDependency)

func WithOpenSearchImage(image string) OpenSearchDependencyOpt {
	return func(dep *OpenSearchDependency) {
		dep.image = image
	}
}

func WithOpenSearchContainerOpts(opts ...testcontainers.ContainerCustomizer) OpenSearchDependencyOpt {
	return func(d *OpenSearchDependency) {
		d.containerOpts = opts
	}
}

func (osd *OpenSearchDependency) Start(ctx context.Context) error {
	c, err := container.Run(ctx, osd.image)
	if err != nil {
		return err
	}

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	osd.container = c
	endpoint, err := c.Address(ctx)
	if err != nil {
		return err
	}

	osd.env = map[string]string{
		"OPENSEARCH_ENDPOINT": endpoint,
		"OPENSEARCH_USERNAME": c.User,
		"OPENSEARCH_PASSWORD": c.Password,
	}

	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Addresses: []string{endpoint},
		Username:  c.User,
		Password:  c.Password,
	})
	if err != nil {
		return err
	}

	osd.client = client

	return nil
}

func (osd *OpenSearchDependency) Client() any {
	return osd.client
}

func (osd *OpenSearchDependency) Env() map[string]string {
	return osd.env
}

func (osd *OpenSearchDependency) Stop(ctx context.Context) error {
	if osd.container != nil {
		return osd.container.Terminate(ctx)
	}
	return nil
}
