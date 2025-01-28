package bigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/testcontainers/testcontainers-go/modules/gcloud"
)

// TODO: Need a version of this image that works on Apple Silicon
const defaultImage = "ghcr.io/goccy/bigquery-emulator:0.6.6"

type Dependency struct {
	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *gcloud.GCloudContainer
	client        *bigquery.Client
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
	c, err := gcloud.RunBigQuery(ctx, dep.image, dep.containerOpts...)
	if err != nil {
		return err
	}

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	dep.container = c

	projectID := c.Settings.ProjectID
	opts := []option.ClientOption{
		option.WithEndpoint(c.URI),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		option.WithoutAuthentication(),
		internaloption.SkipDialSettingsValidation(),
	}

	client, err := bigquery.NewClient(ctx, projectID, opts...)
	if err != nil {
		return err
	}
	dep.client = client

	dep.env = map[string]string{
		"GOOGLE_CLOUD_PROJECT": projectID,
	}

	return err
}

func (dep *Dependency) Env() map[string]string {
	return dep.env
}

func (dep *Dependency) Client() any {
	return dep.client
}

func (dep *Dependency) Stop(ctx context.Context) error {
	if dep.container != nil {
		return dep.container.Terminate(ctx)
	}
	return nil
}
