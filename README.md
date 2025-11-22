# Finnhub Go Explorer

A Go program to explore the Finnhub Stock API using the official finnhub-go client library.

## Setup

1. Get a free API key from [Finnhub.io](https://finnhub.io/register)

2. Set your API key as an environment variable:

   **Windows (PowerShell):**
   ```powershell
   $env:FINNHUB_API_KEY="your_api_key_here"
   ```

   **Windows (Command Prompt):**
   ```cmd
   set FINNHUB_API_KEY=your_api_key_here
   ```

   **Linux/Mac:**
   ```bash
   export FINNHUB_API_KEY=your_api_key_here
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

## Run

```bash
go run main.go
```

## Features

The program demonstrates the following Finnhub API features:

- **Company Profile** - Get company information (name, country, industry, market cap)
- **Real-time Quotes** - Get current stock prices and daily metrics
- **Symbol Search** - Search for stock symbols
- **Company News** - Get recent news articles about a company
- **Market Status** - Check if markets are currently open

## API Documentation

- [Finnhub API Docs](https://finnhub.io/docs/api)
- [finnhub-go GitHub](https://github.com/Finnhub-Stock-API/finnhub-go)
