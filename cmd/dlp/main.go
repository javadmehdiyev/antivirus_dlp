package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"dlpagent/internal/dlp"
)

var (
	checkIntervalDlp time.Duration
)

func main() {
	settingUrl := getDLPURL()

	var files []string

	flag.Func("file", "Path to test file (can be used multiple times)", func(s string) error {
		files = append(files, s)
		return nil
	})

	testURL := flag.String("url", settingUrl, "Target URL for DLP check")
	httpMethod := flag.String("method", "GET", "HTTP method (GET, POST, etc.)")
	jsonFile := flag.String("json", "dlp_results.json", "Path to JSON file to store results")
	flag.Parse()

	// Initialize interval from settings
	checkIntervalDlp = time.Duration(getTimeOutDlp()) * time.Hour

	// Prepare files list
	dlpFiles := prepareDLPFiles(files)

	if *testURL == "" {
		fmt.Println("Usage: dlp -file <path> [-file <path> ...] -url <url> [-method <HTTP_METHOD>] [-json <json_file>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// WaitGroup to wait for goroutine to finish
	var wg sync.WaitGroup

	// Start DLP check goroutine
	wg.Add(1)
	go runDLPCheck(ctx, &wg, *jsonFile, *testURL, *httpMethod, dlpFiles, checkIntervalDlp)

	// Wait for interrupt signal
	<-sigChan
	log.Println("Received interrupt signal, shutting down gracefully...")

	// Cancel context to stop goroutine
	cancel()

	// Wait for goroutine to finish
	wg.Wait()
	log.Println("Shutdown complete")
}

func runDLPCheck(ctx context.Context, wg *sync.WaitGroup, dlpJsonFile, dlpURL, httpMethod string, files []string, interval time.Duration) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run immediately on start
	runDLPCheckOnce(dlpJsonFile, dlpURL, httpMethod, files)

	// Then run on interval
	for {
		select {
		case <-ctx.Done():
			log.Println("DLP check goroutine stopping...")
			return
		case <-ticker.C:
			runDLPCheckOnce(dlpJsonFile, dlpURL, httpMethod, files)
		}
	}
}

func prepareDLPFiles(files []string) []string {
	if len(files) > 0 {
		return files
	}

	defaultFiles := []string{
		"test_credit_card.txt",
		"test_passport.txt",
		"test_dlp_data.csv",
		"test_dlp_data.xlsx",
	}

	allExist := true
	for _, fileName := range defaultFiles {
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			allExist = false
			break
		}
	}

	if !allExist {
		createCreditCardFile("test_credit_card.txt")
		createPassportFile("test_passport.txt")
		createCSVFile("test_dlp_data.csv")
		createXLSXFile("test_dlp_data.xlsx")
	}

	return defaultFiles
}

func runDLPCheckOnce(dlpJsonFile, settingUrl, httpMethod string, files []string) {
	orchestrator := dlp.NewOrchestrator()
	var hasError bool

	for i, file := range files {
		fmt.Printf("\n[%d/%d] Processing file: %s\n", i+1, len(files), file)

		result := orchestrator.RunDLPCheck(file, settingUrl, httpMethod)

		// Save result to JSON file
		if err := orchestrator.SaveResultToJSON(result, dlpJsonFile, file); err != nil {
			fmt.Printf("Warning: Failed to save result to JSON: %v\n", err)
		}

		// send data to dashboard
		saveJsonDlpDashboardResult()

		fmt.Printf("DLP Active: %v\n", result.IsDLPActive)
		fmt.Printf("Status: %s\n", result.StatusText)

		if result.IsDLPActive {
			fmt.Printf("❌ DLP detected in file: %s\n", file)
			hasError = true
		}
	}

	if !hasError {
		fmt.Printf("\n✅ All files processed successfully. No DLP detected.\n")
	}
}
