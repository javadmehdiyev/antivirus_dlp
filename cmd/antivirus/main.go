package main

import (
	"flag"
	"fmt"
	"os"

	"dlpagent/internal/antivirus"
)

func main() {
	settingUrl := getAntivirusURL()

	jsonFile := flag.String("json", "antivirus_results.json", "Path to JSON file to store results")
	flag.Parse()

	orchestrator := antivirus.NewOrchestrator()
	result := orchestrator.RunAntivirusCheck(settingUrl)

	// Save result to JSON file
	if err := orchestrator.SaveResultToJSON(result, *jsonFile); err != nil {
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
		os.Exit(1)
	}
}
