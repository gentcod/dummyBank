sqlc:
	docker run --rm -v $($(pwd)):/src -w /src sqlc/sqlc generate 