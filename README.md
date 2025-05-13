# Telegram Bot

A production-ready Golang Telegram bot using webhooks with PostgreSQL storage.

## Features

- Webhook-based Telegram bot
- PostgreSQL storage
- Database migrations
- Clean architecture
- Graceful shutdown

## Configuration

Edit the `config.toml` file with your settings:

```toml
[server]
webhook_url = "https://yourdomain.com/YOUR_SECRET_TOKEN"
port = 8443

[telegram]
bot_token = "YOUR_BOT_TOKEN"

[postgres]
host = "localhost"
port = 5432
user = "youruser"
password = "yourpassword"
dbname = "yourdb"
sslmode = "disable"
```

## Build and Run

Build the application:

```bash
make build
```

Run the application:

```bash
make run
```

## Database Migrations

Apply migrations:

```bash
make migrate-up
```

Rollback migrations:

```bash
make migrate-down
```

## Development

This project follows clean architecture principles with the following structure:

- `cmd/bot`: Main application entry point
- `internal/config`: Configuration handling
- `internal/db`: Database connection and operations
- `internal/telegram`: Telegram bot service
- `internal/handlers`: Message handlers
- `migrations`: Database migrations

## Requirements

- Go 1.21 or higher
- PostgreSQL
- golang-migrate CLI (for migrations)
