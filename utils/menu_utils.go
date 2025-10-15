package utils

import (
	"context"
	"fmt"
	"gobot/db"
	m "gobot/models"
	"log"
	"time"

	"github.com/go-telegram/bot"
	tbmodels "github.com/go-telegram/bot/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func MenuSubscribe(ctx context.Context, b *bot.Bot, update *tbmodels.Update) {
	coll := db.Client.Database("ftibot").Collection("users")
	chat := update.Message.Chat

	user := m.User{
		ChatID:   chat.ID,
		Username: chat.Username,
	}
	_, err := coll.InsertOne(context.TODO(), user)
	message := "Subscribed!"
	if err != nil {
		message = "Subsribe failed."
		if mongo.IsDuplicateKeyError(err) {
			message = "You already subscribed."
		}
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat.ID,
		Text:   message,
	})
	if err != nil {
		log.Printf("ERROR: Send Message Failed: %v", err)
	}
}

func MenuUnsubscribe(ctx context.Context, b *bot.Bot, update *tbmodels.Update) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	coll := db.Client.Database("ftibot").Collection("users")
	chat := update.Message.Chat
	filter := bson.M{"chatid": chat.ID}
	message := "Unsubscribe from news."
	res, err := coll.DeleteOne(ctxTimeout, filter)
	if err != nil {
		log.Printf("delete gone wrong: %v", err)
		message = "Unsubscribe failed."
	}
	if res.DeletedCount == 0 {
		message = "You not subscribed yet."
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat.ID,
		Text:   message,
	})
	if err != nil {
		log.Printf("Error sending unsubscribe message: %v", err)
	}
}

func MenuLastInfo(ctx context.Context, b *bot.Bot, update *tbmodels.Update) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	coll := db.Client.Database("ftibot").Collection("news")
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})
	var news m.NewsType 
	err := coll.FindOne(ctxTimeout, bson.D{}, opts).Decode(&news)
	if err != nil {
		log.Printf("Error find last news: %v", err)
	}
	message := fmt.Sprintf("%s\n%s", news.Title, news.Link)
	user := update.Message.Chat

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: user.ID,
		Text:   message,
	})

	if err != nil {
		log.Printf("Error sending last info message: %v", err)
	} else {
		log.Printf("INFO: Sent %v - %v", user.Username, news.Title)
	}
}
