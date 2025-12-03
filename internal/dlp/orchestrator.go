package dlp

type Orchestrator struct {
	client *HTTPClient
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		client: NewHTTPClient(),
	}
}

func (o *Orchestrator) RunDLPCheck(testFile, testURL, httpMethod string) *Result {
	req := &CheckRequest{
		TestFile:   testFile,
		TestURL:    testURL,
		HTTPMethod: httpMethod,
	}

	resp, err := o.client.SendRequest(req)
	return EvaluateResult(resp, err)
}
