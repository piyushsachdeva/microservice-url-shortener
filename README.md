# Golang URL Shortener

This URL shortener service, built with Go and Hexagonal Architecture, leverages a serverless approach for efficient scalability and performance. It uses a variety of AWS services to provide a robust, maintainable, and highly available URL shortening service.


- [Prerequisites](#prerequisites)
- [Technologies Used](#technologies-used)
- [System Architecture](#system-architecture)
# URL Shortener Microservices - Complete Repository Guide

A production-ready URL shortener built with **Go microservices** using **Clean Architecture**, featuring containerized deployment, Kubernetes manifests, and comprehensive monitoring. This repository demonstrates modern cloud-native development practices with multiple deployment options.


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
DOCKER_HUB_USERNAME="piyushsachdeva"          # Your Docker Hub username
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

```

**Access after Portainer deployment:**
```bash
# Via Ingress (recommended)
kubectl get ingress -n url-shortener

# Via LoadBalancer
kubectl get svc -n url-shortener | grep LoadBalancer

# Via Port Forward (development)
kubectl port-forward svc/frontend 8080:80 -n url-shortener
```

## 🔗 **Related Documentation**

- **Portainer GitOps Setup**: [k8s/gitopsportainer/README-GITOPS.md](k8s/gitopsportainer/README-GITOPS.md)
- **API Documentation**: Available at `/api/docs` when services are running
- **Database Schema**: See `scripts/init.sql` for complete schema

