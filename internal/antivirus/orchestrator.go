package antivirus

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Orchestrator struct {
	client   *HTTPClient
	fileMap  map[string]bool
	mapMutex sync.RWMutex
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		client:  NewHTTPClient(),
		fileMap: make(map[string]bool),
	}
}

// getLocalIP returns the local IP address of the machine
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func (o *Orchestrator) RunAntivirusCheck(settingUrl string) *Result {
	// Send GET request to http://127.0.0.1:8000/api/antivirus/download?type=file
	req := &CheckRequest{
		TestFile:     "", // No file content for GET
		TestURL:      settingUrl,
		HTTPMethod:   "GET",
		SentFileName: "",
	}

	resp, err := o.client.SendRequest(req)
	result := EvaluateResult(resp, err)

	fileExists := false
	var savedFilePath string
	var fileContent string

	// Get local IP address
	result.IP = getLocalIP()
	result.FileContent = "" // Initialize as empty, will be set if file is received

	// If request succeeded, save the file locally
	if !result.IsVirusDetected && resp != nil && len(resp.Body) > 0 {
		// Use file name from response or generate one with timestamp
		fileName := resp.FileName
		if fileName == "" {
			fileName = time.Now().Format("2006_01_02_15_04_05") + ".txt"
		}

		// Store file content
		fileContent = string(resp.Body)

		// Create uploads directory if it doesn't exist
		uploadsDir := "uploads"
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			return &Result{
				IsVirusDetected: true,
				StatusText:      "Failed to create uploads directory: " + err.Error(),
				IP:              getLocalIP(),
			}
		}

		// Save file to uploads directory
		savedFilePath = filepath.Join(uploadsDir, fileName)
		if err := os.WriteFile(savedFilePath, resp.Body, 0644); err != nil {
			return &Result{
				IsVirusDetected: true,
				StatusText:      "Failed to save file: " + err.Error(),
				IP:              getLocalIP(),
			}
		}

		result.FileName = fileName
		result.FilePath = savedFilePath
		result.FileContent = fileContent

		// Wait 5 seconds
		time.Sleep(5 * time.Second)

		// Check if file still exists
		if _, err := os.Stat(savedFilePath); err == nil {
			fileExists = true
			result.StatusText = fmt.Sprintf("Request succeeded: %s. File exists: %s", resp.StatusText, savedFilePath)
		} else {
			fileExists = false
			result.StatusText = fmt.Sprintf("Request succeeded: %s. File not found: %s", resp.StatusText, savedFilePath)
		}

		result.FileExists = fileExists
	} else if !result.IsVirusDetected && resp != nil {
		// Request succeeded but no file content
		result.FileExists = false
		result.StatusText = fmt.Sprintf("Request succeeded: %s. No file content received", resp.StatusText)
		result.FileContent = ""
	}

	return result
}

// SaveResultToJSON saves the result to JSON file, keeping only last 15 entries
func (o *Orchestrator) SaveResultToJSON(result *Result, jsonFilePath string) error {
	history := &CheckResultsHistory{
		Results: []CheckResultEntry{},
	}

	// Read existing results if file exists
	if _, err := os.Stat(jsonFilePath); err == nil {
		data, err := os.ReadFile(jsonFilePath)
		if err == nil {
			json.Unmarshal(data, history)
		}
	}

	// Create new entry
	entry := CheckResultEntry{
		Timestamp:       time.Now(),
		FileName:        result.FileName,
		StatusText:      result.StatusText,
		IsVirusDetected: result.IsVirusDetected,
		FileExists:      result.FileExists,
		FilePath:        result.FilePath,
		IP:              result.IP,
		FileContent:     result.FileContent,
	}

	// Add new entry
	history.Results = append(history.Results, entry)

	// Keep only last 15 entries
	if len(history.Results) > 15 {
		history.Results = history.Results[len(history.Results)-15:]
	}

	// Save to JSON file
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(jsonFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}
