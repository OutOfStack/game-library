build:
	mkdir -p bin
	go build -o bin/game-library-api cmd/game-library-api/main.go

run:
	go run ./cmd/game-library-api/.

runpg:
	docker-compose up -d

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

.PHONY: build run runpg createdb dropdb migrate rollback seed swaggen