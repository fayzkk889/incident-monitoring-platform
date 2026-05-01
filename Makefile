.PHONY: up down logs test-logs test clean

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f

test-logs:
	bash scripts/send_test_logs.sh 20

test:
	powershell.exe -ExecutionPolicy Bypass -File .\test-api.ps1

clean:
	docker compose down -v
