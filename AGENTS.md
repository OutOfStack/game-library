# Repository Guidelines

## Project Structure & Module Organization

Read `README.md` for service context. This repository contains the Go backend; UI, authentication, and infrastructure live elsewhere.

- `cmd/game-library-api/` and `cmd/game-library-manage/`: service and database-management entry points.
- `internal/api/`: HTTP/gRPC handlers; `internal/facade/`: business logic; `internal/repo/`: PostgreSQL persistence; `internal/model/`: domain models.
- `internal/client/`, `internal/taskprocessor/`, and `internal/appconf/`: external integrations, background jobs, and configuration.
- `pkg/`: shared utilities and generated protobuf code; `api/proto/`: protobuf sources.
- `scripts/migrations/`: numbered SQL migration pairs; `test-data/`: manual test fixtures; `docs/`: generated Swagger assets.
- `.github/workflows/` and `.k8s/`: CI and deployment configuration.

## Build, Test, and Development Commands

Use Go 1.27, Docker, and Make. Copy `app.example.env` to `app.env` and configure local connections before starting.

- `make drunpg` / `make drunredis`: start local dependencies.
- `make migrate`: apply database migrations; `make seed`: load optional sample data.
- `make build` / `make build-mng`: compile binaries into `bin/`.
- `make run`: start the API from source.
- `make test`: run all default tests with the race detector.
- `make cover`: generate `coverage.out` and print coverage.
- `make lint`: run golangci-lint using `.golangci.yaml`.
- `make generate`: regenerate protobufs, Swagger documentation, and mocks.

## Coding Style & Naming Conventions

Use standard Go naming, tab indentation, and configured `gofmt`, `goimports`, and `gci` formatters. Document exported functions and structs. Start comments inside functions with lowercase letters; omit trailing periods. Keep handler Swagger annotations current. Use lowercase filenames such as `get_games.go` and paired `000025_description.up.sql` / `.down.sql` migrations.

Never hand-edit generated files. After protobuf changes, run `make generate-proto` using `buf.yaml` and `buf.gen.yaml`. Extend Makefile generation targets for new mocks. Run `go mod tidy` when dependencies change.

## Testing Guidelines

Use Go's `testing`, Testify, and GoMock. Place `*_test.go` files beside source, use external packages such as `api_test`, and descriptive names such as `Test_GetGames_Success`. Test all exported functions and error paths. Use `t.Context()` and `internal/pkg/td/random.go` for random fixtures. Comment only non-obvious test logic.

Repository tests require Docker and start PostgreSQL on port 5439. `make test-integration` enables live API tests requiring `OPENAI_API_KEY`. Codecov receives CI coverage; no numeric minimum is configured.

For manual API checks, consult `docs/swagger.json` and the [auth API specification](https://github.com/OutOfStack/game-library-auth/blob/main/docs/swagger.json). Existing local test accounts are `aiuser:aiuser__` (user) and `aipublisher:aipublisher` (publisher).

## Commit & Pull Request Guidelines

Follow `CONTRIBUTING.md` and recent Conventional Commits: `fix:`, `feat:`, `docs:`, or `ci:`. Breaking changes use `!` or `BREAKING CHANGE:` and affect release versions.

PRs should describe changes, link issues, report validation, and identify migration/configuration changes. Run `make build`, `make test`, and `make lint` for code changes; fix failures. Regenerate changed definitions and update affected README guidance.

## Configuration

Follow the documented configuration mechanism in `internal/appconf/`; introduce environment variables only when specified. Keep credentials out of commits and document settings in `app.example.env`. gRPC defaults to port 9000 (`APP_GRPC_ADDRESS`).

## Workflow Restrictions

- Do not stage changes or create commits.
- Do not delete files; notify the user when files become redundant.

## File Review Boundaries

- Do not review or analyze generated documentation (`docs/**`), mocks (`**/mocks/**`, `**/*_mock.go`), or vendored dependencies (`vendor/**`).
- Do not review or analyze generated code matching `**/*.gen.go`, `**/*.pb.go`, or `**/*_grpc.pb.go`.
- Do not read sensitive files: `*.pem`, `*.key`, or `app.env`.
