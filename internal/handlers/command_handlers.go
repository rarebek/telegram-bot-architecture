package handlers

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleCommand processes bot commands
func (h *Handler) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		h.handleStart(message)
	case "help":
		h.handleHelp(message)
	case "rules":
		h.handleRules(message)
	case "about":
		h.handleAbout(message)
	case "group":
		h.handleGroup(message)
	case "roadmap":
		h.handleRoadmap(message)
	case "useful":
		h.handleUseful(message)
	case "latest":
		h.handleLatest(message)
	case "version":
		h.handleVersion(message)
	case "warn":
		h.handleWarn(message)
	case "stats":
		h.handleStats(message)
	}
}

func (h *Handler) handleStart(message *tgbotapi.Message) {
	text, err := h.msgSvc.GetCommandText("start")
	if err != nil {
		log.Printf("Error getting start command text: %v", err)
		text = "Assalomu alaykum! GoferUz Golang botiga xush kelibsiz 👋"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending command response:", err)
	}
}

func (h *Handler) handleHelp(message *tgbotapi.Message) {
	// Check if the user is an admin
	isAdmin, err := h.isAdmin(message.Chat.ID, message.From.ID)
	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		isAdmin = false
	}

	// Get help text with optional admin content
	text, err := h.msgSvc.GetMessageForAdmin("help", isAdmin || message.Chat.Type == "private")
	if err != nil {
		log.Printf("Error getting help command text: %v", err)
		text = "Mavjud komandalar ro'yxati..."
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending command response:", err)
	}
}

func (h *Handler) handleRules(message *tgbotapi.Message) {
	text, err := h.msgSvc.GetCommandText("rules")
	if err != nil {
		log.Printf("Error getting rules command text: %v", err)
		text = "GoferUz hamjamiyati qoidalari..."
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending command response:", err)
	}
}

func (h *Handler) handleAbout(message *tgbotapi.Message) {
	text, err := h.msgSvc.GetCommandText("about")
	if err != nil {
		log.Printf("Error getting about command text: %v", err)
		text = "Bu bot Go dasturlash tilida yaratilgan..."
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending command response:", err)
	}
}

func (h *Handler) handleGroup(message *tgbotapi.Message) {
	text, err := h.msgSvc.GetCommandText("group")
	if err != nil {
		log.Printf("Error getting group command text: %v", err)
		text = "Go dasturlash tili bo'yicha guruhlar va hamjamiyatlar..."
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending command response:", err)
	}
}

func (h *Handler) handleRoadmap(message *tgbotapi.Message) {
	text, err := h.msgSvc.GetCommandText("roadmap")
	if err != nil {
		log.Printf("Error getting roadmap command text: %v", err)
		text = "Go dasturlash tilini o'rganish uchun yo'l xaritasi..."
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending command response:", err)
	}
}

func (h *Handler) handleUseful(message *tgbotapi.Message) {
	text, err := h.msgSvc.GetCommandText("useful")
	if err != nil {
		log.Printf("Error getting useful command text: %v", err)
		text = "Go dasturlash tili bo'yicha foydali ma'lumotlar..."
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending command response:", err)
	}
}

func (h *Handler) handleLatest(message *tgbotapi.Message) {
	text, err := h.msgSvc.GetCommandText("latest")
	if err != nil {
		log.Printf("Error getting latest command text: %v", err)
		text = "Go 1.21 - Eng so'nggi Go versiyasi..."
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending command response:", err)
	}
}

func (h *Handler) handleVersion(message *tgbotapi.Message) {
	// By default, show the latest version info
	version := "1.21"

	// If there's a parameter after the command, use that as the version
	if len(message.CommandArguments()) > 0 {
		version = message.CommandArguments()
	}

	// Get version-specific text
	text, err := h.msgSvc.GetVersionText(version)
	if err != nil {
		log.Printf("Error getting version command text: %v", err)
		text = fmt.Sprintf("Kechirasiz, \"%s\" versiyasi haqida ma'lumot mavjud emas.", version)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending command response:", err)
	}
}

// Helper method to check if user is admin
func (h *Handler) isAdmin(chatID, userID int64) (bool, error) {
	config := tgbotapi.ChatConfigWithUser{
		ChatID: chatID,
		UserID: userID,
	}

	chatMember, err := h.bot.API.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: config,
	})
	if err != nil {
		return false, err
	}

	return chatMember.Status == "administrator" || chatMember.Status == "creator", nil
}

