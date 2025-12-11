package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

type SettingsResponse struct {
	Success bool         `json:"success"`
	Data    SettingsData `json:"data"`
}

type SettingsData struct {
	URLDLP       string `json:"url_dlp"`
	URLAntivirus string `json:"url_antivirus"`
}

func getDLPURL() string {
	dlpIP := getIp()
	fmt.Printf("DLP IP: %s\n", dlpIP)

	// Remove trailing slashes and construct proper URL
	dlpIP = strings.TrimSuffix(dlpIP, "/")
	settingsURL := dlpIP + ":8000/api/settings-agent"
	fmt.Printf("Settings request URL: %s\n", settingsURL)

	resp, err := http.Get(settingsURL)
	if err != nil {
		fmt.Printf("Error: Failed to get settings: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read response body: %v\n", err)
		os.Exit(1)
	}

	// Check if response is actually JSON
	if len(body) > 0 && body[0] != '{' && body[0] != '[' {
		fmt.Printf("Error: Server returned non-JSON response (status: %d)\n", resp.StatusCode)
		fmt.Printf("Response preview: %s\n", string(body[:min(200, len(body))]))
		os.Exit(1)
	}

	var settingsResp SettingsResponse
	if err := json.Unmarshal(body, &settingsResp); err != nil {
		fmt.Printf("Error: Failed to parse settings response: %v\n", err)
		fmt.Printf("Response status: %d\n", resp.StatusCode)
		fmt.Printf("Response body: %s\n", string(body))
		os.Exit(1)
	}

	dlpURL := settingsResp.Data.URLDLP
	fmt.Printf("DLP service URL: %s\n", dlpURL)
	return dlpURL
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func createCreditCardFile(fileName string) error {
	content := `John Doe
                12345
                $50000
                john.doe@example.com
                4532-1234-5678-9010
                123-45-6789
                01/15/1980
                `
	return os.WriteFile(fileName, []byte(content), 0644)
}
func createPassportFile(path string) error {
	content := `John Doe
                AB1234567
                01/15/1980
                USA
                `

	return os.WriteFile(path, []byte(content), 0644)
}
func createCSVFile(path string) error {
	content := "Card Number,Expiry,CVV,Name\n" +
		"4532-1234-5678-9010,12/27,123,John Doe\n"

	return os.WriteFile(path, []byte(content), 0644)
}
func createXLSXFile(path string) error {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	f.SetCellValue(sheet, "A1", "Passport Number")
	f.SetCellValue(sheet, "B1", "Name")
	f.SetCellValue(sheet, "C1", "DOB")

	f.SetCellValue(sheet, "A2", "AB1234567")
	f.SetCellValue(sheet, "B2", "John Doe")
	f.SetCellValue(sheet, "C2", "01/15/1980")

	return f.SaveAs(path)
}
func saveJsonDlpDashboardResult() string {
	dlpIP := getIp()
	fmt.Printf("Dlp IP: %s\n", dlpIP)

	// Remove trailing slashes and construct proper URL
	dlpIP = strings.TrimSuffix(dlpIP, "/")
	settingsURL := dlpIP + ":8000/api/dlp/get-data"
	fmt.Printf("Settings import URL: %s\n", settingsURL)

	fileContent, err := os.ReadFile("dlp_results.json")

	if err != nil {
		fmt.Printf("Failed to read file: " + err.Error())
	}
	resp, err := http.Post(settingsURL, "application/json", bytes.NewBuffer(fileContent))
	if err != nil {
		fmt.Printf("Error: Failed to get settings: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading response: " + err.Error()
	}

	return string(bodyBytes)
}
