package integrationtest

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/nrf110/integration-test/migrations"
	"github.com/nrf110/integration-test/pkg/postgres"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTestSystem(t *testing.T) {
	pgConfig := &postgres.Config{
		Database: "test",
		User:     "test",
		Password: "test",
	}

	t.Run("starts all dependencies", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		t.Cleanup(cancel)

		system, err := NewTestSystem(
			WithRedis(),
			WithElasticsearch(),
			WithPubSub(),
			WithGCS(),
			WithPostgres(pgConfig),
			WithPermify())
		assert.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, system.Stop(ctx))
		})

		err = system.Start(ctx)
		assert.NoError(t, err)
	})

	t.Run("runs migrations", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		t.Cleanup(cancel)

		system, err := NewTestSystem(
			WithPostgres(pgConfig),
			WithGooseProviders(func(s *TestSystem) (*goose.Provider, error) {
				conn := s.Postgres.Client().(*pgx.Conn)
				return migrations.NewSQLProvider(conn)
			}),
		)
		assert.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, system.Stop(ctx))
		})

		err = system.Start(ctx)
		assert.NoError(t, err)

		conn := system.Postgres.Client().(*pgx.Conn)

		row := conn.QueryRow(ctx, "select count(*) from test;")
		var count int
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}
