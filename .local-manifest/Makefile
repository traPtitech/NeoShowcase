.DEFAULT_GOAL := help
CLUSTER_NAME = ns-dev
PROJECT_ROOT = $(shell dirname $(shell pwd))

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ---- k3d commands ----

.PHONY: k3d-up
k3d-up: ## Spin up k3d cluster
	k3d cluster create $(CLUSTER_NAME) --image rancher/k3s:latest --k3s-arg '--disable=traefik@server:*' \
		--host-alias 10.43.11.193:registry.local \
 		-v $(PROJECT_ROOT):/work@all \
 		-v $(PROJECT_ROOT)/.local-dev/k3d-local-storage:/var/lib/rancher/k3s/storage@all \
 		-v ~/.ssh/known_hosts:/known_hosts@all \
 		-p 80:80@loadbalancer \
 		-p 443:443@loadbalancer \
 		-p 2201:2201@loadbalancer \
 		-p 8080:8080@loadbalancer

.PHONY: k3d-down
k3d-down: ## Spin down k3d cluster
	k3d cluster delete $(CLUSTER_NAME)

# ---- Startup comands ----

.PHONY: ensure-context
ensure-context:
	kubectl config use-context k3d-$(CLUSTER_NAME)

.PHONY: import
import: ensure-context ## Import ns images into the cluster
	sh -c 'cd ../ && make build'
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-dashboard:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-auth-dev:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-builder:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-controller:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-gateway:main
	# Uncomment if testing gitea integration
	# k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-gitea-integration:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-migrate:main
	k3d image import --cluster $(CLUSTER_NAME) ghcr.io/traptitech/ns-ssgen:main

.PHONY: apply-namespaces
apply-namespaces:
	ls -d */ | tr -d '/' | grep -v 'coredns-patch' | xargs -I{} sh -c 'kubectl create namespace {} --dry-run=client -o yaml | kubectl apply -f -'

.PHONY: apply-crd
apply-crd:
	kustomize build ./crd | kubectl apply -f -

.PHONY: apply
apply: ensure-context apply-namespaces apply-crd ## Apply manifests to k3d cluster
	kustomize build . | kubectl apply -f -
	kubectl rollout restart deployment/coredns -n kube-system # For CoreDNS ConfigMap patch

# ---- Dev commands ----

.PHONY: events
events: ensure-context ## Tail cluster events
	kubectl get events --watch --all-namespaces

NAMESPACE=ns-system
APP=ns-controller
.PHONY: logs
logs: ensure-context ## Tail builder logs
	kubectl logs --prefix --timestamps --all-containers --max-log-requests=100 -n $(NAMESPACE) -l app=$(APP) -f

# ---- all in one commands ----

.PHONY: up
up: k3d-up import apply events

.PHONY: down
down: k3d-down
