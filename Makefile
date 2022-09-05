.PHONY: build
build:
		go build -v ./cmd/service

pub:
		go build -v ./cmd/stan-pub

.PHONY: test
test:
		go test -v -race -timeout 30s ./...

clean:
		rm service stan-pub

.DEFAULT_GOAL := build