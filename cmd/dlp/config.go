package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
