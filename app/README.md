# Platform API Server

A simple, production-ready Go API server for Kubernetes demonstration.

## Features

- **Health Endpoints**: `/health` and `/ready` for Kubernetes probes
- **Prometheus Metrics**: Built-in metrics endpoint at `/metrics`
- **Graceful Shutdown**: Handles SIGINT/SIGTERM gracefully
- **Structured Logging**: JSON logging support
- **CORS Support**: Configurable CORS headers

## Endpoints

| Endpoint | Port | Description |
|----------|------|-------------|
| `/health` | 8080 | Health check (liveness probe) |
| `/ready` | 8080 | Readiness check |
| `/api/v1/status` | 8080 | Application status |
| `/api/v1/info` | 8080 | Pod and environment information |
| `/api/v1/echo` | 8080 | Echo request details (testing) |
| `/metrics` | 9090 | Prometheus metrics |

## Local Development

```bash
# Run locally
go run ./cmd/server

# Build binary
go build -o bin/server ./cmd/server

# Run with custom environment
PORT=3000 METRICS_PORT=9091 ./bin/server
```

## Docker Build

```bash
# Build image
docker build -t platform-api:latest .

# Run container
docker run -p 8080:8080 -p 9090:9090 platform-api:latest
```

## Configuration (Environment Variables)

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | HTTP server port |
| `METRICS_PORT` | 9090 | Metrics server port |
| `APP_NAME` | platform-api | Application name |
| `ENVIRONMENT` | development | Environment name |
| `LOG_LEVEL` | info | Logging level |
| `ENABLE_CORS` | true | Enable CORS headers |

## Kubernetes Integration

The application automatically reads Kubernetes metadata:
- `POD_NAME` - Pod name
- `POD_NAMESPACE` - Pod namespace
- `NODE_NAME` - Node name

These are injected via the Kubernetes Downward API.

## Metrics

The following metrics are exposed:

| Metric | Type | Description |
|--------|------|-------------|
| `platform_http_requests_total` | Counter | Total HTTP requests |
| `platform_http_request_duration_seconds` | Histogram | Request duration |
| `platform_active_connections` | Gauge | Active connections |
| `platform_app_info` | Gauge | Application info |

## License

Apache 2.0
