#!/usr/bin/env bash

set -euo pipefail

# Complete single-file deploy script for the URL Shortener
NAMESPACE="url-shortener"
DIR="$(cd "$(dirname "$0")" && pwd)"

MANIFESTS=(
    00-namespace.yaml
    00-secrets.yaml
    01-postgres-config.yaml
    02-postgres.yaml
    03-link-service.yaml
    04-redirect-service.yaml
    05-stats-service.yaml
    07-frontend.yaml
    simple-ingress.yaml
    ingress-controller.yaml
)

# Check if ingress controller is installed

echo "🔧 Applying manifests from $DIR"
for f in "${MANIFESTS[@]}"; do
    if [ -f "$DIR/$f" ]; then
        echo "  applying $f"
        kubectl apply -f "$DIR/$f"
    else
        echo "  skipping missing $f"
    fi
done

echo "⏳ Waiting for main pods to be ready (short wait)"
kubectl wait --for=condition=ready pod -l app=postgres -n $NAMESPACE --timeout=120s || true
kubectl wait --for=condition=ready pod -l app=link-service -n $NAMESPACE --timeout=120s || true
kubectl wait --for=condition=ready pod -l app=redirect-service -n $NAMESPACE --timeout=120s || true
kubectl wait --for=condition=ready pod -l app=stats-service -n $NAMESPACE --timeout=120s || true
kubectl wait --for=condition=ready pod -l app=frontend -n $NAMESPACE --timeout=120s || true

IP=$(kubectl get ingress -n $NAMESPACE -o jsonpath='{.items[0].status.loadBalancer.ingress[0].ip}' 2>/dev/null || true)
if [ -z "$IP" ]; then
    IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[0].address}' 2>/dev/null || echo "localhost")
fi

echo "✅ Deployment applied. Frontend should be at: http://$IP"
echo "Run: kubectl get pods -n $NAMESPACE; kubectl get svc -n $NAMESPACE; kubectl get ingress -n $NAMESPACE"

exit 0
