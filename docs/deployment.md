# Deployment

## Kubernetes Manifests (KCL)

Manifests are generated using KCL in the `kcl/` directory. The modular structure:

| File               | Resource        |
|--------------------|-----------------|
| `schema.k`         | Config schema with validation |
| `labels.k`         | Common labels and config instantiation |
| `namespace.k`      | Namespace       |
| `serviceaccount.k` | ServiceAccount  |
| `configmap.k`      | ConfigMap       |
| `secret.k`         | Secret (Redis password) |
| `deploy.k`         | Deployment      |
| `main.k`           | Entry point     |

## Render Manifests

```bash
# Using Taskfile (interactive)
task render-manifests

# Using Taskfile (non-interactive)
task render-manifests-quick

# Using KCL directly
kcl kcl/main.k -Y tests/kcl-deploy-profile.yaml
```

## Configuration Options

Override via KCL profile file or CLI options:

```yaml
config.image: ghcr.io/stuttgart-things/homerun2-core-catcher:v0.1.0
config.namespace: homerun2
config.redisAddr: redis-stack.redis-stack.svc.cluster.local
config.redisPort: "6379"
config.redisStream: messages
config.consumerGroup: homerun2-core-catcher
config.redisPassword: my-secret-password
```

## Deploy to Cluster

```bash
# Push kustomize base as OCI artifact
task push-kustomize-base

# Deploy via Dagger blueprint
task deploy-kcl

# Deploy with custom parameters
task deploy-kcl KUBECONFIG=~/.kube/my-cluster NAMESPACE=my-namespace
```

## Kustomize OCI Pipeline

Releases push a kustomize base as an OCI artifact:

```bash
# Pull the base
oras pull ghcr.io/stuttgart-things/homerun2-core-catcher-kustomize:v0.1.0

# Apply with overlays
kubectl apply -k .
```

## Container Image

Built with [ko](https://ko.build/) using a distroless base image (`cgr.dev/chainguard/static:latest`):

```bash
# Build locally
ko build .

# Build via Taskfile
task build-scan-image-ko
```
