migrateup:
	migrate -path db/migration -database "postgresql://postgres:1@localhost:5432/simple-bank?sslmode=disable" -verbose up

migrateupone:
	migrate -path db/migration -database "postgresql://root:postgres@simple-bank.che2e06ko90w.eu-central-1.rds.amazonaws.com:5432/simple-bank?sslmode=require" -verbose up 1
	
migratedown:
	migrate -path db/migration -database "postgresql://root:postgres@simple-bank.che2e06ko90w.eu-central-1.rds.amazonaws.com:5432/simple-bank?sslmode=require" -verbose down
	
migratedownone:
	migrate -path db/migration -database "postgresql://root:postgres@simple-bank.che2e06ko90w.eu-central-1.rds.amazonaws.com:5432/simple-bank?sslmode=require" -verbose down 1
		
sqlc:
	sqlc generate

testDb:
	go test -v -cover ./...

testApi:
	go test -v -cover ./api/...

server:
	go run main.go

createmigrate:
	migrate create -ext sql -dir db/migration -seq <<MIGRATION NAME>>

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/auronvila/simple-bank/db/sqlc Store

.PHONY: migrateup migratedown server mock