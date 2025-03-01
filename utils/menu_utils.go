package utils

import (
	"context"
	"fmt"
	"gobot/models"
	"log"

	"github.com/go-telegram/bot"
	tbmodels "github.com/go-telegram/bot/models"
)

func MenuSubscribe(ctx context.Context, b *bot.Bot, update *tbmodels.Update) {
	users, _ := ReadFile[models.User]("db/users.json")
	chat := update.Message.Chat
	var isExists bool
	for _, user := range users {
		if user.ChatID == chat.ID {
			isExists = true
			break
		}
	}
	var message string
	if isExists {
		message = "You already subscribed."
	} else {
		message = "You have subscribed!"
		newUser := models.User{
			ChatID:   chat.ID,
			Username: chat.Username,
		}
		users = append(users, newUser)
		WriteFile("db/users.json", users)
	}
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat.ID,
		Text:   message,
	})
	if err != nil {
		log.Printf("Error sending subscribe message: %v", err)
	}
}

func MenuUnsubscribe(ctx context.Context, b *bot.Bot, update *tbmodels.Update) {
	users, _ := ReadFile[models.User]("db/users.json")
	chat := update.Message.Chat
	var isExists bool
	for _, user := range users {
		if user.ChatID == chat.ID {
			isExists = true
			break
		}
	}
	var message string
	if !isExists {
		message = "You not subscribed yet."
	} else {
		message = "You have unsubscribed!"
		var filtered []models.User
		for _, user := range users {
			if user.ChatID != chat.ID {
				filtered = append(filtered, user)
			}
		}

		var result []models.User
		if filtered == nil {
			result = []models.User{}
		} else {
			result = filtered
		}
		WriteFile("db/users.json", result)
	}
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat.ID,
		Text:   message,
	})
	if err != nil {
		log.Printf("Error sending unsubscribe message: %v", err)
	}
}

func MenuLastInfo(ctx context.Context, b *bot.Bot, update *tbmodels.Update) {
	news, _ := ReadFile[models.NewsType]("db/news.json")
	message := fmt.Sprintf("%s\n%s", news[0].Title, news[0].Link)
	chatID := update.Message.Chat.ID
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   message,
	})
	if err != nil {
		log.Printf("Error sending last info message: %v", err)
	}
}
