package antivirus

import (
	"bytes"
	"fmt"
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
		body := bytes.NewBuffer([]byte(req.TestFile))
		httpReq, err = http.NewRequest(req.HTTPMethod, req.TestURL, body)
	} else {
		httpReq, err = http.NewRequest(req.HTTPMethod, req.TestURL, nil)
	}

	if err != nil {
		return nil, err
	}

	return httpReq, nil
}

