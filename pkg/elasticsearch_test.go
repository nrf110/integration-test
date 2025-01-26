package integrationtest_test

import (
	"github.com/elastic/go-elasticsearch/v8"
	integrationtest "github.com/nrf110/integration-test/pkg/elasticsearch"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dependency", func() {
	It("can connect", func(ctx SpecContext) {
		es := integrationtest.NewDependency()
		Expect(es.Start(ctx)).To(BeNil())
		defer func() {
			es.Stop(ctx)
		}()

		client := es.Client().(*elasticsearch.TypedClient)
		res, err := client.API.Indices.Create("test").Do(ctx)
		Expect(err).To(BeNil())
		Expect(res.Acknowledged).To(BeTrue())
	})
})
