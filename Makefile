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

dockerrunpg:
	docker compose up -d --no-recreate db

migrate:
	go run ./cmd/game-library-manage/. migrate

rollback:
	go run ./cmd/game-library-manage/. rollback

seed:
	go run ./cmd/game-library-manage/. seed

SWAG_PKG := github.com/swaggo/swag/cmd/swag@v1.16.4
SWAG_BIN := $(shell go env GOPATH)/bin/swag
MOCKGEN_PKG := go.uber.org/mock/mockgen@v0.5.0
MOCKGEN_BIN := $(shell go env GOPATH)/bin/mockgen
generate:
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
	${MOCKGEN_BIN} -destination=internal/taskprocessor/mocks/tx.go -package=taskprocessor_mock github.com/jackc/pgx/v5 Tx

LINT_PKG := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8
LINT_BIN := $(shell go env GOPATH)/bin/golangci-lint
lint:
	@if \[ ! -f ${LINT_BIN} \]; then \
		echo "Installing golangci-lint..."; \
    	go install ${LINT_PKG}; \
  	fi
	@if \[ -f ${LINT_BIN} \]; then \
  		echo "Found golangci-lint at '$(LINT_BIN)', running..."; \
	else \
    	echo "golangci-lint not found or the file does not exist"; \
    	exit 1; \
  	fi
	${LINT_BIN} run

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

dockerrunprom:
	docker compose up -d prometheus

dockerbuildmng:
	docker build -f Dockerfile.mng -t game-library-mng:latest .

