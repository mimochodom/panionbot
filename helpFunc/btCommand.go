package helpFunc

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"panionbot/models"
	"strconv"
	"strings"
	"time"
)

func HandleCommandReg(db *gorm.DB, user models.Users, userID int64, chatID int64, groupName string) string {
	var userGroup models.UsersGroups

	if db.First(&models.UsersGroups{}, "user_id = ? and group_id = ?", userID, chatID).RowsAffected > 0 {
		return "Вы уже участвуете"
	}

	db.FirstOrCreate(&user)

	db.FirstOrCreate(&models.Groups{GroupID: chatID, GroupName: groupName})

	db.FirstOrCreate(&userGroup, &models.UsersGroups{UserID: userID, GroupID: chatID})

	return "Вы успешно зарегистрировались в этой замечательной онлайн-игре \"Зайки-Томатики\""
}

func HandleCommandBunnyTomato(bot *tgbotapi.BotAPI, db *gorm.DB, group models.Groups, chatID int64, groupName string) string {
	var users []models.Users

	if db.First(&group, "group_id = ?", chatID).RowsAffected == 0 {
		return "Сначала нажмите /reg"
	}

	db.Joins("JOIN users_groups on users_groups.user_id = users.user_id").
		Joins("JOIN groups on groups.group_id = users_groups.group_id").
		Where("groups.group_id = ?", chatID).
		Find(&users)

	today := time.Now().Truncate(24 * time.Hour)

	db.Table("groups").Select("group_id, last_game_played").First(&group)
	randomEmoji := SelectRandomItem(models.SmileyList)
	// Selecting random users for the game

	md := tgbotapi.NewDiceWithEmoji(chatID, randomEmoji)
	if group.LastGamePlayed.Before(today) {
		sleep := 500 * time.Millisecond
		bunny := SelectRandomItem(users)
		tomato := SelectRandomItem(users)

		timeNow := time.Now()

		db.Save(&models.Groups{GroupID: chatID, GroupName: groupName, LastGamePlayed: timeNow})
		db.Create(&models.GroupsBTGameResult{GamePlayed: timeNow, GroupID: chatID, UserIDBunny: bunny.UserID, UserIDTomato: tomato.UserID})

		db.Model(&models.UsersGroups{}).Where("user_id = ? AND group_id = ?", bunny.UserID, chatID).UpdateColumn("bunny_count", gorm.Expr("bunny_count+?", 1))
		db.Model(&models.UsersGroups{}).Where("user_id = ? AND group_id = ?", tomato.UserID, chatID).UpdateColumn("tomato_count", gorm.Expr("tomato_count+?", 1))

		if bunny.UserID == tomato.UserID {
			bot.Send(md)
			time.Sleep(sleep * 10)
			SendMessage(bot, chatID, "ПУ-ПУ-ПУ")
			time.Sleep(sleep)
			for i := range users {
				if users[i].UserID == bunny.UserID {
					users[i].BunnyCountGlobal++

				}
				if users[i].UserID == tomato.UserID {
					users[i].TomatoCountGlobal++
				}
			}

			db.Save(&users)
			return "На трон села, на трон села, царь во дворца: " + bunny.UserName

		} else {
			bot.Send(md)
			time.Sleep(sleep * 10)
			SendMessage(bot, chatID, "ПУ-ПУ-ПУ")
			time.Sleep(sleep)
			for i := range users {
				if users[i].UserID == bunny.UserID {
					users[i].BunnyCountGlobal++

				}
				if users[i].UserID == tomato.UserID {
					users[i].TomatoCountGlobal++
				}
			}

			db.Save(&users)
			return "🐰 дня: " + bunny.UserName + " \n" + "🍅 дня: " + tomato.UserName

		}

	} else {
		lastGameResult := models.GroupsBTGameResult{}
		userBunny := models.Users{}
		userTomato := models.Users{}
		db.Table("groups_bt_game_results").Select("user_id_bunny, user_id_tomato").Where("group_id = ?", chatID).Order("id desc").First(&lastGameResult)
		db.Table("users").Select("user_name").Where("user_id = ?", lastGameResult.UserIDBunny).First(&userBunny)
		db.Table("users").Select("user_name").Where("user_id = ?", lastGameResult.UserIDTomato).First(&userTomato)

		if lastGameResult.UserIDBunny == lastGameResult.UserIDTomato {
			return "Уже определили \n" + "Счастливчик, выбил две позиции 🐰🍅: " + userBunny.UserName
		} else {
			return "Уже определили \n" + "🐰 дня: " + userBunny.UserName + " \n" + "🍅 дня: " + userTomato.UserName
		}
	}

	return "Результат игры"
}

func HandleCommandGroupStat(db *gorm.DB, chatID int64) string {
	var users []models.UsersGroups
	var output []string
	var usersR models.Users

	db.Table("users_groups").Find(&users, "group_id =?", chatID)
	db.Table("users_groups").Select("user_id, bunny_count, tomato_count").Order("bunny_count + tomato_count desc").Limit(5).Find(&users, "group_id = ?", chatID)

	for _, user := range users {
		// Получение имени пользователя по ID и другая обработка
		db.Table("users").Select("user_name").First(&usersR, user.UserID)
		info := "Имя пользователя: " + usersR.UserName + "\n" +
			"🐰: " + strconv.Itoa(user.BunnyCount) + " раз(а)\n" +
			"🍅: " + strconv.Itoa(user.TomatoCount) + " раз(а)\n" +
			"---------------------------\n"
		output = append(output, info)
	}

	sentence := strings.Join(output, "")
	return "Топ 5: \n" + sentence + "Из суммарно: " + strconv.Itoa(len(users)) + " человек(а)"
}

func HandleCommandMyStat(db *gorm.DB, userID int, chatID int64) string {
	var user models.Users
	var userGroup models.UsersGroups

	db.Table("users").Select("user_name, bunny_count_global, tomato_count_global").First(&user, userID)

	db.Table("users_groups").Select("bunny_count, tomato_count").First(&userGroup, "user_id = ? AND group_id = ?", userID, chatID)

	return fmt.Sprintf("Вот такая у тебя статистика %s :\n"+
		"В этой группе\n"+
		"- Ты был \"🐰\" %d раз(а)\n"+
		"- и \"🍅\" %d раз(а).\n"+
		"А в общей статистике\n"+
		"- Ты был \"🐰\" %d раз(а)\n"+
		"- и \"🍅\" %d раз(а).",
		user.UserName, userGroup.BunnyCount, userGroup.TomatoCount, user.BunnyCountGlobal, user.TomatoCountGlobal)
}
