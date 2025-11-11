build:
	go build -o build/wami ./cmd/wami/*

test:
	go test ./... -race

coverage:
	go test -cover ./...

.PHONY: build test coverage
