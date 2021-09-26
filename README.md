# game-library
Is an app for exploring and rating games

### Usage:
    make build          builds app
    make run            runs app
    make dockerrunpg    runs postgres server in docker container
    make createdb       creates database on postgres server started by 'make dockerrunpg'
    make dropdb         drops database on postgres server created by 'make dockerrunpg'
    make migrate        applies all migrations to database
    make rollback       rollbacks one last migration on database
    make seed           seeds test data to database
    make swaggen        generates documentation for swagger UI
    make dockerbuildweb builds web app docker image
    make dockerrunweb   runs web app in docker container
    make dockerbuildmng builds manage app docker image

### Migrations
Creating a new migration:

`touch {i}_{name}.up.sql {i}_{name}.down.sql` , where  
`{i}` - consecutive migration ID of length 6 padded with zeroes,  
`{name}` - migration name

### Swagger
Swagger file generation tool is located in `tools/swag`  
Swagger UI url: http://localhost:8000/swagger/index.html