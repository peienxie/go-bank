all: gobank-image gobank

gobank-network:
	docker network inspect gobank-network >/dev/null 2>&1 || docker network create gobank-network
	docker network connect gobank-network postgres || true

gobank: gobank-network
	docker run --network gobank-net --name gobank -p 8080:8080 -e GOBANK_DB_SOURCE="postgresql://root:root@postgres:5432/go-bank?sslmode=disable" -e GOBANK_DB_DRIVER="postgres" -d --rm gobank:latest

gobank-image:
	docker build -t gobank .

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

gen-mockdb:
	mockgen -package mockdb -destination ./db/mock/store.go github.com/peienxie/go-bank/db/sqlc Store

.PHONY: all gobank-network gobank gobank-image postgres createdb dropdb migrateup migratedown gen-sqlc test lint server

