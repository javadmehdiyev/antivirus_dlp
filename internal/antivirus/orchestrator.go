package antivirus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func (o *Orchestrator) RunAntivirusCheck(testFile, testURL, httpMethod, checkUrl string) *Result {
	fileContent, err := os.ReadFile(testFile)
	if err != nil {
		return &Result{
			IsVirusDetected: true,
			StatusText:      "Failed to read file: " + err.Error(),
		}
	}

	// Generate file_name before creating request with file extension
	fileExt := filepath.Ext(testFile)
	fileName := time.Now().Format("2006_01_02_15_04_05") + fileExt

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

			// Check if file exists via HTTP request to checkUrl
			if checkUrl != "" {
				// Parse URL and add file as query parameter
				parsedUrl, err := url.Parse(checkUrl)
				if err == nil {
					query := parsedUrl.Query()
					query.Set("file", usedFileName)
					parsedUrl.RawQuery = query.Encode()
					checkUrlWithFile := parsedUrl.String()
					filePath = checkUrlWithFile

					// Make HTTP GET request to check if file exists
					httpResp, err := http.Get(checkUrlWithFile)
					if err == nil {
						defer httpResp.Body.Close()

						// Read response body
						bodyBytes, err := io.ReadAll(httpResp.Body)
						if err == nil {
							// Parse JSON response
							var jsonResp map[string]interface{}
							if err := json.Unmarshal(bodyBytes, &jsonResp); err == nil {
								// Check if "exists" field is true
								if exists, ok := jsonResp["exists"].(bool); ok && exists {
									fileExists = true
									result.StatusText = fmt.Sprintf("Request succeeded: %s. File exists: %s", resp.StatusText, checkUrlWithFile)
								} else {
									result.StatusText = fmt.Sprintf("Request succeeded: %s. File not found at: %s", resp.StatusText, checkUrlWithFile)
								}
							} else {
								result.StatusText = fmt.Sprintf("Request succeeded: %s. Failed to parse check response: %s", resp.StatusText, checkUrlWithFile)
							}
						} else {
							result.StatusText = fmt.Sprintf("Request succeeded: %s. Failed to read check response: %s", resp.StatusText, checkUrlWithFile)
						}
					} else {
						result.StatusText = fmt.Sprintf("Request succeeded: %s. Failed to check file existence: %s", resp.StatusText, err.Error())
					}
				} else {
					result.StatusText = fmt.Sprintf("Request succeeded: %s. Failed to parse check URL: %s", resp.StatusText, err.Error())
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
