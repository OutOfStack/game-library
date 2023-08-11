build:
	mkdir -p bin
	go build -o bin/game-library-api cmd/game-library-api/main.go

run:
	go run ./cmd/game-library-api/.

test:
	go test -v -race ./...

dockerrunpg:
	docker-compose up -d --no-recreate db

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

generate:
	@if ! command -v swag &> /dev/null; then\
	  	go install github.com/swaggo/swag/cmd/swag@v1.8.12\
  		exit;\
	fi
	swag init -g cmd/game-library-api/main.go

lint:
	@if ! command -v golangci-lint &> /dev/null; then\
	  	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest;\
  		exit;\
	fi
	golangci-lint run

dockerbuildweb:
	docker build -t game-library-web:latest .

dockerrunweb:
	docker compose up -d web

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
