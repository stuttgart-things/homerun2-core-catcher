# homerun2-core-catcher

A Go CLI microservice that consumes messages from Redis Streams using consumer groups and logs them for processing.

[![Build & Test](https://github.com/stuttgart-things/homerun2-core-catcher/actions/workflows/build-test.yaml/badge.svg)](https://github.com/stuttgart-things/homerun2-core-catcher/actions/workflows/build-test.yaml)

## How It Works

The catcher connects to a Redis Stream (default: `messages`) via a consumer group and processes incoming messages. Messages are enqueued by [homerun2-omni-pitcher](https://github.com/stuttgart-things/homerun2-omni-pitcher) (or any producer using the same stream format).

Each stream entry contains a `messageID` referencing a Redis JSON object with the full message payload. The catcher logs every received message with its metadata.

```
omni-pitcher → Redis Stream → core-catcher → slog output
```

## Deployment

<details>
<summary><b>Run locally</b></summary>

```bash
# Start Redis (via Dagger)
task run-redis-as-service

# Run the catcher
export REDIS_ADDR=localhost REDIS_PORT=6379 REDIS_STREAM=messages
go run .
```

</details>

<details>
<summary><b>Container image (ko / ghcr.io)</b></summary>

The container image is built with [ko](https://ko.build) on top of `cgr.dev/chainguard/static` and published to GitHub Container Registry.

```bash
# Pull the image
docker pull ghcr.io/stuttgart-things/homerun2-core-catcher:<tag>

# Run with Docker
docker run \
  -e REDIS_ADDR=redis -e REDIS_PORT=6379 \
  -e REDIS_STREAM=messages \
  ghcr.io/stuttgart-things/homerun2-core-catcher:<tag>
```

</details>

<details>
<summary><b>Deploy Redis (prerequisite)</b></summary>

```bash
helmfile apply -f \
  git::https://github.com/stuttgart-things/helm.git@database/redis-stack.yaml.gotmpl \
  --state-values-set storageClass=openebs-hostpath \
  --state-values-set password="<REPLACE>" \
  --state-values-set namespace=homerun2
```

</details>

<details>
<summary><b>Run locally against a remote Kubernetes Redis</b></summary>

Port-forward the Redis service from your cluster and run the catcher locally:

```bash
# Port-forward Redis (keep running in a separate terminal)
export KUBECONFIG=~/.kube/<your-kubeconfig>
kubectl port-forward -n redis-stack svc/redis-stack 6379:6379

# Get the Redis password from the cluster secret
kubectl get secret -n redis-stack redis-stack \
  -o jsonpath='{.data.redis-password}' | base64 -d

# Run the catcher (use a script to avoid shell escaping issues with special characters in the password)
cat > /tmp/run-catcher.sh << 'EOF'
#!/bin/bash
export REDIS_ADDR=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD='<REPLACE>'
export REDIS_STREAM=messages
export LOG_FORMAT=text
go run .
EOF
bash /tmp/run-catcher.sh
```

> **Note:** Passwords containing `!` or other shell-special characters must be set inside a script with single-quoted heredoc (`<< 'EOF'`). Passing them directly via `export` or inline env vars can cause silent corruption from bash history expansion.

</details>

## Development

<details>
<summary><b>Project structure</b></summary>

```
main.go                    # Entrypoint, consumer setup, graceful shutdown
internal/
  banner/                  # Animated startup banner (Bubble Tea)
  config/                  # Env-based config loading, slog setup
  catcher/                 # Catcher interface (Redis consumer + Mock)
dagger/                    # CI functions (Lint, Build, Test, Scan)
.ko.yaml                   # ko build configuration
Taskfile.yaml              # Task runner
```

</details>

<details>
<summary><b>Configuration reference</b></summary>

| Variable | Description | Default |
|----------|-------------|---------|
| `REDIS_ADDR` | Redis server address | `localhost` |
| `REDIS_PORT` | Redis server port | `6379` |
| `REDIS_PASSWORD` | Redis password | (empty) |
| `REDIS_STREAM` | Redis stream to consume from | `messages` |
| `CONSUMER_GROUP` | Consumer group name | `homerun2-core-catcher` |
| `CONSUMER_NAME` | Consumer name within the group | hostname |
| `LOG_FORMAT` | Log format: `json` or `text` | `json` |
| `LOG_LEVEL` | Log level: `debug`, `info`, `warn`, `error` | `info` |

</details>

## Testing

<details>
<summary><b>Unit tests</b></summary>

Unit tests run without Redis:

```bash
go test ./...
```

</details>

<details>
<summary><b>Integration tests (Dagger + Redis)</b></summary>

Integration tests spin up a Redis service via Dagger:

```bash
task build-test-binary
```

</details>

<details>
<summary><b>Lint</b></summary>

```bash
task lint
```

</details>

<details>
<summary><b>Build and scan container image</b></summary>

```bash
task build-scan-image-ko
```

</details>

## Links

- [Releases](https://github.com/stuttgart-things/homerun2-core-catcher/releases)
- [Container Images](https://github.com/stuttgart-things/homerun2-core-catcher/pkgs/container/homerun2-core-catcher)
- [homerun2-omni-pitcher](https://github.com/stuttgart-things/homerun2-omni-pitcher) (producer)
- [homerun-library](https://github.com/stuttgart-things/homerun-library)

## License

See [LICENSE](LICENSE) file.
