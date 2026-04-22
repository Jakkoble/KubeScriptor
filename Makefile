.PHONY: setup start

setup:
	@if [ ! -d jobs ]; then \
		mkdir -p jobs; \
		printf '%s\n' \
			'name: demo-job' \
			'commands:' \
			'  - echo "Hello from HexaTask"' \
			'  - echo "Working directory:"' \
			'  - pwd' \
			> jobs/demo.yaml; \
	fi

start: setup
	@set -e; \
	trap 'docker compose down' EXIT; \
	docker compose pull commander tui; \
	docker compose up -d commander; \
	echo "Waiting for commander..."; \
	until curl -sf --http2-prior-knowledge http://localhost:5271/health > /dev/null 2>&1; do \
		sleep 0.5; \
	done; \
	echo "Commander ready!"; \
	docker compose run --rm --use-aliases tui
