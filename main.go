package main

import (
	fetcher "Scrapper/Internal/Fetcher"
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	ctx := context.Background()
	f := fetcher.New("LearningScrapper/0.1", 1)

	urls := []string{

		"https://books.toscrape.com/catalogue/page-2.html",
		"https://books.toscrape.com/catalogue/page-3.html",
		"https://books.toscrape.com/catalogue/page-4.html",
	}

	for _, u := range urls {
		start := time.Now()
		body, err := f.Fetch(ctx, u)
		if err != nil {
			log.Printf("failed: %v", err)
			continue
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Printf("parse failed %v", err)
			continue
		}
		doc.Find("article.product_pod").Each(func(i int, s *goquery.Selection) {
			title, _ := s.Find("h3 a").Attr("title")
			price := s.Find(".price_color").Text()
			fmt.Println(title, "-", price)
		})
		log.Printf("done %s in %v", time.Since(start))
	}

}
