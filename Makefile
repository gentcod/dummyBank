current_dir = $(shell pwd)
.PHONY: build
build:
	gofmt -l -s -w .
	go build -o bin/dummybank .

.PHONY: run
run:
	./bin/dummybank

sqlc:
	sqlc generate

.PHONY: sqlc-docker
sqlc-docker:
	docker run --rm -v $(current_dir):/src -w /src sqlc/sqlc generate

mysql:
	docker run --name mysql8 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8

postgres:
	docker run --name postgresdbs -p 5431:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

postgresBash:
	docker exec -it postgresdbs psql -h localhost -U root -d dummy_bank

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root dummy_bank

dropdb:
	docker exec -it postgres12 dropdb dummy_bank

gooseup:
	goose -dir sql/schemas postgres postgres://root:secret@localhost:5431/dummy_bank?sslmode=disable up

goosedown:
	goose -dir sql/schemas postgres postgres://root:secret@localhost:5431/dummy_bank?sslmode=disable down

.PHONY: goosedownall
goosedownall:
	goose -dir sql/schemas postgres postgres://root:secret@localhost:5431/dummy_bank?sslmode=disable down-to 0

migrateCreate:
	migrate create -ext sql -dir sql/migrations -seq add_sessions

test:
	go test -v -cover -short ./...

mock:
	mockgen -package mockdb -destination internal/database/mock/store.go github.com/gentcod/DummyBank/internal/database Store
	mockgen -package mockwk -destination worker/mock/distributo.go github.com/gentcod/DummyBank/worker TaskDistributor

buildimage:
	docker build -t dummybank:latest .

proto:
	rm -f pb/*go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=dummy_bank \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

redis:
	docker run --name redis -p 6379:6379 -d redis:7.4-alpine

.PHONY: sqlc mysql postgres createdb dropdb gooseup goosedown test mock migrateCreate buildimage postgresBash proto redis