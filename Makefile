include .env
export

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Goose
goose-up:
	goose -dir=$(MIGRATIONS_DIR) postgres "$(DB_URL)" up

goose-down:
	goose -dir=$(MIGRATIONS_DIR) postgres "$(DB_URL)" down

# Start
run:
	go run main.go