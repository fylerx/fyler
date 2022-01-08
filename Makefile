.PHONY: docker-build

APP_NAME := fylerx
REPOSITORY := fylerx
VERSION := $(if $(TAG),$(TAG),$(if $(BRANCH_NAME),$(BRANCH_NAME),$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)))
NOCACHE := $(if $(NOCACHE),"--no-cache")
POSTGRESQL_URL='postgres://fylerx:fylerx@localhost:5432/fylerx?sslmode=disable'

docker-build:
	@docker build ${NOCACHE} --pull -f ./Dockerfile -t ${REPOSITORY}/${APP_NAME}:${VERSION} --ssh default .

db-migrate:
	@migrate -database ${POSTGRESQL_URL} -path db/migrations -verbose up
db-rollback:
	@migrate -database ${POSTGRESQL_URL} -path db/migrations -verbose down

psql-exec:
	PGPASSWORD=fylerx docker-compose exec postgres \
	psql -h postgres -p 5432 -U fylerx -d fylerx

run:
	@docker-compose up -d

stop:
	@docker-compose down

test:
	@go test -v -race ./...

cover:
	@go test -coverprofile=cover.out ./... && go tool cover -html=cover.out -o cover.html
