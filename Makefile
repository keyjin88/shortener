.PHONY: build

build:
	go build -v ./cmd/shortener

test:
	go test ./...

.DEFAULT_GOAL:= build