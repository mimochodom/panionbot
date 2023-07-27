package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"panionbot/commandModule"
	"panionbot/helpFunc"
	"panionbot/keyboard"
	"strings"
)

var joke []string
var workerPool = make(chan struct{}, 10)

func main() {
	luceneHost := helpFunc.GetTextFromFile("./token/lucene.txt")
	anek := helpFunc.GetTextFromFile("./token/joke.json")
	_ = json.Unmarshal([]byte(anek), &joke)
	lenArr := len(joke)
	fmt.Println()

	botToken := helpFunc.GetTextFromFile("./token/botToken.txt")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		workerPool <- struct{}{}
		go func(update tgbotapi.Update) {
			defer func() { <-workerPool }()
			processUpdate(bot, update, luceneHost, joke, lenArr)
		}(update)
	}
}

func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, luceneHost string, joke []string, lenArr int) {
	if update.InlineQuery != nil {
		anekdoty := commandModule.FindAnek(update.InlineQuery.Query, luceneHost)

		var articles []interface{}
		for _, anek := range anekdoty {
			article := tgbotapi.NewInlineQueryResultArticle(string(rune(rand.Intn(100000))), " ", anek)
			article.Description = anek

			articles = append(articles, article)
		}
		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    true,
			CacheTime:     0,
			Results:       articles,
		}
		if _, err := bot.Request(inlineConf); err != nil {
			log.Println(err)
		}
	}

	if update.Message != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		if update.Message.IsCommand() {

			switch update.Message.Command() {
			case "start":
				msg.Text = "Я пока ещё жив"
			case "anek":
				msg.Text = joke[rand.Intn(lenArr)-1]
			case "horoscope":
				msg.ReplyMarkup = keyboard.Horoscope

			case "weather_report":

				if update.Message.Chat.Type == "private" {
					msg.ReplyMarkup = keyboard.Weather
					msg.Text = "Взгляните на клавиатуру"

				} else {
					msg.Text = "Данная команда не работает в группах"
				}

			}
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}

		if update.Message.Text == "По названию" {
			msg.Text = "Напишите город в котором хотите узнать погоду"
			msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}

		if update.Message.ReplyToMessage != nil && update.Message.Chat.Type == "private" {
			msg.Text = commandModule.GetWeatherByName(update.Message.Text)

			if update.Message.Location != nil {
				msg.Text = commandModule.GetWeatherByLocation(update.Message.Location.Latitude, update.Message.Location.Longitude)
			}
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
		_, err := bot.Send(msg)
		if err != nil {
			return
		}
		del := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
		_, err = bot.Send(del)
		if err != nil {
			return
		}
	}
}
