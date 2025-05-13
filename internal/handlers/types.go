package handlers

import (
	"bot/internal/db"
	"bot/internal/services"
	"bot/internal/telegram"
	"context"
)

// Handler handles Telegram updates
type Handler struct {
	bot    *telegram.Bot
	db     *db.DB
	ctx    context.Context
	msgSvc *services.MessageService
}

// New creates a new Handler
func New(bot *telegram.Bot, database *db.DB) (*Handler, error) {
	// Initialize message service
	msgSvc, err := services.NewMessageService()
	if err != nil {
		return nil, err
	}

	return &Handler{
		bot:    bot,
		db:     database,
		ctx:    context.Background(),
		msgSvc: msgSvc,
	}, nil
}
