package dlp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Orchestrator struct {
	client *HTTPClient
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		client: NewHTTPClient(),
	}
}

func (o *Orchestrator) RunDLPCheck(testFile, testURL, httpMethod string) *Result {
	fileContent, err := os.ReadFile(testFile)
	if err != nil {
		return &Result{
			IsDLPActive: true,
			StatusText:  "Failed to read file: " + err.Error(),
		}
	}

	// Extract file extension
	fileExt := filepath.Ext(testFile)

	req := &CheckRequest{
		TestFile:      string(fileContent),
		TestURL:       testURL,
		HTTPMethod:    httpMethod,
		FileExtension: fileExt,
	}

	resp, err := o.client.SendRequest(req)
	return EvaluateResult(resp, err)
}

// getCategory determines the category based on file name
func getCategory(fileName string) string {
	baseName := strings.ToLower(filepath.Base(fileName))

	if strings.Contains(baseName, "credit_card") {
		return "credit_card"
	} else if strings.Contains(baseName, "passport") {
		return "passport_number"
	} else if strings.HasSuffix(baseName, ".csv") {
		return "file_upload_csv"
	} else if strings.HasSuffix(baseName, ".xlsx") {
		return "file_upload_xlsx"
	}

	return "unknown"
}

// SaveResultToJSON saves the result to JSON file, keeping only last 15 entries
func (o *Orchestrator) SaveResultToJSON(result *Result, jsonFilePath string, fileName string) error {
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
	category := getCategory(fileName)
	entry := CheckResultEntry{
		Timestamp:   time.Now(),
		StatusText:  result.StatusText,
		IsDLPActive: result.IsDLPActive,
		FileName:    filepath.Base(fileName),
		Category:    category,
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
