package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"stocks/api"
	"stocks/utils"
	"time"
)

func main() {
	args := os.Args[1:]

	if len(args) == 2 {
		switch args[0] {
		case "-refresh":
			getSymbolsAndSaveToFile(args[1])

		case "-fundamentals":
			getFundamentals(args[1])
		}
	}
}

func getSymbolsAndSaveToFile(exchange string) {
	symbols, err := api.GetSymbols(exchange)

	fileName := fmt.Sprintf("./symbols_%s.json", exchange)
	err = utils.WriteJSON(fileName, symbols)
	if err != nil {
		log.Println("Could not save symbols to file", err)
	}
}

func getFundamentals(exchange string) {
	fileName := fmt.Sprintf("./symbols_%s.json", exchange)
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	var symbols []api.StockSymbol
	err = json.Unmarshal(fileData, &symbols)
	if err != nil {
		fmt.Printf("Error deserializing JSON: %s\n", err)
		return
	}

	fmt.Printf("Found %d symbols. Fetching fundamental data...\n\n", len(symbols))

	var fundamentalData []api.FundamentalData
	for _, symbol := range symbols {
		time.Sleep(1100 * time.Millisecond)
		data, err := api.GetFundamentalData(symbol.DisplaySymbol)
		if err != nil {
			fmt.Printf("Error getting fundamental data for %s: %s\n", symbol.DisplaySymbol, err)
			continue
		}

		fundamentalData = append(fundamentalData, *data)
	}

	fileName = fmt.Sprintf("./fundamentals_%s.json", exchange)
	err = utils.WriteJSON(fileName, fundamentalData)
	if err != nil {
		fmt.Printf("Error writing to file: %s\n", err)
	}
}
