This project is a simple example of a REST API using the Go language and the Postgres database.

### Usage
To run the project, you need to have Docker and Docker Compose installed on your machine.

You need to create a .env file in the root of the project with the following variables, for example:
```
DB_USER=my_user
DB_PASSWORD=my_password
DB_NAME=my_database
DB_PORT=5432
TELEGRAM_TOKEN = <your telegram bot api>
```  


To start the project, run the following command:
```shell
docker-compose up
```
next: 
```shell
./migrate.sh up
```





Command for migrate:
Up
```shell 
./migrate.sh up #up migrations
```

Down
```shell
./migrate.sh down #down  migrations
```

Create
```shell
./migrate.sh create <migration_name> #create migration
```


dump database
```shell
docker-compose exec db pg_dump -U my_user my_database > dump.sql
```


[//]: # (docker build -t app:local .)
