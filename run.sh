#!/bin/bash
set -e

# Build and run the bot
go build -o bin/bot ./cmd/bot
./bin/bot
