# Stack

- Persistent Data: Postgres; Redis
- Cache: Redis *(future planned)*

## Postgres
Postgres is used for storing relational data for web applications such as the workbench.

## Redis
Within the context of Brahman (i.e. running Talkative applications) Redis is used for storing commands for compiled Talkative applications. These commands can mutate the current user's state within the application and persist their state, or mutate the output of the application.

# Setup
1. Install Golang
2. [Ensure go is setup properly](https://golang.org/doc/code.html)
3. Install Docker
4. Install docker-compose if not on Mac
5. Download backend services
  - `$ go get github.com/talkative-ai/brahman`
  - `$ go get github.com/talkative-ai/lakshmi`
  - `$ go get github.com/talkative-ai/shiva`
6. Enter the core directory *(this was automatically downloaded from step 3)*
  - `$ cd $GOPATH/src/github.com/talkative-ai/core`
7. Run docker-compose
  - `$ docker-compose up`

The services should now be available on the following ports at localhost (127.0.0.1):

*Talkative Services*
- shiva:8000
- lakshmi:8042
- brahman:9001

*Additional Services*
- postgres:5432
- adminer:8001
- redis:6380

## Database
0. Ensure you have the migration tool installed
  1. # go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
  2. # go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli
1. Enter the core/db directory
  - `$ cd $GOPATH/src/github.com/talkative-ai/core/db`
2. Run the migrations
  - `$ migrate -database "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" -path migrations up`

If you ever need to reset, you can call "migrations down" and then repeat "migrations up" for fresh data

## Google Actions
To use Brahman live on Google Actions:
1. install ngrok
  - `$ ngrok http 9001`
2. Update the api.ai project to point to the new ngrok address

Once it's up and running, then the workbench frontend will work (edited)
via npm install and npm start