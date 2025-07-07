package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type PageChunk struct {
	URL     string
	Title   string
	Content string
}

var MaxChunkLength = 1000
var OverlapLength = 100

var (
	fullDocumentText string // Accumulate all text here
	allTextChunks    []PageChunk
	isSyllabusPage   bool
)

func main() {
	links := 0
	c := colly.NewCollector(
		colly.MaxDepth(4),
		colly.Async(true),
	)
	c.AllowURLRevisit = false
	c.OnHTML("body", func(e *colly.HTMLElement) {
		response := e.Response
		fullDocumentText = ""                 // Reset for each page
		isSyllabusPage = isSyllabus(response) // Set flag
		if isSyllabusPage {
			fmt.Println("Syllabus page detected:", response.Request.URL.String())
			traverseDOMForFullText(e.DOM) // Extract all text from the page
		}
	})

	c.OnScraped(func(r *colly.Response) {
		if isSyllabusPage && fullDocumentText != "" {
			chunkTextByLength(fullDocumentText, r.Request.URL.String()) // Pass URL for tracking
			fmt.Printf("Chunked syllabus: %s - created %d chunks\n",
				r.Request.URL.String(), len(allTextChunks))
		}
	})

	c.OnHTML("a[href]", handleLink)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
		links++
	})
	c.Visit("https://web.cs.toronto.edu/people/faculty-directory")
	c.Wait()
	fmt.Printf("Total links visited: %d\n", links)
	fmt.Printf("Total chunks created: %d\n", len(allTextChunks))
	fmt.Println("Crawling completed.")
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("crawled_chunks_%s.json", timestamp)

	if err := saveChunksToJSON(allTextChunks, filename); err != nil {
		fmt.Printf("Error saving chunks: %v\n", err)
	}

}

func isSyllabus(r *colly.Response) bool {
	content := string(r.Body)
	if strings.Contains(content, "syllabus") || strings.Contains(content, "course outline") || strings.Contains(content, "course description") {
		return true
	}
	return false
}

func handleLink(e *colly.HTMLElement) {
	href := e.Attr("href")
	absoluteURL := e.Request.AbsoluteURL(href)
	e.Request.Visit(absoluteURL)
}

func chunkTextByLength(text string, url string) {
	if text == "" {
		return
	}

	for i := 0; i < len(text); {
		end := i + MaxChunkLength
		if end > len(text) {
			end = len(text)
		}

		chunkContent := text[i:end]
		allTextChunks = append(allTextChunks, PageChunk{Content: chunkContent, URL: url})

		// Move the starting point for the next chunk
		if end == len(text) {
			break // Reached the end of the text
		}

		// Implement overlap
		nextStart := end - OverlapLength
		if nextStart < i { // Ensure nextStart doesn't go backwards past the current i
			nextStart = i + 1 // At least move by one character
		}
		i = nextStart
	}
}

// Get all text on the site recursively
func traverseDOMForFullText(s *goquery.Selection) {
	s.Contents().Each(func(_ int, node *goquery.Selection) {
		nodeName := goquery.NodeName(node)

		// Skip script, style, noscript, and comment nodes
		if nodeName == "script" || nodeName == "style" || nodeName == "noscript" || nodeName == "#comment" {
			return
		}

		if nodeName == "#text" {
			// It's a text node directly inside the current element
			text := node.Text()
			text = strings.TrimSpace(text)
			if text != "" {
				if fullDocumentText != "" {
					fullDocumentText += " "
				}
				fullDocumentText += text
			}
		} else {
			// It's an element node. Recurse into it.
			traverseDOMForFullText(node)
		}
	})
}

func saveChunksToJSON(chunks []PageChunk, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON

	if err := encoder.Encode(chunks); err != nil {
		return fmt.Errorf("failed to encode JSON: %v", err)
	}

	fmt.Printf("Saved %d chunks to %s\n", len(chunks), filename)
	return nil
}
