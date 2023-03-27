package commandModule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func FindAnek(req string, url string) []string {

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

	var dataStr []string
	err2 := json.Unmarshal(b, &dataStr)
	if err2 != nil {

	}
	return dataStr
}
