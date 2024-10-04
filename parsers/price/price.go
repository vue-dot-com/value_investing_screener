package price

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os/exec"
	"sync"
	"time"
)

// Struct to hold the Python script's output
type PythonResult struct {
	Symbol string `json:"symbol"`
	Action string `json:"action"`
	Result string `json:"result,omitempty"` // Result field for valid output
	Error  string `json:"error,omitempty"`  // Error field for errors
}

const python string = "python3.10" // Adjust depending on your environment

func GetStockPrice(tickers []string, action string, maxConcurrency int) map[string]string {
	prices := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	if maxConcurrency <= 0 {
		maxConcurrency = 1 // Ensure at least one concurrent execution
	}
	sem := make(chan struct{}, maxConcurrency)

	for _, ticker := range tickers {
		wg.Add(1)
		go func(ticker string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			log.Printf("Starting processing for ticker: %s\n", ticker)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// done := make(chan bool)
			result := make(chan string)

			go func() {
				cmd := exec.CommandContext(ctx, python, "parsers/price/price.py", ticker, action)

				var outBuf, errBuf bytes.Buffer
				cmd.Stdout = &outBuf
				cmd.Stderr = &errBuf

				err := cmd.Run()
				if err != nil {
					if ctx.Err() == context.DeadlineExceeded {
						log.Printf("Timeout running Python script for ticker %s\n", ticker)
					} else {
						log.Printf("Error running Python script for ticker %s: %v\n", ticker, err)
						log.Printf("Python script stderr: %s\n", errBuf.String())
					}
					result <- ""
					return
				}

				var priceResult PythonResult
				err = json.Unmarshal(outBuf.Bytes(), &priceResult)
				if err != nil {
					log.Printf("Error parsing JSON for ticker %s: %v\n", ticker, err)
					result <- ""
					return
				}

				if priceResult.Error != "" {
					log.Printf("Error for %s: %s\n", priceResult.Symbol, priceResult.Error)
					result <- ""
					return
				}

				result <- priceResult.Result
			}()

			select {
			case price := <-result:
				mu.Lock()
				prices[ticker] = price
				mu.Unlock()
				log.Printf("Processed price for ticker: %s, price: %s\n", ticker, price)
			case <-time.After(11 * time.Second):
				log.Printf("Forcefully terminating processing for ticker: %s\n", ticker)
				mu.Lock()
				prices[ticker] = "" // Set empty price for timed out ticker
				mu.Unlock()
				// Attempt to kill the process if it's still running
				if cmd := ctx.Value("cmd"); cmd != nil {
					if p := cmd.(*exec.Cmd).Process; p != nil {
						p.Kill()
					}
				}
			}

			cancel() // Ensure the context is cancelled
		}(ticker)
	}

	wg.Wait()
	return prices
}
