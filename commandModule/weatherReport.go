package commandModule

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"panionbot/helpFunc"
)

func GetWeatherByName(cityName string) string {
	API_WEATHER_KEY := helpFunc.GetTokenFromFile("./token/weatherTokenAPI.txt")
	url := "https://api.openweathermap.org/data/2.5/weather?q=" + cityName + "&lang=ru&units=metric&appid=" + API_WEATHER_KEY
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

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var w weather
	err2 := json.Unmarshal(b, &w)
	if err2 != nil {
		return "JSON: CANT UNMARSHAL"
	}
	temp := fmt.Sprintf("%.2f", w.Main.Temp)
	feelsLike := fmt.Sprintf("%.2f", w.Main.FeelsLike)
	windSpeed := fmt.Sprintf("%.2f", w.Wind.Speed)
	prognoz := cityName + ": температура " + temp + "°C " + w.Weather[0].Description + " ветер " + windSpeed + " м/с ошущается как " + feelsLike + "°C"
	return prognoz
}

type weather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}
