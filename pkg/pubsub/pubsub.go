package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/gcloud"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const defaultImage = "gcr.io/google.com/cloudsdktool/google-cloud-cli:stable"

var defaultPubSubContainerOpts = []testcontainers.ContainerCustomizer{
	testcontainers.WithEnv(map[string]string{
		"APT_PACKAGES": "curl python3-crcmod lsb-release gnupg bash apt-utils",
		"COMPONENTS":   "google-cloud-cli-pubsub-emulator",
	}),
	testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			WaitingFor: wait.ForLog("started").WithStartupTimeout(2 * time.Minute),
		},
	}),
}

type Dependency struct {
	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *gcloud.GCloudContainer
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
	dep := &Dependency{
		image:         defaultImage,
		containerOpts: defaultPubSubContainerOpts,
	}

	for _, opt := range opts {
		opt(dep)
	}

	return dep
}

func (dep *Dependency) Start(ctx context.Context) error {
	container, err := gcloud.RunPubsub(ctx, dep.image, dep.containerOpts...)
	if err != nil {
		return err
	}
	dep.container = container

	conn, err := grpc.NewClient(dep.container.URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	projectID := container.Settings.ProjectID
	options := []option.ClientOption{option.WithGRPCConn(conn)}
	client, err := pubsub.NewClient(ctx, projectID, options...)
	if err != nil {
		return err
	}

	dep.client = client
	dep.env = map[string]string{
		"GOOGLE_CLOUD_PROJECT": projectID,
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
