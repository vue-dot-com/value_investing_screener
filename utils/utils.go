package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vue-dot-com/value_investing_screener/models"
)

func CalculateTenCap(ownerearnings string) string {

	tenCapString := ""

	if ownerearnings != "" {
		n, _ := strconv.ParseFloat(ownerearnings, 32)
		tenCap := n * 10
		tenCapString := fmt.Sprintf("%f", tenCap)
		return tenCapString
	}

	return tenCapString
}

func CalculateMoS(price, fairvalue string) string {

	marginOfSafetyString := ""

	if price != "" && fairvalue != "" {
		nPrice, _ := strconv.ParseFloat(price, 32)
		nFairValue, _ := strconv.ParseFloat(fairvalue, 32)
		marginOfSafety := (nFairValue - nPrice) / nFairValue * 100
		marginOfSafetyString := fmt.Sprintf("%f%%", marginOfSafety)
		return marginOfSafetyString
	}

	return marginOfSafetyString
}

func Timer(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", name, time.Since(start))
	}
}

// Function to merge two maps
func MergeTickerInfoMaps(map1, map2 map[string]models.TickerInfo) map[string]models.TickerInfo {
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

func ReadCSVFile(filePath string) (map[string]models.TickerInfo, error) {
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

func SaveToCSV(filePath string, data map[string]models.TickerData) error {
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
		"Ticker", "Name", "IPO Year", "Country", "Sector", "Industry", "Market Cap", "Price", "Fair Value", "Margin Of Safety", "Enterprise Value", "ROIC", "Owner Earnings", "Ten Cap", "RevenueGrowth10Y", "RevenueGrowth5Y", "EpsGrowth10Y", "EpsGrowth5Y", "EbitGrowth10Y", "EbitGrowth5Y", "EbitdaGrowth10Y", "EbitdaGrowth5Y", "FcfGrowth10Y", "FcfGrowth5Y", "DividendGrowth10Y", "DividendGrowth5Y", "BvGrowth10Y", "BvGrowth5Y", "StockPriceGrowth10Y", "StockPriceGrowth5Y"}

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
			td.FairValue,
			td.MarginOfSafety,
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
