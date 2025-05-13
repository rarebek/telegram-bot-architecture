package handlers

import (
	"bot/internal/db"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleMessage processes regular messages
func (h *Handler) handleMessage(message *tgbotapi.Message) {
	// Skip messages without a user
	if message.From == nil {
		return
	}

	// Try to get or create user
	user, err := h.db.GetOrCreateUser(
		h.ctx,
		message.From.ID,
		message.From.UserName,
		message.From.FirstName,
		message.From.LastName,
	)
	if err != nil {
		log.Printf("Error getting/creating user: %v", err)
	}

	// Check if this is a command
	isCommand := message.IsCommand()

	// Track user activity
	if user != nil {
		// Update user activity stats
		if err := h.db.UpdateUserActivity(h.ctx, message.From.ID, isCommand); err != nil {
			log.Printf("Error updating user activity: %v", err)
		}

		// Save message to database
		dbMessage := &db.Message{
			UserID:      user.ID, // Using DB user ID, not Telegram ID
			MessageText: message.Text,
			IsCommand:   isCommand,
			MessageType: "text", // Default type is text
			ChatID:      message.Chat.ID,
			ChatType:    message.Chat.Type,
		}

		// Check for other message types
		if message.Photo != nil {
			dbMessage.MessageType = "photo"
			dbMessage.MessageText = message.Caption
		} else if message.Document != nil {
			dbMessage.MessageType = "document"
			dbMessage.MessageText = message.Caption
		} else if message.Audio != nil {
			dbMessage.MessageType = "audio"
			dbMessage.MessageText = message.Caption
		} else if message.Video != nil {
			dbMessage.MessageType = "video"
			dbMessage.MessageText = message.Caption
		} else if message.Voice != nil {
			dbMessage.MessageType = "voice"
			dbMessage.MessageText = message.Caption
		} else if message.Sticker != nil {
			dbMessage.MessageType = "sticker"
		}

		if err := h.db.SaveMessage(h.ctx, dbMessage); err != nil {
			log.Printf("Error saving message: %v", err)
		}
	}

	// Process commands
	if isCommand {
		h.handleCommand(message)
		return
	}

	// Handle private messages differently than group messages
	if message.Chat.Type == "private" {
		h.handlePrivateMessage(message)
	} else {
		// For groups, only respond if mentioned or if it's a Go-related question
		if h.shouldRespondToGroupMessage(message) {
			h.handleGroupMessage(message)
		}
	}
}

// shouldRespondToGroupMessage determines if the bot should respond to a group message
func (h *Handler) shouldRespondToGroupMessage(message *tgbotapi.Message) bool {
	// Always respond if the bot is mentioned
	if message.Chat.Type != "private" && message.Entities != nil {
		for _, entity := range message.Entities {
			// Check for mention entities
			if entity.Type == "mention" {
				// Extract the username from the mention (removing the @ symbol)
				mention := message.Text[entity.Offset+1 : entity.Offset+entity.Length]
				if mention == h.bot.API.Self.UserName {
					return true
				}
			}
		}
	}

	// Use a simpler approach - just check if any keyword exists in the message
	if message.Text != "" {
		lowerText := strings.ToLower(message.Text)
		keywords := []string{
			"golang", "go lang", "go dasturlash", "go dastur",
			"goroutine", "channel", "interface", "struct",
			"go module", "go package", "go compiler", "gofmt",
		}

		for _, keyword := range keywords {
			if strings.Contains(lowerText, strings.ToLower(keyword)) {
				return true
			}
		}
	}

	return false
}

// handlePrivateMessage responds to messages in private chats
func (h *Handler) handlePrivateMessage(message *tgbotapi.Message) {
	// More personalized response in private chats - use the default keyword response
	response := h.msgSvc.GetKeywordResponse(message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending message:", err)
	}
}

// handleGroupMessage responds to relevant messages in group chats
func (h *Handler) handleGroupMessage(message *tgbotapi.Message) {
	// Use the message service to find the appropriate response
	response := h.msgSvc.GetKeywordResponse(message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ReplyToMessageID = message.MessageID
	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending message:", err)
	}
}
