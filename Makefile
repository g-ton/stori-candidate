createdb:
	docker exec -it postgres12 createdb --username=root --owner=root stori-candidate

# Command to create migration entries: migrate create -ext sql -dir db/migration -seq init_schema
migrateup:
	migrate -path db/migration -database "postgresql://root:kuna1234@localhost:5432/stori-candidate?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:kuna1234@localhost:5432/stori-candidate?sslmode=disable" -verbose down

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=kuna1234 -d postgres:12-alpine

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/g-ton/stori-candidate/db/sqlc Store

mock_mail:
	mockgen -package mockmail -destination mail/mock/mail.go github.com/g-ton/stori-candidate/mail Mail

.PHONY: postgres createdb migrateup  migratedown sqlc server mock mock_mail