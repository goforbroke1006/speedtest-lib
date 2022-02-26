.PHONY: all dep test benchmark coverage

all: dep test

dep:
	go mod download

test:
	go test ./... -cover

benchmark:
	go test -gcflags="-N" ./... -bench=.

coverage:
	go test --coverprofile ./.coverage ./...
	go tool cover -html ./.coverage