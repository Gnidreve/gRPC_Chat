.PHONY: build test lint lint-fix fmt run

build:
	go build ./...

test:
	go test -race -shuffle=on ./...

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

fmt:
	golangci-lint fmt ./...

run:
	go run ./cmd/server
