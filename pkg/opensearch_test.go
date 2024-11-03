package integrationtest_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	container "github.com/testcontainers/testcontainers-go/modules/opensearch"
)

var _ = Describe("OpenSearchDependency", func() {
	It("connect", func(ctx SpecContext) {
		c, err := container.Run(ctx, "opensearchproject/opensearch:2.11.1")
		Expect(err).To(BeNil())
		defer func() {
			c.Terminate(ctx)
		}()
	})
})
