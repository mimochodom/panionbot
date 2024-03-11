package commandModule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Aneks struct {
	Items []struct {
		Text string `json:"text"`
	} `json:"items"`
}

func FindAnek(req string, url string) Aneks {
	spl := strings.Replace(req, " ", "~ +", -1)
	spl = "+" + spl + "~"
	req = "select text from telegramAneks where text = '" + spl + "'" +
		"MERGE (select text from anecdoticaAneks where text = '" + spl + "')" +
		"MERGE (select text from anekdotRuAneks where text = '" + spl + "') LIMIT 50 OFFSET 0"
	data := []byte(req)
	r := bytes.NewReader(data)
	res, err := http.Post(url, "application/json", r)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var w Aneks
	err2 := json.Unmarshal(b, &w)
	if err2 != nil {
	}
	return w
}
