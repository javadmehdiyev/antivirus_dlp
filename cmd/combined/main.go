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

	"dlpagent/internal/antivirus"
	"dlpagent/internal/dlp"
)

var (
	checkIntervalDlp       time.Duration
	checkIntervalAntivirus time.Duration
)

func main() {
	var files []string
	flag.Func("file", "Path to test file for DLP (can be used multiple times)", func(s string) error {
		files = append(files, s)
		return nil
	})

	antivirusJsonFile := flag.String("antivirus-json", "antivirus_results.json", "Path to JSON file to store antivirus results")
	dlpJsonFile := flag.String("dlp-json", "dlp_results.json", "Path to JSON file to store DLP results")
	dlpURL := flag.String("dlp-url", "", "Target URL for DLP check (if not provided, will fetch from settings)")
	httpMethod := flag.String("method", "GET", "HTTP method (GET, POST, etc.)")
	skipAntivirus := flag.Bool("skip-antivirus", false, "Skip antivirus check")
	skipDLP := flag.Bool("skip-dlp", false, "Skip DLP check")
	flag.Parse()

	// Initialize intervals from settings
	checkIntervalDlp = time.Duration(getTimeOutDlp()) * time.Minute
	checkIntervalAntivirus = time.Duration(getTimeOutAntivirus()) * time.Minute

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start antivirus check goroutine
	if !*skipAntivirus {
		wg.Add(1)
		go runAntivirusCheck(ctx, &wg, *antivirusJsonFile, checkIntervalAntivirus)
	} else {
		fmt.Println("Skipping antivirus check (--skip-antivirus flag set)")
	}

	// Start DLP check goroutine
	if !*skipDLP {
		wg.Add(1)
		go runDLPCheck(ctx, &wg, *dlpJsonFile, *dlpURL, *httpMethod, files, checkIntervalDlp)
	} else {
		fmt.Println("Skipping DLP check (--skip-dlp flag set)")
	}

	// Wait for interrupt signal
	<-sigChan
	log.Println("Received interrupt signal, shutting down gracefully...")

	// Cancel context to stop all goroutines
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("Shutdown complete")
}

func runAntivirusCheck(ctx context.Context, wg *sync.WaitGroup, antivirusJsonFile string, interval time.Duration) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run immediately on start
	runAntivirusCheckOnce(antivirusJsonFile)

	// Then run on interval
	for {
		select {
		case <-ctx.Done():
			log.Println("Antivirus check goroutine stopping...")
			return
		case <-ticker.C:
			runAntivirusCheckOnce(antivirusJsonFile)
		}
	}
}

func runAntivirusCheckOnce(antivirusJsonFile string) {
	settingUrl := getAntivirusURL()

	orchestrator := antivirus.NewOrchestrator()
	result := orchestrator.RunAntivirusCheck(settingUrl)

	// Save result to JSON file
	if err := orchestrator.SaveResultToJSON(result, antivirusJsonFile); err != nil {
		fmt.Printf("Warning: Failed to save antivirus result to JSON: %v\n", err)
	}

	// send data to dashboard
	saveJsonAntivirusDashboardResult()

	fmt.Printf("Virus Detected: %v\n", result.IsVirusDetected)
	fmt.Printf("Status: %s\n", result.StatusText)
	if result.FileName != "" {
		fmt.Printf("File Name: %s\n", result.FileName)
	}
	if result.FilePath != "" {
		fmt.Printf("File Exists: %v\n", result.FileExists)
		fmt.Printf("File Path: %s\n", result.FilePath)
	}

	if result.IsVirusDetected {
		fmt.Println("\n❌ Antivirus check FAILED: Virus detected!")
	} else {
		fmt.Println("\n✅ Antivirus check PASSED")
	}
}

func runDLPCheck(ctx context.Context, wg *sync.WaitGroup, dlpJsonFile, dlpURL, httpMethod string, files []string, interval time.Duration) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Prepare files list
	dlpFiles := prepareDLPFiles(files)

	// Get DLP URL
	settingUrl := getDLPURL()
	if dlpURL != "" {
		settingUrl = dlpURL
	}

	if settingUrl == "" {
		log.Println("Error: DLP URL is required")
		return
	}

	// Run immediately on start
	runDLPCheckOnce(dlpJsonFile, settingUrl, httpMethod, dlpFiles)

	// Then run on interval
	for {
		select {
		case <-ctx.Done():
			log.Println("DLP check goroutine stopping...")
			return
		case <-ticker.C:
			runDLPCheckOnce(dlpJsonFile, settingUrl, httpMethod, dlpFiles)
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
			fmt.Printf("Warning: Failed to save DLP result to JSON: %v\n", err)
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
		fmt.Printf("\n✅ All DLP files processed successfully. No DLP detected.\n")
	}
}
