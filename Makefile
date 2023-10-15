ifneq (,$(wildcard ./.env))
	include .env
	export
endif

test:
	@go test $(FLGS) -cover ./... -coverprofile=cover.out
	@go tool cover -html=cover.out -o cover.html

lint:
	@golangci-lint run --issues-exit-code 1 --print-issued-lines=true --config .golangci.yml ./...

pre: test lint
	@go mod tidy
	@go mod verify

.PHONY: lint test
