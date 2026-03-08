# CLAUDE.md

## Project

homerun2-core-catcher — Go CLI microservice that consumes messages from Redis Streams using consumer groups and logs/processes them.

## Tech Stack

- **Language**: Go 1.24+
- **Consumer**: Redis Streams via `redisqueue` (consumer groups)
- **Library**: `homerun-library` for shared types and helpers
- **Build**: ko (`.ko.yaml`), no Dockerfile
- **CI**: Dagger modules (`dagger/main.go`), Taskfile
- **Infra**: GitHub Actions, semantic-release, renovate

## Git Workflow

**Branch-per-issue with PR and merge.** Every change gets its own branch, PR, and merge to main.

### Branch naming

- `fix/<issue-number>-<short-description>` for bugs
- `feat/<issue-number>-<short-description>` for features
- `test/<issue-number>-<short-description>` for test-only changes

### Workflow

1. Branch off `main`: `git checkout -b fix/<N>-<desc> main`
2. Make changes, commit with `Closes #<N>` in the message
3. Push: `git push -u origin <branch>`
4. Create PR: `gh pr create --base main`
5. Merge: `gh pr merge <N> --merge --delete-branch`
6. If multiple issues are closely related (e.g., same file), combine into one branch with multiple `Closes #N`

### Commit messages

- Use conventional commits: `fix:`, `feat:`, `test:`, `chore:`, `docs:`
- End with `Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>` when Claude authored
- Include `Closes #<issue-number>` to auto-close issues

## Code Conventions

- No Dockerfile — use ko for image builds
- Config via environment variables, loaded once at startup
- Tests: `go test ./...` — unit tests must not require Redis; integration tests run via Dagger with Redis service
- Catcher interface pattern: pluggable backends (Redis consumer, mock for testing)

## Key Paths

- `main.go` — entrypoint, consumer setup, graceful shutdown
- `internal/catcher/` — Catcher interface, Redis consumer, mock
- `internal/config/` — env-based config loading
- `internal/banner/` — animated TUI startup banner
- `dagger/main.go` — CI functions (Lint, Build, BuildImage, ScanImage, BuildAndTestBinary)
- `Taskfile.yaml` — task runner for build/test/deploy/release
- `.ko.yaml` — ko build configuration
- `.github/workflows/` — CI/CD (build-test, build-scan-image, release, lint-repo)

## Environment Variables

- `REDIS_ADDR` (default: `localhost`) — Redis host
- `REDIS_PORT` (default: `6379`) — Redis port
- `REDIS_PASSWORD` (default: empty) — Redis password
- `REDIS_STREAM` (default: `messages`) — Redis stream to consume from
- `CONSUMER_GROUP` (default: `homerun2-core-catcher`) — Consumer group name
- `CONSUMER_NAME` (default: hostname) — Consumer name within the group
- `LOG_FORMAT` (default: `json`) — Log format (`json` or `text`)
- `LOG_LEVEL` (default: `info`) — Log level (`debug`, `info`, `warn`, `error`)

## Testing

```bash
# Unit tests (no Redis needed)
go test ./...

# Integration test via Dagger (spins up Redis)
task build-test-binary

# Lint
task lint

# Build + scan image
task build-scan-image-ko
```

## Reference Project

`homerun2-omni-pitcher` is the sibling producer service — same patterns for infra, CI, and deployment.
