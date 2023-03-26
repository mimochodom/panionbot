package helpFunc

import (
	"log"
	"os"
)

func GetTokenFromFile(pathFile string) string {
	bytes, err := os.ReadFile(pathFile)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes[:])
}
