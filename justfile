sqlc_version := "v1.30.0"
goose_version := "v3.27.0"

db_host := env_var_or_default("DB_HOST", "localhost")
db_port := env_var_or_default("DB_PORT", "5432")
db_user := env_var_or_default("DB_USER", "postgres")
db_password := env_var_or_default("DB_PASSWORD", "postgres")
db_name := env_var_or_default("DB_NAME", "postgres")
db_sslmode := env_var_or_default("DB_SSLMODE", "disable")
migrations_dir := env_var_or_default("MIGRATIONS_DIR", "db/migrations")
db_url := "postgres://{{db_user}}:{{db_password}}@{{db_host}}:{{db_port}}/{{db_name}}?sslmode={{db_sslmode}}"

# Compiles and runs the application.
run:
    go run cmd/api/main.go

# Generates the SQL code from the SQL files.
sqlc-generate:
    mkdir -p sql sql/queries
    go run github.com/sqlc-dev/sqlc/cmd/sqlc@{{sqlc_version}} generate

# Verifies that the generated code is up to date with the SQL files.
sqlc-verify:
    mkdir -p sql sql/queries
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
    mkdir -p sql sql/queries
    docker compose exec -T -e PGPASSWORD="{{db_password}}" database pg_dump -h localhost -U "{{db_user}}" -d "{{db_name}}" --schema-only --no-owner --no-privileges > sql/schema.sql

# Brings up the database and cache services using Docker Compose.
infra-up:
    docker compose up -d --build database cache --remove-orphans

# Brings down the database and cache services using Docker Compose.
infra-down:
    docker compose down --remove-orphans

# Brings down the database and cache services using Docker Compose, removing volumes.
infra-reset:
    docker compose down -v database cache --remove-orphans