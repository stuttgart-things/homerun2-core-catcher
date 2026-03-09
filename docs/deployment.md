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

## Flux App Deployment

The recommended way to deploy the full homerun2 stack (Redis Stack + omni-pitcher + core-catcher) is via the [homerun2 Flux app](https://github.com/stuttgart-things/flux/tree/main/apps/homerun2). It uses Kustomize Components to deploy all services into a shared namespace with a single Flux Kustomization.

```yaml
---
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: homerun2-flux
  namespace: flux-system
spec:
  interval: 1h
  retryInterval: 1m
  timeout: 5m
  sourceRef:
    kind: GitRepository
    name: flux-apps
  path: ./apps/homerun2
  prune: true
  wait: true
  postBuild:
    substitute:
      HOMERUN2_NAMESPACE: homerun2-flux
      HOMERUN2_CORE_CATCHER_VERSION: v0.5.0
      HOMERUN2_CORE_CATCHER_KUSTOMIZE_VERSION: v0.5.0-web
      HOMERUN2_CORE_CATCHER_HOSTNAME: catcher
      HOMERUN2_OMNI_PITCHER_VERSION: v1.2.0
      HOMERUN2_OMNI_PITCHER_HOSTNAME: pitcher
      GATEWAY_NAME: my-gateway
      GATEWAY_NAMESPACE: default
      DOMAIN: my-cluster.example.com
      HOMERUN2_REDIS_VERSION: "17.1.4"
      HOMERUN2_REDIS_STORAGE_CLASS: nfs4-csi
      HOMERUN2_REDIS_STORAGE_SIZE: 8Gi
    substituteFrom:
      - kind: Secret
        name: homerun2-flux-secrets
```

The core-catcher component patches the KCL base to:

- Set `CATCHER_MODE=web` for the HTMX dashboard
- Override the Redis connection to point to the co-deployed redis-stack
- Patch the Redis password secret with the correct credentials
- Replace the KCL-generated HTTPRoute with a component-level one (custom hostname)

See the [Flux app README](https://github.com/stuttgart-things/flux/tree/main/apps/homerun2) for all substitution variables and a complete example.

## Container Image

Built with [ko](https://ko.build/) using a distroless base image (`cgr.dev/chainguard/static:latest`):

```bash
# Build locally
ko build .

# Build via Taskfile
task build-scan-image-ko
```
