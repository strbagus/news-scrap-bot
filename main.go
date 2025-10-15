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
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	dbc := db.Client.Database("ftibot")
	cNews := dbc.Collection("news")
	cUsers := dbc.Collection("users")

	news := utils.GetData()
	docs := make([]any, len(news))
	for i, n := range news {
		docs[i] = n
	}
	opts := options.InsertMany().SetOrdered(false)
	res, err := cNews.InsertMany(context.Background(), docs, opts)

	if err != nil {
		if bwe, ok := err.(mongo.BulkWriteException); ok {
			for _, e := range bwe.WriteErrors {
				if e.Code != 11000 {
					fmt.Printf("write error: %+v\n", e)
				}
			}
		} else {
			log.Fatalf("InsertMany failed: %v", err)
		}
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	usersCursor, err := cUsers.Find(ctxTimeout, bson.D{})
	if err != nil {
		log.Printf("ERROR: Get list users: %v", err)
		return
	}
	defer usersCursor.Close(ctxTimeout)

	for usersCursor.Next(ctxTimeout) {
		var user models.User

		if err := usersCursor.Decode(&user); err != nil {
			log.Printf("ERROR: Decode user: %v", err)
			continue
		}

		for i := range res.InsertedIDs {
			news := news[i]

			message := fmt.Sprintf("%s\n%s", news.Title, news.Link)
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: user.ChatID,
				Text:   message,
			})
			if err != nil {
				log.Printf("Error sending news: %v", err)
			} else {
				log.Printf("INFO: Sent %v - %v", user.Username, news.Title)
			}
		}

	}

	if err := usersCursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
	}
}

func handler(ctx context.Context, b *bot.Bot, update *tbmodels.Update) {
	// fmt.Println("LOG: Received update")
}
