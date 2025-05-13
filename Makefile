.PHONY: build run migrate-up migrate-down

# Build the application
build:
	go build -o bin/bot ./cmd/bot

# Run the application
run:
	go run ./cmd/bot

# Apply database migrations
migrate-up:
	migrate -path migrations -database "postgres://$(shell grep -A 5 '\[postgres\]' config.toml | grep 'user' | cut -d'"' -f2):$(shell grep -A 5 '\[postgres\]' config.toml | grep 'password' | cut -d'"' -f2)@$(shell grep -A 5 '\[postgres\]' config.toml | grep 'host' | cut -d'"' -f2):$(shell grep -A 5 '\[postgres\]' config.toml | grep 'port' | cut -d'=' -f2 | tr -d ' ')/$(shell grep -A 5 '\[postgres\]' config.toml | grep 'dbname' | cut -d'"' -f2)?sslmode=$(shell grep -A 5 '\[postgres\]' config.toml | grep 'sslmode' | cut -d'"' -f2)" up

# Rollback database migrations
migrate-down:
	migrate -path migrations -database "postgres://$(shell grep -A 5 '\[postgres\]' config.toml | grep 'user' | cut -d'"' -f2):$(shell grep -A 5 '\[postgres\]' config.toml | grep 'password' | cut -d'"' -f2)@$(shell grep -A 5 '\[postgres\]' config.toml | grep 'host' | cut -d'"' -f2):$(shell grep -A 5 '\[postgres\]' config.toml | grep 'port' | cut -d'=' -f2 | tr -d ' ')/$(shell grep -A 5 '\[postgres\]' config.toml | grep 'dbname' | cut -d'"' -f2)?sslmode=$(shell grep -A 5 '\[postgres\]' config.toml | grep 'sslmode' | cut -d'"' -f2)" down
