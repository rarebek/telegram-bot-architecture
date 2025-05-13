package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"bot/internal/config"
	"bot/internal/db"
	"bot/internal/handlers"
	"bot/internal/telegram"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.toml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database
	database, err := db.New(cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer database.Close(context.Background())

	// Setup webhook URL properly
	webhookURL := cfg.Server.WebhookURL + "/bot/" + cfg.Telegram.BotToken
	log.Printf("Using webhook URL: %s", webhookURL)

	// Initialize bot
	bot, err := telegram.New(&cfg.Telegram, webhookURL)
	if err != nil {
		log.Fatalf("Error initializing bot: %v", err)
	}

	// Set the webhook to the URL specified
	if err := bot.SetupWebhook(); err != nil {
		log.Fatalf("Error setting up webhook: %v", err)
	}

	// Get webhook info and display
	info, err := bot.API.GetWebhookInfo()
	if err != nil {
		log.Fatalf("Failed to get webhook info: %v", err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	// Initialize handlers
	handler, err := handlers.New(bot, database)
	if err != nil {
		log.Fatalf("Error initializing handlers: %v", err)
	}

	// Setup HTTP server with the webhook handler
	webhookPath := "/bot/" + cfg.Telegram.BotToken
	http.HandleFunc(webhookPath, func(w http.ResponseWriter, r *http.Request) {
		// Use the bot's webhook handler
		bot.GetWebhookHandler(handler.HandleUpdate)(w, r)
	})

	// Start the HTTP server
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	log.Printf("Starting server on %s", addr)

	// Create server
	server := &http.Server{
		Addr: addr,
		// Use default handler
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give 10 seconds for connections to drain
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
