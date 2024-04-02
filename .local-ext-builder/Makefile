.DEFAULT_GOAL := help
CLUSTER_NAME = ext-builder

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: k3d-up
k3d-up: ## Spin up k3d cluster
	k3d cluster create $(CLUSTER_NAME) --image rancher/k3s:latest --no-lb --k3s-arg "--disable=traefik,servicelb,metrics-server"

.PHONY: ensure-context
ensure-context:
	kubectl config use-context k3d-$(CLUSTER_NAME)

.PHONY: import
import: ensure-context ## Import ns-builder image into the cluster
	sh -c 'cd ../ && make build'
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-builder:main

.PHONY: apply
apply: ensure-context ## Apply manifests to k3d cluster
	kustomize build ./manifest | kubectl apply -f -
	kubectl rollout restart deployment/coredns -n kube-system # For CoreDNS ConfigMap patch

REPLICAS?=1
.PHONY: scale
scale: ensure-context ## Scale the number of builder instances
	kubectl scale deployment/ns-builder --replicas=$(REPLICAS)

.PHONY: events
events: ensure-context ## Tail cluster events
	kubectl get events --watch

.PHONY: logs
logs: ensure-context ## Tail builder logs
	kubectl logs --prefix --timestamps --all-containers --max-log-requests=100 -l app=ns-builder -f

.PHONY: k3d-down
k3d-down: ## Spin down k3d cluster
	k3d cluster delete $(CLUSTER_NAME)

# --- all in one commands ---

.PHONY: up ## All in one up
up: k3d-up import apply events

.PHONY: down ## All in one down
down: k3d-down
