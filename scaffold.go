// Copyright 2025 β ORI Inc.Canada All Rights Reserved.
// Author: Awase Khirni Syed
// Product Name: β ORI Inc. Hexagonal Architecture Generator in Golang
// Golang Onion Architecture/Hexagonal Architecture Scaffolder Project
// This is base project that is used as a scaffolder to automatically generate artefacts by connecting to database
// this would help us render MVC style web api with end-points to render json web api later on.
// this is base project enhanced script.

// scaffold.go
//
// Generates the complete hexagonal architecture project skeleton
// (directories + starter files) exactly as described in the project tree.
//
// Usage:
//
//	go run scaffold.go                                     # ./hexagonal_architecture_golang
//	go run scaffold.go -root myservice                    # custom target dir
//	go run scaffold.go -module github.com/acme/svc        # custom module path
//	go run scaffold.go -dry-run                           # preview only
//	go run scaffold.go -force                             # overwrite existing files
//
// Generated Go files are compilable stubs, so `go build ./...` passes
// immediately after `go mod tidy`.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ─────────────────────────────── flags ────────────────────────────────

var (
	rootDir    = flag.String("root", "hexagonal_architecture_golang", "target directory (created if missing)")
	moduleName = flag.String("module", "hexagonal_architecture_golang", "Go module path written into go.mod")
	force      = flag.Bool("force", false, "overwrite existing files")
	dryRun     = flag.Bool("dry-run", false, "preview without writing anything")
)

// ───────────────────────── path shorthands ────────────────────────────

const (
	cmdAPI     = "cmd/api/"
	cmdWorker  = "cmd/worker/"
	cmdMigrate = "cmd/migrate/"
	cmdHealth  = "cmd/healthcheck/"

	userEntity = "internal/domain/user/entity/"
	userVO     = "internal/domain/user/valueobject/"
	userEvent  = "internal/domain/user/event/"
	userPort   = "internal/domain/user/port/"
	userErr    = "internal/domain/user/error/"

	authEntity = "internal/domain/auth/entity/"
	authEvent  = "internal/domain/auth/event/"
	authPort   = "internal/domain/auth/port/"
	authErr    = "internal/domain/auth/error/"

	sharedVO    = "internal/domain/shared/valueobject/"
	sharedErr   = "internal/domain/shared/error/"
	sharedEvent = "internal/domain/shared/event/"
	sharedPort  = "internal/domain/shared/port/"

	appUserCmd = "internal/application/user/command/"
	appUserQry = "internal/application/user/query/"
	appUserDTO = "internal/application/user/dto/"
	appAuthCmd = "internal/application/auth/command/"
	appAuthQry = "internal/application/auth/query/"
	appAuthDTO = "internal/application/auth/dto/"
	appDecor   = "internal/application/shared/decorator/"
	appHandler = "internal/application/shared/handler/"

	httpUser   = "internal/adapter/inbound/http/v1/user/"
	httpAuth   = "internal/adapter/inbound/http/v1/auth/"
	httpHealth = "internal/adapter/inbound/http/v1/health/"
	httpMw     = "internal/adapter/inbound/http/middleware/"
	httpDocs   = "internal/adapter/inbound/http/docs/"
	httpSpec   = "internal/adapter/inbound/http/docs/spec/"
	httpRoot   = "internal/adapter/inbound/http/"

	grpcRoot  = "internal/adapter/inbound/grpc/"
	grpcProto = "internal/adapter/inbound/grpc/proto/"
	cliRoot   = "internal/adapter/inbound/cli/"

	pgRoot      = "internal/adapter/outbound/persistence/postgres/"
	pgModel     = pgRoot + "model/"
	pgMapper    = pgRoot + "mapper/"
	pgRepo      = pgRoot + "repository/"
	pgMigration = pgRoot + "migration/"

	myRoot      = "internal/adapter/outbound/persistence/mysql/"
	myModel     = myRoot + "model/"
	myMapper    = myRoot + "mapper/"
	myRepo      = myRoot + "repository/"
	myMigration = myRoot + "migration/"

	inMemRepo = "internal/adapter/outbound/persistence/inmemory/"

	cacheRedis = "internal/adapter/outbound/cache/redis/"
	cacheMem   = "internal/adapter/outbound/cache/inmemory/"

	secJWT    = "internal/adapter/outbound/security/jwt/"
	secBcrypt = "internal/adapter/outbound/security/bcrypt/"

	msgKafka       = "internal/adapter/outbound/messaging/kafka/"
	msgRabbit      = "internal/adapter/outbound/messaging/rabbitmq/"
	msgRedisStream = "internal/adapter/outbound/messaging/redis_streams/"
	msgMem         = "internal/adapter/outbound/messaging/inmemory/"

	outLog = "internal/adapter/outbound/logging/"
	outObs = "internal/adapter/outbound/observability/"

	cfg  = "internal/config/"
	cont = "internal/container/"

	pkgHTTP      = "pkg/httputil/"
	pkgValidator = "pkg/validation/"
	pkgLogger    = "pkg/logger/"
	pkgTelemetry = "pkg/telemetry/"

	apiOpenapi = "api/openapi/"
	apiProto   = "api/proto/"

	testIntegration = "test/integration/"
	testFunctional  = "test/functional/"
	testFixture     = "test/fixture/"
	testMock        = "test/mock/"

	scriptsDir = "scripts/"

	dockerDir    = "deployments/docker/"
	k8sDir       = "deployments/kubernetes/"
	helmDir      = "deployments/helm/hexagonal-api/"
	helmTemplate = helmDir + "templates/"

	docsDir = "docs/"
	adrDir  = "docs/adr/"
)

// spec describes one file to generate: its path (forward slashes) and a
// short purpose description (taken from the project tree).
type spec struct {
	path string
	desc string
}

