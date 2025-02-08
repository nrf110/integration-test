package permify

import (
	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	"context"
	"fmt"
	permifygrpc "github.com/Permify/permify-go/grpc"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/testcontainers/testcontainers-go"
	permifytest "github.com/theoriginalstove/testcontainers-permify"
)

const defaultImage = "ghcr.io/permify/permify:v1.2.3"

var defaultOpts = []DependencyOpt{
	WithContainerOpts(
		testcontainers.WithWaitStrategy(wait.ForLog("successfully started"))),
}

type Dependency struct {
	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *permifytest.PermifyContainer
	client        *permifygrpc.Client
	schema        string
	data          []*permifypayload.DataWriteRequest
	env           map[string]string
}

func NewDependency(opts ...DependencyOpt) *Dependency {
	allOpts := append(defaultOpts, opts...)

	dep := &Dependency{
		image: defaultImage,
	}
	for _, opt := range allOpts {
		opt(dep)
	}
	return dep
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

func WithSchema(schema string) DependencyOpt {
	return func(dep *Dependency) {
		dep.schema = schema
	}
}

func WithData(data ...*permifypayload.DataWriteRequest) DependencyOpt {
	return func(dep *Dependency) {
		dep.data = append(dep.data, data...)
	}
}

func (dep *Dependency) Start(ctx context.Context) error {
	c, err := permifytest.Run(ctx, dep.containerOpts...)
	if err != nil {
		return err
	}

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	dep.container = c

	host, err := c.Host(ctx)
	if err != nil {
		return err
	}

	grpcPort, err := c.GRPCPort(ctx)
	if err != nil {
		return err
	}

	cfg := permifygrpc.Config{
		Endpoint: fmt.Sprintf("%s:%d", host, grpcPort),
	}
	client, err := permifygrpc.NewClient(cfg, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	dep.client = client
	dep.env = map[string]string{
		"PERMIFY_ENDPOINT": cfg.Endpoint,
	}

	if dep.schema != "" {
		_, err = client.Schema.Write(ctx, &permifypayload.SchemaWriteRequest{
			TenantId: "t1",
			Schema:   dep.schema,
		})
		if err != nil {
			return err
		}

		for _, data := range dep.data {
			_, err = client.Data.Write(ctx, data)
			if err != nil {
				return err
			}
		}
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
		return dep.container.Terminate(ctx)
	}
	return nil
}
