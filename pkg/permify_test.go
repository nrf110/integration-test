package integrationtest

import (
	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	"context"
	permifygrpc "github.com/Permify/permify-go/grpc"
	integrationtest "github.com/nrf110/integration-test/pkg/permify"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPermifyDependency(t *testing.T) {
	const tenantId = "t1"

	t.Run("can check permissions", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		t.Cleanup(cancel)

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

		assert.NoError(t, r.Start(ctx))
		t.Cleanup(func() {
			assert.NoError(t, r.Stop(ctx))
		})

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
		assert.NoError(t, err)
		assert.Equal(t, permifypayload.CheckResult_CHECK_RESULT_ALLOWED, check1Res.GetCan())

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
		assert.NoError(t, err)
		assert.Equal(t, permifypayload.CheckResult_CHECK_RESULT_DENIED, check2Res.GetCan())
	})
}