// ─────────────────────────── the file list ────────────────────────────

var files = []spec{
	// ── entrypoints ──────────────────────────────────────────────────
	{cmdAPI + "main.go", "Main entry: wires DI, starts HTTP server"},
	{cmdWorker + "main.go", "Consumes events from message bus"},
	{cmdMigrate + "main.go", "Runs golang-migrate migrations"},
	{cmdHealth + "main.go", "CLI health probe"},

	// ── domain / user ────────────────────────────────────────────────
	{userEntity + "user.go", "User entity with behavior methods"},
	{userEntity + "user_test.go", "Unit tests for entity behavior"},
	{userEntity + "role.go", "Role enum + role hierarchy logic"},
	{userVO + "email.go", "Email VO with validation"},
	{userVO + "email_test.go", "Unit tests for Email VO"},
	{userVO + "password.go", "Password VO with complexity rules"},
	{userVO + "password_test.go", "Unit tests for Password VO"},
	{userVO + "user_status.go", "UserStatus enum"},
	{userEvent + "user_created.go", "UserCreated domain event"},
	{userEvent + "user_updated.go", "UserUpdated domain event"},
	{userEvent + "user_deleted.go", "UserDeleted domain event"},
	{userPort + "repository.go", "UserRepository interface"},
	{userPort + "password_service.go", "PasswordService interface"},
	{userPort + "event_publisher.go", "EventPublisher interface"},
	{userErr + "errors.go", "Domain-specific errors (sentinel)"},

	// ── domain / auth ────────────────────────────────────────────────
	{authEntity + "session.go", "Session entity (JWT claims)"},
	{authEvent + "user_logged_in.go", "UserLoggedIn domain event"},
	{authEvent + "user_registered.go", "UserRegistered domain event"},
	{authPort + "auth_repository.go", "AuthRepository interface"},
	{authPort + "token_service.go", "TokenServicePort interface"},
	{authErr + "errors.go", "Auth domain errors (sentinel)"},

	// ── domain / shared ──────────────────────────────────────────────
	{sharedVO + "id.go", "UUID wrapper VO"},
	{sharedVO + "pagination.go", "Pagination VO"},
	{sharedErr + "errors.go", "Base error types"},
	{sharedEvent + "domain_event.go", "Base DomainEvent interface"},
	{sharedPort + "unit_of_work.go", "UnitOfWork interface"},

	// ── application / user ───────────────────────────────────────────
	{appUserCmd + "create_user.go", "CreateUserUseCase"},
	{appUserCmd + "create_user_test.go", "Tests for CreateUserUseCase"},
	{appUserCmd + "update_user.go", "UpdateUserUseCase"},
	{appUserCmd + "update_user_test.go", "Tests for UpdateUserUseCase"},
	{appUserCmd + "delete_user.go", "DeleteUserUseCase"},
	{appUserCmd + "delete_user_test.go", "Tests for DeleteUserUseCase"},
	{appUserQry + "get_user_by_id.go", "GetUserByIdQuery"},
	{appUserQry + "get_user_by_id_test.go", "Tests for GetUserByIdQuery"},
	{appUserQry + "list_users.go", "ListUsersQuery (paginated)"},
	{appUserQry + "list_users_test.go", "Tests for ListUsersQuery"},
	{appUserDTO + "create_user_dto.go", "CreateUser DTO"},
	{appUserDTO + "update_user_dto.go", "UpdateUser DTO"},
	{appUserDTO + "user_response_dto.go", "UserResponse DTO"},

	// ── application / auth ───────────────────────────────────────────
	{appAuthCmd + "register_user.go", "RegisterUserUseCase"},
	{appAuthCmd + "login_user.go", "LoginUserUseCase"},
	{appAuthCmd + "logout_user.go", "LogoutUserUseCase"},
	{appAuthQry + "get_current_user.go", "GetCurrentUserQuery"},
	{appAuthQry + "verify_token.go", "VerifyTokenQuery"},
	{appAuthDTO + "login_dto.go", "Login DTOs"},
	{appAuthDTO + "register_dto.go", "Register DTOs"},

	// ── application / shared ─────────────────────────────────────────
	{appDecor + "cache.go", "CacheResult decorator"},
	{appDecor + "audit.go", "AuditLog decorator"},
	{appDecor + "validate.go", "ValidateInput decorator"},
	{appHandler + "event_handler.go", "Base event handler interface"},

	// ── inbound / http ───────────────────────────────────────────────
	{httpUser + "handler.go", "User HTTP handlers"},
	{httpUser + "handler_test.go", "Tests for user HTTP handlers"},
	{httpUser + "dto.go", "Request/Response DTOs"},
	{httpUser + "mapper.go", "DTO <-> Domain mappers"},
	{httpAuth + "handler.go", "Auth HTTP handlers"},
	{httpAuth + "dto.go", "Auth request/response DTOs"},
	{httpAuth + "mapper.go", "Auth DTO <-> Domain mappers"},
	{httpHealth + "handler.go", "/healthz endpoints"},
	{httpMw + "auth.go", "JWT validation middleware"},
	{httpMw + "error_handler.go", "Global error -> JSON"},
	{httpMw + "logging.go", "Structured request logging"},
	{httpMw + "security_headers.go", "CSP, HSTS, CORS"},
	{httpMw + "rate_limit.go", "Rate limiting middleware"},
	{httpMw + "tracing.go", "OpenTelemetry span injection"},
	{httpMw + "recovery.go", "Panic recovery"},
	{httpDocs + "swagger.go", "Swagger UI handler"},
	{httpDocs + "redoc.go", "ReDoc handler"},
	{httpSpec + "docs.go", "Generated swagger docs"},
	{httpSpec + "swagger.json", "OpenAPI 3.0 spec"},
	{httpSpec + "swagger.yaml", "OpenAPI 3.0 spec (YAML)"},
	{httpRoot + "router.go", "Main router setup"},
	{httpRoot + "server.go", "HTTP server config (graceful shutdown)"},

	// ── inbound / grpc ───────────────────────────────────────────────
	{grpcProto + "user.proto", "gRPC user service definition"},
	{grpcRoot + "handler.go", "gRPC user service handlers"},
	{grpcRoot + "mapper.go", "Proto <-> Domain mappers"},

	// ── inbound / cli ────────────────────────────────────────────────
	{cliRoot + "user_commands.go", "Admin CLI commands for users"},

	// ── outbound / persistence / postgres ────────────────────────────
	{pgModel + "user_model.go", "GORM/sqlc model"},
	{pgMapper + "user_mapper.go", "Model <-> Domain mapper"},
	{pgRepo + "user_repository.go", "Implements domain.UserRepository"},
	{pgRepo + "user_repository_test.go", "Tests for Postgres UserRepository"},
	{pgRepo + "auth_repository.go", "Implements domain.AuthRepository"},
	{pgRepo + "auth_repository_test.go", "Tests for Postgres AuthRepository"},
	{pgMigration + "000001_init.up.sql", "Initial schema"},
	{pgMigration + "000001_init.down.sql", "Rollback initial schema"},
	{pgMigration + "000002_seed.up.sql", "Seed data"},
	{pgRoot + "unit_of_work.go", "Implements UnitOfWork"},
	{pgRoot + "connection.go", "DB connection pool (pgx/sqlx)"},

	// ── outbound / persistence / mysql (mirror of postgres) ──────────
	{myModel + "user_model.go", "MySQL model"},
	{myMapper + "user_mapper.go", "Model <-> Domain mapper"},
	{myRepo + "user_repository.go", "Implements domain.UserRepository (MySQL)"},
	{myRepo + "user_repository_test.go", "Tests for MySQL UserRepository"},
	{myMigration + "000001_init.up.sql", "Initial schema (MySQL)"},
	{myMigration + "000001_init.down.sql", "Rollback initial schema (MySQL)"},

	// ── outbound / persistence / inmemory ────────────────────────────
	{inMemRepo + "user_repository.go", "In-memory UserRepository (for tests)"},

	// ── outbound / cache ─────────────────────────────────────────────
	{cacheRedis + "redis_cache.go", "Implements CachePort"},
	{cacheRedis + "redis_cache_test.go", "Tests for Redis cache"},
	{cacheMem + "lru_cache.go", "In-memory LRU cache (for tests)"},

	// ── outbound / security ──────────────────────────────────────────
	{secJWT + "jwt_service.go", "Implements TokenServicePort"},
	{secJWT + "jwt_service_test.go", "Tests for JWT service"},
	{secBcrypt + "bcrypt_service.go", "Implements PasswordService"},
	{secBcrypt + "bcrypt_service_test.go", "Tests for bcrypt service"},

	// ── outbound / messaging ─────────────────────────────────────────
	{msgKafka + "kafka_publisher.go", "Implements EventPublisher (Kafka)"},
	{msgRabbit + "rabbit_publisher.go", "Implements EventPublisher (RabbitMQ)"},
	{msgRedisStream + "redis_publisher.go", "Implements EventPublisher (Redis Streams)"},
	{msgMem + "event_bus.go", "In-memory event bus (for tests)"},

	// ── outbound / logging ───────────────────────────────────────────
	{outLog + "logger.go", "Structured logger (zerolog/zap)"},
	{outLog + "pii_redactor.go", "PII redaction hook"},
	{outLog + "context_logger.go", "Request-scoped logger with trace_id"},

	// ── outbound / observability ─────────────────────────────────────
	{outObs + "metrics.go", "Prometheus metrics"},
	{outObs + "tracing.go", "OpenTelemetry tracing"},
	{outObs + "health.go", "Health check registry"},

	// ── config & container ───────────────────────────────────────────
	{cfg + "config.go", "Config struct + loader (viper/envconfig)"},
	{cfg + "config_test.go", "Tests for config loading"},
	{cfg + "database.go", "Database config"},
	{cfg + "server.go", "Server config"},
	{cfg + "security.go", "Security config"},
	{cfg + "observability.go", "Observability config"},
	{cont + "container.go", "Wire/DI setup (google/wire or manual)"},
	{cont + "providers.go", "Provider functions for each component"},

	// ── pkg ──────────────────────────────────────────────────────────
	{pkgHTTP + "response.go", "Standard HTTP response helpers"},
	{pkgHTTP + "errors.go", "HTTP error mapping"},
	{pkgHTTP + "pagination.go", "Pagination helpers"},
	{pkgValidator + "validator.go", "Custom validation rules"},
	{pkgLogger + "logger.go", "Logger factory"},
	{pkgTelemetry + "telemetry.go", "Telemetry helpers"},

	// ── api specs ────────────────────────────────────────────────────
	{apiOpenapi + "openapi.yaml", "OpenAPI 3.0 spec (source of truth)"},
	{apiOpenapi + "openapi.json", "Generated JSON"},
	{apiProto + "user.proto", "Public gRPC user service definition"},

	// ── test ─────────────────────────────────────────────────────────
	{testIntegration + "user_repository_test.go", "Integration tests for UserRepository (Testcontainers)"},
	{testIntegration + "setup_test.go", "Testcontainers setup"},
	{testIntegration + "testdata/.gitkeep", ""},
	{testFunctional + "user_api_test.go", "HTTP E2E tests for user API"},
	{testFunctional + "auth_api_test.go", "HTTP E2E tests for auth API"},
	{testFixture + "user_fixture.go", "Test fixtures for users"},
	{testMock + "user_repository_mock.go", "Generated mock (mockery)"},
	{testMock + "event_publisher_mock.go", "Generated mock (mockery)"},

	// ── scripts ──────────────────────────────────────────────────────
	{scriptsDir + "migrate.sh", "Run migrations"},
	{scriptsDir + "seed.sh", "Seed test data"},
	{scriptsDir + "generate-mocks.sh", "Generate mocks (mockery)"},
	{scriptsDir + "generate-swagger.sh", "Generate OpenAPI docs (swaggo)"},
	{scriptsDir + "health-check.sh", "CLI health probe"},

	// ── deployments ──────────────────────────────────────────────────
	{dockerDir + "Dockerfile", "Multi-stage build"},
	{dockerDir + "Dockerfile.dev", "Development Dockerfile"},
	{dockerDir + "docker-compose.yml", "Local stack (api + postgres + redis)"},
	{k8sDir + "deployment.yaml", "K8s Deployment"},
	{k8sDir + "service.yaml", "K8s Service"},
	{k8sDir + "ingress.yaml", "K8s Ingress"},
	{k8sDir + "configmap.yaml", "K8s ConfigMap"},
	{k8sDir + "secret.yaml", "K8s Secret"},
	{k8sDir + "hpa.yaml", "K8s HorizontalPodAutoscaler"},
	{helmDir + "Chart.yaml", "Helm chart metadata"},
	{helmDir + "values.yaml", "Helm chart values"},
	{helmTemplate + ".gitkeep", ""},

	// ── docs ─────────────────────────────────────────────────────────
	{docsDir + "ARCHITECTURE.md", "High-level design"},
	{docsDir + "DEVELOPMENT.md", "Onboarding guide"},
	{docsDir + "API.md", "API usage docs"},
	{adrDir + "0001-hexagonal-go.md", "ADR: hexagonal architecture in Go"},
	{adrDir + "0002-structured-logging.md", "ADR: structured logging"},

	// ── runtime & secrets ────────────────────────────────────────────
	{"logs/.gitkeep", ""},
	{"secrets/.gitignore", "Keep the folder but ignore its contents"},
	{"secrets/.env.example", "Template for local env vars"},

	// ── root files ───────────────────────────────────────────────────
	{".gitignore", "Root gitignore"},
	{".golangci.yml", "Linter config"},
	{".dockerignore", "Docker build context exclusions"},
	{"Makefile", "Common tasks (build, test, lint, migrate)"},
	{"go.mod", "Go module definition"},
	{"go.sum", "Dependency checksums"},
	{"config.yaml", "Default config"},
	{"config.dev.yaml", "Dev overrides"},
	{"config.prod.yaml", "Prod overrides"},
	{"README.md", "Project README"},
}

