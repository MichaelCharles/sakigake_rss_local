package feed

import (
	"log"
	"time"

	"github.com/gorilla/feeds"
)

func BuildFeed(items []*feeds.Item) (string, error) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "Sakigake English",
		Link:        &feeds.Link{Href: "https://www.sakigake.jp/"},
		Description: "Local Akita News in English by Sakigake",
		Author:      &feeds.Author{Name: "Michael Charles Aubrey", Email: "aubrey@michaelcharl.es"},
		Created:     now,
	}

	feed.Items = items

	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	return rss, err
}
