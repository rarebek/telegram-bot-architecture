package telegram

import (
	"fmt"
	"log"

	"bot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents a Telegram bot instance
type Bot struct {
	API        *tgbotapi.BotAPI
	Config     *config.TelegramConfig
	WebhookURL string
}

// New creates a new Telegram bot
func New(cfg *config.TelegramConfig, webhookURL string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Bot{
		API:        bot,
		Config:     cfg,
		WebhookURL: webhookURL,
	}, nil
}
