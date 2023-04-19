# manifest

Manifest files required to deploy NeoShowcase locally

## bootstrap

1. Install k3s (or any other k8s installation)
   - In case of k3s, do not forget to disable default traefik installation, otherwise traefik pod will not be able to bind to port 80, 443
2. Install ArgoCD
   - `kubectl create ns argocd`
   - `kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.6.7/manifests/install.yaml`
     - refer to `./argocd/kustomization.yaml` for the current version
   - `kubectl port-forward svc/argocd-server -n argocd 8080:443`
3. Access localhost:8080
   - Get admin password from ` kubectl get secret -n argocd argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 --decode && echo`
4. Add initial `applications` application
   - Add known hosts and connect repository
   - Add application (path: `applications`)
5. Build and import images with `make build && make k3s-import` at root of the repository
6. Sync other applications (and optionally, change their target revision)
7. Access `cd.local.trapti.tech` and more

## local registry setup

When using local (insecure) registry url which points to `127.0.0.1`,
k8s node needs to be configured to pull from the correct endpoint.

### k3s example

`/etc/rancher/k3s/registries.yaml`
```yaml
mirrors:
  registry.ns-system.svc.cluster.local:
    endpoint:
      # set it to actual registry svc IP
      - "http://10.43.135.226"
```

and `sudo systemctl restart k3s`
