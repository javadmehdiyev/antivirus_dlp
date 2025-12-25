package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"dlpagent/internal/antivirus"
	"dlpagent/internal/dlp"
)

func main() {
	// Parse flags
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

	var hasError bool

	// Step 1: Run Antivirus Check
	if !*skipAntivirus {
		fmt.Println("=" + strings.Repeat("=", 60))
		fmt.Println("STEP 1: Running Antivirus Check")
		fmt.Println("=" + strings.Repeat("=", 60))

		settingUrl := getAntivirusURL()

		orchestrator := antivirus.NewOrchestrator()
		result := orchestrator.RunAntivirusCheck(settingUrl)

		// Save result to JSON file
		if err := orchestrator.SaveResultToJSON(result, *antivirusJsonFile); err != nil {
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
			hasError = true
		} else {
			fmt.Println("\n✅ Antivirus check PASSED")
		}
	} else {
		fmt.Println("Skipping antivirus check (--skip-antivirus flag set)")
	}

	// Step 2: Run DLP Check
	if !*skipDLP {
		fmt.Println("\n" + strings.Repeat("=", 62))
		fmt.Println("STEP 2: Running DLP Check")
		fmt.Println(strings.Repeat("=", 62))

		settingUrl := getDLPURL()
		if *dlpURL != "" {
			settingUrl = *dlpURL
		}

		if settingUrl == "" {
			fmt.Println("Error: DLP URL is required")
			os.Exit(1)
		}

		// Handle default files if none provided
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

		orchestrator := dlp.NewOrchestrator()

		for i, file := range files {
			fmt.Printf("\n[%d/%d] Processing file: %s\n", i+1, len(files), file)

			result := orchestrator.RunDLPCheck(file, settingUrl, *httpMethod)

			// Save result to JSON file
			if err := orchestrator.SaveResultToJSON(result, *dlpJsonFile, file); err != nil {
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
	} else {
		fmt.Println("Skipping DLP check (--skip-dlp flag set)")
	}

	// Final summary
	fmt.Println("\n" + strings.Repeat("=", 62))
	fmt.Println("SUMMARY")
	fmt.Println(strings.Repeat("=", 62))
	if hasError {
		fmt.Println("❌ Overall result: FAILED")
		os.Exit(1)
	} else {
		fmt.Println("✅ Overall result: PASSED - All checks completed successfully")
		os.Exit(0)
	}
}


