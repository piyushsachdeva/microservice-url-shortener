#!/bin/bash

# Docker Hub Push Script for URL Shortener Microservices
# This script builds and pushes Docker images to Docker Hub

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
DOCKER_HUB_USERNAME="piyushsachdeva"
IMAGE_TAG="${IMAGE_TAG:-latest}"
BUILD_PLATFORM="${BUILD_PLATFORM:-linux/amd64,linux/arm64}"

# Services to build and push
SERVICES=(
    "url-shortener-link"
    "url-shortener-redirect" 
    "url-shortener-stats"
    "url-shortener-frontend"
)

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if Docker is available
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed or not in PATH"
        exit 1
    fi
    print_success "Docker is available"
}

# Function to check if Docker daemon is running
check_docker_daemon() {
    if ! docker info &> /dev/null; then
        print_error "Docker daemon is not running"
        exit 1
    fi
    print_success "Docker daemon is running"
}

# Function to check Docker Hub login
check_docker_login() {
    print_info "Checking Docker Hub authentication..."
    
    # Try to get current logged in user
    CURRENT_USER=$(docker system info --format '{{.Username}}' 2>/dev/null || echo "")
    
    if [ "$CURRENT_USER" = "$DOCKER_HUB_USERNAME" ]; then
        print_success "Already logged in as $DOCKER_HUB_USERNAME"
        return 0
    fi
    
    print_info "Please log in to Docker Hub as $DOCKER_HUB_USERNAME"
    if docker login --username "$DOCKER_HUB_USERNAME"; then
        print_success "Successfully logged in to Docker Hub"
    else
        print_error "Failed to log in to Docker Hub"
        exit 1
    fi
}

# Function to build a single service image
build_service() {
    local service=$1
    local dockerfile_path
    local image_name="${DOCKER_HUB_USERNAME}/${service}"
    
    print_info "Building $service..."
    
    # Handle different service types
    if [ "$service" = "url-shortener-frontend" ]; then
        dockerfile_path="frontend/Dockerfile"
    elif [ "$service" = "url-shortener-link" ]; then
        dockerfile_path="services/link-service/Dockerfile"
    elif [ "$service" = "url-shortener-redirect" ]; then
        dockerfile_path="services/redirect-service/Dockerfile"
    elif [ "$service" = "url-shortener-stats" ]; then
        dockerfile_path="services/stats-service/Dockerfile"
    else
        print_error "Unknown service: $service"
        return 1
    fi
    
    if [ ! -f "$dockerfile_path" ]; then
        print_error "Dockerfile not found: $dockerfile_path"
        return 1
    fi
    
    # Build the image
    docker build \
        -t "${image_name}:${IMAGE_TAG}" \
        -t "${image_name}:latest" \
        -f "$dockerfile_path" \
        .
    
    print_success "Built ${image_name}:${IMAGE_TAG}"
}

# Function to build all service images
build_all_images() {
    print_info "Building all Docker images..."
    
    for service in "${SERVICES[@]}"; do
        build_service "$service"
    done
    
    print_success "All images built successfully"
}

# Function to push a single service image
push_service() {
    local service=$1
    local image_name="${DOCKER_HUB_USERNAME}/${service}"
    
    print_info "Pushing $service to Docker Hub..."
    
    # Push both tags
    docker push "${image_name}:${IMAGE_TAG}"
    docker push "${image_name}:latest"
    
    print_success "Pushed ${image_name}"
}

# Function to push all service images
push_all_images() {
    print_info "Pushing all images to Docker Hub..."
    
    for service in "${SERVICES[@]}"; do
        push_service "$service"
    done
    
    print_success "All images pushed successfully"
}

# Function to build and push all images
build_and_push() {
    print_info "Building and pushing all images..."
    
    build_all_images
    push_all_images
    
    print_success "Build and push completed successfully!"
}

# Function to list built images
list_images() {
    print_info "Local Docker images:"
    echo
    
    for service in "${SERVICES[@]}"; do
        image_name="${DOCKER_HUB_USERNAME}/${service}"
        echo "Repository: $image_name"
        docker images "$image_name" --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.Size}}\t{{.CreatedSince}}"
        echo
    done
}

# Function to remove local images
cleanup_local() {
    print_info "Removing local Docker images..."
    
    for service in "${SERVICES[@]}"; do
        image_name="${DOCKER_HUB_USERNAME}/${service}"
        print_info "Removing $image_name images..."
        
        # Remove all tags for this image
        docker rmi "${image_name}:${IMAGE_TAG}" "${image_name}:latest" 2>/dev/null || true
    done
    
    # Clean up dangling images
    docker image prune -f
    
    print_success "Local images cleaned up"
}