// ───────────────────────── full-file templates ────────────────────────

var makefileTmpl = strings.Join([]string{
	"# Makefile — common tasks (build, test, lint, migrate)",
	"",
	"BINDIR := bin",
	"",
	".PHONY: help build run test test-integration lint mock swagger migrate-up migrate-down docker-build tidy clean",
	"",
	"help: ## Show help",
	"\t@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = \":.*?## \"}; {printf \"  \\033[36m%-16s\\033[0m %s\\n\", $$1, $$2}'",
	"",
	"build: ## Build api + worker binaries",
	"\tgo build -o $(BINDIR)/api ./cmd/api",
	"\tgo build -o $(BINDIR)/worker ./cmd/worker",
	"",
	"run: ## Run the API server locally",
	"\tgo run ./cmd/api",
	"",
	"test: ## Run unit tests",
	"\tgo test ./...",
	"",
	"test-integration: ## Run integration tests (needs Docker)",
	"\tgo test ./test/integration/...",
	"",
	"lint: ## Run golangci-lint",
	"\tgolangci-lint run",
	"",
	"mock: ## Regenerate mocks",
	"\t./scripts/generate-mocks.sh",
	"",
	"swagger: ## Regenerate swagger docs",
	"\t./scripts/generate-swagger.sh",
	"",
	"migrate-up: ## Apply database migrations",
	"\t./scripts/migrate.sh up",
	"",
	"migrate-down: ## Rollback last migration",
	"\t./scripts/migrate.sh down 1",
	"",
	"docker-build: ## Build the production Docker image",
	"\tdocker build -f deployments/docker/Dockerfile -t hexagonal-api:dev .",
	"",
	"tidy: ## Tidy go modules",
	"\tgo mod tidy",
	"",
	"clean: ## Remove build artifacts",
	"\trm -rf $(BINDIR)",
	"",
}, "\n") + "\n"

