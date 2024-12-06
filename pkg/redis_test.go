package integrationtest_test

import (
	integrationtest "github.com/nrf110/integration-test/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"
)

var _ = Describe("RedisDependency", func() {
	It("connect", func(ctx SpecContext) {
		r := integrationtest.NewRedisDependency()
		Expect(r.Start(ctx)).To(BeNil())
		defer func() {
			r.Stop(ctx)
		}()

		client := r.Client().(*redis.Client)
		status := client.Ping(ctx)
		Expect(status.Err()).To(BeNil())
	})
})
