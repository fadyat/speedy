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
	@go run cmd/*.go

runc:
	@go run cmd/client/*.go

.PHONY: lint test
