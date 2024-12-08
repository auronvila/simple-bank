migrateup:
	migrate -path db/migration -database "postgresql://postgres:1@localhost:5432/simple-bank?sslmode=disable" -verbose up
	
migratedown:
	migrate -path db/migration -database "postgresql://postgres:1@localhost:5432/simple-bank?sslmode=disable" -verbose down
	
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/auronvila/simple-bank/db/sqlc Store

.PHONY: migrateup migratedown server mock