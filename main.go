package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"stocks/api"
	"stocks/utils"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	args := os.Args[1:]

	if len(args) == 2 {
		switch args[0] {
		case "-refresh":
			getSymbolsAndSaveToFile(args[1])

		case "-fundamentals":
			getFundamentals(args[1])

		case "-process-pd-ratio":
			processPbRatio(args[1])

		case "-value":
			getValueStocks(args[1])

		case "-test":
			numGoroutines, err := strconv.Atoi(args[1])
			if err != nil || numGoroutines < 1 {
				fmt.Printf("Invalid goroutine count. Usage: go run . -test <number>\n")
				return
			}
			testCPUWork(numGoroutines)
		}
	}
}

func getSymbolsAndSaveToFile(exchange string) {
	symbols, err := api.GetSymbols(exchange)

	fileName := fmt.Sprintf("./symbols_%s.json", exchange)
	err = utils.AppendJSON(fileName, symbols)
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

	outFileName := fmt.Sprintf("./fundamentals_%s.json", exchange)

	limiter := rate.NewLimiter(rate.Every(3*time.Second), 1)

	symbolChan := make(chan string, 100)
	resultChan := make(chan *api.FundamentalData, 100)

	ctx := context.Background()

	var workerWg sync.WaitGroup

	for i := 0; i < 3; i++ {
		workerWg.Add(1)
		go runWorker(symbolChan, resultChan, limiter, ctx, &workerWg)
	}

	writerDone := make(chan bool)
	go writeResults(outFileName, resultChan, writerDone)

	for _, symbol := range symbols {
		symbolChan <- symbol.DisplaySymbol
	}

	close(symbolChan)

	workerWg.Wait()

	close(resultChan)

	<-writerDone

	fmt.Printf("\nCompleted! Successfully processed %d symbols.\n", len(symbols))
}

func runWorker(symbolChan <-chan string, resultChan chan<- *api.FundamentalData, limiter *rate.Limiter, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	api.FetchWorker(symbolChan, resultChan, limiter, ctx)
}

func writeResults(fileName string, resultChan <-chan *api.FundamentalData, done chan<- bool) {
	for data := range resultChan {
		err := utils.AppendJSON(fileName, *data)
		if err != nil {
			log.Printf("Error writing to file: %v\n", err)
		}
	}

	done <- true
}

func getValueStocks(exchange string) {
	start := time.Now()

	fileName := fmt.Sprintf("./fundamentals_%s.json", exchange)
	fileData, err := utils.ReadFile[api.FundamentalData](fileName)

	if err != nil {
		log.Printf("Error read from file: %v\n", err)
		return
	}

	var filteredSymbols []api.FundamentalData

	for _, symbol := range fileData {
		if symbol.PBRatio > 0 && symbol.PBRatio < 1.5 &&
			symbol.PERatio > 0 && symbol.PERatio < 15 &&
			symbol.DebtToEquity < 1.0 &&
			symbol.ROE > 15 &&
			symbol.ProfitMargin > 20 &&
			symbol.CurrentPrice > 0 {
			filteredSymbols = append(filteredSymbols, symbol)
		}
	}

	fmt.Printf("%d out of %d passed filters\n", len(filteredSymbols), len(fileData))

	outFileName := fmt.Sprintf("./value_stocks_%s.json", exchange)
	jsonData, err := json.MarshalIndent(filteredSymbols, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	err = os.WriteFile(outFileName, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing to file: %v\n", err)
		return
	}

	elapsed := time.Since(start)

	fmt.Printf("Wrote %d symbols to %s in %v\n", len(filteredSymbols), outFileName, elapsed)
}

func processPbRatio(exchange string) {
	start := time.Now()

	defer func() {
		fmt.Printf("Time: %v\n", time.Since(start))
	}()

	fileName := fmt.Sprintf("./fundamentals_%s.json", exchange)

	var data, err = utils.ReadFile[api.FundamentalData](fileName)
	if err != nil {
		return
	}

	// Process all data and collect results
	results := make([]api.FundamentalData, len(data))

	dataCh := make(chan int, 100)
	var wg sync.WaitGroup

	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()

			for idx := range dataCh {
				d := data[idx]
				if d.BookValuePerShare != 0 {
					d.PBRatio = d.CurrentPrice / d.BookValuePerShare
				} else {
					d.PBRatio = 0
				}
				results[idx] = d
			}
		}()
	}

	// Send indices to workers
	for i := range data {
		dataCh <- i
	}
	close(dataCh)

	wg.Wait()

	// Write all results at once
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("Successfully processed and wrote %d records\n", len(results))
}

func testCPUWork(numGoroutines int) {
	const numTasks = 100000
	const workPerTask = 10000 // iterations of expensive calculation

	fmt.Printf("Testing with %d goroutines processing %d tasks...\n", numGoroutines, numTasks)

	start := time.Now()

	taskCh := make(chan int, 100)
	resultCh := make(chan float64, 100)
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for taskNum := range taskCh {
				result := 0.0
				for j := 0; j < workPerTask; j++ {
					x := float64(taskNum*j + 1)
					result += math.Sqrt(x)
					result += math.Sin(x) * math.Cos(x)
					result /= math.Log(x + 1)
				}

				resultCh <- result
			}
		}(i)
	}

	done := make(chan bool)
	totalResults := 0
	go func() {
		for range resultCh {
			totalResults++
		}

		done <- true
	}()

	for i := 0; i < numTasks; i++ {
		taskCh <- i
	}
	close(taskCh)

	wg.Wait()

	close(resultCh)

	<-done

	elapsed := time.Since(start)

	fmt.Printf("Completed %d tasks in %v\n", totalResults, elapsed)
	fmt.Printf("Average time per task: %v\n", elapsed/time.Duration(numTasks))
	fmt.Printf("Tasks per second: %.2f\n", float64(numTasks)/elapsed.Seconds())
}
