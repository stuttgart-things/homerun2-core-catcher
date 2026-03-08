# Homerun2 Core Catcher

A Go CLI microservice that consumes messages from Redis Streams using consumer groups and logs them for processing.

## Overview

Core Catcher is part of the homerun2 platform. It connects to a Redis Stream via consumer groups and processes incoming messages enqueued by [homerun2-omni-pitcher](https://github.com/stuttgart-things/homerun2-omni-pitcher) or any producer using the same stream format.

Each stream entry contains a `messageID` referencing a Redis JSON object with the full message payload. The catcher logs every received message with its metadata.

```
omni-pitcher → Redis Stream → core-catcher → slog output
```

## Quick Start

```bash
# Set required environment variables
export REDIS_ADDR=localhost
export REDIS_PORT=6379
export REDIS_STREAM=messages

# Run locally
go run .
```

## Architecture

- **Go** - CLI consumer with graceful shutdown
- **Redis Streams** - Message consumption via `redisqueue` consumer groups
- **homerun-library** - Shared types and helpers
- **ko** - Container image builds (distroless)
- **KCL** - Kubernetes manifest generation
- **Dagger** - CI/CD pipeline functions
