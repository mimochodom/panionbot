package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"panionbot/commandModule"
	"panionbot/helpFunc"
	"panionbot/keyboard"
	"panionbot/models"
	"strconv"
	"strings"
	"time"
)

var joke []string
var workerPool = make(chan struct{}, 250000)

func main() {

	luceneHost := helpFunc.GetTextFromFile("./token/lucene.txt")
	anek := helpFunc.GetTextFromFile("./token/joke.json")
	db, err := helpFunc.SetupDatabase()
	_ = json.Unmarshal([]byte(anek), &joke)
	lenArr := len(joke)
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
			processUpdate(bot, db, update, luceneHost, joke, lenArr)
		}(update)
	}
}

func processUpdate(bot *tgbotapi.BotAPI, db *gorm.DB, update tgbotapi.Update, luceneHost string, joke []string, lenArr int) {
	switch {
	case update.InlineQuery != nil:
		handleInlineQuery(bot, update.InlineQuery, luceneHost)
	case update.Message != nil:
		handleMessage(bot, db, update.Message, joke, lenArr)
	case update.CallbackQuery != nil:
		handleCallbackQuery(bot, update.CallbackQuery)

	}
}

func handleInlineQuery(bot *tgbotapi.BotAPI, inlineQuery *tgbotapi.InlineQuery, luceneHost string) {

	anekdoty := commandModule.FindAnek(inlineQuery.Query, luceneHost)

	var articles []interface{}
	for _, anek := range anekdoty {
		article := tgbotapi.NewInlineQueryResultArticle(string(rune(rand.Intn(100000))), " ", anek)
		article.Description = anek

		articles = append(articles, article)
	}
	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       articles,
	}
	if _, err := bot.Request(inlineConf); err != nil {
		log.Println(err)
	}

}

