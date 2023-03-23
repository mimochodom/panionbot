package main

import (
	"fmt"
	"panionbot/commandModule"
)

func main() {
	//botToken := helpFunc.GetTokenFromFile("./token/botToken.txt")
	//fmt.Println(botToken)

	fmt.Print(commandModule.GetWeatherByName("Москва"))
}
