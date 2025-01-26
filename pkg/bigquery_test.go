package integrationtest_test

import (
	"cloud.google.com/go/bigquery"
	bigquerycontainer "github.com/nrf110/integration-test/pkg/bigquery"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go/modules/gcloud"
)

var _ = Describe("Dependency", func() {
	It("can connect", func(ctx SpecContext) {
		Skip("Emulator image not yet available for arm64")
		bq := bigquerycontainer.NewDependency(
			bigquerycontainer.WithContainerOpts(
				gcloud.WithProjectID("test")))
		err := bq.Start(ctx)
		Expect(err).To(BeNil())
		defer func() {
			bq.Stop(ctx)
		}()
		client := bq.Client().(*bigquery.Client)
		it, err := client.Query("select 1;").Read(ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(it.TotalRows).To(Equal(1))
	})
})
