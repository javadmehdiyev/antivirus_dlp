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
	checkPath := flag.String("check-path", "/Users/togrul/Documents/redmantis/go/redmantisupdated/api/uploads", "Path to check if file was saved (optional)")
	jsonFile := flag.String("json", "antivirus_results.json", "Path to JSON file to store results")
	flag.Parse()

	if *testFile == "" || *testURL == "" {
		fmt.Println("Usage: antivirus -file <path> -url <url> [-method <HTTP_METHOD>] [-check-path <path>] [-json <json_file>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	orchestrator := antivirus.NewOrchestrator()
	result := orchestrator.RunAntivirusCheck(*testFile, *testURL, *httpMethod, *checkPath)

	// Save result to JSON file
	if err := orchestrator.SaveResultToJSON(result, *jsonFile); err != nil {
		fmt.Printf("Warning: Failed to save result to JSON: %v\n", err)
	}

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
		os.Exit(1)
	}
}
