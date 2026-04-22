.PHONY: run test lint build clean

run:
	go run cmd/server/main.go

test:
	go test -v ./...

lint:
	golangci-lint run ./...

build:
	go build -o bin/notes-api cmd/server/main.go

clean:
	rm -rf bin/