var templates = map[string]string{
	// root
	"go.mod":           "module __MODULE__\n\ngo 1.22\n",
	"go.sum":           "",
	"Makefile":         makefileTmpl,
	"README.md":        readmeTmpl,
	".gitignore":       gitignoreTmpl,
	".dockerignore":    dockerignoreTmpl,
	".golangci.yml":    golangciTmpl,
	"config.yaml":      configTmpl,
	"config.dev.yaml":  configDevTmpl,
	"config.prod.yaml": configProdTmpl,

	// scripts
	"scripts/migrate.sh":          migrateSh,
	"scripts/seed.sh":             seedSh,
	"scripts/generate-mocks.sh":   genMocksSh,
	"scripts/generate-swagger.sh": genSwaggerSh,
	"scripts/health-check.sh":     healthCheckSh,

	// docker
	"deployments/docker/Dockerfile":         dockerfileTmpl,
	"deployments/docker/Dockerfile.dev":     dockerfileDevTmpl,
	"deployments/docker/docker-compose.yml": composeTmpl,

	// kubernetes
	"deployments/kubernetes/deployment.yaml": k8sDeploymentTmpl,
	"deployments/kubernetes/service.yaml":    k8sServiceTmpl,
	"deployments/kubernetes/ingress.yaml":    k8sIngressTmpl,
	"deployments/kubernetes/configmap.yaml":  k8sConfigMapTmpl,
	"deployments/kubernetes/secret.yaml":     k8sSecretTmpl,
	"deployments/kubernetes/hpa.yaml":        k8sHPATmpl,

	// helm
	helmDir + "Chart.yaml":  helmChartTmpl,
	helmDir + "values.yaml": helmValuesTmpl,

	// docs
	"docs/ARCHITECTURE.md":                archDocTmpl,
	"docs/DEVELOPMENT.md":                 devDocTmpl,
	"docs/API.md":                         apiDocTmpl,
	"docs/adr/0001-hexagonal-go.md":       adr1Tmpl,
	"docs/adr/0002-structured-logging.md": adr2Tmpl,

	// migrations (postgres)
	pgMigration + "000001_init.up.sql":   pgInitUp,
	pgMigration + "000001_init.down.sql": pgInitDown,
	pgMigration + "000002_seed.up.sql":   pgSeedUp,

	// migrations (mysql)
	myMigration + "000001_init.up.sql":   myInitUp,
	myMigration + "000001_init.down.sql": myInitDown,

	// api specs
	apiProto + "user.proto":     apiProtoTmpl,
	grpcProto + "user.proto":    grpcProtoTmpl,
	apiOpenapi + "openapi.yaml": openapiYamlTmpl,
	apiOpenapi + "openapi.json": openapiJSONTmpl,
	httpSpec + "swagger.json":   swaggerJSONTmpl,
	httpSpec + "swagger.yaml":   swaggerYAMLTmpl,

	// runtime & secrets
	"logs/.gitkeep":                       "",
	"secrets/.gitignore":                  secretsGitignore,
	"secrets/.env.example":                envExampleTmpl,
	testIntegration + "testdata/.gitkeep": "",
	helmTemplate + ".gitkeep":             "",
}

