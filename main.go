package main

import (
	"fmt"
	"log"
	"os"
	"panionbot/horoscopeParse"
)

func main() {
	bytes, err := os.ReadFile("./token/token.txt")
	if err != nil {
		log.Fatal(err)
	}
	botToken := string(bytes[:])
	fmt.Println(botToken)

	horoscopeParse.HoroscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=1")
}
