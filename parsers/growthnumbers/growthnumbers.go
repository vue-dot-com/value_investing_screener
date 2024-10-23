package growthnumbers

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/vue-dot-com/value_investing_screener/models"
	"github.com/vue-dot-com/value_investing_screener/utils"
)

const URL string = "https://www.gurufocus.com/stock/{STOCK}/financials"
const TAG string = `tbody[data-v-217973d7]`
const RETRIES int = 1

func GrowthCatcher(ticker string) map[string]models.GrowthData {

	grData := models.GrowthData{}

	parseHTML := func(e *colly.HTMLElement, result map[string]models.GrowthData) {
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

			result[ticker] = grData
		})
	}

	// Replace the ticker in the URL and visit the page
	pageURL := strings.ReplaceAll(URL, "{STOCK}", ticker)

	result := utils.Scraper[models.GrowthData](ticker, pageURL, TAG, parseHTML, RETRIES, grData)

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
