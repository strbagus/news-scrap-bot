package main

import (
	"context"
	"fmt"
	"gobot/db"
	"gobot/models"
	"gobot/utils"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
	tbmodels "github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}

func main() {
	db.InitMongo()
	defer db.CloseMongo()
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set in the environment")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		log.Fatalf("Error initializing bot: %v", err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/subscribe", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *tbmodels.Update) {
		utils.MenuSubscribe(ctx, b, update)
	})
	b.RegisterHandler(bot.HandlerTypeMessageText, "/unsubscribe", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *tbmodels.Update) {
		utils.MenuUnsubscribe(ctx, b, update)
	})
	b.RegisterHandler(bot.HandlerTypeMessageText, "/lastinfo", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *tbmodels.Update) {
		utils.MenuLastInfo(ctx, b, update)
	})

	go startPeriodicTask(ctx, b)

	b.Start(ctx)
}

func startPeriodicTask(ctx context.Context, b *bot.Bot) {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	scrap(ctx, b)

	for {
		select {
		case <-ticker.C:
			scrap(ctx, b)
		case <-ctx.Done():
			fmt.Println("Shutting down periodic task...")
			return
		}
	}
}

func scrap(ctx context.Context, b *bot.Bot) {
	list := utils.GetData()
	users, _ := utils.ReadFile[models.User]("db/users.json")
	oldList, _ := utils.ReadFile[models.NewsType]("db/news.json")
	news := utils.CompareData(oldList, list)
	if len(news) > 0 {
		for _, user := range users {
			for _, item := range news {
				message := fmt.Sprintf("%s\n%s", item.Title, item.Link)
				_, err := b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: user.ChatID,
					Text:   message,
				})
				if err != nil {
					log.Printf("Error sending news: %v", err)
				}
			}
		}
        utils.WriteFile("db/news.json", list)
	} else {
		fmt.Println("No new data to send.")
	}
}

func handler(ctx context.Context, b *bot.Bot, update *tbmodels.Update) {
	// fmt.Println("LOG: Received update")
}
