package utils

import (
	"log"

	"github.com/gocolly/colly"
	"github.com/vue-dot-com/value_investing_screener/models"
)

// Scraper is a shared function that scrapes data based on a given URL, tag, and callback to handle HTML parsing.
func Scraper[T models.ScraperReturnDataType](ticker string, pageURL string, tag string, parseHTML func(e *colly.HTMLElement, result map[string]T), retries int, defaultReturn T) map[string]T {
	result := make(map[string]T)

	// Create Colly collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// Set headers for faster scraping
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	})

	// Handle the HTML response with the provided callback function
	c.OnHTML(tag, func(e *colly.HTMLElement) {
		parseHTML(e, result)
	})

	// Handle request errors
	c.OnError(func(_ *colly.Response, err error) {
		log.Printf("Error occurred while scraping ticker %s: %v", ticker, err)
		result[ticker] = defaultReturn
	})

	// Visit the page with retry logic
	var err error
	for attempts := 0; attempts < retries; attempts++ {
		err = c.Visit(pageURL)
		if err == nil {
			break
		}
		log.Printf("Attempt %d failed for ticker %s: %v. Retrying...", attempts+1, ticker, err)
	}

	if err != nil {
		result[ticker] = defaultReturn
	}

	return result
}
