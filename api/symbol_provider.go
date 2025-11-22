package api

import (
	"context"
	"log"
	"os"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/joho/godotenv"
)

type StockSymbol struct {
	Description   string
	DisplaySymbol string
	MarketIdCode  string
}

func GetSymbols(exchange string) ([]StockSymbol, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("FINNHUB_API_KEY environment variable not set")
	}

	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", apiKey)
	client := finnhub.NewAPIClient(cfg).DefaultApi

	ctx := context.Background()

	symbols, _, err := client.StockSymbols(ctx).Exchange(exchange).Execute()

	if err != nil {
		log.Printf("Error getting symbol list: %v\n", err)
		return nil, err
	}

	symbolSummaries := make([]StockSymbol, 0)
	for _, symbol := range symbols {
		symbolSummaries = append(symbolSummaries, StockSymbol{
			DisplaySymbol: symbol.GetDisplaySymbol(),
			MarketIdCode:  symbol.GetMic(),
			Description:   symbol.GetDescription(),
		})
	}

	return symbolSummaries, nil
}
