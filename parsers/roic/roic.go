package roic

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/vue-dot-com/value_investing_screener/utils"
)

const URL string = "https://www.gurufocus.com/term/roic/{STOCK}"
const TAG string = `font[style="font-size: 24px; font-weight: 700; color: #337ab7"]`
const REGEX string = `[+-]?(?:\d{1,3})(?:\.\d+)?%`
const RETRIES int = 1

// GetOwnerEarnings returns the owner earnings for each ticker
func GetRoic(ticker string) map[string]string {

	regex := regexp.MustCompile(REGEX)

	parseHTML := func(e *colly.HTMLElement, result map[string]string) {
		text := e.Text
		value := regex.FindString(text)
		if value != "" {
			result[ticker] = value
		} else {
			result[ticker] = "" // In case the value is not found, set empty
		}
	}

	// Replace the ticker in the URL and visit the page
	pageURL := strings.ReplaceAll(URL, "{STOCK}", ticker)

	result := utils.Scraper[string](ticker, pageURL, TAG, parseHTML, RETRIES, "")

	return result
}
