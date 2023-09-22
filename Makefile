test:
	go test ./
run:
	go run main.go

docker-compose-up:
	docker compose -f docker-compose.yml up

docker-build:
	docker build -t sber_task:local .
