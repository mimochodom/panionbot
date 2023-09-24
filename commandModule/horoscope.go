package commandModule

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"panionbot/helpFunc"
)

func horoscopeParse(url string) string {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	textHoroscope := doc.Find(".text-link").Contents().Text()
	horoscope := helpFunc.DecoderToWin1251(textHoroscope)

	return horoscope
}

func GetHoroscope(znak string) string {
	switch znak {
	case "овен":
		return horoscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=1")

	case "телец":
		return horoscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=2")

	case "близнецы":
		return horoscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=3")

	case "рак":
		return horoscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=4")

	case "лев":
		return horoscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=5")

	case "дева":
		return horoscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=6")

	case "весы":
		return horoscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=7")

	case "скорпион":
		return horoscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=8")

	case "стрелец":
		return horoscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=9")

	case "козерог":
		return horoscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=10")

	case "водолей":
		return horoscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=11")

	case "рыбы":
		return horoscopeParse("https://www.predskazanie.ru/daily_horoscope/?day=&s=12")
	default:
		fmt.Println("Не повезло")
	}
	return "DEPREACATED MY FRIEND"
}
