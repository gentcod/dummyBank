current_dir = $(shell pwd)
sqlc:
	docker run --rm -v $(current_dir):/src -w /src sqlc/sqlc generate 

test:
	go test -v -cover ./..

.PHONY: sqlc test