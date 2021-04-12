.PHONY: build, test

build:
		go build -v ./cmd/queue_broker 

test: 
		go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build