package dlp

import (
	"os"
	"path/filepath"
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
