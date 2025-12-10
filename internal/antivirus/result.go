package antivirus

import "fmt"

func EvaluateResult(resp *CheckResponse, err error) *Result {
	if err != nil {
		return &Result{
			IsVirusDetected: true,
			StatusText:      fmt.Sprintf("Antivirus blocked request: %v", err),
		}
	}

	return &Result{
		IsVirusDetected: false,
		StatusText:      fmt.Sprintf("Request succeeded: %s", resp.StatusText),
	}
}





