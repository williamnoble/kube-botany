api:
	go run cmd/api/*.go

lint:
	golangci-lint run ./...