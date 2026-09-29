APP=bin/server

ifneq (,$(wildcard .env))
include .env
export
endif

run:
	go run ./cmd/server

dev:
	air -c .air.toml

build:
	go build -o $(APP) ./cmd/server

start: build
	./$(APP)

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-force:
	migrate -path migrations -database "$(DATABASE_URL)" force $(version)

migrate-create:
	migrate create -ext sql -dir migrations $(name)

migrate-reset:
	migrate -path migrations -database "$(DATABASE_URL)" drop -f

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

tidy:
	go mod tidy

install-tools:
	go install github.com/air-verse/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
