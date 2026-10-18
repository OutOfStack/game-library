# Repository Guidelines

## Project Structure & Module Organization

Read `README.md` for service context. This repository contains the Go backend; UI, authentication, and infrastructure live elsewhere.

- `cmd/api/`, `cmd/worker/`, and `cmd/manage/`: API service, background-task worker, and database-management entry points.
- `internal/api/`: HTTP/gRPC handlers; `internal/facade/`: business logic; `internal/repo/`: PostgreSQL persistence; `internal/model/`: domain models.
- `internal/client/`, `internal/taskprocessor/`, and `internal/appconf/`: external integrations, background jobs, and configuration.
- `pkg/`: shared utilities and generated protobuf code; `api/proto/`: protobuf sources.
- `scripts/migrations/`: numbered SQL migration pairs; `test-data/`: manual test fixtures; `api/http/openapi.yaml`: HTTP API spec (source of truth).
- `.github/workflows/` and `.k8s/`: CI and deployment configuration.

## Build, Test, and Development Commands

Use Go 1.27, Docker, and Make. Copy `.example.env` to `.env` and configure local connections before starting.

- `make drunpg` / `make drunredis`: start local dependencies.
- `make migrate`: apply database migrations; `make seed`: load optional sample data.
- `make build` / `make build-mng`: compile binaries into `bin/`.
- `make run`: start the API from source.
- `make test`: run all default tests with the race detector.
- `make cover`: generate `coverage.out` and print coverage.
- `make lint`: run golangci-lint using `.golangci.yaml`.
- `make generate`: regenerate protobufs, OpenAPI server code (`make generate-openapi`), and mocks.

## Coding Style & Naming Conventions

Use standard Go naming, tab indentation, and configured `gofmt`, `goimports`, and `gci` formatters. Document exported functions and structs. Start comments inside functions with lowercase letters; omit trailing periods. Change the HTTP API spec-first: edit `api/http/openapi.yaml`, run `make generate-openapi`, then implement the generated `ServerInterface` methods (path and query params arrive parsed; decode bodies with `decodeRequest`, respond with `web.Respond` or `respondError`). Request/response models are generated from the spec; add behavior such as validation as methods in hand-written files in `internal/api` (see `validate.go`). Use lowercase filenames such as `get_games.go` and paired `000025_description.up.sql` / `.down.sql` migrations.

Never hand-edit generated files. After protobuf changes, run `make generate-proto` using `buf.yaml` and `buf.gen.yaml`. Extend Makefile generation targets for new mocks. Run `go mod tidy` when dependencies change.

## Testing Guidelines

Use Go's `testing`, Testify, and GoMock. Place `*_test.go` files beside source, use external packages such as `api_test`, and descriptive names such as `Test_GetGames_Success`. Test all exported functions and error paths. Use `t.Context()` and `internal/pkg/td/random.go` for random fixtures. Comment only non-obvious test logic.

Repository tests require Docker and start PostgreSQL on port 5439. `make test-integration` enables live API tests requiring `OPENAI_API_KEY`. Codecov receives CI coverage; no numeric minimum is configured.

For manual API checks, consult `api/http/openapi.yaml` and the [auth API specification](https://github.com/OutOfStack/game-library-auth/blob/main/docs/swagger.json). Existing local test accounts are `aiuser:aiuser__` (user) and `aipublisher:aipublisher` (publisher).

## Commit & Pull Request Guidelines

Follow `CONTRIBUTING.md` and recent Conventional Commits: `fix:`, `feat:`, `docs:`, or `ci:`. Breaking changes use `!` or `BREAKING CHANGE:` and affect release versions.

PRs should describe changes, link issues, report validation, and identify migration/configuration changes. Run `make build`, `make test`, and `make lint` for code changes; fix failures. Regenerate changed definitions and update affected README guidance.

## Configuration

Follow the documented configuration mechanism in `internal/appconf/`; introduce environment variables only when specified. Keep credentials out of commits and document settings in `.example.env`. gRPC defaults to port 9000 (`APP_GRPC_ADDRESS`).

## Workflow Restrictions

- Do not stage changes or create commits.
- Do not delete files; notify the user when files become redundant.

## File Review Boundaries

- Do not review or analyze generated documentation (`docs/**`), mocks (`**/mocks/**`, `**/*_mock.go`), or vendored dependencies (`vendor/**`).
- Do not review or analyze generated code matching `**/*.gen.go`, `**/*.pb.go`, or `**/*_grpc.pb.go`.
- Do not read sensitive files: `*.pem`, `*.key`, or `.env`.

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:
- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
