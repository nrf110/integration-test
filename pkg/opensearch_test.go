package integrationtest_test

import (
	integrationtest "github.com/nrf110/integration-test/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/opensearch-project/opensearch-go/v2"
)

var _ = Describe("OpenSearchDependency", func() {
	It("can connect", func(ctx SpecContext) {
		es := integrationtest.NewOpenSearchDependency()
		Expect(es.Start(ctx)).To(BeNil())
		defer func() {
			es.Stop(ctx)
		}()

		client := es.Client().(*opensearch.Client)
		res, err := client.Indices.Create("test")
		Expect(err).To(BeNil())
		Expect(res.StatusCode).To(Equal(200))
	})
})
