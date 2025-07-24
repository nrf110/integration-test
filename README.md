# Integration Test Framework

A Go library that simplifies integration testing by providing a unified interface for managing test dependencies using Testcontainers.

## Features

- 🐳 Automatic container lifecycle management
- 🔌 Pre-built integrations for popular services
- 🔧 Simple, consistent API across all dependencies
- 🏃 Parallel test execution support
- 🔄 Database migration support
- 🌍 Environment variable injection

## Supported Dependencies

- **PostgreSQL** - Full-featured PostgreSQL with migration support
- **Redis** - In-memory data store
- **Elasticsearch** - Search and analytics engine
- **Google Cloud Pub/Sub** - Message queue service
- **Google Cloud Storage (GCS)** - Object storage
- **BigQuery** - Data warehouse
- **Permify** - Authorization service

## Installation

```bash
go get github.com/nrf110/integration-test
```

## Usage

### Basic Example

```go
package myapp_test

import (
    "context"
    "testing"
    "time"
    
    it "github.com/nrf110/integration-test/pkg"
    "github.com/stretchr/testify/assert"
)

func TestMyApp(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    t.Cleanup(cancel)
    
    // Create test system with Redis and PostgreSQL
    system, err := it.NewTestSystem(
        it.WithRedis(),
        it.WithPostgres(&postgres.Config{
            Database: "testdb",
            User:     "testuser",
            Password: "testpass",
        }),
    )
    assert.NoError(t, err)
    
    // Ensure cleanup
    t.Cleanup(func() {
        assert.NoError(t, system.Stop(ctx))
    })
    
    // Start all dependencies
    err = system.Start(ctx)
    assert.NoError(t, err)
    
    // Access clients
    redisClient := system.Redis.Client().(*redis.Client)
    pgConn := system.Postgres.Client().(*pgx.Conn)
    
    // Run your tests...
}
```

### Individual Dependency Examples

#### PostgreSQL with Migrations

```go
import (
    "github.com/nrf110/integration-test/pkg/postgres"
    "github.com/pressly/goose/v3"
)

func TestWithPostgres(t *testing.T) {
    system, err := it.NewTestSystem(
        it.WithPostgres(&postgres.Config{
            Database: "myapp",
            User:     "postgres",
            Password: "postgres",
        }),
        it.WithGooseProviders(func(s *it.TestSystem) (*goose.Provider, error) {
            conn := s.Postgres.Client().(*pgx.Conn)
            // Configure your migrations
            return goose.NewProvider(goose.DialectPostgres, conn, nil)
        }),
    )
    // ... rest of test
}
```

#### Redis

```go
func TestWithRedis(t *testing.T) {
    system, err := it.NewTestSystem(it.WithRedis())
    assert.NoError(t, err)
    
    err = system.Start(ctx)
    assert.NoError(t, err)
    
    client := system.Redis.Client().(*redis.Client)
    
    // Use Redis
    err = client.Set(ctx, "key", "value", 0).Err()
    assert.NoError(t, err)
}
```

#### Elasticsearch

```go
func TestWithElasticsearch(t *testing.T) {
    system, err := it.NewTestSystem(it.WithElasticsearch())
    assert.NoError(t, err)
    
    err = system.Start(ctx)
    assert.NoError(t, err)
    
    esClient := system.Elasticsearch.Client().(*elasticsearch.Client)
    
    // Create index, insert documents, etc.
}
```

#### Google Cloud Services

```go
func TestWithGoogleCloud(t *testing.T) {
    system, err := it.NewTestSystem(
        it.WithPubSub(),
        it.WithGCS(),
    )
    assert.NoError(t, err)
    
    err = system.Start(ctx)
    assert.NoError(t, err)
    
    // Access Pub/Sub client
    pubsubClient := system.PubSub.Client().(*pubsub.Client)
    
    // Access GCS client
    gcsClient := system.GCS.Client().(*storage.Client)
}
```

### Advanced Usage

#### Custom Dependencies

You can add custom dependencies by implementing the `Dependency` interface:

```go
type Dependency interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Client() any
    Env() map[string]string
}

// Then use WithDependency
system, err := it.NewTestSystem(
    it.WithDependency(myCustomDep),
)
```

#### Environment Variables

Each dependency automatically provides environment variables:

```go
system, err := it.NewTestSystem(it.WithPostgres(pgConfig))
err = system.Start(ctx)

// Environment variables are available
// e.g., POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, etc.
```

#### Using Individual Dependencies

You can also use dependencies standalone without the TestSystem:

```go
func TestStandaloneRedis(t *testing.T) {
    ctx := context.Background()
    
    redisDep := redis.NewDependency()
    err := redisDep.Start(ctx)
    assert.NoError(t, err)
    t.Cleanup(func() {
        assert.NoError(t, redisDep.Stop(ctx))
    })
    
    client := redisDep.Client().(*redis.Client)
    // Use client...
}
```

## Best Practices

1. **Always use `t.Cleanup()`** for cleanup instead of `defer`
2. **Set appropriate timeouts** - Container startup can take time
3. **Use the TestSystem** for managing multiple dependencies
4. **Check container logs** if tests fail unexpectedly
5. **Run `make update`** periodically to pull latest container images

## Contributing

We welcome contributions! Here's how you can help:

### Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/integration-test.git`
3. Create a feature branch: `git checkout -b feature/amazing-feature`

### Development Setup

```bash
# Install dependencies
go mod download

# Run tests
make test

# Update dependencies and container images
make update
```

### Adding a New Dependency

1. Create a new package under `pkg/` (e.g., `pkg/mongodb/`)
2. Implement the `Dependency` interface:
   ```go
   type Dependency struct {
       container testcontainers.Container
       client    *mongo.Client
   }
   
   func (d *Dependency) Start(ctx context.Context) error { ... }
   func (d *Dependency) Stop(ctx context.Context) error { ... }
   func (d *Dependency) Client() any { return d.client }
   func (d *Dependency) Env() map[string]string { ... }
   ```
3. Add a constructor function in the main package:
   ```go
   func WithMongoDB(opts ...mongodb.DependencyOpt) Option {
       return func(s *TestSystem) error {
           s.MongoDB = mongodb.NewDependency(opts...)
           return WithDependency(s.MongoDB)(s)
       }
   }
   ```
4. Add the field to `TestSystem` struct
5. Write comprehensive tests

### Submitting Changes

1. Ensure all tests pass: `make test`
2. Commit your changes with clear messages
3. Push to your fork: `git push origin feature/amazing-feature`
4. Create a Pull Request with:
   - Clear description of changes
   - Any breaking changes noted
   - Test coverage for new features

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Add comments for exported functions
- Keep functions focused and testable

### Reporting Issues

- Use GitHub Issues for bug reports and feature requests
- Include reproduction steps for bugs
- Provide context about your use case for feature requests

### Questions?

Feel free to open an issue for any questions about contributing!