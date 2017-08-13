# Database

The database is Postgres.

`$ docker-compose up` from Shiva to run the database and associated servers
adminer can now be accessed via 127.0.0.1:8001, and postgres via 127.0.0.1:5432

In order to conenct to postgres via adminer:
Host: postgres
database: postgres
username: postgres

There is no default password

## Migrations
1. Install https://github.com/mattes/migrate `$ go get github.com/mattes/migrate`
2. To run: `$ cd ./db && migrate -database "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" -path migrations up`