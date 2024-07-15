include .env

deps:
	go install github.com/air-verse/air@v1.52.3
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

up:
	docker compose up -d

down:
	docker compose down

migrate:
	migrate -database ${POSTGRES_URL} -path db/migrations up

migrate-down:
	migrate -database ${POSTGRES_URL} -path db/migrations down
