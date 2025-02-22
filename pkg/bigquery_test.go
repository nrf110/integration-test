package integrationtest

import (
	"cloud.google.com/go/bigquery"
	"context"
	integrationtest "github.com/nrf110/integration-test/pkg/bigquery"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/gcloud"
	"testing"
	"time"
)

type Item struct {
	ID   int    `bigquery:"id"`
	Name string `bigquery:"name"`
}

func TestBigQueryDependency(t *testing.T) {
	t.Run("can connect", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		t.Cleanup(cancel)
		bq := integrationtest.NewDependency(
			integrationtest.WithContainerOpts(
				gcloud.WithProjectID("test")))
		err := bq.Start(ctx)
		assert.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, bq.Stop(ctx))
		})
		client := bq.Client().(*bigquery.Client)

		ds := client.Dataset("test")
		err = ds.Create(ctx, &bigquery.DatasetMetadata{})
		assert.NoError(t, err)

		table := ds.Table("test")
		err = table.Create(ctx, &bigquery.TableMetadata{
			Schema: bigquery.Schema{
				&bigquery.FieldSchema{
					Name:     "id",
					Type:     bigquery.IntegerFieldType,
					Required: true,
				},
				&bigquery.FieldSchema{
					Name:     "name",
					Type:     bigquery.StringFieldType,
					Required: true,
				},
			},
		})
		assert.NoError(t, err)

		items := []*Item{
			{ID: 1, Name: "item1"},
			{ID: 2, Name: "item2"},
		}
		err = table.Inserter().Put(ctx, items)
		assert.NoError(t, err)

		it, err := client.Query("select * from `test`.`test`.`test` order by id").Read(ctx)
		assert.NoError(t, err)
		var row Item

		err = it.Next(&row)
		assert.NoError(t, err)
		assert.Equal(t, 1, row.ID)
		assert.Equal(t, "item1", row.Name)

		err = it.Next(&row)
		assert.NoError(t, err)
		assert.Equal(t, 2, row.ID)
		assert.Equal(t, "item2", row.Name)
	})
}
