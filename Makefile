build:
	mkdir -p bin
	go build -o bin/game-library-api cmd/game-library-api/main.go

run:
	go run ./cmd/game-library-api/.

dockerrunpg:
	docker-compose up -d db

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

swaggen:
	./tools/swag init -g cmd/game-library-api/main.go

dockerbuildweb:
	docker build -f Dockerfile.web -t game-library-web:latest .

dockerrunweb:
	docker-compose up -d web

dockerbuildmng:
	docker build -f Dockerfile.mng -t game-library-mng:latest .