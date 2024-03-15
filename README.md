# game-library

## Introduction

game-library is a web application for exploring and rating games written in Go and TypeScript. It consists of three services: 
- current service is responsible for fetching, storing games data and providing it to UI, 
- [auth service](https://github.com/OutOfStack/game-library-auth) is responsible for user authentication and authorization,
- [ui service](https://github.com/OutOfStack/game-library-ui) is responsible for UI representation.

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Features](#features)
- [Configuration](#configuration)
- [Documentation](#documentation)
- [Examples](#examples)
- [List of Make commands](#list-of-make-commands)
- [Contributors](#contributors)
- [License](#license)

## Installation

Prerequisites: `go`, `Docker`, `Make`. To set up the service, follow these steps:

1. Clone the repository:
    ```bash
    git clone https://github.com/OutOfStack/game-library.git
    cd game-library
    ```

2. Set up the database:
    ```bash
    make dockerrunpg # runs postgres in docker container
    make createdb # creates db
    make migrate # applies migrations
    # optionally
    make seed # applies test data
    ```

3. Install and run dependencies:
    ```bash
    make dockerrunglog # [Optional] runs graylog in docker container
    make dockerrunredis # [Optional] runs redis in docker container
    make dockerrunzipkin # [Optional] runs zipkin
    ```

4. _[Optional]_ Set up fetching games data:
    - Get credentials from [IGDB API](https://api-docs.igdb.com/#account-creation) to run background task that fetches games
    - Get credentials from [Uploadcare API](https://uploadcare.com/api/) for uploading game images

5. _[Optional]_ Install [auth service](https://github.com/OutOfStack/game-library-auth) for using handlers that require authentication

6. Build and run the service:
    ```bash
    make build
    make run
    ```

_Optional_ steps are not required for minimal install but required for full functionality.

Refer to the [List of Make commands](#list-of-make-commands) for a complete list of commands.

## Usage

After installation, you can use the following Make commands to develop the service:

- `make test`: Runs tests.
- `make generate`: Generates documentation for Swagger UI.
- `make lint`: Runs golangci-lint for code analysis.

Refer to the [List of Make commands](#list-of-make-commands) for a complete list of commands.

## Features

- Data storage with PostgreSQL.
- Caching with Redis.
- Tracing with Zipkin.
- Log management with Graylog.
- Background fetching of game data from [IGDB](https://api-docs.igdb.com/).
- Reuploading game images to Uploadcare CDN.
- Code analysis with golangci-lint.
- CI/CD with GitHub Actions and deploy to Kubernetes (microk8s) cluster.

## Configuration

- The service can be configured using [app.env](./app.env) or environment variables, described in [settings.go](./internal/appconf/settings.go)
- CI/CD configs are in _./.github/workflows/_ 
- k8s deployment configs are in _./.k8s/_ 

## Documentation

API documentation is available at [Swagger UI](http://localhost:8000/swagger/index.html). 
For regenerating documentation after handlers change run `make generate`.

## Examples

Endpoint that returns 3 games ordered by release date:
```bash
curl -X GET "http://localhost:8000/api/games?pageSize=3&page=1&orderBy=releaseDate"
```

To see other examples of API endpoints, refer to the [documentation](#documentation).

## List of Make commands:
    build           builds app
    run             runs app
    test            runs tests for the whole project
    generate        generates documentation for swagger UI
    lint            runs golangci-lint

    dockerrunpg     runs postgres server in docker container
    createdb        creates database on postgres server started by 'dockerrunpg'
    dropdb          drops database on postgres server created by 'dockerrunpg'
    migrate         applies all migrations to database
    rollback        rollbacks last migration on database
    seed            seeds test data to database

    dockerbuildapi  builds app docker image
    dockerrunapi    runs app in docker container
    dockerrunzipkin runs zipkin in docker container
    dockerrunglog   runs graylog in docker container
    dockerrunredis  runs redis in docker container
    dockerbuildmng  builds manage app docker image
    dockerrunmng-m  applies migrations to database using docker manage image
    dockerrunmng-r  rollbacks one last migration using docker manage image
    dockerrunmng-s  seeds test data to database using docker manage image

## Contributors

- [OutOfStack](https://github.com/OutOfStack)
- For a complete list of contributors, refer to the GitHub repository.

## License

[MIT License](./LICENSE.md)
