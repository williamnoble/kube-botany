
api:
	go run cmd/api/*.go

index:
	curl http://localhost:8080/ | jq

water:
	id=$(xh localhost:8080 | jq '.ID' -r)
	curl -X POST -H "Content-Type: application/json" -d '{"plant_id": "$id", "amount": 10}' http://localhost:8080/water

ascii:
	curl http://localhost:8080/ascii