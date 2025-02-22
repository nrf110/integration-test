package integrationtest

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	integrationtest "github.com/nrf110/integration-test/pkg/elasticsearch"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestElasticsearchDependency(t *testing.T) {
	t.Run("can connect", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		t.Cleanup(cancel)

		es := integrationtest.NewDependency()
		assert.NoError(t, es.Start(ctx))
		t.Cleanup(func() {
			assert.NoError(t, es.Stop(ctx))
		})

		client := es.Client().(*elasticsearch.TypedClient)
		res, err := client.API.Indices.Create("test").Do(ctx)
		assert.NoError(t, err)
		assert.True(t, res.Acknowledged)
	})
}