const gitignoreTmpl = `# Binaries
bin/
dist/
*.exe

# Test / coverage
coverage.out
coverage.html
*.test

# Runtime logs
logs/*
!logs/.gitkeep
*.log

# Local secrets — never commit
.env
secrets/*
!secrets/.gitignore
!secrets/.env.example

# IDE / editor
.idea/
.vscode/
*.swp

# OS
.DS_Store
Thumbs.db
`

const dockerignoreTmpl = `.git
.gitignore
bin
logs
secrets
test
docs
scripts
deployments
*.md
.dockerignore
`

const golangciTmpl = `run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - unused
    - gosimple
    - revive
    - gofmt
    - goimports
    - misspell

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
`

const configTmpl = `# Default configuration (overridden by config.<env>.yaml)

server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 10s
  write_timeout: 10s
  shutdown_timeout: 30s

database:
  driver: "postgres"           # postgres | mysql
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"         # override via env in production!
  name: "hexagonal"
  ssl_mode: "disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

cache:
  driver: "redis"              # redis | inmemory
  addr: "localhost:6379"
  default_ttl: 15m

messaging:
  driver: "inmemory"           # kafka | rabbitmq | redis_streams | inmemory
  kafka:
    brokers: ["localhost:9092"]
    topic: "domain-events"

security:
  jwt:
    secret: "change-me"        # override via env in production!
    access_ttl: 15m
    refresh_ttl: 720h
  bcrypt:
    cost: 12

observability:
  metrics_enabled: true
  tracing_enabled: true
  log_level: "info"
  log_format: "json"

rate_limit:
  requests_per_second: 100
  burst: 200
`

const configDevTmpl = `# Development overrides

server:
  port: 8081

database:
  ssl_mode: "disable"

observability:
  log_level: "debug"
  log_format: "console"
`

const configProdTmpl = `# Production overrides — secrets must come from env / secret manager

server:
  port: 8080

database:
  ssl_mode: "require"

observability:
  log_level: "info"
  log_format: "json"
`

const readmeTmpl = `# hexagonal_architecture_golang

Go service built with Hexagonal Architecture (Ports & Adapters).

## Layout

| Path | Purpose |
| --- | --- |
| cmd/ | Entrypoints: api, worker, migrate, healthcheck |
| internal/domain/ | Business core: entities, value objects, events, ports |
| internal/application/ | Use cases: commands/queries + DTOs |
| internal/adapter/inbound/ | Driving adapters: HTTP, gRPC, CLI |
| internal/adapter/outbound/ | Driven adapters: postgres, redis, kafka, jwt, bcrypt, ... |
| internal/container/ | Dependency injection (composition root) |
| pkg/ | Reusable public helpers |
| test/ | Integration, functional and fixture tests |

## Quick start

    go mod tidy
    go run ./cmd/api

    # or the full local stack (api + postgres + redis):
    docker compose -f deployments/docker/docker-compose.yml up --build

## Docs

- docs/ARCHITECTURE.md
- docs/DEVELOPMENT.md
- docs/API.md
- docs/adr/ — Architecture Decision Records

## Make targets

    make help
`

const migrateSh = `#!/usr/bin/env bash
# Run golang-migrate migrations.
# Usage: ./scripts/migrate.sh [up|down [N]|version|force VERSION]
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIGRATIONS_DIR="${ROOT_DIR}/internal/adapter/outbound/persistence/postgres/migration"
DB_URL="${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/hexagonal?sslmode=disable}"

ACTION="${1:-up}"
shift || true

exec migrate -path "${MIGRATIONS_DIR}" -database "${DB_URL}" "${ACTION}" "$@"
`

const seedSh = `#!/usr/bin/env bash
# Seed the database with test data (dev / test only).
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SEED_FILE="${ROOT_DIR}/internal/adapter/outbound/persistence/postgres/migration/000002_seed.up.sql"
DB_URL="${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/hexagonal?sslmode=disable}"

psql "${DB_URL}" -f "${SEED_FILE}"
echo "seed data applied"
`

