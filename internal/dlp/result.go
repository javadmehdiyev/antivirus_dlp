package dlp

import "fmt"

func EvaluateResult(resp *CheckResponse, err error) *Result {
	if err != nil {
		return &Result{
			IsDLPActive: true,
			StatusText:  fmt.Sprintf("DLP blocked request: %v", err),
			IP:          "",
			FileContent: "",
		}
	}

	return &Result{
		IsDLPActive: false,
		StatusText:  fmt.Sprintf("Request succeeded: %s", resp.StatusText),
		IP:          "",
		FileContent: "",
	}
}
