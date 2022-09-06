TESTS = $(addprefix github.com/Gr1LyA/L0_golang/internal/app/, server storage stan)

.PHONY: build
build:
		go build -v ./cmd/service

pub:
		go build -v ./cmd/stan-pub

.PHONY: test
test:
		go test -v -timeout 30s $(TESTS)

clean:
		rm service stan-pub

.DEFAULT_GOAL := build