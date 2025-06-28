package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type PageChunk struct {
	URL   string
	Title string
	Body  string
}

func main() {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
	)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	c.OnHTML("html", handleHtml)
	c.Visit("https://web.cs.toronto.edu/")

	c.Wait()

	fmt.Println("Crawling completed.")

}

func handleHtml(e *colly.HTMLElement) {
	content := e.Text
	fmt.Println("Page Text: ", content)
}
