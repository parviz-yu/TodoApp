.PHONY: build
build:
	go build -v ./cmd/todoapp

.PHONY: test
test:
	echo 'mode: atomic' > coverage.txt && go test -covermode=atomic -coverprofile=coverage.txt -v -race -timeout=30s ./...

.PHONY: cover
cover: test
	go tool cover -html=coverage.txt

.PHONY: create-migration
create-migration:
	migrate create -ext sql -dir schema/ -seq init

.PHONY: migrate
migrate:
	migrate -path schema/ -database "postgresql://localhost/todoapp?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	migrate -path schema/ -database "postgresql://localhost/todoapp?sslmode=disable" down 1

.PHONY: migrate-drop
migrate-drop:
	migrate -path schema/ -database "postgresql://localhost/todoapp?sslmode=disable" drop

.DEFAULT_GOAL := build