ifneq (,$(wildcard ./.env))
	include .env
	export
endif

lint:
	@golangci-lint run --issues-exit-code 1 --print-issued-lines=true --config .golangci.yml ./...

.PHONY: lint
