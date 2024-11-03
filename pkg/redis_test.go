package integrationtest_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	container "github.com/testcontainers/testcontainers-go/modules/redis"
)

var _ = Describe("RedisDependency", func() {
	It("connect", func(ctx SpecContext) {
		c, err := container.Run(ctx, "redis:7")
		Expect(err).To(BeNil())
		defer func() {
			c.Terminate(ctx)
		}()
	})
})
