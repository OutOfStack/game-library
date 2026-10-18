# game-library

[![Build](https://img.shields.io/github/actions/workflow/status/OutOfStack/game-library/main.yml?branch=main&label=build)](https://github.com/OutOfStack/game-library/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/OutOfStack/game-library/graph/badge.svg?token=NRIKDH6OFR)](https://codecov.io/gh/OutOfStack/game-library)

## Introduction

`game-library` is a web application for exploring and rating games.

`game-library` consists of:
- current service is responsible for fetching, storing games data and providing it to UI
- [game-library-auth](https://github.com/OutOfStack/game-library-auth) - authentication and authorization
- [game-library-ui](https://github.com/OutOfStack/game-library-ui) - UI service
- [game-library-infra](https://github.com/OutOfStack/game-library-infra) - infrastructure and deployment configurations

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Tech Stack and Integrations](#tech-stack-and-integrations)
- [Configuration](#configuration)
- [Documentation](#documentation)
- [Examples](#examples)
- [Contributing](#contributing)
- [List of Make commands](#list-of-make-commands)
- [License](#license)

## Installation

Prerequisites: [`go`](https://go.dev/doc/install), [`Docker`](https://docs.docker.com/desktop/), [`Make`](https://www.gnu.org/software/make/)

To set up the service, follow these steps:

1. Clone the repository:
    ```bash
    git clone https://github.com/OutOfStack/game-library.git
    cd game-library
    ```

2. Set up the database:
    ```bash
    make drunpg # runs postgres in docker container
    make migrate # applies migrations
    # optionally
    make seed # applies test data
    ```

3. Install and run dependencies:
    ```bash
    make drunredis # runs redis in docker container
    ```

4. _[Optional]_ Set up fetching games data:
    - Get credentials from [IGDB API](https://api-docs.igdb.com/#account-creation) to run background tasks that fetch and update games
    - Get S3 compatible storage, for example [AWS S3](https://aws.amazon.com/s3/) or [Cloudflare R2](https://www.cloudflare.com/en-gb/developer-platform/products/r2/) for uploading game images
    - Get an API key from [OpenAI](https://platform.openai.com/api-keys) for automatic game moderation functionality

5. _[Optional]_ Launch [auth service](https://github.com/OutOfStack/game-library-auth) for using handlers that require authentication

6. Create the `.env` file based on [`./.example.env`](./.example.env) and update it with your local configuration settings.

7.  Build and run the service:
    ```bash
    make build
    make run
    ```
    Background tasks (IGDB fetch, trending index, game info updates, moderation) run in a separate worker:
    ```bash
    make run-worker
    ```
    `make run-worker` overrides `APP_HTTP_ADDRESS` to `localhost:8001` for the worker's `/metrics` endpoint, so it can run alongside the API.

_Optional_ steps are not required for minimal install but required for full functionality.

Refer to the [List of Make commands](#list-of-make-commands) for a complete list of commands.

## Usage

After installation, you can use the following Make commands to develop the service:

- `make test`: Runs tests.
- `make generate`: Generates proto files, HTTP server code from the OpenAPI spec and mocks for testing.
- `make lint`: Runs golangci-lint for code analysis.

Refer to the [List of Make commands](#list-of-make-commands) for a complete list of commands.

## Tech Stack and Integrations

- Data storage with PostgreSQL.
- Caching with Redis.
- Background task for fetching and updating games data using IGDB API.
- Game image upload and storage with S3-compatible services (Cloudflare R2).
- Automatic game moderation using OpenAI API.
- gRPC for internal service-to-service communication.
- Tracing with OTLP exporter (Jaeger).
- Log management with Graylog.
- Metrics collection with Prometheus.
- Code analysis with golangci-lint.
- CI/CD with GitHub Actions and deploy to Kubernetes (microk8s) cluster.

## Configuration

- The service can be configured using `.env` or environment variables, described in [`settings.go`](./internal/appconf/settings.go)
- CI/CD configs are in [`./github/workflows/`](./.github/workflows/)
- k8s deployment configs are in [`./k8s`](./.k8s/)

The worker exports traces to `JAEGER_OTLP_ENDPOINT` under service name `game-library-worker`. Its `/metrics` endpoint 
exposes the worker's own Prometheus registry. The `game-library-infra` Prometheus configuration discovers opted-in pods 
in the `game-library` namespace using `prometheus.io/scrape`, `prometheus.io/port`, and `prometheus.io/path` annotations. 
The pod `app` label becomes the scrape job name (`game-library` or `game-library-worker`). New components only need these 
annotations and a matching declared TCP container port; no infra edit is needed. During migration, deploy the annotated 
API, auth, and worker manifests before switching Prometheus to generic pod discovery.

## Documentation

The HTTP API is defined spec-first in [`api/http/openapi.yaml`](api/http/openapi.yaml) (OpenAPI 3).
Server interfaces and routing helpers in `internal/api/openapi.gen.go` are generated from it with
[oapi-codegen](https://github.com/oapi-codegen/oapi-codegen); handlers implement the generated `ServerInterface` and respond with `web.Respond`/`respondError`.
After changing the spec, regenerate the code:

```bash
make generate-openapi
```

API documentation is available via [Swagger UI](http://localhost:8000/swagger/index.html), which serves the embedded spec
from `/swagger/openapi.yaml`.

### gRPC

The service exposes a gRPC endpoint for internal service-to-service communication.

**Endpoint:** `localhost:9000` (`APP_GRPC_ADDRESS` environment variable in [`.example.env`](./.example.env))

**Protobuf Definition:** [`api/proto/infoapi/v1/infoapi.proto`](api/proto/infoapi/v1/infoapi.proto)

**Testing with grpcurl:**
```bash
# list services
grpcurl -plaintext localhost:9000 list

# describe service
grpcurl -plaintext localhost:9000 describe infoapi.v1.InfoApiService

# call CompanyExists method
grpcurl -plaintext -d '{"company_name": "Nintendo"}' -emit-defaults localhost:9000 infoapi.v1.InfoApiService/CompanyExists
```

## Examples

Endpoint that returns 3 games ordered by release date:
```bash
curl -X GET "http://localhost:8000/api/games?pageSize=3&page=1&orderBy=releaseDate"
```

To see other examples of API endpoints, refer to the [documentation](#documentation).

## Contributing

Commit message and release rules are documented in [CONTRIBUTING.md](./CONTRIBUTING.md).

## List of Make commands:

#### Main Commands
    build         builds app
    build-mng     build manage app
    build-worker  builds background tasks worker
    run           runs app
    run-worker    runs background tasks worker
    test          runs tests for the whole project
    generate      generates proto files, OpenAPI server code and mocks for testing
    lint          runs golangci-lint
    cover         outputs tests coverage

#### Database Commands
    drunpg        runs postgres server in docker container
    migrate       applies all migrations to database (reads from config file)
    rollback      rollbacks last migration on database (reads from config file)
    seed          seeds test data to database (reads from config file)

#### Docker Commands
    dbuildapi     builds app docker image
    dbuildmng     builds manage app docker image
    drunapi       runs app in docker container
    drunredis     runs redis in docker container

## License

[MIT License](./LICENSE.md)
