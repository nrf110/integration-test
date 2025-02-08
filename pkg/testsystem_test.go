package integrationtest_test

import (
	"github.com/jackc/pgx/v5"
	"github.com/nrf110/integration-test/migrations"
	integrationtest "github.com/nrf110/integration-test/pkg"
	"github.com/nrf110/integration-test/pkg/postgres"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pressly/goose/v3"
)

var _ = Describe("TestSystem", func() {
	Describe("Start", func() {
		pgConfig := &postgres.Config{
			Database: "test",
			User:     "test",
			Password: "test",
		}

		It("should start all dependencies", func(ctx SpecContext) {
			system, err := integrationtest.NewTestSystem(
				integrationtest.WithRedis(),
				integrationtest.WithElasticsearch(),
				integrationtest.WithPubSub(),
				integrationtest.WithGCS(),
				integrationtest.WithPostgres(pgConfig),
				integrationtest.WithPermify())
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				err = system.Stop(ctx)
				Expect(err).NotTo(HaveOccurred())
			}()

			err = system.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should run migrations", func(ctx SpecContext) {
			system, err := integrationtest.NewTestSystem(
				integrationtest.WithPostgres(pgConfig),
				integrationtest.WithGooseProviders(func(s *integrationtest.TestSystem) (*goose.Provider, error) {
					conn := s.Postgres.Client().(*pgx.Conn)
					return migrations.NewSQLProvider(conn)
				}),
			)
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				err = system.Stop(ctx)
				Expect(err).NotTo(HaveOccurred())
			}()

			err = system.Start(ctx)
			Expect(err).NotTo(HaveOccurred())

			conn := system.Postgres.Client().(*pgx.Conn)

			row := conn.QueryRow(ctx, "select count(*) from test;")
			var count int
			err = row.Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})
})
