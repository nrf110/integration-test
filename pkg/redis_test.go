package integrationtest

import (
	"context"
	integrationtest "github.com/nrf110/integration-test/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisDependency(t *testing.T) {
	t.Run("can connect", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		t.Cleanup(cancel)
		r := integrationtest.NewDependency()
		assert.NoError(t, r.Start(ctx))
		t.Cleanup(func() {
			assert.NoError(t, r.Stop(ctx))
		})

		client := r.Client().(*redis.Client)
		status := client.Ping(ctx)
		assert.NoError(t, status.Err())
	})
}
