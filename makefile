postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:latest

createdb:
	docker exec -it postgres createdb --username=root --owner=root go-bank

dropdb:
	docker exec -it postgres dropdb go-bank

migrateup:
	migrate -path ./db/schema -database "postgresql://root:root@localhost:5432/go-bank?sslmode=disable" -verbose up

migratedown:
	migrate -path ./db/schema -database "postgresql://root:root@localhost:5432/go-bank?sslmode=disable" -verbose down

gen-sqlc:
	sqlc generate

test:
	go test ./... -v -cover

lint:
	golangci-lint run

server:
	go run main.go


.PHONY: postgres createdb dropdb migrateup migratedown gen-sqlc test lint server

