package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"github.com/vue-dot-com/value_investing_screener/models"
	"github.com/vue-dot-com/value_investing_screener/parsers/enterprisevalue"
	"github.com/vue-dot-com/value_investing_screener/parsers/growthnumbers"
	"github.com/vue-dot-com/value_investing_screener/parsers/ownerearnings"
	"github.com/vue-dot-com/value_investing_screener/parsers/price"
	"github.com/vue-dot-com/value_investing_screener/parsers/roic"
)

func main() {
	// Map to store results for each ticker
	tickerResults := make(map[string]models.TickerData)

	// Read csv and merge results
	tickerInfoNasdaq, _ := readCSVFile("NASDAQ.csv")
	tickerInfoNyse, _ := readCSVFile("NYSE.csv")

	tickerInfo := mergeTickerInfoMaps(tickerInfoNasdaq, tickerInfoNyse)

	// List of tickers to scrape data for
	var tickers []string
	for ticker := range tickerInfo {
		tickers = append(tickers, ticker)
	}

	tickers = tickers[:5]

	log.Print("tickers are", tickers)

	// Launch the browser once for all scraping functions
	u := launcher.New().Headless(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	// Set max concurrency (e.g., 10 concurrent requests)
	maxConcurrency := 2

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Use channels to collect results from goroutines
	priceChan := make(chan map[string]string)
	enterpriseValueChan := make(chan map[string]string)
	roicChan := make(chan map[string]string)
	ownerEarningsChan := make(chan map[string]string)
	growthDataChan := make(chan map[string]growthnumbers.GrowthData)

	// Start goroutines for each scraping function
	wg.Add(5)

	go func() {
		wg.Wait()
		close(priceChan)
		close(enterpriseValueChan)
		close(roicChan)
		close(ownerEarningsChan)
		close(growthDataChan)
	}()

	go func() {
		defer wg.Done()
		log.Print("In Price catcher")
		action := "fast_info[lastPrice]"
		priceChan <- price.GetStockPrice(tickers, action, maxConcurrency)
	}()

	go func() {
		defer wg.Done()
		log.Print("In enterprise value catcher")
		enterpriseValueChan <- enterprisevalue.GetEnterpriseValue(browser, tickers, maxConcurrency)
	}()

	go func() {
		defer wg.Done()
		log.Print("In roic catcher")
		roicChan <- roic.GetRoic(browser, tickers, maxConcurrency)

	}()

	go func() {
		defer wg.Done()
		log.Print("In owner earnings catcher")
		ownerEarningsChan <- ownerearnings.GetOwnerEarnings(browser, tickers, maxConcurrency)

	}()

	go func() {
		defer wg.Done()
		log.Print("In growth data catcher")
		growthDataChan <- growthnumbers.GrowthCatcher(browser, tickers, maxConcurrency)
	}()

	// Collect results from channels
	price := <-priceChan
	enterpriseValue := <-enterpriseValueChan
	roicData := <-roicChan
	ownerEarnings := <-ownerEarningsChan
	growthData := <-growthDataChan

	// Populate the struct for each ticker
	for _, ticker := range tickers {
		tickerResults[ticker] = models.TickerData{
			TickerInfo:      tickerInfo[ticker],
			LastPrice:       price[ticker],
			EnterpriseValue: enterpriseValue[ticker],
			Roic:            roicData[ticker],
			OwnerEarnings:   ownerEarnings[ticker],
			GrowthData:      growthData[ticker],
		}
	}

	// Save to CSV
	log.Print("Saving data to csv file")
	filePath := "ticker_data.csv"
	if err := saveToCSV(filePath, tickerResults); err != nil {
		fmt.Printf("Error saving to CSV: %v\n", err)
	} else {
		fmt.Printf("Data successfully saved to %s\n", filePath)
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
		"Ticker",
		"Name",
		"IPO Year",
		"Country",
		"Sector",
		"Industry",
		"Price",
		"Enterprise Value",
		"ROIC",
		"Owner Earnings",
		"RevenueGrowth10Y",
		"RevenueGrowth5Y",
		"EpsGrowth10Y",
		"EpsGrowth5Y",
		"EbitGrowth10Y",
		"EbitGrowth5Y",
		"EbitdaGrowth10Y",
		"EbitdaGrowth5Y",
		"FcfGrowth10Y",
		"FcfGrowth5Y",
		"DividendGrowth10Y",
		"DividendGrowth5Y",
		"BvGrowth10Y",
		"BvGrowth5Y",
		"StockPriceGrowth10Y",
		"StockPriceGrowth5Y"}

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
			td.LastPrice,
			td.EnterpriseValue,
			td.Roic,
			td.OwnerEarnings,
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
