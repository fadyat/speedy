ifneq (,$(wildcard ./.env))
	include .env
	export
endif

test:
	@go test $(FLGS) -cover ./... -coverprofile=cover.out
	@go tool cover -html=cover.out -o cover.html

lint:
	@golangci-lint run --issues-exit-code 1 --print-issued-lines=true --config .golangci.yml ./...

pre:
	@go mod tidy
	@go mod verify
	@make lint
	@make test

proto:
	@protoc --proto_path api \
		--go_out api \
		--go_opt paths=source_relative \
		--go-grpc_out api \
		--go-grpc_opt paths=source_relative \
		api/*.proto

run:
	@go run cmd/*.go


.PHONY: lint test