const genMocksSh = `#!/usr/bin/env bash
# Generate mocks with mockery.
set -euo pipefail

mockery --recursive --inpackage --dir internal
echo "mocks generated (configure .mockery.yaml as needed)"
`

const genSwaggerSh = `#!/usr/bin/env bash
# Generate OpenAPI docs with swaggo/swag.
set -euo pipefail

swag init -g cmd/api/main.go --output internal/adapter/inbound/http/docs/spec --parseDependency
echo "swagger docs generated"
`

const healthCheckSh = `#!/usr/bin/env bash
# CLI health probe (useful for K8s exec probes and local checks).
set -euo pipefail

URL="${HEALTHCHECK_URL:-http://localhost:8080/healthz}"

if curl -fsS --max-time 3 "${URL}" > /dev/null; then
  echo "healthy"
else
  echo "unhealthy" >&2
  exit 1
fi
`

const dockerfileTmpl = `# syntax=docker/dockerfile:1

# ── Build stage ─────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/api ./cmd/api

# ── Runtime stage ───────────────────────────────────────────────
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 appuser
USER appuser
COPY --from=builder /bin/api /bin/api
EXPOSE 8080
ENTRYPOINT ["/bin/api"]
`

const dockerfileDevTmpl = `# Development image — mount the source, run with hot reload tooling.
FROM golang:1.22-alpine
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8080
CMD ["go", "run", "./cmd/api"]
`

const composeTmpl = `services:
  api:
    build:
      context: ../..
      dockerfile: deployments/docker/Dockerfile
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://postgres:postgres@postgres:5432/hexagonal?sslmode=disable
      REDIS_ADDR: redis:6379
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: hexagonal
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 3s
      retries: 10

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  pgdata:
`

const k8sDeploymentTmpl = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: hexagonal-api
  labels:
    app: hexagonal-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: hexagonal-api
  template:
    metadata:
      labels:
        app: hexagonal-api
    spec:
      containers:
        - name: api
          image: hexagonal-api:latest
          ports:
            - containerPort: 8080
          envFrom:
            - configMapRef:
                name: hexagonal-api-config
            - secretRef:
                name: hexagonal-api-secret
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8080
            initialDelaySeconds: 5
          readinessProbe:
            httpGet:
              path: /readyz
              port: 8080
            initialDelaySeconds: 3
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 256Mi
`

const k8sServiceTmpl = `apiVersion: v1
kind: Service
metadata:
  name: hexagonal-api
spec:
  type: ClusterIP
  selector:
    app: hexagonal-api
  ports:
    - port: 80
      targetPort: 8080
`

const k8sIngressTmpl = `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: hexagonal-api
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: hexagonal-api
                port:
                  number: 80
`

const k8sConfigMapTmpl = `apiVersion: v1
kind: ConfigMap
metadata:
  name: hexagonal-api-config
data:
  LOG_LEVEL: "info"
  SERVER_PORT: "8080"
`

const k8sSecretTmpl = `apiVersion: v1
kind: Secret
metadata:
  name: hexagonal-api-secret
type: Opaque
stringData:
  DATABASE_URL: "postgres://user:pass@host:5432/hexagonal?sslmode=require"
  JWT_SECRET: "change-me"
`

const k8sHPATmpl = `apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: hexagonal-api
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: hexagonal-api
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
`

const helmChartTmpl = `apiVersion: v2
name: hexagonal-api
description: Hexagonal architecture Go API service
type: application
version: 0.1.0
appVersion: "1.0.0"
`

const helmValuesTmpl = `replicaCount: 2

image:
  repository: hexagonal-api
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false
  host: api.example.com

resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 256Mi
`

const archDocTmpl = `# Architecture

This service follows Hexagonal Architecture (a.k.a. Ports & Adapters) with a
strict inward dependency rule.

## Layer overview

    cmd/                              entrypoints (thin mains)
    internal/container/               composition root (DI wiring)
    internal/application/             use cases (commands / queries / DTOs)
    internal/domain/                  business core (entities, VOs, events, ports)
    internal/adapter/inbound/         DRIVING adapters: http, grpc, cli
    internal/adapter/outbound/        DRIVEN adapters: postgres, redis, kafka, jwt, bcrypt
    pkg/                              reusable public helpers

## The dependency rule

- internal/domain depends on nothing but the standard library.
- internal/application depends only on internal/domain.
- Adapters depend on application/domain and implement their ports.
- Nothing inside domain or application ever imports an adapter.

## Patterns in use

- CQRS-lite: command/ vs query/ packages per module.
- Domain events published through the EventPublisher port.
- Unit of Work port for transactional use cases.
- Decorators (cache, audit, validation) around use cases.

See docs/adr/ for the decision records.
`

const devDocTmpl = `# Development Guide

## Prerequisites

- Go 1.22+
- Docker + docker compose (for Postgres / Redis)
- make

Optional: golang-migrate, mockery, swag, golangci-lint.

## Getting started

    make tidy          # download deps
    make run           # start the API locally

    # or run the whole local stack (api + postgres + redis):
    docker compose -f deployments/docker/docker-compose.yml up --build

## Database

    make migrate-up    # apply migrations
    ./scripts/seed.sh  # seed dev data

## Testing

    make test              # unit tests
    make test-integration  # integration tests (Testcontainers)

## Quality

    make lint

## Layout

See docs/ARCHITECTURE.md for the big picture.
`

const apiDocTmpl = `# API Usage

Base URL (local): http://localhost:8080

## Health

    GET /healthz   liveness
    GET /readyz    readiness

## Auth

    POST /v1/auth/register
    POST /v1/auth/login
    POST /v1/auth/logout

