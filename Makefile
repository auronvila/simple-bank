DB_URL="postgresql://postgres:1@localhost:5432/simple-bank?sslmode=disable"

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateupone:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
	
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
	
migratedownone:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
		
sqlc:
	sqlc generate

testDb:
	go test -v -cover ./...

testApi:
	go test -v -cover ./api/...

server:
	go run main.go

# when using air if the files are updated the server restarts automatically
devServer:
	air

createmigrate:
	migrate create -ext sql -dir db/migration -seq <<MIGRATION NAME>>

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/auronvila/simple-bank/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
        --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
        proto/*.proto 
		statik -src=./doc/swagger -dest=./doc

.PHONY: migrateup migratedown server mock proto