# Function to verify images on Docker Hub
verify_hub_images() {
    print_info "Verifying images on Docker Hub..."
    
    for service in "${SERVICES[@]}"; do
        image_name="${DOCKER_HUB_USERNAME}/${service}"
        print_info "Checking $image_name..."
        
        # Try to pull the image info
        if docker manifest inspect "${image_name}:${IMAGE_TAG}" &>/dev/null; then
            print_success "✓ ${image_name}:${IMAGE_TAG} is available on Docker Hub"
        else
            print_warning "✗ ${image_name}:${IMAGE_TAG} not found on Docker Hub"
        fi
    done
}

# Function to update Kubernetes manifests with new image names
update_k8s_manifests() {
    print_info "Updating Kubernetes manifests with Docker Hub image names..."
    
    # Backup original files
    backup_dir="k8s/backup-$(date +%Y%m%d-%H%M%S)"
    mkdir -p "$backup_dir"
    
    for service in "${SERVICES[@]}"; do
        k8s_file="k8s/0*-${service}.yaml"
        if ls $k8s_file 1> /dev/null 2>&1; then
            cp $k8s_file "$backup_dir/"
            
            # Update image name in the manifest
            sed -i.bak "s|image: url-shortener/${service}:latest|image: ${DOCKER_HUB_USERNAME}/${service}:${IMAGE_TAG}|g" $k8s_file
            sed -i.bak "s|image: ${service}:latest|image: ${DOCKER_HUB_USERNAME}/${service}:${IMAGE_TAG}|g" $k8s_file
            
            # Remove backup files created by sed
            rm -f ${k8s_file}.bak
            
            print_success "Updated $(ls $k8s_file)"
        fi
    done
    
    print_info "Original files backed up to: $backup_dir"
    print_success "Kubernetes manifests updated with Docker Hub image names"
}

# Function to show Docker Hub URLs
show_docker_hub_info() {
    print_info "Docker Hub Repository URLs:"
    echo
    
    for service in "${SERVICES[@]}"; do
        echo "🐳 $service: https://hub.docker.com/r/${DOCKER_HUB_USERNAME}/${service}"
    done
    
    echo
    print_info "Image Pull Commands:"
    echo
    
    for service in "${SERVICES[@]}"; do
        echo "docker pull ${DOCKER_HUB_USERNAME}/${service}:${IMAGE_TAG}"
    done
}

# Function to show usage
show_help() {
    echo "Docker Hub Push Script for URL Shortener Microservices"
    echo
    echo "Usage: $0 [command] [options]"
    echo
    echo "Commands:"
    echo "  build              Build all Docker images locally"
    echo "  push               Push images to Docker Hub (requires build first)"
    echo "  deploy             Build and push all images (default)"
    echo "  list               List local Docker images"
    echo "  login              Login to Docker Hub"
    echo "  verify             Verify images exist on Docker Hub"
    echo "  update-k8s         Update Kubernetes manifests with Docker Hub image names"
    echo "  cleanup            Remove local Docker images"
    echo "  info               Show Docker Hub repository information"
    echo "  help               Show this help message"
    echo
    echo "Environment Variables:"
    echo "  IMAGE_TAG          Tag for the images (default: latest)"
    echo "  BUILD_PLATFORM     Target platform for builds (default: linux/amd64,linux/arm64)"
    echo
    echo "Examples:"
    echo "  $0 deploy                    # Build and push all images"
    echo "  $0 build                     # Build images only"
    echo "  $0 push                      # Push existing images" 
    echo "  $0 update-k8s                # Update K8s manifests"
    echo "  IMAGE_TAG=v1.0.0 $0 deploy   # Use specific tag"
    echo
    echo "Docker Hub Username: $DOCKER_HUB_USERNAME"
    echo "Image Tag: $IMAGE_TAG"
}

# Main function
main() {
    local command=${1:-deploy}
    
    # Always check Docker availability
    check_docker
    check_docker_daemon
    
    case $command in
        "build")
            print_info "Building Docker images..."
            build_all_images
            list_images
            ;;
        "push")
            print_info "Pushing images to Docker Hub..."
            check_docker_login
            push_all_images
            verify_hub_images
            show_docker_hub_info
            ;;
        "deploy")
            print_info "Building and pushing all images to Docker Hub..."
            check_docker_login
            build_and_push
            verify_hub_images
            show_docker_hub_info
            print_warning "Don't forget to run: $0 update-k8s"
            ;;
        "list")
            list_images
            ;;
        "login")
            check_docker_login
            ;;
        "verify")
            verify_hub_images
            ;;
        "update-k8s")
            update_k8s_manifests
            ;;
        "cleanup")
            cleanup_local
            ;;
        "info")
            show_docker_hub_info
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            echo "Use '$0 help' for available commands"
            exit 1
            ;;
    esac
}

# Check if script is being run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
