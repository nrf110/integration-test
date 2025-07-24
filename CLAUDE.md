# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go integration testing framework that provides utilities for testing applications with various dependencies using Testcontainers.

### Core Architecture

The system is built around the `TestSystem` struct (pkg/testsystem.go:16-26) which orchestrates multiple dependencies:
- **Dependency Interface** (pkg/dependency.go:5-10): Core abstraction for all test dependencies requiring Start/Stop lifecycle, Client access, and environment variable management
- **Test Dependencies**: Pre-built integrations for Elasticsearch, Redis, PostgreSQL, BigQuery, Pub/Sub, GCS, and Permify
- **Testcontainers Integration**: All dependencies use testcontainers-go for running services in Docker containers during tests

### Key Patterns

1. **Option Pattern**: All dependencies use functional options for configuration (e.g., `WithElasticsearch`, `WithPostgres`)
2. **Environment Variables**: Each dependency provides its connection details via `Env()` map
3. **Database Migrations**: PostgreSQL integration supports Goose migrations via `WithGooseProviders`

## Development Commands

```bash
# Run all tests
make test

# Clean build artifacts
make clean

# Update dependencies and pull latest container images
make update

# Run specific package tests
go test ./pkg

# Run a single test
go test -run TestName ./pkg
```

## Testing Best Practices

When adding new tests:
1. Use the `TestSystem` to manage dependencies - don't create containers manually
2. Always use `t.Cleanup()` instead of `defer` for cleanup operations (as per commit 5a30374)
3. Dependencies start automatically when calling `TestSystem.Start(ctx)`
4. Access clients via the dependency fields (e.g., `ts.Redis.Client()`)

## Common Tasks

### Adding a New Dependency
1. Create a new package under `pkg/` with your dependency implementation
2. Implement the `Dependency` interface
3. Add a field to `TestSystem` struct
4. Create a `With<Dependency>` option function
5. Write tests using the pattern from existing test files

### Running Database Migrations
Use `WithGooseProviders` option when creating the TestSystem with Postgres to automatically run migrations during startup.