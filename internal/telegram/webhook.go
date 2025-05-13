package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SetupWebhook configures the Telegram webhook
func (b *Bot) SetupWebhook() error {
	wh, err := tgbotapi.NewWebhook(b.WebhookURL)
	if err != nil {
		return fmt.Errorf("failed to create webhook: %w", err)
	}

	_, err = b.API.Request(wh)
	if err != nil {
		return fmt.Errorf("webhook registration failed: %w", err)
	}

	info, err := b.API.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("failed to get webhook info: %w", err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return nil
}

// GetWebhookHandler returns a handler for the webhook endpoint
func (b *Bot) GetWebhookHandler(updateHandler func(update tgbotapi.Update)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Println("Error decoding update:", err)
			return
		}

		// Log the entire update to debug
		log.Printf("Received update: %+v\n", update)

		// Process the update
		updateHandler(update)
	}
}
