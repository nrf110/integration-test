package postgres

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	container "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const defaultPostgresImage = "postgres:16"

type Config struct {
	User     string
	Password string
	Database string
}

type Dependency struct {
	config        *Config
	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *container.PostgresContainer
	client        *pgx.Conn
	env           map[string]string
}

func NewDependency(config *Config, opts ...DependencyOpt) *Dependency {
	dep := &Dependency{
		config: config,
		image:  defaultPostgresImage,
	}
	for _, opt := range opts {
		opt(dep)
	}
	return dep
}

type DependencyOpt func(*Dependency)

func WithImage(image string) DependencyOpt {
	return func(dep *Dependency) {
		dep.image = image
	}
}

func WithContainerOpts(opts ...testcontainers.ContainerCustomizer) DependencyOpt {
	return func(d *Dependency) {
		d.containerOpts = opts
	}
}

func (dep *Dependency) Start(ctx context.Context) error {
	c, err := container.Run(ctx, dep.image,
		container.WithDatabase(dep.config.Database),
		container.WithUsername(dep.config.User),
		container.WithPassword(dep.config.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
		container.WithSQLDriver("pgx"))

	if err != nil {
		return err
	}

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	dep.container = c

	connectionString, err := c.ConnectionString(ctx)
	if err != nil {
		return err
	}

	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return err
	}

	dep.env = map[string]string{
		"PG_HOST":     config.Host,
		"PG_PORT":     strconv.Itoa(int(config.Port)),
		"PG_DATABASE": dep.config.Database,
		"PG_USER":     dep.config.User,
		"PG_PASSWORD": dep.config.Password,
	}

	dep.client, err = pgx.Connect(ctx, connectionString)

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
