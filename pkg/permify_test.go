package integrationtest

import (
	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	permifygrpc "github.com/Permify/permify-go/grpc"
	integrationtest "github.com/nrf110/integration-test/pkg/permify"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("permify.Dependency", func() {
	const tenantId = "t1"

	It("can check permissions", func(ctx SpecContext) {
		r := integrationtest.NewDependency(
			integrationtest.WithTenantId(tenantId),
			integrationtest.WithSchema(`
				entity user {
					relation self @user

					permission read = self
				}
			`),
			integrationtest.WithTuples(&permifypayload.Tuple{
				Entity: &permifypayload.Entity{
					Type: "user",
					Id:   "1",
				},
				Relation: "self",
				Subject: &permifypayload.Subject{
					Type: "user",
					Id:   "1",
				},
			}))

		Expect(r.Start(ctx)).To(BeNil())
		defer func() {
			r.Stop(ctx)
		}()

		client := r.Client().(*permifygrpc.Client)

		check1Res, err := client.Permission.Check(ctx, &permifypayload.PermissionCheckRequest{
			TenantId: tenantId,
			Subject: &permifypayload.Subject{
				Type: "user",
				Id:   "1",
			},
			Entity: &permifypayload.Entity{
				Type: "user",
				Id:   "1",
			},
			Permission: "read",
			Metadata: &permifypayload.PermissionCheckRequestMetadata{
				Depth: 10,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(check1Res.GetCan()).To(Equal(permifypayload.CheckResult_CHECK_RESULT_ALLOWED))

		check2Res, err := client.Permission.Check(ctx, &permifypayload.PermissionCheckRequest{
			TenantId: tenantId,
			Subject: &permifypayload.Subject{
				Type: "user",
				Id:   "2",
			},
			Entity: &permifypayload.Entity{
				Type: "user",
				Id:   "1",
			},
			Permission: "read",
			Metadata: &permifypayload.PermissionCheckRequestMetadata{
				Depth: 10,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(check2Res.GetCan()).To(Equal(permifypayload.CheckResult_CHECK_RESULT_DENIED))
	})
})
