.PHONY: up down build logs migrate

up:
	docker compose up --build

down:
	docker compose down --volumes

build:
	docker compose build

logs:
	docker compose logs -f

migrate:
	docker compose run --rm migrate
