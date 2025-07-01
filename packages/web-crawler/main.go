package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type PageChunk struct {
	URL   string
	Title string
	Body  string
}

func main() {
	links := 0
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)
	c.AllowURLRevisit = false
	c.OnResponse(func(r *colly.Response) {
		if isSyllabus(r) {
			fmt.Println("Found syllabus at:", r.Request.URL.String())
		}
	})
	// c.OnHTML("table", handleTable)
	c.OnHTML("body", func(e *colly.HTMLElement) {
		chunks := getChunks(e)
		fmt.Println(chunks)
	})
	c.OnHTML("a[href]", handleLink)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
		links++
	})
	c.Visit("https://web.cs.toronto.edu/people/faculty-directory")
	c.Wait()

	fmt.Println("Crawling completed.")
	fmt.Printf("Total links visited: %d\n", links)

}

func isSyllabus(r *colly.Response) bool {
	content := string(r.Body)
	if strings.Contains(content, "syllabus") || strings.Contains(content, "course outline") {
		return true
	}
	return false
}

func handleLink(e *colly.HTMLElement) {
	href := e.Attr("href")
	absoluteURL := e.Request.AbsoluteURL(href)
	e.Request.Visit(absoluteURL)
}

func getChunks(e *colly.HTMLElement) []PageChunk {
	chunks := []PageChunk{}
	currChunk := PageChunk{
		URL:   e.Request.URL.String(),
		Title: e.ChildText("title"),
	}
	currChunkContent := ""
	e.ForEach("p, h1, h2, h3, h4, h5, h6", func(_ int, el *colly.HTMLElement) {
		text := strings.TrimSpace(el.Text)
		if text != "" {
			if len(currChunkContent)+len(text) > 1000 { // Limit chunk size
				chunks = append(chunks, currChunk)
				currChunkContent = ""
				currChunk = PageChunk{
					URL:   e.Request.URL.String(),
					Title: e.ChildText("title"),
				}
			}
			currChunkContent += text + " "
		}
	})
	if currChunkContent != "" {
		currChunk.Body = strings.TrimSpace(currChunkContent)
		chunks = append(chunks, currChunk)
	}
	return chunks
}
