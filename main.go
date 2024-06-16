package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set viewport to 1920x1080
	err := chromedp.Run(ctx, chromedp.EmulateViewport(1920, 1080))
	if err != nil {
		log.Fatalf("Failed to set viewport: %v", err)
	}

	// Create a directory to save the screenshots if it doesn't exist
	savePath := "screenshots"
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		err = os.Mkdir(savePath, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}

	// Define the timeframes and corresponding intervals
	timeframes := map[string]string{
		"1D": "1D",
		"1W": "1W",
		"1M": "1M",
	}

	// Start timer
	start := time.Now()
	log.Println("Starting the script")

	// Navigate to TradingView and take screenshots
	for label, interval := range timeframes {
		log.Printf("Processing interval: %s", label)
		var screenshotBuf []byte

		// Start timer for each interval
		intervalStart := time.Now()

		// Navigate to the TradingView page with the interval as a query parameter
		url := fmt.Sprintf(`https://www.tradingview.com/chart/?symbol=BSE:HDFCBANK&interval=%s`, interval)
		log.Printf("Navigating to URL: %s", url)
		err := chromedp.Run(ctx, chromedp.Navigate(url))
		if err != nil {
			log.Fatalf("Failed to navigate to TradingView: %v", err)
		}

		// Wait for the page to load
		log.Println("Waiting for the page to load")
		time.Sleep(5 * time.Second)

		// Check if the watchlist is open by evaluating the aria-pressed attribute
		var isWatchlistOpen bool
		log.Println("Checking if the watchlist is open")
		err = chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector('button[data-tooltip="Watchlist, details and news"]').getAttribute('aria-pressed') === 'true'`, &isWatchlistOpen))
		if err != nil {
			log.Fatalf("Failed to check watchlist status: %v", err)
		}

		if isWatchlistOpen {
			log.Println("Watchlist is open, attempting to close it")
			err = chromedp.Run(ctx, chromedp.Click(`button[data-tooltip="Watchlist, details and news"]`, chromedp.NodeVisible))
			if err != nil {
				log.Printf("Failed to close the watchlist: %v", err)
			} else {
				log.Println("Watchlist closed successfully")
			}
		} else {
			log.Println("Watchlist is not open, no action needed")
		}

		// Take a screenshot of the chart
		log.Println("Taking screenshot")
		err = chromedp.Run(ctx, chromedp.FullScreenshot(&screenshotBuf, 90))
		if err != nil {
			log.Fatalf("Failed to take screenshot: %v", err)
		}

		// Save the screenshot
		screenshotPath := filepath.Join(savePath, fmt.Sprintf("hdfc_stock_chart_tradingview_%s.png", label))
		err = os.WriteFile(screenshotPath, screenshotBuf, 0644)
		if err != nil {
			log.Fatalf("Failed to save screenshot: %v", err)
		}

		log.Printf("Screenshot saved to %s", screenshotPath)
		log.Printf("Time taken for %s interval: %v", label, time.Since(intervalStart))
	}

	log.Printf("Total time taken: %v", time.Since(start))
}
