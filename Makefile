.PHONY: run build test clean

run:
    go run cmd/api/main.go

build:
    go build -o bin/api cmd/api/main.go

test:
    go test -v ./...

clean:
    rm -rf bin/
