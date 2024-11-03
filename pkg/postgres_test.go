package integrationtest_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	container "github.com/testcontainers/testcontainers-go/modules/postgres"
)

var _ = Describe("PostgresDependency", func() {
	It("connect", func(ctx SpecContext) {
		c, err := container.Run(ctx, "postgres:16",
			container.WithDatabase("postgres"),
			container.WithUsername("postgres"),
			container.WithPassword("postgres"),
			container.WithSQLDriver("pgx"))
		Expect(err).To(BeNil())
		defer func() {
			c.Terminate(ctx)
		}()
	})
})
