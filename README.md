# Go Microservice Skeleton

A demonstration project showcasing how to build a complete Go microservice using modern technologies. This skeleton serves as both a learning example and a practical starting point for developers who want to build similar services.

## Technologies Demonstrated

- **Go**: Modern Go service with proper project structure
- **Docker**: Containerization with multi-stage builds
- **Kubernetes**: Complete deployment manifests (Deployment/StatefulSet options)
- **Kafka**: Message broker integration with producers/consumers
- **BoltDB**: Local embedded database for persistent storage
- **HTTP Services**: RESTful API with standard endpoints
- **Monitoring**: Prometheus metrics and Sentry error tracking
- **Security**: Comprehensive vulnerability scanning and license management

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

## Testing

Uses Ginkgo BDD framework with Gomega matchers. Run tests with `make test` or `ginkgo run ./...`.  
