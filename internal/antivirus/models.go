package antivirus

import "time"

type CheckRequest struct {
	TestFile     string
	TestURL      string
	HTTPMethod   string
	SentFileName string // file_name that we're sending in the request
}

type CheckResponse struct {
	StatusCode int
	StatusText string
	FileName   string // file_name from server response
	Body       []byte // response body for GET requests (file content)
}

type Result struct {
	IsVirusDetected bool
	StatusText      string
	FileName        string // file_name for file check
	FileExists      bool   // whether file was found after check
	FilePath        string // path where file was checked
	IP              string // IP address of the computer sending the request
	FileContent     string // content of the file
}

// CheckResultEntry represents a single result entry stored in JSON
type CheckResultEntry struct {
	Timestamp       time.Time `json:"timestamp"`
	FileName        string    `json:"file_name"`
	StatusText      string    `json:"status_text"`
	IsVirusDetected bool      `json:"is_virus_detected"`
	FileExists      bool      `json:"file_exists"`
	FilePath        string    `json:"file_path"`
	IP              string    `json:"ip"`
	FileContent     string    `json:"file_content"`
}

// CheckResultsHistory stores the history of check results
type CheckResultsHistory struct {
	Results []CheckResultEntry `json:"results"`
}

// AntivirusAPIResponse represents the response from /api/antivirus endpoint
type AntivirusAPIResponse struct {
	Success bool             `json:"success"`
	Data    AntivirusAPIData `json:"data"`
}

// AntivirusAPIData represents the data field in the API response
type AntivirusAPIData struct {
	File        string `json:"file"`
	FileContent string `json:"file_content"`
	URL         string `json:"url"`
	Method      string `json:"method"`
	JSON        string `json:"json"`
}
