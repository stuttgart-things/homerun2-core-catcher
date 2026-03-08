# CLAUDE.md

## Project

${{ values.fullName }} — ${{ values.description }}

## Tech Stack

- **Language**: Go 1.24+
{%- if values.serviceType == "pitcher" %}
- **HTTP**: stdlib `net/http` (no framework)
- **Queue**: Redis Streams via `homerun-library`
{%- endif %}
{%- if values.serviceType == "catcher" %}
- **Consumer**: Redis Streams via `redisqueue` (consumer groups)
- **Library**: `homerun-library` for shared types and helpers
{%- endif %}
- **Build**: ko (`.ko.yaml`), no Dockerfile
- **CI**: Dagger modules (`dagger/main.go`), Taskfile
- **Infra**: GitHub Actions, semantic-release, renovate

## Git Workflow

**Branch-per-issue with PR and merge.**

### Branch naming

- `fix/<issue-number>-<short-description>` for bugs
- `feat/<issue-number>-<short-description>` for features
- `test/<issue-number>-<short-description>` for test-only changes

### Commit messages

- Use conventional commits: `fix:`, `feat:`, `test:`, `chore:`, `docs:`
- End with `Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>` when Claude authored
- Include `Closes #<issue-number>` to auto-close issues

## Code Conventions

- No Dockerfile — use ko for image builds
- Config via environment variables, loaded once at startup
- Tests: `go test ./...` — unit tests must not require Redis

## Testing

```bash
go test ./...
task build-test-binary
task lint
task build-scan-image-ko
```
