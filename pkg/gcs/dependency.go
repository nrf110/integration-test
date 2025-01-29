package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"net/url"
	"os"
)

const defaultImage = "fsouza/fake-gcs-server:1"

type roundTripper url.URL

func (rt roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Host = rt.Host
	req.URL.Host = rt.Host
	req.URL.Scheme = rt.Scheme
	return http.DefaultTransport.RoundTrip(req)
}

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

	endpoint, err := c.Url(ctx)
	if err != nil {
		return err
	}

	dep.env = map[string]string{
		"STORAGE_EMULATOR_HOST": fmt.Sprintf("%s/storage/v1/", endpoint),
	}

	os.Setenv("STORAGE_EMULATOR_HOST", dep.env["STORAGE_EMULATOR_HOST"])

	//u, _ := url.Parse(endpoint)
	//httpClient := &http.Client{
	//	Transport: roundTripper(*u),
	//}
	client, err := storage.NewClient(ctx,
		option.WithoutAuthentication(),
		//option.WithHTTPClient(httpClient))
		option.WithEndpoint(fmt.Sprintf("%s/storage/v1/", endpoint)))
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
