.DEFAULT_GOAL := help
CLUSTER_NAME = ns-dev
PROJECT_ROOT = $(shell dirname $(shell pwd))

test:
	echo $(PROJECT_ROOT)

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: k3d-up
k3d-up: ## Spin up k3d cluster
	k3d cluster create $(CLUSTER_NAME) --image rancher/k3s:latest --no-lb --k3s-arg "--disable=traefik" \
 		-v $(PROJECT_ROOT):/work@all \
 		-v $(PROJECT_ROOT)/.local-dev/k3d-local-storage:/var/lib/rancher/k3s/storage@all \
 		-p 80:80@agent:0 \
 		-p 2201:2201@agent:0

.PHONY: ensure-context
ensure-context:
	kubectl config use-context k3d-$(CLUSTER_NAME)

.PHONY: import
import: ensure-context ## Import ns images into the cluster
	sh -c 'cd ../../ && make build'
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-dashboard:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-auth-dev:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-builder:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-controller:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-gateway:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-gitea-integration:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-migrate:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-ssgen:main

.PHONY: apply
apply: ensure-context ## Apply manifests to k3d cluster
	kustomize build . | kubectl apply -f -
	kubectl rollout restart deployment/coredns -n kube-system # For CoreDNS ConfigMap patch

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

.PHONY: up
up: k3d-up import apply events

.PHONY: down
down: k3d-down
