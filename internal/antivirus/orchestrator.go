package antivirus

import (
	"encoding/json"
	"fmt"
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

func (o *Orchestrator) RunAntivirusCheck(testFile, testURL, httpMethod, checkPath string) *Result {
	fileContent, err := os.ReadFile(testFile)
	if err != nil {
		return &Result{
			IsVirusDetected: true,
			StatusText:      "Failed to read file: " + err.Error(),
		}
	}

	// Generate file_name before creating request
	fileName := time.Now().Format("2006_01_02_15_04_05")

	req := &CheckRequest{
		TestFile:     string(fileContent),
		TestURL:      testURL,
		HTTPMethod:   httpMethod,
		SentFileName: fileName,
	}

	resp, err := o.client.SendRequest(req)
	result := EvaluateResult(resp, err)

	fileExists := false
	filePath := ""

	// Use file_name from response or from request (we always send one)
	usedFileName := fileName
	if resp != nil && resp.FileName != "" {
		usedFileName = resp.FileName
	}

	// If request succeeded, save file_name and check file
	if !result.IsVirusDetected && resp != nil {
		// Save file_name to map
		if usedFileName != "" {
			o.mapMutex.Lock()
			o.fileMap[usedFileName] = true
			o.mapMutex.Unlock()

			result.FileName = usedFileName

			// Wait 5 seconds
			time.Sleep(5 * time.Second)

			// Check if file exists at checkPath with name file_name
			if checkPath != "" {
				filePath = filepath.Join(checkPath, usedFileName)
				if _, err := os.Stat(filePath); err == nil {
					fileExists = true
					result.StatusText = fmt.Sprintf("Request succeeded: %s. File saved: %s", resp.StatusText, filePath)
				} else {
					result.StatusText = fmt.Sprintf("Request succeeded: %s. File not found at: %s", resp.StatusText, filePath)
				}
			}
		}
		result.FileExists = fileExists
		result.FilePath = filePath
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
