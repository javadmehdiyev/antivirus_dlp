package antivirus

import (
	"os"
)

type Orchestrator struct {
	client *HTTPClient
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		client: NewHTTPClient(),
	}
}

func (o *Orchestrator) RunAntivirusCheck(testFile, testURL, httpMethod string) *Result {
	fileContent, err := os.ReadFile(testFile)
	if err != nil {
		return &Result{
			IsVirusDetected: true,
			StatusText:      "Failed to read file: " + err.Error(),
		}
	}

	req := &CheckRequest{
		TestFile:   string(fileContent),
		TestURL:    testURL,
		HTTPMethod: httpMethod,
	}

	resp, err := o.client.SendRequest(req)
	return EvaluateResult(resp, err)
}

