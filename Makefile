.PHONY: all dep test lint benchmark benchcmp coverage setup

all: dep test lint

dep:
	go mod download

gen:
	go generate ./...

test:
	go test ./... -cover

lint:
	golangci-lint run

benchmark:
	go test -gcflags="-N" -bench=. ./...

benchcmp:
	go test -gcflags="-N" -bench=. -benchmem ./... > new.bench.txt
	@git stash
	go test -gcflags="-N" -bench=. -benchmem ./... > old.bench.txt
	@git stash pop
	benchcmp old.bench.txt new.bench.txt

coverage:
	go test --coverprofile ./.coverage ./...
	go tool cover -html ./.coverage

setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2
	golangci-lint --version
	go install github.com/golang/mock/mockgen@v1.6.0
	mockgen --version
	go install golang.org/x/tools/cmd/benchcmp@v0.1.9
