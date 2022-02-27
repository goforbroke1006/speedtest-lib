.PHONY: all dep test lint benchmark coverage

all: dep test lint

dep:
	go mod download

test:
	go test ./... -cover

lint:
	golangci-lint run

benchmark:
	go test -gcflags="-N" ./... -bench=.

coverage:
	go test --coverprofile ./.coverage ./...
	go tool cover -html ./.coverage