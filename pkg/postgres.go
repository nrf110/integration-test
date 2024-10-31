package integrationtest

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"strconv"

	"github.com/jackc/pgx/v5"
	container "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const defaultPostgresImage = "postgres:16"

type PostgresConfig struct {
	User     string
	Password string
	Database string
}

type PostgresDependency struct {
	Dependency

	config        *PostgresConfig
	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *container.PostgresContainer
	client        *pgx.Conn
	env           map[string]string
}

func NewPostgresDependency(config *PostgresConfig, opts ...PostgresDependencyOpt) *PostgresDependency {
	dep := &PostgresDependency{
		config: config,
		image:  defaultPostgresImage,
	}
	for _, opt := range opts {
		opt(dep)
	}
	return dep
}

type PostgresDependencyOpt func(*PostgresDependency)

func WithPostgresImage(image string) PostgresDependencyOpt {
	return func(dep *PostgresDependency) {
		dep.image = image
	}
}

func WithPostgresContainerOpts(opts ...testcontainers.ContainerCustomizer) PostgresDependencyOpt {
	return func(d *PostgresDependency) {
		d.containerOpts = opts
	}
}

func (pg *PostgresDependency) Start(ctx context.Context) error {
	c, err := container.Run(ctx, pg.image,
		container.WithDatabase(pg.config.Database),
		container.WithUsername(pg.config.User),
		container.WithPassword(pg.config.Password),
		container.WithSQLDriver("pgx"))

	if err != nil {
		return err
	}

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	pg.container = c

	connectionString, err := c.ConnectionString(ctx)
	if err != nil {
		return err
	}

	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return err
	}

	pg.env = map[string]string{
		"PG_HOST":     config.Host,
		"PG_PORT":     strconv.Itoa(int(config.Port)),
		"PG_DATABASE": pg.config.Database,
		"PG_USER":     pg.config.User,
		"PG_PASSWORD": pg.config.Password,
	}

	pg.client, err = pgx.Connect(ctx, connectionString)

	return err
}

func (pg *PostgresDependency) Env() map[string]string {
	return pg.env
}

func (pg *PostgresDependency) Client() any {
	return pg.client
}

func (pg *PostgresDependency) Stop(ctx context.Context) error {
	if pg.container != nil {
		return pg.container.Terminate(ctx)
	}
	return nil
}
