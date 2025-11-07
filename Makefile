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

SWAG_VERSION := v1.16
SWAG_PKG := github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)
generate-swag:
	@swag --version >/dev/null 2>&1 || { echo "Installing swag..."; go install ${SWAG_PKG}; }
	@echo "Found swag, generating documentation..."
	swag init \
	-d cmd/game-library-api,internal/api,internal/api/model,internal/web

MOCKGEN_VERSION := v0.6
MOCKGEN_PKG := go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)
generate-mocks:
	@mockgen -version >/dev/null 2>&1 || { echo "Installing mockgen..."; go install ${MOCKGEN_PKG}; }
	@echo "Found mockgen, generating mocks..."
	mockgen -source=internal/api/provider.go -destination=internal/api/mocks/provider.go -package=api_mock
	mockgen -source=internal/pkg/cache/redis.go -destination=internal/pkg/cache/mocks/redis.go -package=cache_mock
	mockgen -source=internal/facade/provider.go -destination=internal/facade/mocks/provider.go -package=facade_mock
	mockgen -source=internal/auth/auth.go -destination=internal/auth/mocks/auth.go -package=auth_mock
	mockgen -source=internal/middleware/auth.go -destination=internal/middleware/mocks/auth.go -package=middleware_mock
	mockgen -source=internal/taskprocessor/task.go -destination=internal/taskprocessor/mocks/task.go -package=taskprocessor_mock
	mockgen -destination=internal/repo/mocks/tx.go -package=repo_mock github.com/jackc/pgx/v5 Tx
	mockgen -source=internal/api/grpc/infoapi/service.go -destination=internal/api/grpc/infoapi/mocks/service.go -package=infoapi_mock

BUF_VERSION := v1.59
PROTOC_GEN_GO_VERSION := v1.36.10
PROTOC_GEN_GO_GRPC_VERSION := v1.5.1
BUF_PKG := github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
PROTOC_GEN_GO_PKG := google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION}
PROTOC_GEN_GO_GRPC_PKG := google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION}
generate-proto:
	@buf --version >/dev/null 2>&1 || { echo "Installing buf..."; go install ${BUF_PKG}; }
	@protoc-gen-go --version >/dev/null 2>&1 || { echo "Installing protoc-gen-go..."; go install ${PROTOC_GEN_GO_PKG}; }
	@protoc-gen-go-grpc --version >/dev/null 2>&1 || { echo "Installing protoc-gen-go-grpc..."; go install ${PROTOC_GEN_GO_GRPC_PKG}; }
	@echo "Generating protobuf code with buf..."; \
	buf generate

generate: generate-proto generate-swag generate-mocks

LINT_VERSION := v2.6
LINT_PKG := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(LINT_VERSION)
lint:
	@golangci-lint version >/dev/null 2>&1 || { echo "Installing golangci-lint..."; go install ${LINT_PKG}; }
	@echo "Found golangci-lint, running..."
	golangci-lint run

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
