package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	settingsURL := dlpIP + ":8000/api/settings-agent"

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

	var settingsResp SettingsResponse
	if err := json.Unmarshal(body, &settingsResp); err != nil {
		fmt.Printf("Error: Failed to parse settings response: %v\n", err)
		os.Exit(1)
	}

	dlpURL := settingsResp.Data.URLDLP
	fmt.Printf("Settings URL: %s\n", dlpURL)
	return dlpURL
}
