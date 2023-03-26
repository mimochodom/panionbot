package main

import (
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

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			if update.Message.Command() == "weather_report" && update.Message.Chat.Type == "private" {
				msg.ReplyMarkup = keyboard.Weather

				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}

			if update.Message.Text == "По названию" {
				msg.Text = "Напишите город в котором хотите узнать погоду"
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
			}

			if update.Message.ReplyToMessage != nil {
				msg.Text = commandModule.GetWeatherByName(update.Message.Text)

				if update.Message.Location != nil {
					msg.Text = commandModule.GetWeatherByLocation(update.Message.Location.Latitude, update.Message.Location.Longitude)
				}
			}

			if update.Message.IsCommand() {

				switch update.Message.Command() {
				case "start":
					msg.Text = "Я пока ещё жив"
				case "anek":
					msg.Text = commandModule.GetAnek()
				case "horoscope":
					msg.ReplyMarkup = keyboard.Horoscope
				}
			}

			if update.Message.Command() != "weather_report" {
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}
		} else if update.CallbackQuery != nil {
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
