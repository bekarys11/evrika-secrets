dev:
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose --file docker-compose.dev.yml up --build



## up: starts all containers in the background without forcing build
up:
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"
	make test

## up_build: stops docker compose (if running), builds all projects and starts docker compose
up_build:
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## start: starts the server
start:
	@echo "Running test"
	go test -count=1 ./test/...

	@echo "Starting server"
	go run cmd/api/main.go

swagger:
	@echo "Generating swagger documentation"
	swag init -g pkg/config/routes.go
	@echo "Swagger documentation is available at: http://localhost:8888/swagger/index.html"

test:
	@echo "Running test"
	go test ./test/...
