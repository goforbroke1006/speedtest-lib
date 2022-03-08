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
	go test -bench=. -benchmem bench_test.go > new.bench.txt
	git stash
	go test -bench=. -benchmem bench_test.go > old.bench.txt
	benchcmp old.bench.txt new.bench.txt

coverage:
	go test --coverprofile ./.coverage ./...
	go tool cover -html ./.coverage

setup:
	go install golang.org/x/tools/cmd/benchcmp