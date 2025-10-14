build:
	mkdir -p bin
	go build -o bin/game-library-api cmd/game-library-api/main.go

build-mng:
	mkdir -p bin
	go build -o bin/game-library-manage cmd/game-library-manage/main.go

run:
	go run ./cmd/game-library-api/.

test:
	go test -v -race ./...

cover:
	go test -cover -coverpkg=./... -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

drunpg:
	docker compose up -d --no-recreate db

migrate:
	go run ./cmd/game-library-manage/. -from-file migrate

rollback:
	go run ./cmd/game-library-manage/. -from-file rollback

seed:
	go run ./cmd/game-library-manage/. -from-file seed

SWAG_PKG := github.com/swaggo/swag/cmd/swag@v1.16.4
SWAG_BIN := $(shell go env GOPATH)/bin/swag
generate-swag:
	@if \[ ! -f ${SWAG_BIN} \]; then \
		echo "Installing swag..."; \
    	go install ${SWAG_PKG}; \
  	fi
	@if \[ -f ${SWAG_BIN} \]; then \
  		echo "Found swag at '$(SWAG_BIN)', generating documentation..."; \
	else \
    	echo "swag not found or the file does not exist"; \
    	exit 1; \
  	fi
	${SWAG_BIN} init \
	-d cmd/game-library-api,internal/app/game-library-api/api,internal/app/game-library-api/api/model,internal/app/game-library-api/web

MOCKGEN_PKG := go.uber.org/mock/mockgen@v0.6
MOCKGEN_BIN := $(shell go env GOPATH)/bin/mockgen
generate-mocks:
	@if \[ ! -f ${MOCKGEN_BIN} \]; then \
		echo "Installing mockgen..."; \
		go install ${MOCKGEN_PKG}; \
	fi
	@if \[ -f ${MOCKGEN_BIN} \]; then \
		echo "Found mockgen at '$(MOCKGEN_BIN)', generating mocks..."; \
	else \
		echo "mockgen not found or the file does not exist"; \
		exit 1; \
  	fi
	${MOCKGEN_BIN} -source=internal/app/game-library-api/api/provider.go -destination=internal/app/game-library-api/api/mocks/provider.go -package=api_mock
	${MOCKGEN_BIN} -source=internal/pkg/cache/redis.go -destination=internal/pkg/cache/mocks/redis.go -package=cache_mock
	${MOCKGEN_BIN} -source=internal/app/game-library-api/facade/provider.go -destination=internal/app/game-library-api/facade/mocks/provider.go -package=facade_mock
	${MOCKGEN_BIN} -source=internal/auth/auth.go -destination=internal/auth/mocks/auth.go -package=auth_mock
	${MOCKGEN_BIN} -source=internal/middleware/auth.go -destination=internal/middleware/mocks/auth.go -package=middleware_mock
	${MOCKGEN_BIN} -source=internal/taskprocessor/task.go -destination=internal/taskprocessor/mocks/task.go -package=taskprocessor_mock
	${MOCKGEN_BIN} -destination=internal/app/game-library-api/repo/mocks/tx.go -package=repo_mock github.com/jackc/pgx/v5 Tx

generate: generate-swag generate-mocks

LINT_PKG := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5
LINT_BIN := $(shell command -v golangci-lint 2>/dev/null || echo $(shell go env GOPATH)/bin/golangci-lint)
lint:
	@if [ ! -f ${LINT_BIN} ]; then \
		echo "Installing golangci-lint..."; \
    	go install ${LINT_PKG}; \
		LINT_BIN=$(shell go env GOPATH)/bin/golangci-lint; \
  	fi
	@echo "Found golangci-lint at '$(LINT_BIN)', running..."; \
	${LINT_BIN} run

dbuildapi:
	docker build -t game-library:latest .

drunapi:
	docker compose up -d api

drunzipkin:
	docker compose up -d zipkin

drunredis:
	docker compose up -d redis

drunglog:
	docker compose up -d graylog

drunprom:
	docker compose up -d prometheus

dbuildmng:
	docker build -f Dockerfile.mng -t game-library-mng:latest .

