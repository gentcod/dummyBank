current_dir = $(shell pwd)
sqlc:
	docker run --rm -v $(current_dir):/src -w /src sqlc/sqlc generate 

mysql:
	docker run --name mysql8 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root dummy_bank

dropdb:
	docker exex -it postgres12 dropdb dummy_bank

migrateup:
	migrate -path ./sql/schemas -database postgres://root:secret@localhost:5432/dummy_bank?sslmode=disable -verbose up

migratedown:
		migrate -path sql/schemas -database postgres://root:secret@localhost:5432/dummy_bank?sslmode=disable -verbose down

gooseup:
	goose postgres postgres://root:secret@localhost:5432/dummy_bank?sslmode=disable up

test:
	go test -v -cover ./...

.PHONY: sqlc postgres createdb dropdb migrateup migratedown test