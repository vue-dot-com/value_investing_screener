package enterprisevalue

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/vue-dot-com/value_investing_screener/utils"
)

const URL string = "https://www.gurufocus.com/term/enterprise-value/{STOCK}"
const TAG string = `font[style="font-size: 24px; font-weight: 700; color: #337ab7"]`
const REGEX string = `(\$\d+([,\.]\d+)?\d+([,\.]\d+)?k?)`
const RETRIES int = 1

func GetEnterpriseValue(ticker string) map[string]string {

	regex := regexp.MustCompile(REGEX)

	parseHTML := func(e *colly.HTMLElement, result map[string]string) {
		text := e.Text
		value := regex.FindString(text)
		if value != "" {
			result[ticker] = sanitize(value)
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
	sanitizedValue := strings.Trim(value, ",$")
	sanitizedValueFinal := strings.ReplaceAll(sanitizedValue, ",", "")

	return sanitizedValueFinal
}
