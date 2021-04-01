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

	// Write to use POST
	// https://stackoverflow.com/questions/24455147/how-do-i-send-a-json-string-in-a-post-request-in-go
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
		opt, _ := s.Find(".mdCMN03Share a").Attr("data-opt")

		optArr := strings.SplitAfter(opt, "'")

		spos := len(optArr[3]) - 1
		share := optArr[3][:spos]

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
			Link:        &feeds.Link{Href: share},
			Description: content,
			Created:     time.Now(),
		}

		f = append(f, &entry)
	})

	writeToFile(fmt.Sprintln(feed.BuildFeed(f)), "feed.rss")

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

func writeToFile(content string, filename string) {

	f, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(content)

	if err2 != nil {
		log.Fatal(err2)
	}
}
