package main

import (
	"flag"
	"fmt"
	"os"

	"dlpagent/internal/dlp"
)

func main() {
	testFile := flag.String("file", "test_dlp_data.txt", "Path to test file")
	testURL := flag.String("url", "http://127.0.0.1:8080/antivirus/load-data", "Target URL for DLP check")
	httpMethod := flag.String("method", "POST", "HTTP method (GET, POST, etc.)")
	flag.Parse()

	if *testFile == "" || *testURL == "" {
		fmt.Println("Usage: dlp -file <path> -url <url> [-method <HTTP_METHOD>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	orchestrator := dlp.NewOrchestrator()
	result := orchestrator.RunDLPCheck(*testFile, *testURL, *httpMethod)

	fmt.Printf("DLP Active: %v\n", result.IsDLPActive)
	fmt.Printf("Status: %s\n", result.StatusText)

	if result.IsDLPActive {
		os.Exit(1)
	}
}
