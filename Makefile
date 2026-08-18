test:
	go test ./...

build:
	go build -o bin/kinc ./cmd/kinc

.PHONY: test build
