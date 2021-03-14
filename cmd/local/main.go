package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/mcaubrey/sakigake_rss_local/internal/feed"
	"github.com/mcaubrey/sakigake_rss_local/pkg/deepl"
)

func main() {
	// database.InitDatabase()

	url, err := getURLFromArgs()
	if err != nil {
		log.Fatal("Please specify a URL.")
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var f []*feeds.Item
	doc.Find("doc-content[type=lineNewsPC]").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h2").Text()
		content, _ := s.Find("p").Html()

		title, err = deepl.Translate(title, "ja", "en")
		if err != nil {
			panic(err)
		}
		content, err = deepl.Translate(content, "ja", "en")
		if err != nil {
			panic(err)
		}

		entry := feeds.Item{
			Title:       properTitle(title),
			Link:        &feeds.Link{Href: "https://example.com"},
			Description: content,
			Created:     time.Now(),
		}

		f = append(f, &entry)
	})

	fmt.Println(feed.BuildFeed(f))
}

func getURLFromArgs() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("not enough arguments")
	}
	return os.Args[1], nil
}

func properTitle(input string) string {
	words := strings.Fields(input)
	smallwords := " a an the for nor and but or yet so as at around by after along for from of on to with without in"

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}
