.PHONY: build run lint migrate

name ?= default

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
	goose -dir migrations postgres "postgres://user:password@localhost:5000/url_shortener?sslmode=disable" up

migrate-create:
ifndef name
	$(error Переменная name не задана. Используй: make migrate-create name=имя_миграции)
endif
	goose -dir migrations create $(name) sql