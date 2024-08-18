build:
	mkdir -p bin
	go build -o bin/game-library-api cmd/game-library-api/main.go

run:
	go run ./cmd/game-library-api/.

test:
	go test -v -race ./...

dockerrunpg:
	docker compose up -d --no-recreate db

createdb:
	docker exec -it games_db createdb --username=postgres --owner=postgres games

dropdb:
	docker exec -it games_db dropdb --username=postgres games

migrate:
	go run ./cmd/game-library-manage/. migrate

rollback:
	go run ./cmd/game-library-manage/. rollback

seed:
	go run ./cmd/game-library-manage/. seed

SWAG_PKG := github.com/swaggo/swag/cmd/swag@v1.16.3
SWAG_BIN := $(shell go env GOPATH)/bin/swag
generate:
	@if \[ ! -f ${SWAG_BIN} \]; then \
		echo "Installing swag..."; \
    go install ${SWAG_PKG}; \
  fi
	@if \[ -f ${SWAG_BIN} \]; then \
  	echo "Found swag at '$(SWAG_BIN)', generating documentation..."; \
    ${SWAG_BIN} init --parseDependency -g cmd/game-library-api/main.go; \
	else \
    echo "swag not found or the file does not exist"; \
    exit 1; \
  fi

LINT_PKG := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.1
LINT_BIN := $(shell go env GOPATH)/bin/golangci-lint
lint:
	@if \[ ! -f ${LINT_BIN} \]; then \
		echo "Installing golangci-lint..."; \
    go install ${LINT_PKG}; \
  fi
	@if \[ -f ${LINT_BIN} \]; then \
  	echo "Found golangci-lint at '$(LINT_BIN)', running..."; \
    ${LINT_BIN} run; \
	else \
    echo "golangci-lint not found or the file does not exist"; \
    exit 1; \
  fi

dockerbuildapi:
	docker build -t game-library:latest .

dockerrunapi:
	docker compose up -d api

dockerrunzipkin:
	docker compose up -d zipkin

dockerrunredis:
	docker compose up -d redis

dockerrunglog:
	docker compose up -d graylog

dockerbuildmng:
	docker build -f Dockerfile.mng -t game-library-mng:latest .

dockerrunmng-m:
	docker compose run --rm mng migrate

dockerrunmng-r:
	docker compose run --rm mng rollback

dockerrunmng-s:
	docker compose run --rm mng seed
