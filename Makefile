.PHONY: all dep test coverage

all: dep test

dep:
	go mod download

test:
	go test ./... -cover

coverage:
	go test --coverprofile ./.coverage ./...
	go tool cover -html ./.coverage