## Users

    GET    /v1/users?page=1&page_size=20
    POST   /v1/users
    GET    /v1/users/{id}
    PUT    /v1/users/{id}
    DELETE /v1/users/{id}

## Specs

- OpenAPI 3.0 (source of truth): api/openapi/openapi.yaml
- Swagger UI: /swagger/index.html
- ReDoc: /redoc
`

const adr1Tmpl = `# ADR-0001: Hexagonal architecture in Go

- Status: Accepted
- Date: 2025-01-01

## Context

The business core must stay independent of frameworks, databases and
delivery mechanisms so it remains testable and swappable.

## Decision

Adopt Hexagonal Architecture (Ports & Adapters):

- internal/domain          pure business core + ports (interfaces)
- internal/application     use cases orchestrating the domain
- internal/adapter/inbound     driving adapters (HTTP, gRPC, CLI)
- internal/adapter/outbound    driven adapters (persistence, cache, messaging)

## Consequences

- The domain has no framework imports; unit tests are fast and cheap.
- Adapters can be swapped (postgres -> mysql, kafka -> rabbitmq) without
  touching the core.
- More files and mapping boilerplate (DTO <-> domain <-> persistence model).
`

const adr2Tmpl = `# ADR-0002: Structured logging

- Status: Accepted
- Date: 2025-01-01

## Context

Logs must be machine-parsable, correlated per request (trace_id) and free of
PII, across all environments.

## Decision

Use structured JSON logging (zerolog or zap) behind a thin logging adapter in
internal/adapter/outbound/logging, with:

- request-scoped logger carrying trace_id via context
- PII redaction hook
- env-based level/format (dev: console/debug, prod: json/info)

## Consequences

- One logging API everywhere; the library can be swapped behind the adapter.
- Slight overhead for redaction, acceptable for compliance.
`

const pgInitUp = `-- 000001_init.up.sql — initial schema

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id            UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(50)  NOT NULL DEFAULT 'user',
    status        VARCHAR(50)  NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ  NULL
);

CREATE INDEX idx_users_deleted_at ON users (deleted_at);
`

const pgInitDown = `-- 000001_init.down.sql — rollback initial schema

DROP INDEX IF EXISTS idx_users_deleted_at;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";
`

const pgSeedUp = `-- 000002_seed.up.sql — seed data (dev / test only)
-- Both passwords are "password" (bcrypt, cost 10).

