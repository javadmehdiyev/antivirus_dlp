package main

import (
	"flag"
	"fmt"
	"os"

	"dlpagent/internal/dlp"
)

func main() {
	settingUrl := getDLPURL()

	testFile := flag.String("file", "test_dlp_data.txt", "Path to test file")
	testURL := flag.String("url", settingUrl, "Target URL for DLP check")
	httpMethod := flag.String("method", "POST", "HTTP method (GET, POST, etc.)")
	jsonFile := flag.String("json", "dlp_results.json", "Path to JSON file to store results")
	flag.Parse()

	if *testFile == "" || *testURL == "" {
		fmt.Println("Usage: dlp -file <path> -url <url> [-method <HTTP_METHOD>] [-json <json_file>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	orchestrator := dlp.NewOrchestrator()
	result := orchestrator.RunDLPCheck(*testFile, *testURL, *httpMethod)

	// Save result to JSON file
	if err := orchestrator.SaveResultToJSON(result, *jsonFile); err != nil {
		fmt.Printf("Warning: Failed to save result to JSON: %v\n", err)
	}

	fmt.Printf("DLP Active: %v\n", result.IsDLPActive)
	fmt.Printf("Status: %s\n", result.StatusText)

	if result.IsDLPActive {
		os.Exit(1)
	}
}
