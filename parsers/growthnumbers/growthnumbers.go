package growthnumbers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

const URL string = "https://www.gurufocus.com/stock/{STOCK}/financials"
const TAG string = `tbody[data-v-217973d7]`

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

func GrowthCatcher(browser *rod.Browser, pool rod.Pool[*rod.Page], ticker string) map[string]GrowthData {
	result := make(map[string]GrowthData)

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
		var grData GrowthData
		var err error
		for attempts := 0; attempts < 3; attempts++ {
			grData, err = scrapeGrowthData(*page, ticker)
			if err == nil {
				break
			}
			log.Printf("Attempt %d failed for ticker %s: %v. Retrying...", attempts+1, ticker, err)
			// time.Sleep(time.Duration(attempts+1) * time.Second) // Exponential backoff
		}

		if err != nil {
			log.Printf("Error processing ticker %s: %v", ticker, err)
			result[ticker] = GrowthData{} // Empty GrowthData in case of error
		} else {
			result[ticker] = grData
		}

	}

	job(ticker)
	pool.Cleanup(func(p **rod.Page) {
		(*p).MustClose() // Dereference **rod.Page to get *rod.Page and call MustClose
	})
	return result
}

func scrapeGrowthData(page *rod.Page, ticker string) (GrowthData, error) {
	grData := GrowthData{}

	pageURL := strings.ReplaceAll(URL, "{STOCK}", ticker)
	err := rod.Try(func() {
		page.Timeout(60 * time.Second).MustNavigate(pageURL).MustWaitLoad().MustElement(TAG)
	})
	if err != nil {
		return grData, fmt.Errorf("failed to navigate to %s: %v", pageURL, err)
	}

	tbody, err := page.Timeout(30 * time.Second).Element(TAG)
	if err != nil {
		return grData, fmt.Errorf("no matching element found for ticker %s: %v", ticker, err)
	}

	rows, err := tbody.Elements(`tr[data-v-217973d7]`)
	if err != nil {
		return grData, nil // No need to return error if tbody doesn't contain element we continue
	}

	for i, row := range rows {
		tdElements, err := row.Elements(`td[data-v-217973d7]`)
		if err != nil {
			log.Printf("Failed to find td elements in row: %v", err)
			continue
		}

		if len(tdElements) < 4 {
			log.Println("Not enough <td> elements found in row")
			continue
		}

		gr10Y, err := tdElements[2].Text()
		if err != nil {
			log.Printf("Failed to get text from first td: %v", err)
			gr10Y = ""
		}
		if gr10Y == "-" {
			gr10Y = ""
		}

		gr5Y, err := tdElements[3].Text()
		if err != nil {
			log.Printf("Failed to get text from second td: %v", err)
			gr5Y = ""
		}
		if gr5Y == "-" {
			gr5Y = ""
		}

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
	}

	return grData, nil
}
