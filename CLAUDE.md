# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

Uses `just` as the command runner (see `justfile`).

```bash
just run          # Compile and run (default port 8080)
just dev          # Run with hot reload (Air)
just test         # go test -race -count=1 ./...
just lint         # golangci-lint
just vet          # go vet
just check        # vet + lint + test + sqlc-verify

just migrate-up                # Apply pending migrations
just migrate-down              # Rollback last migration
just migrate-status            # Show migration state
just migrate-create <name>     # Create new migration file

just sqlc-generate             # Regenerate Go code from SQL queries
just sqlc-verify               # Verify generated code is current
just swagger-init              # Regenerate Swagger docs from annotations

just infra-up                  # Start PostgreSQL, Redis, pgAdmin (Docker)
just infra-down                # Stop containers
just infra-reset               # Stop and wipe volumes
```

## Architecture

Layered architecture with: **Handler → Service → Repository**, one package per domain.

```
cmd/api/main.go           # Entry point: wires config, DB, Redis, Gin router, services
internal/
  auth/                   # Session-based auth (sign-up, sign-in, sign-out, middleware)
  user/                   # User entity and password management
  project/                # Project CRUD
  step/                   # Steps within projects (with position ordering)
  item/                   # Items within steps (with position ordering + priority enum)
  platform/
    config/               # Env config loading (.env via godotenv)
    crypto/               # Argon2id password hashing
    db/                   # SQLC-generated models and queries (do not edit manually)
    http/                 # ErrorMapper, request logger middleware
db/
  migrations/             # Goose SQL migrations
  queries/                # SQLC source SQL queries
  schema.sql              # Current schema snapshot
docs/swagger.yaml         # Generated Swagger spec (do not edit manually)
```

Each domain package contains: `handler.go`, `service.go`, `repository.go`, `{entity}.go` (domain model), `dto.go`, `errors.go`.

## Key Patterns

**Error handling:** Each domain defines sentinel errors in `errors.go`. Handlers use `ErrorMapper` from `internal/platform/http/apierror.go` to map domain errors to HTTP status codes:
```go
var errMapper = platformHTTP.NewErrorMapper(
    platformHTTP.E(ErrProjectNotFound, http.StatusNotFound, "project not found"),
)
// In handler:
errMapper.Respond(c, err, "fallback message")
```

**Database access:** All SQL is written in `db/queries/*.sql` and compiled via SQLC into `internal/platform/db/`. Never write raw SQL in Go files. After modifying queries, run `just sqlc-generate`.

**Migrations:** Use Goose format. After creating a migration with `just migrate-create <name>`, apply with `just migrate-up`. The schema has two PostgreSQL schemas: `auth` (users) and `todo` (projects, steps, items).

**Authentication:** Session-based with Redis. The `auth.Middleware` reads the `session_id` cookie, validates against Redis, refreshes TTL, and injects the session into the Gin context. Protected routes are grouped under `/v1/*` behind this middleware.

**Position ordering:** Steps and Items have integer `position` fields with `DEFERRABLE INITIALLY DEFERRED` unique constraints on `(parent_id, position)` to allow bulk repositioning within a transaction.

**Authorization:** Services validate ownership — users can only access resources they own. The repository layer provides `IsStepInOwnedProject` / `IsItemInOwnedStep` style helpers used by services before mutations.

## Environment

Copy `.env.example` to `.env`. Key variables: `DB_HOST/PORT/USER/PASSWORD/NAME`, `CACHE_HOST/PORT`, `SERVER_PORT`, `COOKIE_SECURE`, `CORS_ALLOWED_ORIGIN`.

