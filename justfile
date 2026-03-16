set windows-shell := ["C:\\Program Files\\Git\\bin\\sh.exe", "-c"]
set dotenv-load

sqlc_version := "v1.30.0"
goose_version := "v3.27.0"
air_version := "v1.61.7"
golangci_lint_version := "v2.10.1"
swag_version := "v1.16.4"

db_host := env_var_or_default("DB_HOST", "localhost")
db_port := env_var_or_default("DB_PORT", "5432")
db_user := env_var_or_default("DB_USER", "postgres")
db_password := env_var_or_default("DB_PASSWORD", "postgres")
db_name := env_var_or_default("DB_NAME", "postgres")
db_sslmode := env_var_or_default("DB_SSLMODE", "disable")
migrations_dir := env_var_or_default("MIGRATIONS_DIR", "db/migrations")
db_url := "postgres://" + db_user + ":" + db_password + "@" + db_host + ":" + db_port + "/" + db_name + "?sslmode=" + db_sslmode

# Compiles and runs the application.
run port="8080":
    go run cmd/api/main.go -port={{port}}

# Runs all Go tests in the project.
test:
    go test ./...

# Runs go vet checks for suspicious constructs.
vet:
    go vet ./...

# Runs golangci-lint across the project.
lint:
    go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@{{golangci_lint_version}} run ./...

# Runs the standard local quality checks.
check: vet lint test sqlc-verify

# Prints the version of golangci-lint being used.
lint-version:
    go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@{{golangci_lint_version}} version

# Runs the API with hot reload using Air.
dev port="8080":
    mkdir -p tmp
    SERVER_PORT={{port}} go run github.com/air-verse/air@{{air_version}} -c .air.toml

# Prints the version of Air being used.
dev-version:
    go run github.com/air-verse/air@{{air_version}} -v

# Generates Swagger docs from code annotations.
swagger-init:
    mkdir -p docs
    go run github.com/swaggo/swag/cmd/swag@{{swag_version}} init -g main.go -d cmd/api,internal -o docs

# Formats Swagger annotations in Go files.
swagger-fmt:
    go run github.com/swaggo/swag/cmd/swag@{{swag_version}} fmt -g cmd/api/main.go

# Prints the version of swag being used.
swagger-version:
    go run github.com/swaggo/swag/cmd/swag@{{swag_version}} --version

# Compiles the application and outputs the binary to the bin directory.
build:
    mkdir -p bin
    go build -o bin/api cmd/api/main.go

# Compiles the application for Windows and outputs the binary to the bin directory.
build-windows:
    mkdir -p bin
    GOOS=windows GOARCH=amd64 go build -o bin/api.exe cmd/api/main.go

# Compiles the application for Linux and outputs the binary to the bin directory.
build-linux:
    mkdir -p bin
    GOOS=linux GOARCH=amd64 go build -o bin/api cmd/api/main.go

# Compiles the application for macOS and outputs the binary to the bin directory.
build-macos:
    mkdir -p bin
    GOOS=darwin GOARCH=amd64 go build -o bin/api cmd/api/main.go

# Generates the SQL code from the SQL files.
sqlc-generate:
    mkdir -p db/queries
    go run github.com/sqlc-dev/sqlc/cmd/sqlc@{{sqlc_version}} generate

# Verifies that the generated code is up to date with the SQL files.
sqlc-verify:
    mkdir -p db/queries
    go run github.com/sqlc-dev/sqlc/cmd/sqlc@{{sqlc_version}} verify

# Prints the version of sqlc being used.
sqlc-version:
    go run github.com/sqlc-dev/sqlc/cmd/sqlc@{{sqlc_version}} version

# Runs the database migrations.
migrate-up:
    mkdir -p "{{migrations_dir}}"
    GOOSE_DRIVER=postgres GOOSE_DBSTRING="{{db_url}}" go run github.com/pressly/goose/v3/cmd/goose@{{goose_version}} -dir "{{migrations_dir}}" up

# Rolls back the last database migration.
migrate-down:
    mkdir -p "{{migrations_dir}}"
    GOOSE_DRIVER=postgres GOOSE_DBSTRING="{{db_url}}" go run github.com/pressly/goose/v3/cmd/goose@{{goose_version}} -dir "{{migrations_dir}}" down

# Prints the status of the database migrations.
migrate-status:
    mkdir -p "{{migrations_dir}}"
    GOOSE_DRIVER=postgres GOOSE_DBSTRING="{{db_url}}" go run github.com/pressly/goose/v3/cmd/goose@{{goose_version}} -dir "{{migrations_dir}}" status

# Migrates up to a specific migration version.
migrate-up-to version:
    mkdir -p "{{migrations_dir}}"
    GOOSE_DRIVER=postgres GOOSE_DBSTRING="{{db_url}}" go run github.com/pressly/goose/v3/cmd/goose@{{goose_version}} -dir "{{migrations_dir}}" up-to {{version}}

# Creates a new database migration with the given name.
migrate-create name:
    mkdir -p "{{migrations_dir}}"
    GOOSE_DRIVER=postgres GOOSE_DBSTRING="{{db_url}}" go run github.com/pressly/goose/v3/cmd/goose@{{goose_version}} -dir "{{migrations_dir}}" create {{name}} sql

# Dumps the current database schema used by sqlc.
schema-dump:
    mkdir -p db
    docker compose exec -T -e PGPASSWORD="{{db_password}}" database pg_dump -h localhost -U "{{db_user}}" -d "{{db_name}}" --schema-only --no-owner --no-privileges | grep -v '^[\\]' > db/schema.sql

# Brings up the database and cache services using Docker Compose.
infra-up:
    docker compose up -d --build database cache pgadmin --remove-orphans

# Brings down the database and cache services using Docker Compose.
infra-down:
    docker compose down --remove-orphans

# Brings down the database and cache services using Docker Compose, removing volumes.
infra-reset:
    docker compose down -v database cache --remove-orphans