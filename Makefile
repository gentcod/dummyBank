current_dir = $(shell pwd)
sqlc:
	docker run --rm -v $(current_dir):/src -w /src sqlc/sqlc generate

mysql:
	docker run --name mysql8 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8

postgres:
	docker run --name --network bank-network  postgres12 -p 5431:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root dummy_bank

dropdb:
	docker exec -it postgres12 dropdb dummy_bank

gooseup:
	goose -dir sql/schemas postgres postgres://root:secret@localhost:5431/dummy_bank?sslmode=disable up

goosedown:
	goose -dir sql/schemas postgres postgres://root:secret@localhost:5431/dummy_bank?sslmode=disable down

test:
	go test -v -cover -short ./...

mock:
	mockgen -package mockdb -destination internal/database/mock/store.go github.com/gentcod/DummyBank/internal/database Store

buildimage:
	docker build -t dummybank:latest .

.PHONY: sqlc mysql postgres createdb dropdb gooseup goosedown test mock 