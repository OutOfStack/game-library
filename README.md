# game-library

## Introduction

game-library is a web application for exploring and rating games. 
It consists of three services:
- current service is responsible for fetching, storing games data and providing it to UI,
- [auth service](https://github.com/OutOfStack/game-library-auth) is responsible for user authentication and authorization,
- [ui service](https://github.com/OutOfStack/game-library-ui) is responsible for UI representation.

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Tech Stack and Integrations](#tech-stack-and-integrations)
- [Configuration](#configuration)
- [Documentation](#documentation)
- [Examples](#examples)
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
    make drunredis # [Optional] runs redis in docker container
    make drunzipkin # [Optional] runs zipkin in docker container
    make drunglog # [Optional] runs graylog in docker container
    make drunprom # [Optional] runs prometheus in docker container
    ```

4. _[Optional]_ Set up fetching games data:
    - Get credentials from [IGDB API](https://api-docs.igdb.com/#account-creation) to run background tasks that fetch and update games
    - Get S3 compatible storage, for example [AWS S3](https://aws.amazon.com/s3/) or [Cloudflare R2](https://www.cloudflare.com/en-gb/developer-platform/products/r2/) for uploading game images
    - Get an API key from [OpenAI](https://platform.openai.com/api-keys) for automatic game moderation functionality

5. _[Optional]_ Launch [auth service](https://github.com/OutOfStack/game-library-auth) for using handlers that require authentication

6. Create the `app.env` file based on [`./app.example.env`](./app.example.env) and update it with your local configuration settings.

7.  Build and run the service:
    ```bash
    make build
    make run
    ```

_Optional_ steps are not required for minimal install but required for full functionality.

Refer to the [List of Make commands](#list-of-make-commands) for a complete list of commands.

## Usage

After installation, you can use the following Make commands to develop the service:

- `make test`: Runs tests.
- `make generate`: Generates proto files, documentation for Swagger UI and mocks for testing.
- `make lint`: Runs golangci-lint for code analysis.

Refer to the [List of Make commands](#list-of-make-commands) for a complete list of commands.

## Tech Stack and Integrations

- Data storage with PostgreSQL.
- Caching with Redis.
- Tracing with Zipkin.
- Log management with Graylog.
- Background task for fetching and updating games data using IGDB API.
- Game image upload and storage with S3-compatible services (Cloudflare R2).
- Automatic game moderation using OpenAI API.
- gRPC service for internal service-to-service communication.
- Code analysis with golangci-lint.
- CI/CD with GitHub Actions and deploy to Kubernetes (microk8s) cluster.

## Configuration

- The service can be configured using `app.env` or environment variables, described in [`settings.go`](./internal/appconf/settings.go)
- CI/CD configs are in [`./github/workflows/`](./.github/workflows/)
- k8s deployment configs are in [`./k8s`](./.k8s/)

## Documentation

API documentation is available via [Swagger UI](http://localhost:8000/swagger/index.html).
For regenerating documentation after swagger description change run:
```bash
make generate
```

### gRPC Service

The service exposes a gRPC endpoint for internal service-to-service communication.

**Endpoint:** `localhost:9000` (`APP_GRPC_ADDRESS` environment variable in [`app.example.env`](./app.example.env))

**Testing with grpcurl:**
```bash
# list services
grpcurl -plaintext localhost:9000 list

# describe service
grpcurl -plaintext localhost:9000 describe infoapi.InfoApiService

# call method
grpcurl -plaintext -d '{"company_name": "Nintendo"}' -emit-defaults localhost:9000 infoapi.InfoApiService/CompanyExists
```

**Protobuf Definition:**
The protobuf schema is located in [`api/proto/infoapi.proto`](./api/proto/infoapi.proto)

## Examples

Endpoint that returns 3 games ordered by release date:
```bash
curl -X GET "http://localhost:8000/api/games?pageSize=3&page=1&orderBy=releaseDate"
```

To see other examples of API endpoints, refer to the [documentation](#documentation).

## List of Make commands:

#### Main Commands
    build         builds app
    build-mng     build manage app
    run           runs app
    test          runs tests for the whole project
    generate      generates proto files, docs for swagger UI and mocks for testing
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
    drunzipkin    runs zipkin in docker container
    drunredis     runs redis in docker container
    drunglog      runs graylog in docker container
    drunprom      runs prometheus in docker container

## License

[MIT License](./LICENSE.md)
