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
	tenantId      string
	schema        string
	tuples        []*permifypayload.Tuple
	env           map[string]string
}

func NewDependency(opts ...DependencyOpt) *Dependency {
	allOpts := append(defaultOpts, opts...)

	dep := &Dependency{
		image:    defaultImage,
		tenantId: "t1",
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

func WithTenantId(tenantId string) DependencyOpt {
	return func(dep *Dependency) {
		dep.tenantId = tenantId
	}
}

func WithTuples(tuples ...*permifypayload.Tuple) DependencyOpt {
	return func(dep *Dependency) {
		dep.tuples = append(dep.tuples, tuples...)
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
		schemaRes, err := client.Schema.Write(ctx, &permifypayload.SchemaWriteRequest{
			TenantId: dep.tenantId,
			Schema:   dep.schema,
		})
		if err != nil {
			return err
		}

		_, err = client.Data.Write(ctx, &permifypayload.DataWriteRequest{
			TenantId: "t1",
			Metadata: &permifypayload.DataWriteRequestMetadata{
				SchemaVersion: schemaRes.SchemaVersion,
			},
			Tuples: dep.tuples,
		})
		if err != nil {
			return err
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
