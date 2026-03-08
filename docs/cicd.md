# CI/CD

## GitHub Actions Workflows

| Workflow | Trigger | Description |
|----------|---------|-------------|
| `build-test.yaml` | PR / push to main | Dagger lint + build + test |
| `build-scan-image.yaml` | Push to main | ko build + Trivy scan |
| `release.yaml` | After image build / manual | Semantic release + stage image + push kustomize OCI |
| `pages.yaml` | After release / manual | Deploy MkDocs to GitHub Pages |
| `lint-repo.yaml` | PR / push to main | Repository linting |

## Dagger Functions

The `dagger/` module provides:

| Function | Description |
|----------|-------------|
| `Lint` | Go linting via golangci-lint |
| `Build` | Build Go binary |
| `BuildImage` | Build container image with ko |
| `ScanImage` | Trivy vulnerability scan |
| `BuildAndTestBinary` | Build + Redis integration test |
| `IntegrationTest` | Full e2e test with pitcher + Redis + catcher |

## Taskfile

Common tasks available via `task`:

```bash
task lint                  # Run golangci-lint
task build-test-binary     # Build + test with Redis
task integration-test      # Full e2e test with pitcher
task render-manifests      # Render KCL manifests
task build-scan-image-ko   # Build + scan with ko
task deploy-kcl            # Deploy to cluster
```

## Release Process

Releases are automated via semantic-release:

1. Push to `main` triggers build + image workflow
2. On success, release workflow runs semantic-release
3. If releasable commits exist, a new version tag is created
4. GoReleaser builds binaries for linux/darwin (amd64/arm64)
5. Container image is staged from `:main` to `:vX.Y.Z`
6. Kustomize base is pushed as OCI artifact to GHCR
7. GitHub Pages documentation is deployed
