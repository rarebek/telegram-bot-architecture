package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleUpdate processes a Telegram update
func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	// Handle Callback Queries
	if update.CallbackQuery != nil {
		h.handleCallback(update.CallbackQuery)
		return
	}

	// Handle regular messages
	if update.Message != nil {
		h.handleMessage(update.Message)
		return
	}
}

