.PHONY: all dep test lint benchmark benchcmp coverage setup

all: dep test lint

dep:
	go mod download

test:
	go test ./... -cover

lint:
	golangci-lint run

benchmark:
	go test -gcflags="-N" ./... -bench=.

benchcmp:
	go test -gcflags="-N" -bench=. -benchmem ./... > new.bench.txt
	git stash
	go test -gcflags="-N" -bench=. -benchmem ./... > old.bench.txt
	git stash pop
	benchcmp old.bench.txt new.bench.txt

coverage:
	go test --coverprofile ./.coverage ./...
	go tool cover -html ./.coverage

setup:
	go install golang.org/x/tools/cmd/benchcmp