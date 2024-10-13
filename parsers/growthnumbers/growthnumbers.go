package growthnumbers

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

const URL string = "https://www.gurufocus.com/stock/{STOCK}/financials"
const TAG string = `tbody[data-v-217973d7]`
const RETRIES int = 1

type GrowthData struct {
	RevenueGrowth10Y    string
	RevenueGrowth5Y     string
	EpsGrowth10Y        string
	EpsGrowth5Y         string
	EbitGrowth10Y       string
	EbitGrowth5Y        string
	EbitdaGrowth10Y     string
	EbitdaGrowth5Y      string
	FcfGrowth10Y        string
	FcfGrowth5Y         string
	DividendGrowth10Y   string
	DividendGrowth5Y    string
	BvGrowth10Y         string
	BvGrowth5Y          string
	StockPriceGrowth10Y string
	StockPriceGrowth5Y  string
}

func GrowthCatcher(ticker string) map[string]GrowthData {
	result := make(map[string]GrowthData)
	grData := GrowthData{}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	})

	// Match the growth table
	c.OnHTML(TAG, func(e *colly.HTMLElement) {
		// Iterate over the table rows with tr[data-v-217973d7]
		e.DOM.Find(`tr[data-v-217973d7]`).Each(func(i int, row *goquery.Selection) {
			// Get all the td elements in the row
			tdElements := row.Find(`td[data-v-217973d7]`)

			// Check if there are at least 4 elements in the row
			if tdElements.Length() < 4 {
				log.Println("Not enough <td> elements found in row")
				return
			}

			// Extract the 10Y and 5Y growth data
			gr10Y := sanitize(tdElements.Eq(2).Text())
			gr5Y := sanitize(tdElements.Eq(3).Text())

			// Fill the appropriate growth fields based on the row index
			switch i {
			case 0:
				grData.RevenueGrowth10Y, grData.RevenueGrowth5Y = gr10Y, gr5Y
			case 1:
				grData.EpsGrowth10Y, grData.EpsGrowth5Y = gr10Y, gr5Y
			case 2:
				grData.EbitGrowth10Y, grData.EbitGrowth5Y = gr10Y, gr5Y
			case 3:
				grData.EbitdaGrowth10Y, grData.EbitdaGrowth5Y = gr10Y, gr5Y
			case 4:
				grData.FcfGrowth10Y, grData.FcfGrowth5Y = gr10Y, gr5Y
			case 5:
				grData.DividendGrowth10Y, grData.DividendGrowth5Y = gr10Y, gr5Y
			case 6:
				grData.BvGrowth10Y, grData.BvGrowth5Y = gr10Y, gr5Y
			case 7:
				grData.StockPriceGrowth10Y, grData.StockPriceGrowth5Y = gr10Y, gr5Y
			}
		})
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

	result[ticker] = grData

	if err != nil {
		result[ticker] = grData
	}

	return result
}

// Helper function to sanitize and handle "-" as empty values
func sanitize(value string) string {
	sanitizedValue := strings.TrimSpace(strings.Join(strings.Fields(value), " "))
	if sanitizedValue == "-" {
		return ""
	}
	return sanitizedValue
}
