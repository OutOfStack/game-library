# game-library

### Makefile:
    make build      builds app
    make run        runs app
    make runpg      starts postgres server in container
    make createdb   creates database on postgres server started by 'make runpg'
    make dropdb     drops database on postgres server created by 'make runpg'
    make migrate    applies all migrations to database
    make rollback   rollbacks one last migration on database
    make seed       seeds test data to database


### migrate CLI
https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

Creating migrations:
migrate create -ext sql -dir /migrations -seq migration_name