postgres:
	docker run --name postgres -d -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root postgres:14-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "mysql://root:root@tcp(localhost:3306)/simple_bank" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://root:root@tcp(localhost:3306)/simple_bank" -verbose down

migrateuptest:
	migrate -path db/migration -database "mysql://root:root@tcp(localhost:3306)/simple_bank_test" -verbose up

migratedowntest:
	migrate -path db/migration -database "mysql://root:root@tcp(localhost:3306)/simple_bank_test" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test migrateuptest migratedowntest server