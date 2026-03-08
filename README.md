# Kubernetes Platform

A production-ready Kubernetes platform featuring autoscaling, observability, GitOps deployments, comprehensive security policies, and a functional Go API server.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Application](#application)
- [Project Structure](#project-structure)
- [Deployment Strategies](#deployment-strategies)
- [Observability](#observability)
- [Security](#security)
- [Autoscaling](#autoscaling)
- [GitOps](#gitops)
- [Development](#development)
- [Testing](#testing)
- [Makefile Reference](#makefile-reference)

## Overview

This project provides a complete Kubernetes platform infrastructure with:

- **Multi-environment support** (dev, staging, production)
- **GitOps workflows** with ArgoCD
- **Observability stack** with Prometheus and Grafana
- **Security policies** with Network Policies, PSS, and OPA/Kyverno
- **Advanced autoscaling** with HPA and VPA
- **Deployment strategies** including Canary and Blue-Green
- **Functional Go API** with health checks and metrics
- **Comprehensive testing** and validation

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              KUBERNETES CLUSTER                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                           INGRESS LAYER                              │   │
│  │  ┌──────────────┐    ┌──────────────┐    ┌──────────────────────┐  │   │
│  │  │   NGINX      │    │   CERT-      │    │    EXTERNAL DNS      │  │   │
│  │  │   INGRESS    │◄───│   MANAGER    │    │                      │  │   │
│  │  │   CLASS      │    │   (TLS)      │    │                      │  │   │
│  │  └──────┬───────┘    └──────────────┘    └──────────────────────┘  │   │
│  └─────────┼───────────────────────────────────────────────────────────┘   │
│            │                                                                │
│  ┌─────────▼───────────────────────────────────────────────────────────┐   │
│  │                        PLATFORM NAMESPACE                            │   │
│  │                                                                      │   │
│  │  ┌──────────────────────────┐    ┌──────────────────────────┐      │   │
│  │  │      API SERVER (Go)     │    │       FRONTEND           │      │   │
│  │  │  ┌────────────────────┐  │    │  ┌────────────────────┐  │      │   │
│  │  │  │  Pod 1  │  Pod 2  │  │    │  │  Pod 1  │  Pod 2  │  │      │   │
│  │  │  └────────────────────┘  │    │  └────────────────────┘  │      │   │
│  │  │                          │    │                          │      │   │
│  │  │  :8080 HTTP API          │    │  :8080 HTTP              │      │   │
│  │  │  :9090 Prometheus        │    │                          │      │   │
│  │  └──────────────────────────┘    └──────────────────────────┘      │   │
│  │                                                                      │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                        MONITORING NAMESPACE                           │   │
│  │                                                                       │   │
│  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────────────┐ │   │
│  │  │   PROMETHEUS   │  │  ALERTMANAGER  │  │       GRAFANA          │ │   │
│  │  └────────────────┘  └────────────────┘  └────────────────────────┘ │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                         ARGOCD NAMESPACE                              │   │
│  │  ┌────────────────────────────────────────────────────────────────┐ │   │
│  │  │                    ARGOCD CONTROLLER                            │ │   │
│  │  └────────────────────────────────────────────────────────────────┘ │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Features

| Feature | Description | Status |
|---------|-------------|--------|
| **Go API Server** | Functional REST API with health checks and metrics | ✅ |
| **Multi-environment** | Dev, staging, production with Kustomize overlays | ✅ |
| **Helm Charts** | Parameterized, reusable deployment templates | ✅ |
| **GitOps** | ArgoCD with App of Apps pattern | ✅ |
| **HPA** | CPU, memory, and custom metrics | ✅ |
| **VPA** | Vertical scaling recommendations | ✅ |
| **Canary Deployments** | Argo Rollouts with analysis | ✅ |
| **Blue-Green Deployments** | Zero-downtime deployments | ✅ |
| **Network Policies** | Default deny, explicit allow | ✅ |
| **Pod Security Standards** | Restricted policy enforcement | ✅ |
| **Prometheus** | Metrics collection and alerting | ✅ |
| **Grafana Dashboards** | Platform and cluster visibility | ✅ |

## Prerequisites

- Kubernetes cluster v1.28+
- kubectl v1.28+
- Helm v3.14+
- Go 1.21+ (for local development)
- Docker (for building images)
- (Optional) ArgoCD v2.8+
- (Optional) Prometheus Operator

## Quick Start

```bash
# Clone the repository
git clone https://github.com/example/platform.git
cd platform

# Build the Go application
make build-app

# Build Docker image
make docker-build

# Deploy to local development environment
make dev

# Check deployment status
make status

# Test the API
kubectl port-forward svc/api-server 8080:80 -n platform
curl http://localhost:8080/health
```

## Application

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Liveness probe |
| `/ready` | GET | Readiness probe |
| `/api/v1/status` | GET | Application status and uptime |
| `/api/v1/info` | GET | Pod, runtime, and Kubernetes info |
| `/api/v1/echo` | GET | Echo request details (testing) |
| `/metrics` | GET | Prometheus metrics |

### Example Responses

**Health Check:**
```json
{
  "status": "healthy",
  "time": "2024-01-15T10:30:00Z"
}
```

**Application Info:**
```json
{
  "application": {
    "name": "platform-api",
    "environment": "production",
    "version": "1.0.0"
  },
  "runtime": {
    "goVersion": "go1.21.0",
    "os": "linux",
    "arch": "amd64",
    "goroutines": 8
  },
  "kubernetes": {
    "podName": "api-server-abc123",
    "podNamespace": "platform",
    "nodeName": "worker-1"
  }
}
```

### Prometheus Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `platform_http_requests_total` | Counter | Total HTTP requests by method, path, status |
| `platform_http_request_duration_seconds` | Histogram | Request duration in seconds |
| `platform_active_connections` | Gauge | Current active connections |
| `platform_app_info` | Gauge | Application version and environment |

## Project Structure

```
k8s-platform/
├── app/                          # Go API application
│   ├── cmd/server/main.go        # Application entry point
│   ├── internal/
│   │   ├── config/config.go      # Configuration management
│   │   └── handlers/             # HTTP handlers
│   │       ├── response.go       # JSON response helpers
│   │       ├── health.go         # Health/Readiness probes
│   │       ├── api.go            # API endpoints
│   │       └── metrics.go        # Prometheus metrics
│   ├── go.mod
│   ├── Dockerfile
│   └── README.md
│
├── k8s/                          # Kubernetes manifests (Kustomize)
│   ├── base/                     # Base resources
│   │   ├── namespace.yaml
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── ingress.yaml
│   │   ├── configmap.yaml
│   │   ├── secrets.yaml
│   │   └── rbac.yaml
│   └── overlays/                 # Environment-specific overlays
│       ├── dev/
│       ├── staging/
│       └── prod/
│
├── helm-charts/                  # Helm charts
│   └── platform/
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
│
├── gitops/                       # GitOps configuration
│   └── argocd/
│       ├── apps/
│       ├── projects/
│       └── bootstrap/
│
├── security/                     # Security policies
│   ├── network-policies/
│   ├── pod-security/
│   └── policies/
│
├── scripts/                      # Utility scripts
│   ├── validate-policies.sh
│   └── test-deployment.sh
│
├── Makefile
└── README.md
```

## Deployment Strategies

### Canary (Recommended for production)

```bash
kubectl apply -f deployments/canary/rollout.yaml
kubectl argo rollouts get rollout api-server -n platform
```

### Blue-Green

```bash
kubectl apply -f deployments/blue-green/rollout.yaml
kubectl argo rollouts promote api-server -n platform
```

## Observability

### SLOs (Service Level Objectives)

| SLO | Target | Measurement |
|-----|--------|-------------|
| Availability | 99.9% | Successful requests / Total requests |
| Latency (P99) | < 500ms | 99th percentile response time |
| Error Budget | > 10% | Remaining budget before breach |

### Grafana Dashboards

1. **Platform Overview** - Request rate, latency, error budget
2. **Kubernetes Cluster** - Node and pod metrics
3. **SLO Dashboard** - Availability and latency SLOs

## Security

### Network Policies

- Default deny all traffic
- Explicit allow lists for ingress and egress
- DNS resolution allowed

### Pod Security Standards

```yaml
# Production: Restricted
pod-security.kubernetes.io/enforce: restricted
```

## Autoscaling

### HPA Configuration

```yaml
minReplicas: 3
maxReplicas: 20
metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

## GitOps

### ArgoCD Setup

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl apply -f gitops/argocd/projects/
kubectl apply -f gitops/argocd/apps/
```

## Development

### Local Development

```bash
# Run Go app locally
make run-local

# Build binary
make build-app

# Build Docker image
make docker-build

# Create Kind cluster
make kind-create

# Deploy to Kind
make deploy

# Cleanup
make kind-delete
```

### Testing

```bash
# Run all tests
make test

# Run specific tests
make test-helm
make test-policies

# Validate manifests
make lint
```

## Makefile Reference

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make build-app` | Build Go application binary |
| `make docker-build` | Build Docker image |
| `make run-local` | Run application locally |
| `make dev` | Deploy development environment |
| `make deploy` | Deploy with Helm |
| `make deploy-prod` | Deploy to production |
| `make destroy` | Remove all deployments |
| `make test` | Run all tests |
| `make lint` | Run all linters |
| `make status` | Show deployment status |
| `make logs` | View application logs |
| `make port-forward` | Port-forward services |

## Troubleshooting

```bash
# Check pod status
kubectl get pods -n platform
kubectl describe pod <pod-name> -n platform

# View logs
kubectl logs -f deployment/api-server -n platform

# Test API locally
kubectl port-forward svc/api-server 8080:80 -n platform
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/info
curl http://localhost:8080/metrics
```

## License

Apache 2.0

---

**Key Skills Demonstrated:**
- Go HTTP server development
- Kubernetes workload management
- Deployment strategies (Rolling Update, Blue-Green, Canary)
- Autoscaling (HPA, VPA, Custom Metrics)
- Observability (Prometheus, Grafana, Alerting)
- GitOps with ArgoCD
- Security hardening (Network Policies, PSS, OPA/Kyverno)
- Infrastructure as Code (Helm, Kustomize)
