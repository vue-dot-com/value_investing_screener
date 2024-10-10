package price

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

// Struct to hold the Python script's output
type PythonResult struct {
	Symbol string `json:"symbol"`
	Action string `json:"action"`
	Result string `json:"result,omitempty"` // Result field for valid output
	Error  string `json:"error,omitempty"`  // Error field for errors
}

const python string = "python3.10" // Adjust depending on your environment

func GetStockPrice(ticker string, action string) map[string]string {
	result := make(map[string]string)
	script := "parsers/price/price.py"

	cmd := exec.Command(python, script, ticker, action)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running Python script for ticker %s: %v\n", ticker, err)
		log.Printf("Python script stderr: %s\n", errBuf.String())

	}

	var priceResult PythonResult
	err = json.Unmarshal(outBuf.Bytes(), &priceResult)
	if err != nil {
		log.Printf("Error parsing JSON from stdout for ticker %s: %v\n", ticker, err)
		log.Printf("Output was: %s\n", outBuf.String()) // Print raw output for further debugging

	}

	if priceResult.Error != "" {
		log.Printf("Error for %s: %s\n", priceResult.Symbol, priceResult.Error)

	}

	result[ticker] = priceResult.Result
	// fmt.Printf("Result for %s (%s): %s\n", priceResult.Symbol, priceResult.Action, priceResult.Result)

	return result
}
