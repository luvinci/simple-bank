migrateup:
	migrate -path db/migration -database "postgresql://postgres:123456@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:123456@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

.PHONY: migrateup migratedown
