package main

import (
	"flag"
	"fmt"
	"os"

	"dlpagent/internal/antivirus"
)

func main() {
	testFile := flag.String("file", "test_antivirus.js", "Path to test file")
	testURL := flag.String("url", "http://127.0.0.1:8080/antivirus/load-data", "Target URL for antivirus check")
	httpMethod := flag.String("method", "POST", "HTTP method (GET, POST, etc.)")
	flag.Parse()

	if *testFile == "" || *testURL == "" {
		fmt.Println("Usage: antivirus -file <path> -url <url> [-method <HTTP_METHOD>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	orchestrator := antivirus.NewOrchestrator()
	result := orchestrator.RunAntivirusCheck(*testFile, *testURL, *httpMethod)

	fmt.Printf("Virus Detected: %v\n", result.IsVirusDetected)
	fmt.Printf("Status: %s\n", result.StatusText)

	if result.IsVirusDetected {
		os.Exit(1)
	}
}
