package main

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/vue-dot-com/value_investing_screener/config"
	"github.com/vue-dot-com/value_investing_screener/models"
	"github.com/vue-dot-com/value_investing_screener/parsers/enterprisevalue"
	"github.com/vue-dot-com/value_investing_screener/parsers/growthnumbers"
	"github.com/vue-dot-com/value_investing_screener/parsers/ownerearnings"
	"github.com/vue-dot-com/value_investing_screener/parsers/price"
	"github.com/vue-dot-com/value_investing_screener/parsers/roic"
	"github.com/vue-dot-com/value_investing_screener/utils"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No configuration .env file found")
	}
}

func main() {
	defer utils.Timer("main")()

	// Load configuration variables
	conf := config.New()
	// Set python version. Default is python3.10
	pythonVersion := conf.Python
	// Set verbosity for python script execution
	verbose := conf.Verbose
	// Set maxConcurrency. Default is 20
	maxConcurrency := conf.MaxConcurrency
	// List of tickers selected from the environment variable
	selectedTicker := conf.Tickers

	// Map to store results for each ticker
	tickerResults := make(map[string]models.TickerData)
	var mu sync.Mutex // Mutex to protect shared access to tickerResults
	// Initialize counter
	var counter int32

	// Read csv and merge results
	tickerInfoNasdaq, _ := utils.ReadCSVFile("NASDAQ.csv")
	tickerInfoNyse, _ := utils.ReadCSVFile("NYSE.csv")

	// Merge ticker information and apply some filtering
	tickerInfo := utils.MergeTickerInfoMaps(tickerInfoNasdaq, tickerInfoNyse)

	// List of tickers to scrape data for
	var tickers []string
	for ticker := range tickerInfo {
		tickers = append(tickers, ticker)
	}

	// If configuration has specific tickers use these to find values
	if len(selectedTicker) > 0 {
		tickers = selectedTicker
	}

	// Create a semaphore with a buffer to limit concurrency (e.g., 5)
	semaphore := make(chan struct{}, maxConcurrency)

	// Use an outer WaitGroup to wait for all tickers routines to finish in this way we can increment the counter in the for loop
	var outerWg sync.WaitGroup
	outerWg.Add(len(tickers))

	for _, ticker := range tickers {

		// Use a WaitGroup to wait for all goroutines to finish
		var wg sync.WaitGroup
		wg.Add(6) // Add 6 for each ticker since you're spawning 5 goroutines

		go func(ticker string) {
			defer wg.Done()

			mu.Lock()
			defer mu.Unlock()
			result := tickerResults[ticker]
			result.TickerInfo = tickerInfo[ticker]
			tickerResults[ticker] = result

		}(ticker)

		go func(ticker string) {
			defer wg.Done()

			// Acquire a slot in the semaphore
			semaphore <- struct{}{}
			defer func() {
				// Release the slot
				<-semaphore
			}()

			action := "fast_info[lastPrice]"
			priceData := price.GetStockPrice(ticker, action, pythonVersion, verbose)
			mu.Lock()
			defer mu.Unlock()
			result := tickerResults[ticker]

			result.LastPrice = priceData[ticker] // Update the last price
			tickerResults[ticker] = result       // Write back the updated struct

		}(ticker)

		go func(ticker string) {
			defer wg.Done()

			// Acquire a slot in the semaphore
			semaphore <- struct{}{}
			defer func() {
				// Release the slot
				<-semaphore
			}()

			// Fetch and update enterprise value data
			enterpriseValueData := enterprisevalue.GetEnterpriseValue(ticker)

			mu.Lock()
			defer mu.Unlock()
			result := tickerResults[ticker] // Retrieve the current value of the ticker

			result.EnterpriseValue = enterpriseValueData[ticker] // Update the ROIC
			tickerResults[ticker] = result                       // Write back the updated struct

		}(ticker)

		go func(ticker string) {
			defer wg.Done()

			// Acquire a slot in the semaphore
			semaphore <- struct{}{}
			defer func() {
				// Release the slot
				<-semaphore
			}()

			// Fetch and update ROIC data
			roicData := roic.GetRoic(ticker)

			mu.Lock()
			defer mu.Unlock()
			result := tickerResults[ticker] // Retrieve the current value of the ticker

			result.Roic = roicData[ticker] // Update the ROIC
			tickerResults[ticker] = result // Write back the updated struct
		}(ticker)

		go func(ticker string) {
			defer wg.Done()

			// Acquire a slot in the semaphore
			semaphore <- struct{}{}
			defer func() {
				// Release the slot
				<-semaphore
			}()

			// Fetch and update owner earnings data
			ownerEarningsData := ownerearnings.GetOwnerEarnings(ticker)

			mu.Lock()
			defer mu.Unlock()
			result := tickerResults[ticker] // Retrieve the current value of the ticker

			result.OwnerEarnings = ownerEarningsData[ticker]                 // Update the owner earnings
			result.TenCap = utils.CalculateTenCap(ownerEarningsData[ticker]) // Update the ten cap
			tickerResults[ticker] = result                                   // Write back the updated struct
		}(ticker)

		go func(ticker string) {
			defer wg.Done()

			// Acquire a slot in the semaphore
			semaphore <- struct{}{}
			defer func() {
				// Release the slot
				<-semaphore
			}()

			// Fetch and update growth data
			growthData := growthnumbers.GrowthCatcher(ticker)

			mu.Lock()
			defer mu.Unlock()
			result := tickerResults[ticker] // Retrieve the current value of the ticker

			result.GrowthData = growthData[ticker] // Update the growth data
			tickerResults[ticker] = result         // Write back the updated struct
		}(ticker)

		go func() {
			wg.Wait()
			atomic.AddInt32(&counter, 1)
			log.Printf("Processed %d/%d tickers\n", atomic.LoadInt32(&counter), len(tickers))
			// Mark this iteration as done in the outer WaitGroup
			outerWg.Done()
		}()

	}

	// Wait for all iterations (and their goroutines) to complete
	outerWg.Wait()

	// Save to CSV
	log.Print("Saving data to csv file")
	filePath := conf.OutputFile
	if err := utils.SaveToCSV(filePath, tickerResults); err != nil {
		log.Printf("Error saving to CSV: %v\n", err)
	} else {
		log.Printf("Data successfully saved to %s\n", filePath)
	}
}
