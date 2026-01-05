package dlp

import "time"

type CheckRequest struct {
	TestFile      string
	TestURL       string
	HTTPMethod    string
	FileExtension string // file extension to add to file_name
}

type CheckResponse struct {
	StatusCode int
	StatusText string
}

type Result struct {
	IsDLPActive bool
	StatusText  string
	IP          string // IP address of the computer sending the request
	FileContent string // content of the file
}

// CheckResultEntry represents a single result entry stored in JSON
type CheckResultEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	StatusText  string    `json:"status_text"`
	IsDLPActive bool      `json:"is_dlp_active"`
	FileName    string    `json:"file_name"`
	Category    string    `json:"category"`
	IP          string    `json:"ip"`
	FileContent string    `json:"file_content"`
}

// CheckResultsHistory stores the history of check results
type CheckResultsHistory struct {
	Results []CheckResultEntry `json:"results"`
}
