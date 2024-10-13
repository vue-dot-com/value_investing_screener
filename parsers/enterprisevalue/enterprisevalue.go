package enterprisevalue

import (
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

const URL string = "https://www.gurufocus.com/term/enterprise-value/{STOCK}"
const TAG string = `font[style="font-size: 24px; font-weight: 700; color: #337ab7"]`
const REGEX string = `(\$\d+([,\.]\d+)?\d+([,\.]\d+)?k?)`
const RETRIES int = 1

func GetEnterpriseValue(ticker string) map[string]string {
	result := make(map[string]string)
	regex := regexp.MustCompile(REGEX)

	// Create a Colly collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	//c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: maxConcurrency})

	// Disable CSS and images for faster scraping
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	})

	// Handle the HTML response and scrape the required element
	c.OnHTML(TAG, func(e *colly.HTMLElement) {
		text := e.Text

		value := regex.FindString(text)

		if value != "" {
			result[ticker] = value
		} else {
			result[ticker] = "" // In case the value is not found, set empty
		}
	})

	// Handle request errors
	c.OnError(func(_ *colly.Response, err error) {
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
