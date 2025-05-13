package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleCallback processes callback queries
func (h *Handler) handleCallback(callback *tgbotapi.CallbackQuery) {
	// Process based on callback data
	data := callback.Data
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "")

	switch data {
	case "get_information":
		msg.Text = "Here is some information about our services."
	case "start_action":
		msg.Text = "Action started successfully!"
	default:
		msg.Text = "Unknown action: " + data
	}

	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending callback message:", err)
	}
}
