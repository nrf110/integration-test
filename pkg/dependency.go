package integrationtest

import "context"

type Dependency interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Client() any
	Env() map[string]string
}
