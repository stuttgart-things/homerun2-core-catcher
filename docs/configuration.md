# Configuration

## Environment Variables

All configuration is via environment variables, loaded once at startup.

| Variable         | Default                    | Description                          |
|------------------|----------------------------|--------------------------------------|
| `REDIS_ADDR`     | `localhost`                | Redis server address                 |
| `REDIS_PORT`     | `6379`                     | Redis server port                    |
| `REDIS_PASSWORD` | (empty)                    | Redis password                       |
| `REDIS_STREAM`   | `messages`                 | Redis stream to consume from         |
| `CONSUMER_GROUP`| `homerun2-core-catcher`     | Consumer group name                  |
| `CONSUMER_NAME` | hostname                    | Consumer name within the group       |
| `REDIS_STARTUP_TIMEOUT` | `120s`              | How long startup retries Redis before exiting (Go duration) |
| `BASE_PATH`     | *(empty)*                  | URL prefix of the web UI behind a reverse proxy, e.g. `/core` behind homerun2-portal; `/health` stays at the root (web mode) |
| `LOG_FORMAT`     | `json`                     | Log format: `json` or `text`         |
| `LOG_LEVEL`      | `info`                     | Log level: `debug`, `info`, `warn`, `error` |

## Serving Under a Path Prefix

In web mode, `BASE_PATH=/core` serves the dashboard under `/core/`: for
[homerun2-portal](https://github.com/stuttgart-things/homerun2-portal), which
forwards `/core/...` to the service unchanged, or for one gateway hostname
shared with other services.

- Every URL the dashboard hands to the browser (htmx requests, export links,
  static files) carries the prefix.
- `/` and `/core` redirect to `/core/`.
- `/health` answers at the root as well, so the probes and the portal's
  status check need no prefix.
- The value is normalised (`core/` → `/core`, `/` → root). Anything but plain
  path segments (letters, digits, `.` `_` `~` `-`) fails startup.
- In Kubernetes, `config.basePath` sets `BASE_PATH` and the HTTPRoute's
  `PathPrefix` (see [Deployment](deployment.md)).

## Consumer Behavior

The catcher uses Redis consumer groups (`XREADGROUP`) to ensure reliable message delivery:

- Messages are acknowledged after processing
- Multiple catcher instances can share the same consumer group for load balancing
- Each instance should have a unique `CONSUMER_NAME`

## Running Locally Against Remote Redis

Port-forward Redis from a Kubernetes cluster:

```bash
kubectl port-forward -n redis-stack svc/redis-stack 6379:6379

# Use a script to avoid shell escaping issues with special characters
cat > /tmp/run-catcher.sh << 'EOF'
#!/bin/bash
export REDIS_ADDR=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD='<your-password>'
export REDIS_STREAM=messages
export LOG_FORMAT=text
go run .
EOF
bash /tmp/run-catcher.sh
```

!!! note
    Passwords containing `!` or other shell-special characters must be set inside a script with single-quoted heredoc to avoid bash history expansion.
