build:
	go build -o build/wami ./cmd/wami/*

test:
	go test ./... -race

coverage:
	mkdir -p tmp
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: build test coverage
