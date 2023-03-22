package horoscopeParse

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
)

func HoroscopeParse(url string) {
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

	// Find the review items
	doc.Find(".content-wrapper").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find(".text-link").Text()
		dec := charmap.Windows1251.NewDecoder()
		// Разные кодировки = разные длины символов.
		newBody := make([]byte, len(title)*2)
		n, _, err := dec.Transform(newBody, []byte(title), false)
		if err != nil {
			panic(err)
		}
		newBody = newBody[:n]
		fmt.Printf(string(newBody))
	})
}
