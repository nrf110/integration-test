package integrationtest

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/gcloud"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const defaultPubSubImage = "gcr.io/google.com/cloudsdktool/google-cloud-cli:stable"

type PubSubDependency struct {
	Dependency

	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *gcloud.GCloudContainer
	env           map[string]string
	client        *pubsub.Client
}

var defaultContainerOpts = []testcontainers.ContainerCustomizer{
	testcontainers.WithEnv(map[string]string{
		"APT_PACKAGES": "curl python3-crcmod lsb-release gnupg bash",
		"COMPONENTS":   "google-cloud-cli-pubsub-emulator",
	}),
}

func NewPubSubDependency(opts ...testcontainers.ContainerCustomizer) *PubSubDependency {
	dep := &PubSubDependency{
		image:         defaultPubSubImage,
		containerOpts: append(defaultContainerOpts, opts...),
	}

	return dep
}

func (pub *PubSubDependency) Start(ctx context.Context) error {
	container, err := gcloud.RunPubsub(ctx, pub.image, pub.containerOpts...)
	if err != nil {
		return err
	}
	pub.container = container

	conn, err := grpc.NewClient(pub.container.URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	projectID := container.Settings.ProjectID
	options := []option.ClientOption{option.WithGRPCConn(conn)}
	client, err := pubsub.NewClient(ctx, projectID, options...)
	if err != nil {
		return err
	}

	pub.client = client
	pub.env = map[string]string{
		"GOOGLE_CLOUD_PROJECT": projectID,
	}

	return nil
}

func (pub *PubSubDependency) Client() any {
	return pub.client
}

func (pub *PubSubDependency) Env() map[string]string {
	return pub.env
}

func (pub *PubSubDependency) Stop(ctx context.Context) error {
	if pub.container != nil {
		err := pub.container.Terminate(ctx)
		if err != nil {
			log.Fatalf("failed to stop pubsub container: %v", err)
		}
	}
	return nil
}
