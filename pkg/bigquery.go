package integrationtest

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
const defaultBiqQueryImage = "ghcr.io/goccy/bigquery-emulator:0.6.6"

type BigQueryDependency struct {
	Dependency

	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *gcloud.GCloudContainer
	client        *bigquery.Client
	env           map[string]string
}

func NewBigQueryDependency(opts ...BigQueryDependencyOpt) *BigQueryDependency {
	dep := &BigQueryDependency{
		image: defaultBiqQueryImage,
	}
	for _, opt := range opts {
		opt(dep)
	}
	return dep
}

type BigQueryDependencyOpt func(dependency *BigQueryDependency)

func WithBigQueryImage(image string) BigQueryDependencyOpt {
	return func(dep *BigQueryDependency) {
		dep.image = image
	}
}

func WithBigQueryContainerOpts(opts ...testcontainers.ContainerCustomizer) BigQueryDependencyOpt {
	return func(d *BigQueryDependency) {
		d.containerOpts = append(d.containerOpts, opts...)
	}
}

func (bq *BigQueryDependency) Start(ctx context.Context) error {
	c, err := gcloud.RunBigQuery(ctx, bq.image, bq.containerOpts...)
	if err != nil {
		return err
	}

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	bq.container = c

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
	bq.client = client

	bq.env = map[string]string{
		"GOOGLE_CLOUD_PROJECT": projectID,
	}

	return err
}

func (bq *BigQueryDependency) Env() map[string]string {
	return bq.env
}

func (bq *BigQueryDependency) Client() any {
	return bq.client
}

func (bq *BigQueryDependency) Stop(ctx context.Context) error {
	if bq.container != nil {
		return bq.container.Terminate(ctx)
	}
	return nil
}
