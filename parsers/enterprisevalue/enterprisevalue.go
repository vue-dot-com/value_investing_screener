package enterprisevalue

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

const URL string = "https://www.gurufocus.com/term/enterprise-value/{STOCK}"
const TAG string = `font[style="font-size: 24px; font-weight: 700; color: #337ab7"]`
const REGEX string = `(\$\d+([,\.]\d+)?\d+([,\.]\d+)?k?)`
const FUNCTION string = "Enterprise Value"
const RETRIES int = 1

func GetEnterpriseValue(browser *rod.Browser, pool rod.Pool[*rod.Page], ticker string) map[string]string {
	result := make(map[string]string)
	regex := regexp.MustCompile(REGEX)

	create := func() **rod.Page {
		page := browser.MustIncognito().MustPage()

		// Disable CSS and images
		page.MustSetExtraHeaders(
			"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Accept-Encoding", "gzip, deflate, br",
		)
		page.MustEval(`() => {
			Object.defineProperty(window, 'Image', {
				get: () => function() { return {}; }
			});
			Object.defineProperty(window, 'matchMedia', {
				get: () => function() { return {matches: false, addEventListener: function(){}}; }
			});
		}`)

		// Set custom user agent
		page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		})

		return &page
	}

	job := func(ticker string) {
		page := pool.MustGet(create)
		defer pool.Put(page)
		var value string
		var err error
		for attempts := 0; attempts < RETRIES; attempts++ {
			value, err = scrapeEnterpriseValue(*page, ticker, regex)
			if err == nil {
				break
			}
			log.Printf("Attempt %d failed for ticker %s: %v. Retrying...", attempts+1, ticker, err)
			// time.Sleep(time.Duration(attempts+1) * time.Second) // Exponential backoff
		}

		if err != nil {
			result[ticker] = ""
		} else {
			result[ticker] = value
		}
	}

	job(ticker)
	pool.Cleanup(func(p **rod.Page) {
		(*p).MustClose() // Dereference **rod.Page to get *rod.Page and call MustClose
	})
	return result
}

func scrapeEnterpriseValue(page *rod.Page, ticker string, regex *regexp.Regexp) (string, error) {
	pageURL := strings.ReplaceAll(URL, "{STOCK}", ticker)
	err := rod.Try(func() {
		page.MustNavigate(pageURL).MustWaitLoad().MustWaitStable()
	})
	if err != nil {
		return "", fmt.Errorf("failed to navigate to %s: %v", pageURL, err)
	}

	element, err := page.Timeout(10 * time.Second).Element(TAG)
	if err != nil {
		return "", fmt.Errorf("no matching element found: %v", err)
	}

	text, err := element.Text()
	if err != nil {
		return "", fmt.Errorf("error getting text from element %v: %v", text, err)
	}

	value := regex.FindString(text)
	if value == "" {
		return "", nil
	}

	return value, nil
}
