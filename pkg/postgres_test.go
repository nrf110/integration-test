package integrationtest_test

import (
	"github.com/jackc/pgx/v5"
	integrationtest "github.com/nrf110/integration-test/pkg/postgres"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dependency", func() {
	It("can connect", func(ctx SpecContext) {
		pg := integrationtest.NewDependency(&integrationtest.Config{
			Database: "postgres",
			User:     "postgres",
			Password: "postgres",
		})
		err := pg.Start(ctx)
		Expect(err).To(BeNil())
		defer func() {
			pg.Stop(ctx)
		}()
		conn := pg.Client().(*pgx.Conn)
		Expect(conn.Ping(ctx)).To(BeNil())
	})
})
