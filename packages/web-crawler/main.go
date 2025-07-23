package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
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
	pageTitle        string // Store the page title
	allTextChunks    []PageChunk
)

func main() {
	links := 0
	c := colly.NewCollector(
		colly.MaxDepth(3),
		colly.Async(true),
		colly.URLFilters(
			regexp.MustCompile(`^https?://([a-z0-9-]+\.)?(cs\.toronto\.edu|cs\.utoronto\.edu)/.*$`),
		))
	c.AllowURLRevisit = false
	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("Found page title:", e.Text)
		pageTitle = strings.TrimSpace(e.Text)
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		fullDocumentText = ""         // Reset for each page
		traverseDOMForFullText(e.DOM) // Extract all text from the page
	})

	c.OnScraped(func(r *colly.Response) {
		if fullDocumentText == "" {
			return // Skip if no content
		}

		oldChunkCount := len(allTextChunks)
		chunkTextByLength(fullDocumentText, r.Request.URL.String(), pageTitle) // Pass URL and title for tracking
		newChunkCount := len(allTextChunks)

		fmt.Printf("Chunked page: %s - created %d new chunks (total: %d)\n",
			r.Request.URL.String(), newChunkCount-oldChunkCount, newChunkCount)

		for i := oldChunkCount; i < newChunkCount; i++ {
			chunk := allTextChunks[i]
			payload := map[string]interface{}{
				"content":         chunk.Content,
				"url":             chunk.URL,
				"title":           chunk.Title,
				"position":        i, // Position in the chunk list
				"collection_name": "website_chunks",
			}

			fmt.Printf("Sending chunk %d: %s\n", i, chunk.Title)

			jsonData, err := json.Marshal(payload)
			if err != nil {
				fmt.Printf("Error marshaling chunk %d: %v\n", i, err)
				continue
			}

			resp, err := http.Post("http://localhost:8080/add_embedding", "application/json", strings.NewReader(string(jsonData)))
			if err != nil {
				fmt.Printf("Error sending chunk %d to server: %v\n", i, err)
				continue
			}

			if resp.StatusCode != 200 {
				fmt.Printf("Server returned error for chunk %d: Status %d\n", i, resp.StatusCode)
			} else {
				fmt.Printf("Successfully sent chunk %d\n", i)
			}

			resp.Body.Close()
			time.Sleep(100 * time.Millisecond) // Rate limiting
		}
	})

	c.OnHTML("a[href]", handleLink)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
		pageTitle = "" // Reset title for each new page
		links++
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*cs.toronto.edu",
		Parallelism: 1,
		Delay:       1 * time.Second,
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

	fmt.Println("All chunks have been sent to vector database in real-time.")
}

func handleLink(e *colly.HTMLElement) {
	href := e.Attr("href")
	absoluteURL := e.Request.AbsoluteURL(href)
	e.Request.Visit(absoluteURL)
}

func chunkTextByLength(text string, url string, title string) {
	if text == "" {
		return
	}

	for i := 0; i < len(text); {
		end := i + MaxChunkLength
		if end > len(text) {
			end = len(text)
		}

		chunkContent := text[i:end]
		allTextChunks = append(allTextChunks, PageChunk{Content: chunkContent, URL: url, Title: title})

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
