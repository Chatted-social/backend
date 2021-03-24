# Backend for Chatted social network


## Development

### Run app via docker-compose:

First, run `docker-compose -f docker-compose-dev.yml up -d`,
this is will setup db for you.


Then u need to connect to your container using psql and create database `chatted`

`docker exec -ti backend_postgres_1 psql -U postgres`

note that `backend_postgres_1` might be not actual name of your container

`create database chatted;`.

Then apply migrations `docker-compose -f docker-compose-dev.yml up -d migrate`.

Then run `go build ./cmd/apiserver` and `./apiserver`.

### How to run swagger:
First, install it
`go get -u github.com/go-swagger/go-swagger/cmd/swagger`
then go to docs directory `cd ./docs`, then run 

`swagger serve -F=swagger swagger.yaml`
    
