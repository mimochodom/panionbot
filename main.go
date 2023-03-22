package main

import (
	"fmt"
	"log"
	"os"
	"panionbot/commandModule"
)

func main() {
	bytes, err := os.ReadFile("./token/token.txt")
	if err != nil {
		log.Fatal(err)
	}
	botToken := string(bytes[:])
	fmt.Println(botToken)

	fmt.Printf(commandModule.GetHoroscope("водолей"))
	fmt.Println(commandModule.Help())
}
