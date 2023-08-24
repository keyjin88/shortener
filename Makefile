.PHONY: build

build:
	go build -v ./cmd/shortener

test:
	go test ./...

test-covered:
	  go test -covermode=count ./...

.DEFAULT_GOAL:= build