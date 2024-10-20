package ownerearnings

import (
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

const URL string = "https://www.gurufocus.com/term/owner-earnings/{STOCK}"
const TAG string = `font[style="font-size: 24px; font-weight: 700; color: #337ab7"]`
const REGEX string = `\:\s*(.*?)\(`
const FUNCTION string = "Owner Earnings"
const RETRIES int = 1

// GetOwnerEarnings returns the owner earnings for each ticker
func GetOwnerEarnings(ticker string) map[string]string {
	result := make(map[string]string)
	regex := regexp.MustCompile(REGEX)

	// Create a Colly collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// Disable CSS and images for faster scraping
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	})

	// Handle the HTML response and scrape the required element
	c.OnHTML(TAG, func(e *colly.HTMLElement) {
		text := e.Text
		value := regex.FindStringSubmatch(text)
		if len(value) > 1 {
			result[ticker] = sanitize(value[1])
		} else {
			result[ticker] = "" // In case the value is not found, set empty
		}
	})

	// Handle request errors
	c.OnError(func(_ *colly.Response, err error) {
		log.Printf("Error occurred while scraping ticker %s: %v", ticker, err)
		result[ticker] = ""
	})

	// Replace the ticker in the URL and visit the page
	pageURL := strings.ReplaceAll(URL, "{STOCK}", ticker)

	var err error
	for attempts := 0; attempts < RETRIES; attempts++ {
		err = c.Visit(pageURL)
		if err == nil {
			break
		}
		log.Printf("Attempt %d failed for ticker %s: %v. Retrying...", attempts+1, ticker, err)
	}

	if err != nil {
		result[ticker] = ""
	}

	return result
}

// Helper function to sanitize return value. The return value is a string like $1,233.456 we want a 1233.456 value
func sanitize(value string) string {
	sanitizedValue := strings.TrimSpace(value)

	return sanitizedValue
}
