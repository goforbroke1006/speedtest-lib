.PHONY: all dep test

all: dep test

dep:
	go mod download

test:
	go test ./... -cover