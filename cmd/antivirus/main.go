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
)

var (
	checkIntervalAntivirus time.Duration
)

func main() {
	jsonFile := flag.String("json", "antivirus_results.json", "Path to JSON file to store results")
	flag.Parse()

	// Initialize interval from settings
	checkIntervalAntivirus = time.Duration(getTimeOutAntivirus()) * time.Hour

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// WaitGroup to wait for goroutine to finish
	var wg sync.WaitGroup

	// Start antivirus check goroutine
	wg.Add(1)
	go runAntivirusCheck(ctx, &wg, *jsonFile, checkIntervalAntivirus)

	// Wait for interrupt signal
	<-sigChan
	log.Println("Received interrupt signal, shutting down gracefully...")

	// Cancel context to stop goroutine
	cancel()

	// Wait for goroutine to finish
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
		fmt.Printf("Warning: Failed to save result to JSON: %v\n", err)
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
