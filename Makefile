.PHONY: build run lint migrate

dev:
	go run ./cmd/app/main.go --config=./configs/local.yaml

build:
	go build -o bin/app ./cmd/app/main.go

run: build
	./bin/app

test:
	go test -v ./...

lint:
	golangci-lint run

migrate:
	goose -dir migrations postgres "postgres://user:password@localhost:5432/url_shortener?sslmode=disable" up