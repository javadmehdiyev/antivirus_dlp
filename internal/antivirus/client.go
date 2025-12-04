package antivirus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{},
	}
}

func (c *HTTPClient) SendRequest(req *CheckRequest) (*CheckResponse, error) {
	// Get file_name that we're sending in the request
	sentFileName := req.SentFileName

	httpReq, err := c.buildRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body to parse file_name
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try to parse JSON response for file_name
	checkResp := &CheckResponse{
		StatusCode: resp.StatusCode,
		StatusText: resp.Status,
		FileName:   sentFileName, // Use sent file_name as default
	}

	// Parse JSON if response is JSON
	var jsonResp map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonResp); err == nil {
		// Try different possible keys for file_name
		if fileName, ok := jsonResp["file_name"].(string); ok && fileName != "" {
			checkResp.FileName = fileName
		} else if fileName, ok := jsonResp["fileName"].(string); ok && fileName != "" {
			checkResp.FileName = fileName
		} else if fileName, ok := jsonResp["filename"].(string); ok && fileName != "" {
			checkResp.FileName = fileName
		}
		// If none found, use the one we sent
	}

	return checkResp, nil
}

func (c *HTTPClient) buildRequest(req *CheckRequest) (*http.Request, error) {
	var httpReq *http.Request
	var err error

	if req.HTTPMethod == "POST" || req.HTTPMethod == "PUT" {
		// Create multipart form-data with "file" and "file_name" fields
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add file field
		fileField, err := writer.CreateFormFile("file", "test.txt")
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}

		_, err = io.WriteString(fileField, req.TestFile)
		if err != nil {
			return nil, fmt.Errorf("failed to write file content: %w", err)
		}

		// Add file_name field with current date/time in format: 2025_11_22_13_00_45 (with seconds)
		fileName := time.Now().Format("2006_01_02_15_04_05")
		if req.SentFileName != "" {
			fileName = req.SentFileName
		}
		err = writer.WriteField("file_name", fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to write file_name field: %w", err)
		}

		// Store sent file_name in request for later use
		req.SentFileName = fileName

		err = writer.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close multipart writer: %w", err)
		}

		httpReq, err = http.NewRequest(req.HTTPMethod, req.TestURL, body)
		if err != nil {
			return nil, err
		}

		httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		httpReq, err = http.NewRequest(req.HTTPMethod, req.TestURL, nil)
	}

	if err != nil {
		return nil, err
	}

	return httpReq, nil
}
