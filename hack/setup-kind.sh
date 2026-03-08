#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="${1:-gitops-dev}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if ! command -v kind &>/dev/null; then
    echo "Error: kind is not installed. Install it with: go install sigs.k8s.io/kind@latest"
    exit 1
fi

if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
    echo "Cluster '${CLUSTER_NAME}' already exists."
    kubectl cluster-info --context "kind-${CLUSTER_NAME}"
    exit 0
fi

echo "Creating kind cluster '${CLUSTER_NAME}'..."
kind create cluster \
    --name "${CLUSTER_NAME}" \
    --config "${SCRIPT_DIR}/kind-config.yaml" \
    --wait 120s

echo ""
echo "Cluster '${CLUSTER_NAME}' is ready."
echo "Context: kind-${CLUSTER_NAME}"
kubectl cluster-info --context "kind-${CLUSTER_NAME}"
