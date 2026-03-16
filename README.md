# Godo List Backend

REST API built with Go, Gin, PostgreSQL, and Redis. Session-based authentication with Argon2id password hashing.

## Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Just](https://github.com/casey/just) — command runner
- [Docker](https://www.docker.com/) — for infrastructure services

## Getting Started

### 1. Clone the repository

```bash
git clone git@github.com:LuizHVicari/godo-list-backend.git
cd godo-list-backend
```

### 2. Configure environment

Copy `.env.example` to `.env` and fill in the values:

```bash
cp .env.example .env
```

### 3. Start infrastructure

```bash
just infra-up
```

This starts PostgreSQL, Redis, and pgAdmin via Docker Compose.

### 4. Run migrations

```bash
just migrate-up
```

### 5. Start the server

```bash
just run
# or with hot reload:
just dev
```

The API will be available at `http://localhost:8080`.

Swagger UI: `http://localhost:8080/swagger/index.html`

pgAdmin: `http://localhost:5050`

For all available commands, run `just --list`.
