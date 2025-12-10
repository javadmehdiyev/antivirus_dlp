package dlp

import (
	"bytes"
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

func (c *HTTPClient) SendRequest(req *CheckRequest) (*CheckResponse, error) { // bu gedecek EvaluateRequest funksiyasina
	httpReq, err := c.buildRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	return &CheckResponse{
		StatusCode: resp.StatusCode,
		StatusText: resp.Status,
	}, nil
}

func (c *HTTPClient) buildRequest(req *CheckRequest) (*http.Request, error) {
	var httpReq *http.Request
	var err error

	// Create multipart form-data with "file" and "file_name" fields for POST, PUT, and GET
	// Laravel can handle multipart form-data even for GET requests
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

	// Add file_name field with current date/time in format: 2025_11_22_13_00_45 (with seconds) + extension
	fileName := time.Now().Format("2006_01_02_15_04_05")
	if req.FileExtension != "" {
		fileName = fileName + req.FileExtension
	}
	err = writer.WriteField("file_name", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to write file_name field: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	httpReq, err = http.NewRequest(req.HTTPMethod, req.TestURL, body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	return httpReq, nil
}
