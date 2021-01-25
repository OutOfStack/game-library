build:
	mkdir -p bin
	go build -o bin/game-library cmd/game-library/main.go

run:
	go run ./cmd/game-library/.

runpg:
	docker-compose up -d

createdb:
	docker exec -it games_db createdb --username=postgres --owner=postgres games

dropdb:
	docker exec -it games_db dropdb --username=postgres games

migrate:
	go run ./cmd/game-library/. migrate

rollback:
	go run ./cmd/game-library/. rollback

seed:
	go run ./cmd/game-library/. seed

.PHONY: build run runpg createdb dropdb migrate rollback seed