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
}

type Result struct {
	IsVirusDetected bool
	StatusText      string
	FileName        string // file_name for file check
	FileExists      bool   // whether file was found after check
	FilePath        string // path where file was checked
}

// CheckResultEntry represents a single result entry stored in JSON
type CheckResultEntry struct {
	Timestamp       time.Time `json:"timestamp"`
	FileName        string    `json:"file_name"`
	StatusText      string    `json:"status_text"`
	IsVirusDetected bool      `json:"is_virus_detected"`
	FileExists      bool      `json:"file_exists"`
	FilePath        string    `json:"file_path"`
}

// CheckResultsHistory stores the history of check results
type CheckResultsHistory struct {
	Results []CheckResultEntry `json:"results"`
}
