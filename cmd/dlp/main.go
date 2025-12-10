package main

import (
	"flag"
	"fmt"
	"os"

	"dlpagent/internal/dlp"
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

	if len(files) == 0 {
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

		if allExist {
			files = defaultFiles
		} else {
			createCreditCardFile("test_credit_card.txt")
			createPassportFile("test_passport.txt")
			createCSVFile("test_dlp_data.csv")
			createXLSXFile("test_dlp_data.xlsx")
			files = defaultFiles
		}
	}

	if *testURL == "" {
		fmt.Println("Usage: dlp -file <path> [-file <path> ...] -url <url> [-method <HTTP_METHOD>] [-json <json_file>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	orchestrator := dlp.NewOrchestrator()

	for i, file := range files {
		fmt.Printf("\n[%d/%d] Processing file: %s\n", i+1, len(files), file)

		result := orchestrator.RunDLPCheck(file, *testURL, *httpMethod)

		// Save result to JSON file
		if err := orchestrator.SaveResultToJSON(result, *jsonFile); err != nil {
			fmt.Printf("Warning: Failed to save result to JSON: %v\n", err)
		}

		fmt.Printf("DLP Active: %v\n", result.IsDLPActive)
		fmt.Printf("Status: %s\n", result.StatusText)

		if result.IsDLPActive {
			fmt.Printf("DLP detected in file: %s\n", file)
			os.Exit(1)
		}
	}

	fmt.Printf("\nAll files processed successfully. No DLP detected.\n")
}
