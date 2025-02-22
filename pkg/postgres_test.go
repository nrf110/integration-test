package integrationtest

import (
	"context"
	"github.com/jackc/pgx/v5"
	integrationtest "github.com/nrf110/integration-test/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPostgresDependency(t *testing.T) {
	t.Run("can connect", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		t.Cleanup(cancel)

		pg := integrationtest.NewDependency(&integrationtest.Config{
			Database: "postgres",
			User:     "postgres",
			Password: "postgres",
		})
		err := pg.Start(ctx)
		assert.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, pg.Stop(ctx))
		})
		conn := pg.Client().(*pgx.Conn)
		assert.NoError(t, conn.Ping(ctx))
	})
}
