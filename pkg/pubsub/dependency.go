package pubsub

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/testcontainers/testcontainers-go"
	tcpubsub "github.com/testcontainers/testcontainers-go/modules/gcloud/pubsub"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Dependency struct {
	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *tcpubsub.Container
	env           map[string]string
	client        *pubsub.Client
}

type DependencyOpt func(d *Dependency)

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

func NewDependency(opts ...DependencyOpt) *Dependency {
	dep := &Dependency{}

	for _, opt := range opts {
		opt(dep)
	}

	return dep
}

func (dep *Dependency) Start(ctx context.Context) error {
	container, err := tcpubsub.Run(ctx, "gcr.io/google.com/cloudsdktool/cloud-sdk:367.0.0-emulators", dep.containerOpts...)
	if err != nil {
		return err
	}
	container.Start(ctx)
	dep.container = container

	conn, err := grpc.NewClient(dep.container.URI(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	projectID := container.ProjectID()
	options := []option.ClientOption{option.WithGRPCConn(conn)}
	client, err := pubsub.NewClient(ctx, projectID, options...)
	if err != nil {
		return err
	}

	dep.client = client
	dep.env = map[string]string{
		"GOOGLE_CLOUD_PROJECT": projectID,
		"PUBSUB_EMULATOR_HOST": container.URI(),
	}

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
			log.Fatalf("failed to stop pubsub container: %v", err)
		}
	}
	return nil
}
