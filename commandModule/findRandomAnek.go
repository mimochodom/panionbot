package commandModule

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type TotalCountAneks struct {
	QueryTotalItems int `json:"query_total_items"`
}

func FindRandomAnek(src int, url string) string {
	var srcTable string
	if src == 0 {
		srcTable = "telegramAneks"
	}
	if src == 1 {
		srcTable = "anecdoticaAneks"
	}
	if src == 2 {
		srcTable = "anekdotRuAneks"
	}

	req := "select count(*) from " + srcTable
	data := []byte(req)
	r := bytes.NewReader(data)
	res, err := http.Post(url, "application/json", r)
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var w TotalCountAneks
	err2 := json.Unmarshal(b, &w)
	if err2 != nil {
	}
	rnd := rand.Intn(w.QueryTotalItems)

	req = "select text from " + srcTable + " where id=" + strconv.Itoa(rnd)
	data = []byte(req)
	r = bytes.NewReader(data)
	res, err = http.Post(url, "application/json", r)
	b, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var ww Aneks
	err2 = json.Unmarshal(b, &ww)
	if err2 != nil {
	}

	return ww.Items[0].Text
}
