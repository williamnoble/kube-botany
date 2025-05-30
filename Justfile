id := `curl http://localhost:8080/ -s | jq '.[0].Id' -r`

api:
	go run cmd/api/*.go

gen:
	go run cmd/gen/main.go

index:
	@curl http://localhost:8080/ -s | jq

water:
	@curl -s -X POST -H "Content-Type: application/json" -d '{"id": "{{ id }}"}' http://localhost:8080/water | jq

ascii:
	curl http://localhost:8080/ascii

lint:
	golangci-lint run ./...