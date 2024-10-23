package ownerearnings

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/vue-dot-com/value_investing_screener/utils"
)

const URL string = "https://www.gurufocus.com/term/owner-earnings/{STOCK}"
const TAG string = `font[style="font-size: 24px; font-weight: 700; color: #337ab7"]`
const REGEX string = `\:\s*(.*?)\(`
const FUNCTION string = "Owner Earnings"
const RETRIES int = 1

// GetOwnerEarnings returns the owner earnings for each ticker
func GetOwnerEarnings(ticker string) map[string]string {

	regex := regexp.MustCompile(REGEX)

	parseHTML := func(e *colly.HTMLElement, result map[string]string) {
		text := e.Text
		value := regex.FindStringSubmatch(text)
		if len(value) > 1 {
			result[ticker] = sanitize(value[1])
		} else {
			result[ticker] = "" // In case the value is not found, set empty
		}
	}

	// Replace the ticker in the URL and visit the page
	pageURL := strings.ReplaceAll(URL, "{STOCK}", ticker)

	result := utils.Scraper[string](ticker, pageURL, TAG, parseHTML, RETRIES, "")

	return result
}

// Helper function to sanitize return value. The return value is a string like $1,233.456 we want a 1233.456 value
func sanitize(value string) string {
	sanitizedValue := strings.TrimSpace(value)

	return sanitizedValue
}
