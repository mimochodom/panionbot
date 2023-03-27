package helpFunc

import (
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
