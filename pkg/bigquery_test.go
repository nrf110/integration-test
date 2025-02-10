package integrationtest_test

import (
	"cloud.google.com/go/bigquery"
	integrationtest "github.com/nrf110/integration-test/pkg/bigquery"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go/modules/gcloud"
)

type Item struct {
	ID   int    `bigquery:"id"`
	Name string `bigquery:"name"`
}

var _ = Describe("bigquery.Dependency", func() {
	It("can connect", func(ctx SpecContext) {
		bq := integrationtest.NewDependency(
			integrationtest.WithContainerOpts(
				gcloud.WithProjectID("test")))
		err := bq.Start(ctx)
		Expect(err).To(BeNil())
		defer func() {
			bq.Stop(ctx)
		}()
		client := bq.Client().(*bigquery.Client)

		ds := client.Dataset("test")
		err = ds.Create(ctx, &bigquery.DatasetMetadata{})
		Expect(err).To(BeNil())

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
		Expect(err).To(BeNil())

		items := []*Item{
			{ID: 1, Name: "item1"},
			{ID: 2, Name: "item2"},
		}
		err = table.Inserter().Put(ctx, items)
		Expect(err).To(BeNil())

		it, err := client.Query("select * from `test`.`test`.`test` order by id").Read(ctx)
		Expect(err).ToNot(HaveOccurred())
		var row Item

		err = it.Next(&row)
		Expect(err).To(BeNil())
		Expect(row.ID).To(Equal(1))
		Expect(row.Name).To(Equal("item1"))

		err = it.Next(&row)
		Expect(err).To(BeNil())
		Expect(row.ID).To(Equal(2))
		Expect(row.Name).To(Equal("item2"))
	})
})