func (h *Handler) handleWarn(message *tgbotapi.Message) {
	// Check if the user is replying to another message
	if message.ReplyToMessage == nil {
		errorText, _ := h.msgSvc.GetCommandTemplate("warn")
		text, _ := errorText["error_no_reply"].(string)
		if text == "" {
			text = "Bu buyruqni ishlatish uchun biror xabarga javob sifatida yuboring."
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		if _, err := h.bot.API.Send(msg); err != nil {
			log.Println("Error sending command response:", err)
		}
		return
	}

	// Check if the user is an admin
	isAdmin, err := h.isAdmin(message.Chat.ID, message.From.ID)
	if err != nil {
		errorText, _ := h.msgSvc.GetCommandTemplate("warn")
		text, _ := errorText["error_status"].(string)
		if text == "" {
			text = "Foydalanuvchi statusini tekshirishda xatolik yuz berdi."
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		if _, err := h.bot.API.Send(msg); err != nil {
			log.Println("Error sending error message:", err)
		}
		return
	}

	// Check if the user is an admin or creator
	if !isAdmin && message.Chat.Type != "private" {
		errorText, _ := h.msgSvc.GetCommandTemplate("warn")
		text, _ := errorText["error_permission"].(string)
		if text == "" {
			text = "Bu buyruqni faqat adminlar ishlatishi mumkin."
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		if _, err := h.bot.API.Send(msg); err != nil {
			log.Println("Error sending command response:", err)
		}
		return
	}

	// Send warning message
	text, err := h.msgSvc.GetCommandText("warn")
	if err != nil {
		log.Printf("Error getting warn command text: %v", err)
		text = "⚠️ Ogohlantirish ⚠️\n\nIltimos, faqat Go dasturlash tiliga oid mavzularda suhbatlashing..."
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.ReplyToMessage.MessageID

	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending command response:", err)
	}

	// Log the warning action
	log.Printf("Admin %s (ID: %d) issued a warning in chat %d in response to message %d",
		message.From.UserName, message.From.ID, message.Chat.ID, message.ReplyToMessage.MessageID)
}

func (h *Handler) handleStats(message *tgbotapi.Message) {
	// Check if the user is an admin
	isAdmin, err := h.isAdmin(message.Chat.ID, message.From.ID)
	if err != nil {
		errorText, _ := h.msgSvc.GetCommandTemplate("stats")
		text, _ := errorText["error_status"].(string)
		if text == "" {
			text = "Foydalanuvchi statusini tekshirishda xatolik yuz berdi."
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		if _, err := h.bot.API.Send(msg); err != nil {
			log.Println("Error sending error message:", err)
		}
		return
	}

	// Check if the user is an admin or creator, or if it's a private chat
	if !isAdmin && message.Chat.Type != "private" {
		errorText, _ := h.msgSvc.GetCommandTemplate("stats")
		text, _ := errorText["error_permission"].(string)
		if text == "" {
			text = "Bu buyruqni faqat adminlar ishlatishi mumkin."
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		if _, err := h.bot.API.Send(msg); err != nil {
			log.Println("Error sending command response:", err)
		}
		return
	}

	// Update daily stats
	if err := h.db.UpdateDailyStats(h.ctx); err != nil {
		log.Println("Error updating daily stats:", err)
	}

	// Get message statistics
	messageStats, err := h.db.GetMessageStats(h.ctx)
	if err != nil {
		log.Println("Error getting message stats:", err)
	}

	// Get top 10 most active users
	topUsers, err := h.db.GetTopUsers(h.ctx, 10)
	if err != nil {
		log.Println("Error getting top users:", err)
		errorText, _ := h.msgSvc.GetCommandTemplate("stats")
		text, _ := errorText["error_data"].(string)
		if text == "" {
			text = "Foydalanuvchilar statistikasini olishda xatolik yuz berdi."
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		if _, err := h.bot.API.Send(msg); err != nil {
			log.Println("Error sending error message:", err)
		}
		return
	}

	// Get daily stats for the last 7 days
	dailyStats, err := h.db.GetDailyStats(h.ctx, 7)
	if err != nil {
		log.Println("Error getting daily stats:", err)
	}

	// Build statistics message
	var statsText strings.Builder
	statsText.WriteString("📊 Bot Statistikasi:\n\n")

	// Overall message stats
	if messageStats != nil {
		statsText.WriteString(fmt.Sprintf("📝 Umumiy xabarlar: %d\n", messageStats["total_messages"]))
		statsText.WriteString(fmt.Sprintf("👥 Noyob foydalanuvchilar: %d\n", messageStats["unique_users"]))
		statsText.WriteString(fmt.Sprintf("🔍 Buyruqlar: %d\n", messageStats["command_count"]))
		statsText.WriteString(fmt.Sprintf("💬 Oddiy xabarlar: %d\n", messageStats["regular_msg_count"]))
		statsText.WriteString(fmt.Sprintf("👤 Shaxsiy xabarlar: %d\n", messageStats["private_messages"]))
		statsText.WriteString(fmt.Sprintf("👪 Guruh xabarlar: %d\n\n", messageStats["group_messages"]))
	}

	// Daily stats
	if len(dailyStats) > 0 {
		statsText.WriteString("📅 So'nggi kunlar statistikasi:\n")
		for _, stat := range dailyStats {
			date, _ := stat["date"].(string)
			messages, _ := stat["total_messages"].(int)
			commands, _ := stat["total_commands"].(int)
			activeUsers, _ := stat["active_users"].(int)
			newUsers, _ := stat["new_users"].(int)

			statsText.WriteString(fmt.Sprintf("%s: 👥 %d faol, 🆕 %d yangi, 📝 %d xabar, 🔍 %d buyruq\n",
				date, activeUsers, newUsers, messages, commands))
		}
		statsText.WriteString("\n")
	}

	// Top users section
	statsText.WriteString("🏆 Eng faol foydalanuvchilar:\n")
	if len(topUsers) == 0 {
		statsText.WriteString("Hozircha foydalanuvchilar yo'q\n")
	} else {
		for i, user := range topUsers {
			username := user.Username
			if username == "" {
				username = user.FirstName
			}
			if username == "" {
				username = "Anonymous"
			}

			statsText.WriteString(fmt.Sprintf("%d. @%s - %d xabar, %d buyruq\n",
				i+1, username, user.MessageCount, user.CommandCount))
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, statsText.String())
	if _, err := h.bot.API.Send(msg); err != nil {
		log.Println("Error sending command response:", err)
	}
}