func handleMessage(bot *tgbotapi.BotAPI, db *gorm.DB, message *tgbotapi.Message, joke []string, lenArr int) {
	// Extracting relevant information from the update
	user := models.Users{}
	group := models.Groups{}
	userGroup := models.UsersGroups{}

	userID := message.From.ID
	userName := message.From.UserName
	groupName := message.Chat.Title

	chatID := message.Chat.ID
	user.UserID = userID
	user.UserName = userName
	group.GroupName = groupName
	group.GroupID = chatID
	fmt.Println(message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)

	if db.First(&user, "user_id = ?", userID).RowsAffected > 0 {
		if user.UserName != userName {
			db.Model(&user).Update("user_name", userName)
		}
	}

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			msg.Text = "Я пока ещё жив"
		case "anek":
			msg.Text = joke[rand.Intn(lenArr)-1]
		case "horoscope":
			msg.ReplyMarkup = keyboard.Horoscope

		case "weather_report":

			if message.Chat.Type == "private" {
				msg.ReplyMarkup = keyboard.Weather
				msg.Text = "Взгляните на клавиатуру"

			} else {
				msg.Text = "Данная команда не работает в группах"
			}

		case "reg":
			if helpFunc.IsGroupChat(message.Chat.Type) {
				//The time when it all started
				//timeStart := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

				// Checking if the user is already registered
				if db.First(&models.UsersGroups{}, "user_id = ? and group_id = ?", userID, chatID).RowsAffected > 0 {
					msg.Text = "Вы уже участвуете"
					break
				}

				// Registering the user and group
				db.FirstOrCreate(&user)
				db.FirstOrCreate(&models.Groups{GroupID: chatID, GroupName: groupName})
				db.FirstOrCreate(&userGroup, &models.UsersGroups{UserID: userID, GroupID: chatID})
				msg.Text = "Вы успешно зарегистрировались в этой замечательной онлайн-игре \"Зайки-Томатики\""
			} else {
				msg.Text = "Данная команда работает только в группах"
			}

		case "bunny_tomato":
			if helpFunc.IsGroupChat(message.Chat.Type) {
				// Selecting random users for the game
				md := tgbotapi.NewDiceWithEmoji(chatID, "🏀")

				var users []models.Users

				db.Joins("JOIN users_groups on users_groups.user_id = users.user_id").
					Joins("JOIN groups on groups.group_id = users_groups.group_id").
					Where("groups.group_id = ?", chatID).
					Find(&users)

				today := time.Now().Truncate(24 * time.Hour)

				db.Table("groups").Select("group_id, last_game_played").First(&group)

				if group.LastGamePlayed.Before(today) {
					sleep := 500 * time.Millisecond
					bunny, bunny_id := helpFunc.SelectRandomBunnyTomatoPerson(users)
					tomato, tomato_id := helpFunc.SelectRandomBunnyTomatoPerson(users)
					timeNow := time.Now()

					db.Save(&models.Groups{GroupID: chatID, GroupName: groupName, LastGamePlayed: timeNow})
					db.Create(&models.GroupsBTGameResult{GamePlayed: timeNow, GroupID: chatID, UserIDBunny: bunny_id, UserIDTomato: tomato_id})

					db.Model(&models.UsersGroups{}).Where("user_id = ? AND group_id = ?", bunny_id, chatID).UpdateColumn("bunny_count", gorm.Expr("bunny_count+?", 1))
					db.Model(&models.UsersGroups{}).Where("user_id = ? AND group_id = ?", tomato_id, chatID).UpdateColumn("tomato_count", gorm.Expr("tomato_count+?", 1))

					if bunny == tomato {
						bot.Send(md)
						time.Sleep(sleep * 10)
						msg.Text = "ПУ-ПУ-ПУ"
						bot.Send(msg)
						time.Sleep(sleep)
						msg.Text = "Повезло тебе, ты сегодня никакой: " + bunny

					} else {
						bot.Send(md)
						time.Sleep(sleep * 10)
						msg.Text = "ПУ-ПУ-ПУ"
						bot.Send(msg)
						time.Sleep(sleep)
						msg.Text = "🐰 дня: " + bunny + " \n" + "🍅 дня: " + tomato

					}

					for i := range users {
						if users[i].UserName == bunny {
							users[i].BunnyCountGlobal++

						}
						if users[i].UserName == tomato {
							users[i].TomatoCountGlobal++
						}
					}

					db.Save(&users)

				} else {
					lastGameResult := models.GroupsBTGameResult{}
					userBunny := models.Users{}
					userTomato := models.Users{}
					db.Table("groups_bt_game_results").Select("user_id_bunny, user_id_tomato").Where("group_id = ?", chatID).Order("id desc").First(&lastGameResult)
					db.Table("users").Select("user_name").Where("user_id = ?", lastGameResult.UserIDBunny).First(&userBunny)
					db.Table("users").Select("user_name").Where("user_id = ?", lastGameResult.UserIDTomato).First(&userTomato)

					if lastGameResult.UserIDBunny == lastGameResult.UserIDTomato {
						msg.Text = "Уже определили \n" + "Счастливчик, выбил две позиции 🐰🍅: " + userBunny.UserName
					} else {
						msg.Text = "Уже определили \n" + "🐰 дня: " + userBunny.UserName + " \n" + "🍅 дня: " + userTomato.UserName
					}
				}
			} else {
				msg.Text = "Данная команда работает только в группах"
			}
		case "group_stat":

			if helpFunc.IsGroupChat(message.Chat.Type) {
				// Getting the statistics for all users in the group
				var users []models.UsersGroups
				var usersR models.Users
				var output []string
				db.Table("users_groups").Find(&users, "group_id =?", chatID)

				//db.Table("users_groups").Select("bunny_count, tomato_count").First(&userGroup, userID, chatID)
				db.Table("users_groups").Select("user_id, bunny_count, tomato_count").Order("bunny_count + tomato_count desc").Limit(5).Find(&users, "group_id = ?", chatID)

				for _, user := range users {
					db.Table("users").Select("user_name").First(&usersR, user.UserID)
					info := "Имя пользователя: " + usersR.UserName + "\n" +
						"🐰: " + strconv.Itoa(user.BunnyCount) + " раз(а)\n" +
						"🍅: " + strconv.Itoa(user.TomatoCount) + " раз(а)\n" +
						"---------------------------\n"
					output = append(output, info)
				}
				sentence := strings.Join(output, "")
				msg.Text = "Топ 5: \n" + sentence + "Из суммарно: " + strconv.Itoa(len(users)) + " человек(а)"
			} else {
				msg.Text = "Данная команда работает только в группах"
			}
		case "my_stat":
			if helpFunc.IsGroupChat(message.Chat.Type) {
				if db.Table("users").Select("user_name, bunny_count_global, tomato_count_global").First(&user, userID).RowsAffected > 0 {
					db.Table("users_groups").Select("bunny_count, tomato_count").First(&userGroup, userID, chatID)

					msg.Text = "Вот такая у тебя статистика " + user.UserName + " :\n" +
						"В этой группе\n" +
						"- Ты был \"🐰\" " + strconv.Itoa(userGroup.BunnyCount) + " раз(а)\n" +
						"- и \"🍅\" " + strconv.Itoa(userGroup.TomatoCount) + " раз(а).\n" +
						"А в общей статистике\n" +
						"- Ты был \"🐰\" " + strconv.Itoa(user.BunnyCountGlobal) + " раз(а)\n" +
						"- и \"🍅\" " + strconv.Itoa(user.TomatoCountGlobal) + " раз(а)."

				} else {
					msg.Text = "Вы не зарегистрировались"
				}

			} else {
				msg.Text = "Данная команда работает только в группах"
			}
		case "bot_time":
			msg.Text = time.Now().String()
		default:
			imgPath := "./token/What.png"
			helpFunc.SendImage(bot, chatID, imgPath)
			msg.Text = ""

		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered from panic:", r)
				// Выполните здесь необходимые действия после паники
			}
		}()
	}
	if message.Text == "По названию" {
		msg.Text = "Напишите город в котором хотите узнать погоду"
		msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}

	if message.ReplyToMessage != nil && message.Chat.Type == "private" {
		msg.Text = commandModule.GetWeatherByName(message.Text)

		if message.Location != nil {
			msg.Text = commandModule.GetWeatherByLocation(message.Location.Latitude, message.Location.Longitude)
		}
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	callback := tgbotapi.NewCallback(callbackQuery.ID, callbackQuery.Data)

	if _, err := bot.Request(callback); err != nil {
		panic(err)
	}

	horoscopeText := strings.ToUpper(callbackQuery.Data) + ": " + commandModule.GetHoroscope(callbackQuery.Data)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, horoscopeText)
	_, err := bot.Send(msg)
	if err != nil {
		return
	}
	del := tgbotapi.NewDeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
	_, err = bot.Send(del)
	if err != nil {
		return
	}
}
