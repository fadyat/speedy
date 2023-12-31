PORT=8080

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

recreate-tag: delete-tag create-tag

create-tag:
	@echo "Tagging version $(VERSION)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)

delete-tag:
	@echo "Deleting tag $(VERSION)"
	@git tag -d $(VERSION)
	@git push origin --delete $(VERSION)

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
	@export GRPC_PORT=${PORT} && go run cmd/*.go

local:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-X 'main.Version=local' -s -w" \
		-o ./bin/main cmd/*.go
	@docker compose --file docker-compose.local.yml up --build && docker compose rm -fsv

play:
	@cd playground/client-config-sync && go run *.go

.PHONY: lint test
