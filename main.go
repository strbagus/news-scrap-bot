package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	// "time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}

func main() {
	token := os.Getenv("BOT_TOKEN")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}
	/* go func() {
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()

		// Run getData immediately at startup
		// getData()

		// Periodically run getData every 15 minutes
		for {
			select {
			case <-ticker.C:
				// getData()
			case <-ctx.Done():
				fmt.Println("Shutting down periodic task...")
				return
			}
		}
	}() */

	b.Start(ctx)
}
func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("LOG: ")
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
