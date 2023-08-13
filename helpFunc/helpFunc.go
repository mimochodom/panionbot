package helpFunc

import (
	"golang.org/x/text/encoding/charmap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
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
