test:
	go test ./app/internal/test

run:
	go run main.go

docker-build:
	docker build -t sber_task:local .

docker-compose-up:
	docker compose -f docker-compose.yml up

