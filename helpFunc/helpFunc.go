package helpFunc

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
	"unicode"
)

func GetTextFromFile(pathFile string) string {
	file, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}

	butt, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file) // Закрытие файла при выходе из функции

	return string(butt)
}

func DecoderToWin1251(title string) string {
	dec := charmap.Windows1251.NewDecoder()

	// Преобразуем строку в байтовый поток в кодировке Windows-1251.
	input := []byte(title)
	windows1251Reader := transform.NewReader(bytes.NewReader(input), dec)

	// Чтобы нормализовать UTF-8, прочтем байты из windows1251Reader и пропустим их через нормализатор.
	utf8Buffer := &bytes.Buffer{}
	utf8Writer := transform.NewWriter(utf8Buffer, transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC))

	if _, err := io.Copy(utf8Writer, windows1251Reader); err != nil {
		log.Fatal(err)
	}
	utf8Writer.Close()

	return utf8Buffer.String()
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func SetupDatabase() (*gorm.DB, error) {
	dsn := GetTextFromFile("./token/dbConfig.txt")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func IsGroupChat(chatType string) bool {
	return chatType == "group" || chatType == "supergroup"
}

func SendImage(bot *tgbotapi.BotAPI, chatID int64, imagePath string, caption string) {

	image := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(imagePath))
	image.Caption = caption
	if _, err := bot.Send(image); err != nil {
		log.Panic(err)
	}
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func SelectRandomItem[T any](items []T) T {

	randomIndex := random.Intn(len(items))
	return items[randomIndex]
}

func GenerateUniqueID(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending message:", err)
	}
}
