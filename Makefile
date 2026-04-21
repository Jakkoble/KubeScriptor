.PHONY: up pull

up:
	@set -e; \
	trap 'docker compose down' EXIT; \
	docker compose pull commander tui; \
	docker compose up -d commander; \
	docker compose run --rm --use-aliases tui

pull:
	docker compose pull commander tui
