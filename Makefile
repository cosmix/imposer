.PHONY: build test lint clean

build:
	go build -o bin/imposer ./cmd/impose

test:
	go test -v ./...

lint:
	go vet ./...

clean:
	rm -rf bin/

all: lint test build
