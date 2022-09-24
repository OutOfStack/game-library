# game-library
An app for exploring and rating games

### Usage with `Make`:
    build          builds app
    run            runs app
    test           runs tests for the whole project
    dockerrunpg    runs postgres server in docker container
    createdb       creates database on postgres server started by 'make dockerrunpg'
    dropdb         drops database on postgres server created by 'make dockerrunpg'
    migrate        applies all migrations to database
    rollback       rollbacks one last migration on database
    seed           seeds test data to database
    generate       generates documentation for swagger UI
    dockerbuildweb builds web app docker image
    dockerrunweb   runs web app in docker container
    dockerbuildmng builds manage app docker image
    dockerrunmng-m applies migrations to database using docker manage image
    dockerrunmng-r rollbacks one last migration using docker manage image
    dockerrunmng-s applies migrations to database using docker manage image

### Migrations
Creating a new migration:

`touch {i}_{name}.up.sql {i}_{name}.down.sql` , where  
`{i}` - consecutive migration ID of length 6 padded with zeroes,  
`{name}` - migration name

### Swagger 
Swagger UI url: http://localhost:8000/swagger/index.html