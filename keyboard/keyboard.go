package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Horoscope = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("♈овен", "овен"),
		tgbotapi.NewInlineKeyboardButtonData("♉телец", "телец"),
		tgbotapi.NewInlineKeyboardButtonData("♊близнецы", "близнецы"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("♋рак", "рак"),
		tgbotapi.NewInlineKeyboardButtonData("♌лев", "лев"),
		tgbotapi.NewInlineKeyboardButtonData("♍дева", "дева"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("♎весы", "весы"),
		tgbotapi.NewInlineKeyboardButtonData("♏скорпион", "скорпион"),
		tgbotapi.NewInlineKeyboardButtonData("♐стрелец", "стрелец"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("♑козерог", "козерог"),
		tgbotapi.NewInlineKeyboardButtonData("♒воделей", "водолей"),
		tgbotapi.NewInlineKeyboardButtonData("♓рыбы", "рыбы"),
	),
)