INSERT INTO users (email, password_hash, role, status)
VALUES
    ('admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', 'active'),
    ('user@example.com',  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user',  'active')
ON CONFLICT (email) DO NOTHING;
`

const myInitUp = `-- 000001_init.up.sql — initial schema (MySQL)

CREATE TABLE users (
    id            CHAR(36)     NOT NULL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(50)  NOT NULL DEFAULT 'user',
    status        VARCHAR(50)  NOT NULL DEFAULT 'active',
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    DATETIME     NULL,
    KEY idx_users_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
`

const myInitDown = `-- 000001_init.down.sql — rollback initial schema (MySQL)

DROP TABLE IF EXISTS users;
`

const grpcProtoTmpl = `syntax = "proto3";

package user.v1;

option go_package = "__MODULE__/internal/adapter/inbound/grpc/proto;userpb";

// UserService exposes user management over gRPC.
service UserService {
  rpc GetUser(GetUserRequest) returns (User);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}

message GetUserRequest {
  string id = 1;
}

message ListUsersRequest {
  int32 page = 1;
  int32 page_size = 2;
}

message ListUsersResponse {
  repeated User users = 1;
  int64 total = 2;
}

message User {
  string id = 1;
  string email = 2;
  string role = 3;
  string status = 4;
  string created_at = 5;
  string updated_at = 6;
}
`

const apiProtoTmpl = `syntax = "proto3";

package user.v1;

option go_package = "__MODULE__/api/proto;userpb";

// Public definition of the user service (mirrors the gRPC adapter contract).
service UserService {
  rpc GetUser(GetUserRequest) returns (User);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}

message GetUserRequest {
  string id = 1;
}

message ListUsersRequest {
  int32 page = 1;
  int32 page_size = 2;
}

message ListUsersResponse {
  repeated User users = 1;
  int64 total = 2;
}

message User {
  string id = 1;
  string email = 2;
  string role = 3;
  string status = 4;
  string created_at = 5;
  string updated_at = 6;
}
`

const openapiYamlTmpl = `openapi: 3.0.3
info:
  title: Hexagonal Architecture API
  description: >
    OpenAPI 3.0 specification — source of truth for the HTTP API.
    Derived copies live in internal/adapter/inbound/http/docs/spec.
  version: 1.0.0
servers:
  - url: http://localhost:8080
tags:
  - name: health
  - name: auth
  - name: users
paths:
  /healthz:
    get:
      tags: [health]
      summary: Liveness probe
      responses:
        "200":
          description: Service is alive
  /readyz:
    get:
      tags: [health]
      summary: Readiness probe
      responses:
        "200":
          description: Service is ready
  /v1/auth/register:
    post:
      tags: [auth]
      summary: Register a new user
      responses:
        "201":
          description: Registered
  /v1/auth/login:
    post:
      tags: [auth]
      summary: Login and receive tokens
      responses:
        "200":
          description: Token pair
  /v1/auth/logout:
    post:
      tags: [auth]
      summary: Invalidate the current session
      responses:
        "204":
          description: Logged out
  /v1/users:
    get:
      tags: [users]
      summary: List users (paginated)
      parameters:
        - in: query
          name: page
          schema: {type: integer, default: 1}
        - in: query
          name: page_size
          schema: {type: integer, default: 20}
      responses:
        "200":
          description: A page of users
    post:
      tags: [users]
      summary: Create a user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/CreateUserRequest"
      responses:
        "201":
          description: User created
  /v1/users/{id}:
    get:
      tags: [users]
      summary: Get a user by id
      parameters:
        - in: path
          name: id
          required: true
          schema: {type: string, format: uuid}
      responses:
        "200":
          description: The user
        "404":
          description: Not found
components:
  schemas:
    CreateUserRequest:
      type: object
      required: [email, password]
      properties:
        email:
          type: string
          format: email
        password:
          type: string
          minLength: 8
    User:
      type: object
      properties:
        id:
          type: string
          format: uuid
        email:
          type: string
        role:
          type: string
        status:
          type: string
        created_at:
          type: string
          format: date-time
`

const openapiJSONTmpl = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Hexagonal Architecture API",
    "version": "1.0.0"
  },
  "paths": {
    "/healthz": {
      "get": {
        "summary": "Liveness probe",
        "responses": {
          "200": {
            "description": "Service is alive"
          }
        }
      }
    }
  }
}
`

const swaggerJSONTmpl = `{
  "swagger": "2.0",
  "info": {
    "title": "Hexagonal Architecture API",
    "version": "1.0.0"
  },
  "basePath": "/",
  "paths": {}
}
`

const swaggerYAMLTmpl = `swagger: "2.0"
info:
  title: Hexagonal Architecture API
  version: "1.0.0"
basePath: /
paths: {}
`

const secretsGitignore = `*
!.gitignore
!.env.example
`

const envExampleTmpl = `# Copy to secrets/.env (gitignored) and fill in real values.
# All values can also be provided as real environment variables.

APP_ENV=dev

# Postgres
DATABASE_URL=postgres://postgres:postgres@localhost:5432/hexagonal?sslmode=disable

# Redis
REDIS_ADDR=localhost:6379

# Auth
JWT_SECRET=change-me-in-production
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h

# Messaging (optional)
KAFKA_BROKERS=localhost:9092

# Observability
LOG_LEVEL=debug
LOG_FORMAT=console
`

// ───────────────────────── content generation ─────────────────────────

// render decides the content and file mode for one spec.
func render(s spec) (string, os.FileMode) {
	base := filepath.Base(s.path)

	if tmpl, ok := templates[s.path]; ok {
		mode := os.FileMode(0o644)
		if strings.HasSuffix(base, ".sh") {
			mode = 0o755 // make scripts executable
		}
		return tmpl, mode
	}

	switch {
	case base == "main.go":
		return mainStub(s), 0o644
	case strings.HasSuffix(base, "_test.go"):
		return testStub(s), 0o644
	case strings.HasSuffix(base, ".go"):
		return goStub(s), 0o644
	default:
		// Fallback stub for anything not covered by an explicit template.
		return fmt.Sprintf("# %s - %s\n\nTODO: implement.\n", base, desc(s.desc)), 0o644
	}
}

// goStub generates a compilable package stub (package name = folder name).
func goStub(s spec) string {
	return fmt.Sprintf("// %s - %s\n// TODO: implement.\n\npackage %s\n",
		filepath.Base(s.path), desc(s.desc), pkgName(s.path))
}

// testStub generates a skipping placeholder test.
func testStub(s spec) string {
	name := testName(s.path)
	return fmt.Sprintf(`package %s

import "testing"

// %s - %s
// TODO: implement.
func %s(t *testing.T) {
    t.Skip("TODO: implement")
}
`, pkgName(s.path), name, desc(s.desc), name)
}

// mainStub generates a runnable cmd entrypoint stub.
func mainStub(s spec) string {
	name := filepath.Base(filepath.Dir(s.path)) // api | worker | migrate | healthcheck
	return fmt.Sprintf(`package main

import "fmt"

// main - %s
// TODO: load config, wire dependencies, start.
func main() {
    fmt.Println("%s: starting... (TODO: implement)")
}
`, desc(s.desc), name)
}

func pkgName(path string) string { return filepath.Base(filepath.Dir(path)) }

// testName derives a unique test name from the file name,
// e.g. get_user_by_id_test.go -> TestGetUserById.
func testName(path string) string {
	base := strings.TrimSuffix(filepath.Base(path), "_test.go")
	var b strings.Builder
	b.WriteString("Test")
	for _, part := range strings.Split(base, "_") {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		b.WriteString(part[1:])
	}
	return b.String()
}

func desc(d string) string {
	if strings.TrimSpace(d) == "" {
		return "TODO: describe the purpose of this file."
	}
	return d
}

// ─────────────────────────────── main ─────────────────────────────────

func main() {
	flag.Parse()

	var created, skipped, dirs int
	seenDirs := map[string]bool{}

	for _, s := range files {
		dst := filepath.Join(*rootDir, filepath.FromSlash(s.path))

		// Never clobber existing files unless -force was given.
		if !*force {
			if _, err := os.Stat(dst); err == nil {
				skipped++
				continue
			}
		}

		content, mode := render(s)
		content = strings.ReplaceAll(content, "__MODULE__", *moduleName)

		if d := filepath.Dir(s.path); !seenDirs[d] {
			seenDirs[d] = true
			dirs++
		}

		if !*dryRun {
			if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
				fatal("create dir for %s: %v", s.path, err)
			}
			if err := os.WriteFile(dst, []byte(content), mode); err != nil {
				fatal("write %s: %v", s.path, err)
			}
		}
		created++
	}

	verb := "created"
	if *dryRun {
		verb = "would create"
	}
	fmt.Printf("\n[ok] %s %d files across %d directories (%d skipped: already exist)\n",
		verb, created, dirs, skipped)

	if !*dryRun && created > 0 {
		fmt.Printf("\nNext steps:\n  cd %s\n  go mod tidy\n  go build ./...\n  go test ./...\n  go run ./cmd/api\n", *rootDir)
	}
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
