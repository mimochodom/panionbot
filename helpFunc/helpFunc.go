package helpFunc

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/text/encoding/charmap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"os"
	"time"
)

func GetTextFromFile(pathFile string) string {
	bytes, err := os.ReadFile(pathFile)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes[:])
}

func DecoderToWin1251(title string) string {
	dec := charmap.Windows1251.NewDecoder()
	// Разные кодировки = разные длины символов.
	newBody := make([]byte, len(title)*2)
	n, _, err := dec.Transform(newBody, []byte(title), false)
	if err != nil {
		panic(err)
	}
	newBody = newBody[:n]
	return string(newBody)
}

func SetupDatabase() (*gorm.DB, error) {
	dsn := GetTextFromFile("./token/dbConfig.txt")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func IsGroupChat(chatType string) bool {
	return chatType == "group" || chatType == "supergroup"
}

func SendImage(bot *tgbotapi.BotAPI, chatID int64, imagePath string) {

	image := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(imagePath))
	if _, err := bot.Send(image); err != nil {
		log.Panic(err)
	}
}

func SelectRandomItem[T any](items []T) T {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rand.Intn(len(items))
	return items[randomIndex]
}

//func SelectRandomBunnyTomatoPerson(users []models.Users) (string, int64) {
//	rand.New(rand.NewSource(time.Now().UnixNano()))
//	randomIndex := rand.Intn(len(users))
//	return users[randomIndex].UserName, users[randomIndex].UserID
//}
