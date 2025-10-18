# Go Microservice Skeleton

[![Go Reference](https://pkg.go.dev/badge/github.com/bborbe/skeleton.svg)](https://pkg.go.dev/github.com/bborbe/skeleton)
[![CI](https://github.com/bborbe/skeleton/actions/workflows/ci.yml/badge.svg)](https://github.com/bborbe/skeleton/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bborbe/skeleton)](https://goreportcard.com/report/github.com/bborbe/skeleton)

**A template/skeleton project for creating new Go microservices.** This project demonstrates patterns and integration of common technologies, serving as a copy-paste starting point for new services rather than a complete production implementation.

---

## Technologies Demonstrated

- **Go**: Modern Go service with proper project structure
- **Docker**: Containerization with multi-stage builds
- **Kubernetes**: Complete deployment manifests (Deployment/StatefulSet options)
- **Kafka**: Message broker integration with producers/consumers
- **BoltDB**: Local embedded database for persistent storage
- **HTTP Services**: RESTful API with standard endpoints
- **Monitoring**: Prometheus metrics and Sentry error tracking
- **Security**: Comprehensive vulnerability scanning and license management

---

## Installation

As a template project, you typically copy this repository rather than installing it as a dependency. However, if you need to reference it:

```bash
go get github.com/bborbe/skeleton
```

## Quick Start

```bash
# Install dependencies
make deps

# Run the service locally
make run

# Run full development workflow
make precommit
```

## Standard Endpoints

The service includes these standard endpoints:
- `/healthz` - Health check
- `/readiness` - Readiness check
- `/metrics` - Prometheus metrics
- `/resetdb` - Database reset
- `/setloglevel/{level}` - Dynamic log level adjustment
- `/testloglevel` - Test logging at all levels
- `/sentryalert` - Test Sentry integration

## Development Commands

```bash
make ensure    # Dependency management (go mod tidy, verify)
make format    # Code formatting (gofmt, goimports-reviser, golines)
make generate  # Code generation
make test      # Run tests with coverage
make check     # Security and quality checks (vet, errcheck, vulncheck, osv-scanner, trivy)
make precommit # Complete development workflow
```

## Configuration

Configuration is managed through environment variables (see `example.env`):
- `KAFKA_BROKERS` - Kafka broker addresses
- `SKELETON_PORT` - HTTP server port
- `SENTRY_DSN` - Sentry error tracking DSN
- `DATADIR` - Database storage directory

## Deployment

### Kubernetes
The `k8s/` directory contains complete Kubernetes manifests:
- Choose between Deployment (`skeleton-deploy.yaml`) or StatefulSet (`skeleton-sts.yaml`)
- Service, Ingress, Secret, and User configurations included
- **To activate**: Change replicas from 0 → 1
- **For Deployment**: Delete StatefulSet manifest if using Deployment

### Docker
Multi-stage Docker build included. See `Dockerfile` and `Makefile.docker`.

## Architecture

Standard Go microservice architecture:
- `main.go` - Application entry point with dependency injection
- `pkg/factory/` - Factory functions for handler creation
- `pkg/handler/` - HTTP handlers for various endpoints
- Uses proven libraries for service framework, HTTP handling, Kafka, database, and monitoring

## API Documentation

For complete API documentation of the handlers and packages, visit [pkg.go.dev](https://pkg.go.dev/github.com/bborbe/skeleton).

---

## Testing

Uses Ginkgo BDD framework with Gomega matchers. Run tests with `make test` or `ginkgo run ./...`.

Example tests are included to demonstrate testing patterns - copy and adapt these when building your service.

## Usage as Template

1. Copy this repository as a starting point for your new service
2. Update service name throughout (replace "skeleton")
3. Modify handlers and add your business logic
4. Update configuration in `example.env` and K8s manifests
5. Implement your specific Kafka consumers/producers and database schema

---

## License

BSD-style license. See [LICENSE](LICENSE) file for details.
