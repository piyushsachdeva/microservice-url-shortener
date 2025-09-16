# Golang URL Shortener

This URL shortener service, built with Go and Hexagonal Architecture, leverages a serverless approach for efficient scalability and performance. It uses a variety of AWS services to provide a robust, maintainable, and highly available URL shortening service.


- [Prerequisites](#prerequisites)
- [Technologies Used](#technologies-used)
- [System Architecture](#system-architecture)
# URL Shortener Microservices - Complete Repository Guide

A production-ready URL shortener built with **Go microservices** using **Clean Architecture**, featuring containerized deployment, Kubernetes manifests, and comprehensive monitoring. This repository demonstrates modern cloud-native development practices with multiple deployment options.

## 📁 Repository Structure Deep Dive

This repository is organized into clearly separated concerns, following **Clean Architecture principles** with **microservices patterns**:

### 🏗️ **Core Application Structure**

```
golang-url-shortener/
├── internal/                    # Shared Clean Architecture Library
│   ├── adapters/               # Interface Adapters Layer
│   │   ├── cache/             # Redis caching implementation
│   │   ├── functions/         # AWS Lambda functions (serverless deployment)
│   │   ├── handlers/          # HTTP/Lambda request handlers
│   │   ├── messaging/         # RabbitMQ message queue implementation
│   │   └── repository/        # Data persistence implementations
│   ├── config/                # Configuration management
│   ├── core/                  # Business Logic Layer (Enterprise + Application)
│   │   ├── domain/           # Domain entities (Link, Stats) - Enterprise Layer
│   │   ├── ports/            # Interface definitions/contracts - Application Layer
│   │   └── services/         # Business logic implementations - Application Layer
│   └── tests/                # Comprehensive testing suite
│       ├── benchmark/        # Performance benchmarks
│       ├── mock/            # Test doubles and mocks
│       └── unit/            # Unit tests
├── services/                  # Independent Microservices
│   ├── link-service/         # URL generation microservice
│   ├── redirect-service/     # URL redirection microservice  
│   └── stats-service/        # Analytics and statistics microservice
├── frontend/                 # Web User Interface (Framework Layer)
├── k8s/                      # Kubernetes Manifests
│   ├── base/                # Base Kubernetes resources
│   └── gitopsportainer/     # Portainer GitOps deployment
├── configs/                  # Runtime configurations
├── scripts/                  # Database initialization
└── Infrastructure files      # Docker, CI/CD, deployment scripts
```

### 🏗️ **Clean Architecture Implementation**

This project follows **Clean Architecture principles** organized into distinct layers with clear dependency rules:

#### **Enterprise Business Rules Layer (`internal/core/domain/`)**

- **Pure Domain Entities**: `Link`, `Stats` with no external dependencies
- **Business Rules**: Core business logic independent of frameworks or databases
- **Domain Models**: Represent the most general and high-level rules

#### **Application Business Rules Layer (`internal/core/`)**

- **`ports/`**: Interface contracts that define application needs (e.g., `LinkPort`, `CachePort`)
- **`services/`**: Application-specific business logic that orchestrates domain entities
- **Use Cases**: Coordinate between domain entities and external interfaces

#### **Interface Adapters Layer (`internal/adapters/`)**

- **`repository/postgres/`**: Database persistence adapters implementing business interfaces
- **`cache/redis.go`**: Caching layer implementing cache interfaces
- **`handlers/`**: HTTP/API adapters that convert web requests to application use cases
- **`functions/`**: AWS Lambda adapters for serverless deployment
- **`messaging/rabbitmq/`**: Message queue adapters for async processing

#### **Frameworks & Drivers Layer**

- **`services/`**: Independent microservices that use the shared internal library
- **`frontend/`**: User interface framework
- **`configs/`**: External framework configurations (nginx, prometheus)
- **`k8s/`**: Infrastructure and deployment configurations

### 🔀 **Microservices + Clean Architecture Benefits**

- **Shared Business Logic**: All microservices use the same `internal/` clean architecture library
- **Independent Deployment**: Each service in `services/` can be deployed independently
- **Technology Flexibility**: Easy to swap databases, caches, or message queues without changing business logic
- **Testability**: Clean separation enables comprehensive unit testing of business logic
- **Scalability**: Services can scale independently based on their specific load patterns



### **Prerequisites**
- **Docker** & **Docker Compose** (for local development)
- **Go 1.22+** (see `go.mod` for exact version)
- **kubectl** (for Kubernetes deployment)
- **PostgreSQL** (for local database development)

### **Local Development with Docker Compose**

The `docker-compose.yaml` sets up the complete stack locally:

```bash
# Start all services with hot reload
docker-compose up --build

# Access points:
# Frontend: http://localhost:3000
# API Gateway: http://localhost:8080  
# Direct services: 8001, 8002, 8003
# PostgreSQL: localhost:5432
```

### **Manual Service Development**

```bash
# Run database migrations
psql -h localhost -U postgres -d urlshortener -f scripts/init.sql

# Start individual services
cd services/link-service && go run main.go
cd services/redirect-service && go run main.go  
cd services/stats-service && go run main.go
```

## 🐳 **Docker Hub Build & Push Process**

The `push-to-dockerhub.sh` script automates building and pushing all service images to Docker Hub.

### **Script Configuration**

```bash
# Default configuration (edit in script)
DOCKER_HUB_USERNAME="itsbaivab"          # Your Docker Hub username
IMAGE_TAG="${IMAGE_TAG:-latest}"         # Configurable via environment
BUILD_PLATFORM="linux/amd64,linux/arm64" # Multi-architecture support

# Services built:
- url-shortener-link      (services/link-service/Dockerfile)
- url-shortener-redirect  (services/redirect-service/Dockerfile)  
- url-shortener-stats     (services/stats-service/Dockerfile)
- url-shortener-frontend  (frontend/Dockerfile)
```

### **Complete Build & Push Workflow**

```bash
# 1. Build and push all images (recommended)
./push-to-dockerhub.sh deploy

# 2. Using custom tag
IMAGE_TAG=v1.2.0 ./push-to-dockerhub.sh deploy

# 3. Update Kubernetes manifests with new image tags
./push-to-dockerhub.sh update-k8s

# 4. Verify images on Docker Hub
./push-to-dockerhub.sh verify
```

### **Individual Operations**

```bash
# Build images locally only
./push-to-dockerhub.sh build

# Push existing images to Docker Hub
./push-to-dockerhub.sh login  # First-time setup
./push-to-dockerhub.sh push

# List local images
./push-to-dockerhub.sh list

# Clean up local images
./push-to-dockerhub.sh cleanup

# Get Docker Hub repository info
./push-to-dockerhub.sh info
```

### **Environment Variables**

```bash
# Custom image tag
export IMAGE_TAG="v2.1.0"

# Custom Docker Hub username (overrides script default)
export DOCKER_HUB_USERNAME="myusername"

# Custom build platform
export BUILD_PLATFORM="linux/amd64"
```

## ⚓ **Kubernetes Deployment Options**

### **Option 1: Standard Kubernetes** (`k8s/base/`)

```bash
# Deploy to any Kubernetes cluster
kubectl apply -f k8s/base/

# Check deployment status
kubectl get pods -n url-shortener
kubectl get services -n url-shortener
```

### **Option 2: Portainer GitOps** (`k8s/gitopsportainer/`)

Complete GitOps deployment with Portainer management interface. **For detailed Portainer setup instructions, see:**

📖 **[k8s/gitopsportainer/README-GITOPS.md](k8s/gitopsportainer/README-GITOPS.md)**

**Quick Portainer Deployment:**

```bash
# Navigate to GitOps directory
cd k8s/gitopsportainer/

# Option A: Automatic deployment (recommended)
./deploy.sh

# Option B: Using Kustomize
kubectl apply -k .

# Option C: Manual step-by-step
./install-ingress-controller.sh
kubectl apply -f .
```

**Portainer Demo Features:**
- **Visual Kubernetes management** through Portainer UI
- **GitOps workflow** with automatic deployments
- **Ingress controller** setup with external access
- **LoadBalancer services** for cloud environments
- **Monitoring dashboards** integration
- **Secrets management** for database credentials

**Access after Portainer deployment:**
```bash
# Via Ingress (recommended)
kubectl get ingress -n url-shortener

# Via LoadBalancer
kubectl get svc -n url-shortener | grep LoadBalancer

# Via Port Forward (development)
kubectl port-forward svc/frontend 8080:80 -n url-shortener
```

## 🚀 **Production Deployment Considerations**

### **Environment Variables**
```bash
# Database configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secure_password
DB_NAME=urlshortener

# Redis configuration (optional)
REDIS_HOST=redis
REDIS_PORT=6379

# Service configuration  
SERVICE_PORT=8001
LOG_LEVEL=info
```


- **Link Service**: CPU-intensive (short I

## 🔗 **Related Documentation**

- **Portainer GitOps Setup**: [k8s/gitopsportainer/README-GITOPS.md](k8s/gitopsportainer/README-GITOPS.md)
- **API Documentation**: Available at `/api/docs` when services are running
- **Database Schema**: See `scripts/init.sql` for complete schema

