package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vue-dot-com/value_investing_screener/models"
	"github.com/vue-dot-com/value_investing_screener/parsers/enterprisevalue"
	"github.com/vue-dot-com/value_investing_screener/parsers/growthnumbers"
	"github.com/vue-dot-com/value_investing_screener/parsers/ownerearnings"
	"github.com/vue-dot-com/value_investing_screener/parsers/price"
	"github.com/vue-dot-com/value_investing_screener/parsers/roic"
	"github.com/vue-dot-com/value_investing_screener/utils"
)

func main() {
	defer timer("main")()

	// Map to store results for each ticker
	tickerResults := make(map[string]models.TickerData)
	var mu sync.Mutex // Mutex to protect shared access to tickerResults
	// Initialize counter
	var counter int32
	// Read csv and merge results
	tickerInfoNasdaq, _ := readCSVFile("NASDAQ.csv")
	tickerInfoNyse, _ := readCSVFile("NYSE.csv")

	tickerInfo := mergeTickerInfoMaps(tickerInfoNasdaq, tickerInfoNyse)

	// List of tickers to scrape data for
	var tickers []string
	for ticker := range tickerInfo {
		tickers = append(tickers, ticker)
	}

	maxConcurrency := 20
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
			priceData := price.GetStockPrice(ticker, action)
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
	filePath := "Screener.csv"
	if err := saveToCSV(filePath, tickerResults); err != nil {
		log.Printf("Error saving to CSV: %v\n", err)
	} else {
		log.Printf("Data successfully saved to %s\n", filePath)
	}
}

func readCSVFile(filePath string) (map[string]models.TickerInfo, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Initialize the map to store the TickerInfo
	tickerInfoMap := make(map[string]models.TickerInfo)

	// Skip the header row
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Read the CSV records
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		// Extract the relevant fields from the record
		symbol := record[0]
		name := record[1]
		marketCap := record[5]
		country := record[6]
		ipoYear := record[7]
		sector := record[9]
		industry := record[10]

		// Create a TickerInfo object
		tickerInfo := models.TickerInfo{
			Name:      name,
			MarketCap: marketCap,
			Country:   country,
			IpoYear:   ipoYear,
			Sector:    sector,
			Industry:  industry,
		}

		// Add the TickerInfo to the map using the symbol as the key
		tickerInfoMap[symbol] = tickerInfo
	}

	return tickerInfoMap, nil
}

// Function to merge two maps
func mergeTickerInfoMaps(map1, map2 map[string]models.TickerInfo) map[string]models.TickerInfo {
	// Create a new map to store the merged result
	mergedMap := make(map[string]models.TickerInfo)

	// Copy all key-value pairs from map1 into mergedMap
	for key, value := range map1 {
		mergedMap[key] = value
	}

	// Copy key-value pairs from map2 into mergedMap
	for key, value := range map2 {
		mergedMap[key] = value // This will overwrite the key if it exists in both maps
	}

	// Delete any ticker that contains "^" since it might indicate preferred shares or trusts etc.
	// Delete any ticker that contains "warrant" since it might indicate warrants.
	// Delete any ticker that contains "preferred" since it might indicate preferred stock and underlying company is covered in main ticker.
	for symbol := range mergedMap {
		if strings.Contains(symbol, "^") {
			delete(mergedMap, symbol)
		}
		if strings.Contains(strings.ToLower(mergedMap[symbol].Name), "warrant") {
			delete(mergedMap, symbol)
		}
		if strings.Contains(strings.ToLower(mergedMap[symbol].Name), "preferred") {
			delete(mergedMap, symbol)
		}
	}

	return mergedMap
}

func saveToCSV(filePath string, data map[string]models.TickerData) error {
	// Create or open the CSV file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	header := []string{
		"Ticker", "Name", "IPO Year", "Country", "Sector", "Industry", "Market Cap", "Price", "Enterprise Value", "ROIC", "Owner Earnings", "Ten Cap", "RevenueGrowth10Y", "RevenueGrowth5Y", "EpsGrowth10Y", "EpsGrowth5Y", "EbitGrowth10Y", "EbitGrowth5Y", "EbitdaGrowth10Y", "EbitdaGrowth5Y", "FcfGrowth10Y", "FcfGrowth5Y", "DividendGrowth10Y", "DividendGrowth5Y", "BvGrowth10Y", "BvGrowth5Y", "StockPriceGrowth10Y", "StockPriceGrowth5Y"}

	if err := writer.Write(header); err != nil {
		return err
	}

	// Write the data for each ticker
	for ticker, td := range data {
		record := []string{
			ticker,
			td.Name,
			td.IpoYear,
			td.Country,
			td.Sector,
			td.Industry,
			td.MarketCap,
			td.LastPrice,
			td.EnterpriseValue,
			td.Roic,
			td.OwnerEarnings,
			td.TenCap,
			td.GrowthData.RevenueGrowth10Y,
			td.GrowthData.RevenueGrowth5Y,
			td.GrowthData.EpsGrowth10Y,
			td.GrowthData.EpsGrowth5Y,
			td.GrowthData.EbitGrowth10Y,
			td.GrowthData.EbitGrowth5Y,
			td.GrowthData.EbitdaGrowth10Y,
			td.GrowthData.EbitdaGrowth5Y,
			td.GrowthData.FcfGrowth10Y,
			td.GrowthData.FcfGrowth5Y,
			td.GrowthData.DividendGrowth10Y,
			td.GrowthData.DividendGrowth5Y,
			td.GrowthData.BvGrowth10Y,
			td.GrowthData.BvGrowth5Y,
			td.GrowthData.StockPriceGrowth10Y,
			td.GrowthData.StockPriceGrowth5Y,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", name, time.Since(start))
	}
}
