package main

import (
	"fmt"
	"log"
	"panionbot/commandModule"
	"panionbot/helpFunc"
	"panionbot/keyboard"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botToken := helpFunc.GetTokenFromFile("./token/botToken.txt")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil {
			// Create a new MessageConfig. We don't have text yet,
			// so we leave it empty.
			flag := false
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			if update.Message.IsCommand() { // ignore any non-command Messages
				// Extract the command from the Message.
				switch update.Message.Command() {
				case "start":
					msg.Text = "Я пока ещё жив"
				case "anek":
					msg.Text = commandModule.GetAnek()
				case "horoscope":
					msg.ReplyMarkup = keyboard.Horoscope
				case "weather_report":
					msg.Text = "Напишите город в котором хотите узнать погоду"
					msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
					flag = true
					goto L
				default:
					msg.Text = "Ну извини, не могу"
				}
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}

			}
		L:
			if flag == true {
				fmt.Print(commandModule.GetWeatherByName(update.Message.Text))
				flag = false
			}

		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}
			horoscopeText := strings.ToUpper(update.CallbackQuery.Data) + ": " + commandModule.GetHoroscope(update.CallbackQuery.Data)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, horoscopeText)
			bot.Send(msg)
			del := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
			bot.Send(del)
		}

	}

}
