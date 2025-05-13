# Handlers Package

This package contains the Telegram bot handlers for processing various types of updates.

## Files:

- `types.go`: Core handler types and initialization
- `update_handlers.go`: Main update routing logic
- `message_handlers.go`: Message processing handlers
- `command_handlers.go`: Command processing handlers
- `callback_handlers.go`: Callback query handling

## Usage:

Initialize the handler using the `New` function from `types.go`, then use the `HandleUpdate` method to process Telegram updates.
