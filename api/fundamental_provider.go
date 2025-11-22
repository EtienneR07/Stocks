package api

import (
	"context"
	"log"
	"os"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/joho/godotenv"
)

type FundamentalData struct {
	Symbol            string
	CompanyName       string
	PERatio           float64
	EPS               float64
	RevenuePerShare   float64
	ProfitMargin      float64
	ROE               float64
	ROA               float64
	DebtToEquity      float64
	MarketCap         float64
	DividendYield     float64
	Beta              float64
	BookValuePerShare float64
	CurrentPrice      float64
}

func GetFundamentalData(symbol string) (*FundamentalData, error) {
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

	// Get company profile for name and market cap
	profile, _, err := client.CompanyProfile2(ctx).Symbol(symbol).Execute()
	if err != nil {
		return nil, err
	}

	// Get basic financials
	financials, _, err := client.CompanyBasicFinancials(ctx).Symbol(symbol).Metric("all").Execute()
	if err != nil {
		return nil, err
	}

	// Get current quote for price
	quote, _, err := client.Quote(ctx).Symbol(symbol).Execute()
	if err != nil {
		return nil, err
	}

	metric := financials.GetMetric()

	data := &FundamentalData{
		Symbol:            symbol,
		CompanyName:       profile.GetName(),
		MarketCap:         float64(profile.GetMarketCapitalization()),
		CurrentPrice:      float64(quote.GetC()),
		PERatio:           getMetricFloat(metric, "peBasicExclExtraTTM"),
		EPS:               getMetricFloat(metric, "epsExclExtraItemsTTM"),
		RevenuePerShare:   getMetricFloat(metric, "revenuePerShareTTM"),
		ProfitMargin:      getMetricFloat(metric, "netProfitMarginTTM"),
		ROE:               getMetricFloat(metric, "roeTTM"),
		ROA:               getMetricFloat(metric, "roaTTM"),
		DebtToEquity:      getMetricFloat(metric, "totalDebt/totalEquityQuarterly"),
		DividendYield:     getMetricFloat(metric, "currentDividendYieldTTM"),
		Beta:              getMetricFloat(metric, "beta"),
		BookValuePerShare: getMetricFloat(metric, "bookValuePerShareQuarterly"),
	}

	return data, nil
}

// Helper function to safely extract float values from metric map
func getMetricFloat(metric map[string]interface{}, key string) float64 {
	if val, ok := metric[key]; ok {
		if floatVal, ok := val.(float64); ok {
			return floatVal
		}
	}
	return 0.0
}
