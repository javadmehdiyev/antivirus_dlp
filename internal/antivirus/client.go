package antivirus

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

	if req.HTTPMethod == "POST" || req.HTTPMethod == "PUT" {
		// Create multipart form-data with "file" field
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		
		fileField, err := writer.CreateFormFile("file", "test.txt")
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}
		
		_, err = io.WriteString(fileField, req.TestFile)
		if err != nil {
			return nil, fmt.Errorf("failed to write file content: %w", err)
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
	} else {
		httpReq, err = http.NewRequest(req.HTTPMethod, req.TestURL, nil)
	}

	if err != nil {
		return nil, err
	}

	return httpReq, nil
}

